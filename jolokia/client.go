package jolokia

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/valyala/fastjson"
)

type jolokiaClient struct {
	endpoint   string
	httpClient *http.Client
}

var defaultHTTPClient = &http.Client{
	Timeout: 3 * time.Second,
}

// Init initalizes and returns a Client ready for calls. The endpoint should
// consist of <protocol>://<host>:<port>. (example: http://localhost:8778)
func Init(endpoint string) Client {
	c := &jolokiaClient{
		endpoint:   endpoint,
		httpClient: defaultHTTPClient,
	}
	return c
}

// Version gives the running agent version of Jolokia
func (c *jolokiaClient) Version() (string, error) {
	v, err := c.get("/jolokia/version")
	if err != nil {
		return "", fmt.Errorf("err calling /version: %v", err)
	}
	return string(v.Get("value", "agent").GetStringBytes()), nil
}

// Tables gets the list of tables from Cassandra
func (c *jolokiaClient) Tables() ([]Table, error) {
	// We use LiveDiskSpaceUsed as a placeholder name so we aren't grabbing all
	// the table level metrics at once (because for a large cluster, that takes
	// a ton of time and CPU usage)
	v, err := c.read("org.apache.cassandra.metrics", "type=Table", "name=LiveDiskSpaceUsed", "*")
	if err != nil {
		return nil, fmt.Errorf("err reading tables: %v", err)
	}

	tables := []Table{}
	v.Get("value").GetObject().Visit(func(key []byte, _ *fastjson.Value) {
		attributes := extractAttributes(string(key))

		keyspace, _ := attributes["keyspace"]
		table, _ := attributes["scope"] // JMX exposes the table name as scope
		attributeType, _ := attributes["type"]

		if attributeType == "Table" && keyspace != "" && table != "" {
			tables = append(tables, Table{KeyspaceName: keyspace, TableName: table})
		}
	})
	return tables, nil
}

// TableStats gets all the stats for a given Table within Cassandra
func (c *jolokiaClient) TableStats(table Table) (TableStats, error) {
	metricItems := []string{
		"CoordinatorReadLatency",
		"CoordinatorWriteLatency",
		"CoordinatorScanLatency",
		"ReadLatency",
		"WriteLatency",
		"RangeLatency",
		"EstimatedPartitionCount",
		"PendingCompactions",
		"MaxPartitionSize",
		"MeanPartitionSize",
		"BloomFilterFalseRatio",
		"KeyCacheHitRate",
		"PercentRepaired",
	}

	mbeanGroups := make([][]string, 0, len(metricItems))
	for _, name := range metricItems {
		mbeanGroups = append(mbeanGroups, []string{
			"type=Table",
			fmt.Sprintf("keyspace=%s", table.KeyspaceName),
			fmt.Sprintf("scope=%s", table.TableName),
			fmt.Sprintf("name=%s", name),
		})
	}

	v, err := c.bulkRequest("org.apache.cassandra.metrics", mbeanGroups)
	if err != nil {
		return TableStats{}, fmt.Errorf("err reading table: %v", err)
	}

	stats := TableStats{Table: table}
	for _, item := range v.GetArray() {
		if item.Get("status").GetInt64() != http.StatusOK {
			continue
		}

		attributes := extractAttributes(string(item.Get("request", "mbean").GetStringBytes()))
		val := item.Get("value")
		switch attributes["name"] {
		// Latency stats
		case "CoordinatorReadLatency":
			stats.CoordinatorRead = parseLatency(val)
		case "CoordinatorWriteLatency":
			stats.CoordinatorWrite = parseLatency(val)
		case "CoordinatorScanLatency":
			stats.CoordinatorScan = parseLatency(val)
		case "ReadLatency":
			stats.ReadLatency = parseLatency(val)
		case "WriteLatency":
			stats.WriteLatency = parseLatency(val)
		case "RangeLatency":
			stats.RangeLatency = parseLatency(val)

		// Table specific stats
		case "EstimatedPartitionCount":
			stats.EstimatedPartitionCount = Gauge(val.Get("Value").GetInt64())
		case "PendingCompactions":
			stats.PendingCompactions = Gauge(val.Get("Value").GetInt64())
		case "MaxPartitionSize":
			stats.MaxPartitionSize = BytesGauge(val.Get("Value").GetInt64())
		case "MeanPartitionSize":
			stats.MeanPartitionSize = BytesGauge(val.Get("Value").GetInt64())
		case "BloomFilterFalseRatio":
			stats.BloomFilterFalseRatio = FloatGauge(val.Get("Value").GetInt64())
		case "KeyCacheHitRate":
			stats.KeyCacheHitRate = FloatGauge(val.Get("Value").GetInt64())
		case "PercentRepaired":
			stats.PercentRepaired = FloatGauge(val.Get("Value").GetInt64())
		}
	}
	return stats, nil
}

