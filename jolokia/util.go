package jolokia

import "strings"

// extractAttributeMap takes in a string consisting of (possibly) a metric
// name followed by some key value pairs and turns that into a structured
// map of strings
//
//
// example: org.apache.cassandra.metrics:keyspace=system,name=LiveDiskSpaceUsed,scope=IndexInfo,type=Table
// turns into:
//   keyspace: "system"
//   name: 	   "LiveDiskSpaceUsed"
//   scope:    "IndexInfo"
//   type:     "Table"
//
func extractAttributeMap(tag string) map[string]string {
	idx := strings.IndexByte(tag, ':')
	tag = tag[idx+1:]

	commaSplit := strings.Split(tag, ",")
	out := make(map[string]string, len(commaSplit))
	for _, pair := range commaSplit {
		kvSplit := strings.Split(pair, "=")
		if len(kvSplit) != 2 {
			continue
		}
		out[kvSplit[0]] = kvSplit[1]
	}
	return out
}
