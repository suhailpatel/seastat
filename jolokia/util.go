package jolokia

import (
	"strings"
	"time"

	"github.com/valyala/fastjson"
)

// parseHistogram takes a histogram map and converts the various fields
// into a Histogram struct object
//
//    "StdDev": 0,
//    "75thPercentile": 0,
//    "Mean": null,
//    "98thPercentile": 0,
//    "Min": 0,
//    "95thPercentile": 0,
//    "99thPercentile": 0,
//    "Max": 0,
//    "999thPercentile": 0,
//    "Count": 5,
//    "50thPercentile": 0
//
func parseHistogram(val *fastjson.Value) Histogram {
	return Histogram{
		Minimum:       FloatGauge(val.Get("Min").GetFloat64()),
		Maximum:       FloatGauge(val.Get("Max").GetFloat64()),
		Percentile75:  FloatGauge(val.Get("75thPercentile").GetFloat64()),
		Percentile95:  FloatGauge(val.Get("95thPercentile").GetFloat64()),
		Percentile99:  FloatGauge(val.Get("99thPercentile").GetFloat64()),
		Percentile999: FloatGauge(val.Get("999thPercentile").GetFloat64()),
		Mean:          FloatGauge(val.Get("Mean").GetFloat64()),
		Count:         Counter(val.Get("Count").GetInt64()),
	}
}

// parseLatency takes a latency map and converts the various fields into
// a Latency struct object so it's easier to work with
//
//   "StdDev": 0,
//   "75thPercentile": 0,
//   "Mean": null,
//   "98thPercentile": 0,
//   "RateUnit": "events/second",
//   "95thPercentile": 0,
//   "99thPercentile": 0,
//   "Max": 0,
//   "Count": 7,
//   "FiveMinuteRate": 7.900892347061689e-10,
//   "50thPercentile": 0,
//   "MeanRate": 0.0008518263461541576,
//   "Min": 0,
//   "OneMinuteRate": 9.689141333518686e-39,
//   "DurationUnit": "microseconds",
//   "999thPercentile": 0,
//   "FifteenMinuteRate": 0.00002282562138178788
//
func parseLatency(val *fastjson.Value) Latency {
	durationUnitString := string(val.Get("DurationUnit").GetStringBytes())
	durationUnit := parseDurationString(durationUnitString)

	return Latency{
		Minimum:       time.Duration(val.Get("Min").GetFloat64()) * durationUnit,
		Maximum:       time.Duration(val.Get("Max").GetFloat64()) * durationUnit,
		Percentile75:  time.Duration(val.Get("75thPercentile").GetFloat64()) * durationUnit,
		Percentile95:  time.Duration(val.Get("95thPercentile").GetFloat64()) * durationUnit,
		Percentile99:  time.Duration(val.Get("99thPercentile").GetFloat64()) * durationUnit,
		Percentile999: time.Duration(val.Get("999thPercentile").GetFloat64()) * durationUnit,
		Mean:          time.Duration(val.Get("Mean").GetFloat64()) * durationUnit,
		Count:         Counter(val.Get("Count").GetInt64()),
	}
}

func parseDurationString(in string) time.Duration {
	switch strings.ToLower(in) {
	case "nanosecond", "nanoseconds", "ns", "nsec":
		return time.Nanosecond
	case "microsecond", "microseconds", "μs":
		return time.Microsecond
	case "millisecond", "milliseconds", "ms":
		return time.Millisecond
	case "second", "seconds", "sec", "secs":
		return time.Second
	case "minute", "minutes", "min", "mins":
		return time.Minute
	case "hour", "hours", "hr", "hrs":
		return time.Hour
	default:
		return time.Microsecond
	}
}

// extractAttributes takes in a string consisting of (possibly) a metric
// name followed by some key value pairs and turns that into a structured
// map of key value strings
//
//
// example: org.apache.cassandra.metrics:keyspace=system,name=LiveDiskSpaceUsed,scope=IndexInfo,type=Table
// turns into:
//   keyspace: "system"
//   name: 	   "LiveDiskSpaceUsed"
//   scope:    "IndexInfo"
//   type:     "Table"
//
func extractAttributes(tag string) map[string]string {
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

// valueToStringArray takes in an array of fastjson value types
// and converts the ones which are a string value to output an
// array of strings
func valueToStringArray(in []*fastjson.Value) []string {
	out := make([]string, 0, len(in))
	for _, val := range in {
		str := string(val.GetStringBytes())
		if str != "" {
			out = append(out, str)
		}
	}
	return out
}

// valueObjectToStringMap takes in a fastjson objects and extracts
// the keys and values as a map of string to string
func valueObjectToStringMap(in *fastjson.Object) map[string]string {
	out := make(map[string]string, 8)
	in.Visit(func(key []byte, val *fastjson.Value) {
		str := string(val.GetStringBytes())
		if len(key) > 0 && str != "" {
			out[string(key)] = string(val.GetStringBytes())
		}
	})
	return out
}
