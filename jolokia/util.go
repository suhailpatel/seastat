package jolokia

import "strings"

// extractAttributeMap takes in a string consisting of (possibly) a metric
// name followed by some key value pairs and turns that into a structured
// map of strings
//
//
// example: org.apache.cassandra.metrics:keyspace=system,name=LiveDiskSpaceUsed,scope=IndexInfo,type=Table
// turns into:
//   keyspace: system
//   name: LiveDiskSpaceUsed
//   scope: IndexInfo
//   type: Table
//
func extractAttributeMap(tag string) map[string]string {
	idx := strings.IndexByte(tag, ':')
	tag = tag[idx+1:]

	split := strings.Split(tag, ",")
	m := make(map[string]string, len(split))
	for _, pair := range strings.Split(tag, ",") {
		items := strings.Split(pair, "=")
		if len(items) != 2 {
			continue
		}
		m[items[0]] = items[1]
	}
	return m
}
