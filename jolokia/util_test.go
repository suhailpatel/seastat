package jolokia

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractAttributes(t *testing.T) {
	attr1 := extractAttributes("org.apache.cassandra.metrics:keyspace=system,name=LiveDiskSpaceUsed,scope=IndexInfo,type=Table")
	assert.Equal(t, map[string]string{
		"keyspace": "system",
		"name":     "LiveDiskSpaceUsed",
		"scope":    "IndexInfo",
		"type":     "Table",
	}, attr1)

	attr2 := extractAttributes("keyspace=system,name=LiveDiskSpaceUsed,scope=IndexInfo,type=Table")
	assert.Equal(t, map[string]string{
		"keyspace": "system",
		"name":     "LiveDiskSpaceUsed",
		"scope":    "IndexInfo",
		"type":     "Table",
	}, attr2)

	attr3 := extractAttributes("org.apache.cassandra.metrics")
	assert.Equal(t, map[string]string{}, attr3)
}
