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

// TableStats
var (
	PromTableCoordinatorRead = prometheus.NewDesc(
		"seastat_table_coordinator_read_latency_seconds",
		"Coordinator table read latency",
		[]string{"keyspace", "table"}, nil,
	)

	PromTableCoordinatorWrite = prometheus.NewDesc(
		"seastat_table_coordinator_write_latency_seconds",
		"Coordinator table write latency",
		[]string{"keyspace", "table"}, nil,
	)

	PromTableCoordinatorRangeScan = prometheus.NewDesc(
		"seastat_table_coordinator_range_scan_latency_seconds",
		"Coordinator table range scan latency",
		[]string{"keyspace", "table"}, nil,
	)

	PromTableRead = prometheus.NewDesc(
		"seastat_table_read_latency_seconds",
		"Table read latency",
		[]string{"keyspace", "table"}, nil,
	)

	PromTableWrite = prometheus.NewDesc(
		"seastat_table_write_latency_seconds",
		"Table write latency",
		[]string{"keyspace", "table"}, nil,
	)

	PromTableRangeScan = prometheus.NewDesc(
		"seastat_table_range_scan_latency_seconds",
		"Table range scan latency",
		[]string{"keyspace", "table"}, nil,
	)

	PromTableCASPropose = prometheus.NewDesc(
		"seastat_table_cas_propose_latency_seconds",
		"Paxos propose latency",
		[]string{"keyspace", "table"}, nil,
	)

	PromTableCASCommit = prometheus.NewDesc(
		"seastat_table_cas_commit_latency_seconds",
		"Paxos commit latency",
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

	PromTableLiveDiskSpaceUsed = prometheus.NewDesc(
		"seastat_table_live_disk_space_used_bytes",
		"Disk space used for live cells in bytes",
		[]string{"keyspace", "table"}, nil,
	)

	PromTableTotalDiskSpaceUsed = prometheus.NewDesc(
		"seastat_table_total_disk_space_used_bytes",
		"Disk space used for all data in bytes",
		[]string{"keyspace", "table"}, nil,
	)

	PromTableLiveSSTables = prometheus.NewDesc(
		"seastat_table_live_sstables",
		"Number of live SSTables",
		[]string{"keyspace", "table"}, nil,
	)

	PromTableSSTablesPerRead = prometheus.NewDesc(
		"seastat_table_sstables_per_read",
		"Number of SSTables consulted per read query",
		[]string{"keyspace", "table"}, nil,
	)

	PromTableMaxPartitionSize = prometheus.NewDesc(
		"seastat_table_max_partition_size_bytes",
		"Max Partition Size in bytes",
		[]string{"keyspace", "table"}, nil,
	)

	PromTableMeanPartitionSize = prometheus.NewDesc(
		"seastat_table_mean_partition_size_bytes",
		"Mean Partition Size in bytes",
		[]string{"keyspace", "table"}, nil,
	)

	PromTableBloomFilterFalseRatio = prometheus.NewDesc(
		"seastat_table_bloom_filter_false_ratio",
		"False positive ratio of table’s bloom filter",
		[]string{"keyspace", "table"}, nil,
	)

	PromTableTombstonesScanned = prometheus.NewDesc(
		"seastat_table_tombstones_scanned",
		"Number of tombstones scanned per read query",
		[]string{"keyspace", "table"}, nil,
	)

	PromTableLiveCellsScanned = prometheus.NewDesc(
		"seastat_table_live_cells_scanned",
		"Number of live cells scanned per read query",
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

	PromTableSpeculativeRetries = prometheus.NewDesc(
		"seastat_table_speculative_retries_total",
		"Total amount of speculative retries",
		[]string{"keyspace", "table"}, nil,
	)

	PromTableSpeculativeFailedRetries = prometheus.NewDesc(
		"seastat_table_speculative_failed_retries_total",
		"Total amount of speculative failed retries",
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

// CompactionStats
var (
	PromCompactionBytesCompacted = prometheus.NewDesc(
		"seastat_compaction_bytes_compacted_total",
		"Total amount of bytes compacted across all compactions",
		[]string{}, nil,
	)

	PromCompactionPendingTasks = prometheus.NewDesc(
		"seastat_compaction_pending_tasks",
		"Number of pending compaction tasks",
		[]string{}, nil,
	)

	PromCompactionCompletedTasks = prometheus.NewDesc(
		"seastat_compaction_completed_tasks_total",
		"Number of completed compaction tasks",
		[]string{}, nil,
	)
)

// ClientRequestStats
var (
	PromClientRequestLatency = prometheus.NewDesc(
		"seastat_client_request_latency_seconds",
		"Coordinator request latency",
		[]string{"request_type"}, nil,
	)

	PromClientRequestTimeouts = prometheus.NewDesc(
		"seastat_client_request_timeout_total",
		"Total number of coordinated request timeouts",
		[]string{"request_type"}, nil,
	)

	PromClientRequestFailures = prometheus.NewDesc(
		"seastat_client_request_failure_total",
		"Total number of coordinated request failures",
		[]string{"request_type"}, nil,
	)

	PromClientRequestUnavailable = prometheus.NewDesc(
		"seastat_client_request_unavailable_total",
		"Total number of coordinated request unavailable",
		[]string{"request_type"}, nil,
	)
)

// ConnectedClientStats
var (
	PromConnectedClients = prometheus.NewDesc(
		"seastat_connected_clients",
		"Number of connected clients",
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

// StorageStats
var (
	PromStorageKeyspaces = prometheus.NewDesc(
		"seastat_storage_keyspaces",
		"Number of keyspaces",
		[]string{}, nil,
	)

	PromStorageTokens = prometheus.NewDesc(
		"seastat_storage_tokens",
		"Number of tokens",
		[]string{}, nil,
	)

	PromStorageNodeStatus = prometheus.NewDesc(
		"seastat_storage_node_status",
		"Status of the other nodes from Cassandra's point of view",
		[]string{"node", "status", "state"}, nil,
	)
)

var (
	PromStorageInternalExceptions = prometheus.NewDesc(
		"seastat_internal_exceptions",
		"Number of internal exceptions caught",
		[]string{}, nil,
	)

	PromTotalHintsInProgress = prometheus.NewDesc(
		"seastat_hints_in_progress",
		"Number of hints attempting to be handed off since Cassandra started",
		[]string{}, nil,
	)

	PromTotalHints = prometheus.NewDesc(
		"seastat_hints_total",
		"Number of hint messages written to this node since Cassandra started",
		[]string{}, nil,
	)
)
