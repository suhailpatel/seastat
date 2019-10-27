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

// // TableStats
// type TableStats struct {
// 	Table Table
//
// 	// Latency stats
// 	CoordinatorRead  Latency
// 	CoordinatorWrite Latency
// 	CoordinatorScan  Latency
// 	ReadLatency      Latency
// 	WriteLatency     Latency
// 	RangeLatency     Latency
// }

// TableStats
var (
	PromTableCoordinatorRead = prometheus.NewDesc(
		"seastat_table_coordinator_read_latency",
		"Coordinator table read latency",
		[]string{"keyspace", "table"}, nil,
	)

	PromTableCoordinatorWrite = prometheus.NewDesc(
		"seastat_table_coordinator_write_latency",
		"Coordinator table write latency",
		[]string{"keyspace", "table"}, nil,
	)

	PromTableCoordinatorRangeScan = prometheus.NewDesc(
		"seastat_table_coordinator_range_scan_latency",
		"Coordinator table range scan latency",
		[]string{"keyspace", "table"}, nil,
	)

	PromTableRead = prometheus.NewDesc(
		"seastat_table_read_latency",
		"Table read latency",
		[]string{"keyspace", "table"}, nil,
	)

	PromTableWrite = prometheus.NewDesc(
		"seastat_table_write_latency",
		"Table write latency",
		[]string{"keyspace", "table"}, nil,
	)

	PromTableRangeScan = prometheus.NewDesc(
		"seastat_table_range_scan_latency",
		"Table range scan latency",
		[]string{"keyspace", "table"}, nil,
	)

	PromTableEstimatedPartitionCount = prometheus.NewDesc(
		"seastat_table_estimated_partitions",
		"Number of partitions in this table (estimated)",
		[]string{"keyspace", "table"}, nil,
	)

	PromTablePendingCompactions = prometheus.NewDesc(
		"seastat_table_pending_compactions",
		"Number of pending compactions on this table",
		[]string{"keyspace", "table"}, nil,
	)

	PromTableMaxPartitionSize = prometheus.NewDesc(
		"seastat_table_max_partition_size_bytes",
		"Partition Size Max in bytes",
		[]string{"keyspace", "table"}, nil,
	)

	PromTableMeanPartitionSize = prometheus.NewDesc(
		"seastat_table_mean_partition_size_bytes",
		"Partition Size Mean in bytes",
		[]string{"keyspace", "table"}, nil,
	)

	PromTableBloomFilterFalseRatio = prometheus.NewDesc(
		"seastat_table_bloom_filter_false_ratio",
		"Percent of bloom filter false",
		[]string{"keyspace", "table"}, nil,
	)

	PromTableKeyCacheHitRate = prometheus.NewDesc(
		"seastat_table_key_cache_hit_percent",
		"Percent of key cache hits",
		[]string{"keyspace", "table"}, nil,
	)

	PromTablePercentRepaired = prometheus.NewDesc(
		"seastat_table_repaired_percent",
		"Percent of table repaired",
		[]string{"keyspace", "table"}, nil,
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
