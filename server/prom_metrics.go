package server

import "github.com/prometheus/client_golang/prometheus"

// Scrape Stats
var (
	PromScrapeTimestamp = prometheus.NewDesc(
		"seastat_last_scrape_timestamp",
		"Timestamp of the last scrape",
		[]string{}, nil,
	)

	PromScrapeDuration = prometheus.NewDesc(
		"seastat_last_scrape_duration",
		"Duration of the last scrape",
		[]string{}, nil,
	)
)

// MemoryStats
var (
	PromMemoryStatsHeapUsed = prometheus.NewDesc(
		"seastat_memory_heap_used_bytes",
		"Bytes representing the used memory heap size",
		[]string{}, nil,
	)

	PromMemoryStatsNonHeapUsed = prometheus.NewDesc(
		"seastat_memory_nonheap_used_bytes",
		"Bytes representing the used memory non-heap size",
		[]string{}, nil,
	)
)

// GCStats
var (
	PromGCStatsCountTotal = prometheus.NewDesc(
		"seastat_gc_total",
		"Total number of Garbage Collections",
		[]string{"name"}, nil,
	)

	PromGCStatsLastGC = prometheus.NewDesc(
		"seastat_gc_last_duration_seconds",
		"Duration of Last GC",
		[]string{"name"}, nil,
	)

	PromGCStatsAccumulatedGC = prometheus.NewDesc(
		"seastat_gc_accumulated_duration_seconds",
		"Accumulated durations of GC",
		[]string{"name"}, nil,
	)
)
