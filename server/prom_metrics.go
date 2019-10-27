package server

import "github.com/prometheus/client_golang/prometheus"

// ScrapeStats
var (
	PromScrapeTimestamp = prometheus.NewDesc(
		"seastat_last_scrape_timestamp",
		"Timestamp of the last scrape",
		[]string{}, nil,
	)

	PromScrapeDuration = prometheus.NewDesc(
		"seastat_last_scrape_duration_seconds",
		"Duration of the last scrape",
		[]string{}, nil,
	)
)

// CQLStats
var (
	PromCQLPreparedStatementsCount = prometheus.NewDesc(
		"seastat_cql_prepared_statements",
		"Number of prepared statements",
		[]string{}, nil,
	)

	PromCQLPreparedStatementsEvicted = prometheus.NewDesc(
		"seastat_cql_prepared_statements_evicted_total",
		"Number of evicted prepared statements",
		[]string{}, nil,
	)

	PromCQLPreparedStatementsExecuted = prometheus.NewDesc(
		"seastat_cql_prepared_statements_executed_total",
		"Number of executed prepared statements",
		[]string{}, nil,
	)

	PromCQLRegularStatementsExecuted = prometheus.NewDesc(
		"seastat_cql_regular_statements_executed_total",
		"Number of executed regular statements",
		[]string{}, nil,
	)

	PromCQLPreparedStatementsRatio = prometheus.NewDesc(
		"seastat_cql_prepared_statements_ratio",
		"Ratio of prepared statements",
		[]string{}, nil,
	)
)

// ThreadPoolStats
var (
	PromThreadPoolActiveTasks = prometheus.NewDesc(
		"seastat_thread_pool_active_tasks",
		"Number of active tasks in this thread pool",
		[]string{"name"}, nil,
	)

	PromThreadPoolPendingTasks = prometheus.NewDesc(
		"seastat_thread_pool_pending_tasks",
		"Number of pending tasks in this thread pool",
		[]string{"name"}, nil,
	)

	PromThreadPoolCompletedTasks = prometheus.NewDesc(
		"seastat_thread_pool_completed_tasks_total",
		"Number of completed tasks in this thread pool",
		[]string{"name"}, nil,
	)

	PromThreadPoolTotalBlockedTasks = prometheus.NewDesc(
		"seastat_thread_pool_blocked_tasks_total",
		"Number of total blocked tasks in this thread pool",
		[]string{"name"}, nil,
	)

	PromThreadPoolCurrentlyBlockedTasks = prometheus.NewDesc(
		"seastat_thread_pool_currently_blocked_tasks",
		"Number of currently blocked tasks in this thread pool",
		[]string{"name"}, nil,
	)

	PromThreadPoolMaxPoolSize = prometheus.NewDesc(
		"seastat_thread_pool_max_pool_size",
		"Largest thread pool size",
		[]string{"name"}, nil,
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
