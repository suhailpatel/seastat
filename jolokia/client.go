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
		return "", fmt.Errorf("error whilst calling /version: %v", err)
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
		return nil, fmt.Errorf("error whilst calling /read for tables: %v", err)
	}

	tables := []Table{}
	fmt.Println(v.Get("value").String())
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

// MemoryStats returns memory information about the Java process
func (c *jolokiaClient) MemoryStats() (*MemoryStats, error) {
	v, err := c.read("java.lang", "type=Memory/*")
	if err != nil {
		return nil, fmt.Errorf("error whilst calling /read for memory stats: %v", err)
	}

	return &MemoryStats{
		HeapUsed:    v.Get("value", "HeapMemoryUsage", "used").GetInt64(),
		NonHeapUsed: v.Get("value", "NonHeapMemoryUsage", "used").GetInt64(),
	}, nil
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
