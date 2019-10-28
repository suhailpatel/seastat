package jolokia

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseDurationString(t *testing.T) {
	cases := []struct {
		in  string
		out time.Duration
	}{
		{in: "ns", out: time.Nanosecond},
		{in: "nsec", out: time.Nanosecond},
		{in: "μs", out: time.Microsecond},
		{in: "MiCROSecONDS", out: time.Microsecond},
		{in: "ms", out: time.Millisecond},
		{in: "MILLISECONDS", out: time.Millisecond},
		{in: "sec", out: time.Second},
		{in: "MINUTES", out: time.Minute},
		{in: "hRs", out: time.Hour},
	}

	for _, tc := range cases {
		assert.Equal(t, tc.out, parseDurationString(tc.in))
	}
}

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