// CQLStats returns info about the kinds of CQL statements being processed and
// how many were prepared vs non-prepared. It also gives some insight into the
// Prepared Statement cache
func (c *jolokiaClient) CQLStats() (CQLStats, error) {
	v, err := c.read("org.apache.cassandra.metrics", "type=CQL", "name=*")
	if err != nil {
		return CQLStats{}, fmt.Errorf("err reading CQL stats: %v", err)
	}

	stats := CQLStats{}
	v.Get("value").GetObject().Visit(func(key []byte, val *fastjson.Value) {
		attributes := extractAttributes(string(key))
		switch attributes["name"] {
		case "PreparedStatementsCount":
			stats.PreparedStatementsCount = Gauge(val.Get("Count").GetInt64())
		case "PreparedStatementsEvicted":
			stats.PreparedStatementsEvicted = Counter(val.Get("Count").GetInt64())
		case "PreparedStatementsExecuted":
			stats.PreparedStatementsExecuted = Counter(val.Get("Count").GetInt64())
		case "RegularStatementsExecuted":
			stats.RegularStatementsExecuted = Counter(val.Get("Count").GetInt64())
		case "PreparedStatementsRatio":
			stats.PreparedStatementsRatio = FloatGauge(val.Get("Value").GetFloat64())
		}
	})
	return stats, nil
}

// ThreadPoolStats returns info about each of the Thread Pools running
// in Cassandra
func (c *jolokiaClient) ThreadPoolStats() ([]ThreadPoolStats, error) {
	v, err := c.read("org.apache.cassandra.metrics", "type=ThreadPools", "*")
	if err != nil {
		return []ThreadPoolStats{}, fmt.Errorf("err reading ThreadPool stats: %v", err)
	}

	// The structure of this response is slightly weird because is just a flat
	// list of stats, to keep on top of this, we use a map which we'll convert
	// to a list later on
	pools := map[string]*ThreadPoolStats{}
	v.Get("value").GetObject().Visit(func(key []byte, val *fastjson.Value) {
		attributes := extractAttributes(string(key))
		poolName := attributes["scope"] // pool name is embedded as scope
		pool, ok := pools[poolName]
		if !ok {
			pool = &ThreadPoolStats{PoolName: poolName}
			pools[poolName] = pool
		}

		switch attributes["name"] {
		case "ActiveTasks":
			pool.ActiveTasks = Gauge(val.Get("Value").GetInt64())
		case "PendingTasks":
			pool.PendingTasks = Gauge(val.Get("Value").GetInt64())
		case "CompletedTasks":
			// TODO(suhail): This feels like a Counter but has a value rather
			// than a Count which is odd?
			pool.CompletedTasks = Counter(val.Get("Value").GetInt64())
		case "TotalBlockedTasks":
			pool.TotalBlockedTasks = Counter(val.Get("Count").GetInt64())
		case "CurrentlyBlockedTasks":
			// TODO(suhail): This feels like a gauge but is exposed as a Counter
			pool.CurrentlyBlockedTasks = Counter(val.Get("Count").GetInt64())
		case "MaxPoolSize":
			pool.MaxPoolSize = Gauge(val.Get("Value").GetInt64())
		}
	})

	// We want this function to be determinstic output given two calls and
	// assuming the response from Jolokia is consistent. Thus, we sort our
	// pools in the output by Pool Name
	names := make([]string, 0, len(pools))
	for poolName := range pools {
		names = append(names, poolName)
	}
	sort.Strings(names)

	out := make([]ThreadPoolStats, 0, len(names))
	for _, poolName := range names {
		out = append(out, *pools[poolName])
	}
	return out, nil
}

// ConnectedClients returns the number of connected clients via the Native
// Protocol in Cassandra
func (c *jolokiaClient) ConnectedClients() (Gauge, error) {
	// We want to be very specific with our query here because otherwise we'll
	// get a list of all connected clients which might be huge if there are lots
	// of them!
	v, err := c.read("org.apache.cassandra.metrics", "type=Client", "name=connectedNativeClients")
	if err != nil {
		return 0, fmt.Errorf("err reading clients: %v", err)
	}
	return Gauge(v.Get("value", "Value").GetInt64()), nil
}

