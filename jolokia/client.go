package jolokia

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/valyala/fastjson"
)

type jolokiaClient struct {
	endpoint   string
	httpClient *http.Client
}

var defaultHTTPClient = &http.Client{
	Timeout: 15 * time.Second,
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
	return v.Get("value", "agent").String(), nil
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
		attributes := extractAttributeMap(string(key))

		keyspace, _ := attributes["keyspace"]
		table, _ := attributes["scope"] // JMX exposes the table name as scope
		attributeType, _ := attributes["type"]

		if attributeType == "Table" && keyspace != "" && table != "" {
			tables = append(tables, Table{KeyspaceName: keyspace, TableName: table})
		}
	})
	return tables, nil
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
		attributes := extractAttributeMap(string(key))
		switch attributes["name"] {
		case "PreparedStatementsCount":
			stats.PreparedStatementsCount = val.Get("Count").GetInt64()
		case "PreparedStatementsEvicted":
			stats.PreparedStatementsEvicted = val.Get("Count").GetInt64()
		case "PreparedStatementsExecuted":
			stats.PreparedStatementsExecuted = val.Get("Count").GetInt64()
		case "RegularStatementsExecuted":
			stats.RegularStatementsExecuted = val.Get("Count").GetInt64()
		case "PreparedStatementsRatio":
			stats.PreparedStatementsRatio = val.Get("Value").GetFloat64()
		}
	})
	return stats, nil
}

// ConnectedClients returns the number of connected clients via the Native
// Protocol in Cassandra
func (c *jolokiaClient) ConnectedClients() (int64, error) {
	// We want to be very specific with our query here because otherwise we'll
	// get a list of all connected clients which might be huge if there are lots
	// of them!
	v, err := c.read("org.apache.cassandra.metrics", "type=Client", "name=connectedNativeClients")
	if err != nil {
		return 0, fmt.Errorf("err reading clients: %v", err)
	}
	return v.Get("value", "Value").GetInt64(), nil
}

// MemoryStats returns memory information about the Java process
func (c *jolokiaClient) MemoryStats() (MemoryStats, error) {
	v, err := c.read("java.lang", "type=Memory/*")
	if err != nil {
		return MemoryStats{}, fmt.Errorf("err reading memory stats: %v", err)
	}

	return MemoryStats{
		HeapUsed:    v.Get("value", "HeapMemoryUsage", "used").GetInt64(),
		NonHeapUsed: v.Get("value", "NonHeapMemoryUsage", "used").GetInt64(),
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
			Name:        val.Get("Name").String(),
			Count:       val.Get("CollectionCount").GetInt64(),
			LastGC:      time.Duration(val.Get("LastGcInfo", "duration").GetInt64()) * time.Millisecond,
			Accumulated: time.Duration(val.Get("CollectionTime").GetInt64()) * time.Millisecond,
		})
	})
	return stats, nil
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