// MemoryStats returns memory information about the Java process
func (c *jolokiaClient) MemoryStats() (MemoryStats, error) {
	v, err := c.read("java.lang", "type=Memory/*")
	if err != nil {
		return MemoryStats{}, fmt.Errorf("err reading memory stats: %v", err)
	}

	return MemoryStats{
		HeapUsed:    BytesGauge(v.Get("value", "HeapMemoryUsage", "used").GetInt64()),
		NonHeapUsed: BytesGauge(v.Get("value", "NonHeapMemoryUsage", "used").GetInt64()),
	}, nil
}

// GarbageCollectorStatus returns information about Garbage Collections that
// occur in the process. Since there are different kinds of GC processes
// occurring, the stats are returned as a list with an item for each kind
// of GC step
func (c *jolokiaClient) GarbageCollectionStats() ([]GCStats, error) {
	v, err := c.read("java.lang", "type=GarbageCollector,*")
	if err != nil {
		return []GCStats{}, fmt.Errorf("err reading GC stats: %v", err)
	}

	stats := []GCStats{}
	v.Get("value").GetObject().Visit(func(_ []byte, val *fastjson.Value) {
		stats = append(stats, GCStats{
			Name:        string(val.Get("Name").GetStringBytes()),
			Count:       Counter(val.Get("CollectionCount").GetInt64()),
			LastGC:      time.Duration(val.Get("LastGcInfo", "duration").GetInt64()) * time.Millisecond,
			Accumulated: time.Duration(val.Get("CollectionTime").GetInt64()) * time.Millisecond,
		})
	})
	return stats, nil
}

// get makes a GET request to the targetPath and returns the contents of the
// body as a JSON value ready for items to be plucked. If any part of the
// request pipeline fails, an err is returned
func (c *jolokiaClient) get(targetPath string) (*fastjson.Value, error) {
	u, err := url.Parse(fmt.Sprintf("%v", c.endpoint))
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, targetPath)

	rsp, err := c.httpClient.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	// We do a quick sanity check to see if the response was OK. Note that
	// this isn't much use because Jolokia has a response code embedded in
	// the response body
	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected 200 OK, got %v", rsp.StatusCode)
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	var p fastjson.Parser
	v, err := p.ParseBytes(body)
	if err != nil {
		return nil, fmt.Errorf("error whilst decoding: %v", err)
	}

	// Jolokia helpfully gives a response code in the inner body which you
	// also need to check (the HTTP request might be a 200 OK but the
	// Jolokia response code might be a 404 for example)
	rspStatus := v.Get("status").GetInt()
	if rspStatus != http.StatusOK {
		return nil, fmt.Errorf("expected 200 response from Jolokia, got %v", rspStatus)
	}

	return v, nil
}

// bulkRequest does a Jolokia bulk request. You pass in a list of groups of
// mbeans (one per request). Responses are provided in order of mbeanGroups
// queried
func (c *jolokiaClient) bulkRequest(metricName string, mbeanGroups [][]string) (*fastjson.Value, error) {
	bodyBytes, err := buildBulkRequestBody(metricName, mbeanGroups)
	if err != nil {
		return nil, fmt.Errorf("could not build bulkRequest body: %v", err)
	}
	reader := bytes.NewReader(bodyBytes)

	u, err := url.Parse(fmt.Sprintf("%v", c.endpoint))
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, "/jolokia/read")

	rsp, err := c.httpClient.Post(u.String(), "application/json", reader)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	var p fastjson.Parser
	v, err := p.ParseBytes(body)
	if err != nil {
		return nil, fmt.Errorf("error whilst decoding: %v", err)
	}
	return v, nil
}

// read is a convinience method around get. It takes in a metric name and a
// series of key=value strings and constructs a query to /jolokia/read
func (c *jolokiaClient) read(metricName string, kv ...string) (*fastjson.Value, error) {
	var targetPath string
	if len(kv) == 0 {
		targetPath = metricName
	} else {
		targetPath = fmt.Sprintf("/jolokia/read/%v:%v", metricName, strings.Join(kv, ","))
	}
	return c.get(targetPath)
}

func buildBulkRequestBody(metricName string, mbeanGroups [][]string) ([]byte, error) {
	queries := make([]map[string]string, 0, len(mbeanGroups))
	for _, group := range mbeanGroups {
		mbean := fmt.Sprintf("%s:%s", metricName, strings.Join(group, ","))
		queries = append(queries, map[string]string{
			"type":  "read",
			"mbean": mbean,
		})
	}
	return json.Marshal(queries)
}
