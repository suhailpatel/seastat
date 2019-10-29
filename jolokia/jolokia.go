package jolokia

import "time"

// Types of metrics (mostly for documention)
type (
	// Counter represents a monotonically increasing value
	Counter int64
	// Gauge represents a point-in-time value
	Gauge int64
	// BytesGauge is a gauge representing bytes
	BytesGauge Gauge
	// FloatGauge represents a point-in-time (float) value
	FloatGauge float64
	// Histogram represents a Histogram distribution
	Histogram struct {
		Minimum       FloatGauge
		Maximum       FloatGauge
		Percentile75  FloatGauge
		Percentile95  FloatGauge
		Percentile99  FloatGauge
		Percentile999 FloatGauge
		Mean          FloatGauge
		Count         Counter
	}
	// Latency represents a latency distribution. It's similar
	// to a Histogram with the caveat that the fields are more
	// based around durations than gauge values
	Latency struct {
		Minimum       time.Duration
		Maximum       time.Duration
		Percentile75  time.Duration
		Percentile95  time.Duration
		Percentile99  time.Duration
		Percentile999 time.Duration
		Mean          time.Duration
		Count         Counter
	}
)

// Client embeds all the methods which can be called by a Jolokia client
// running alongside Cassandra
type Client interface {
	// Version gives the running agent version of Jolokia
	Version() (string, error)

	// Tables returns the list of tables from Cassandra
	Tables() ([]Table, error)

	// TableStats returns all the stats for a given Table from Cassandra
	TableStats(table Table) (TableStats, error)

	// CQLStats returns info about the kinds of CQL statements being processed
	// and how many were prepared vs non-prepared. It also gives some insight
	// into the Prepared Statement cache
	CQLStats() (CQLStats, error)

	// ThreadPoolStats returns info about each of the Thread Pools running
	// in Cassandra
	ThreadPoolStats() ([]ThreadPoolStats, error)

	// CompactionStats returns info about compactions which have happened
	// or are waiting in Cassandra
	CompactionStats() (CompactionStats, error)

	// ClientRequestStats returns info about client requests which happen
	// at the coordinator level
	ClientRequestStats() ([]ClientRequestStats, error)

	// ConnectedClients returns the number of connected clients via the
	// Native Protocol in Cassandra
	ConnectedClients() (Gauge, error)

	// MemoryStats returns memory information about the Java process
	MemoryStats() (MemoryStats, error)

	// GarbageCollectorStatus returns information about Garbage Collections
	// that occur in the process. Since there are different kinds of GC
	// processes occurring, the stats are returned as a list with an item for
	// each kind of GC step
	GarbageCollectionStats() ([]GCStats, error)
}

// Table embeds information about a Keyspace and Table that exists in
// Cassandra
type Table struct {
	KeyspaceName string
	TableName    string
}

// TableStats embeds all the stats associated with a table
type TableStats struct {
	Table Table

	// Latency stats
	CoordinatorRead  Latency
	CoordinatorWrite Latency
	CoordinatorScan  Latency
	ReadLatency      Latency
	WriteLatency     Latency
	RangeLatency     Latency

	// Table specific stats
	EstimatedPartitionCount  Gauge
	PendingCompactions       Gauge
	LiveDiskSpaceUsed        Gauge
	TotalDiskSpaceUsed       Gauge
	LiveSSTables             Gauge
	SSTablesPerRead          Histogram
	MaxPartitionSize         BytesGauge
	MeanPartitionSize        BytesGauge
	BloomFilterFalseRatio    FloatGauge
	TombstonesScanned        Histogram
	LiveCellsScanned         Histogram
	KeyCacheHitRate          FloatGauge
	PercentRepaired          FloatGauge
	SpeculativeRetries       Counter
	SpeculativeFailedRetries Counter
}

// CQLStats embeds stats about Prepared and Regular CQL statements
// to give insight of the kinds of queries hitting the cluster
type CQLStats struct {
	PreparedStatementsCount    Gauge
	PreparedStatementsEvicted  Counter
	PreparedStatementsExecuted Counter
	RegularStatementsExecuted  Counter
	PreparedStatementsRatio    FloatGauge
}

// ThreadPoolStats embeds stats for a type of Thread Pool
type ThreadPoolStats struct {
	PoolName              string
	ActiveTasks           Gauge
	PendingTasks          Gauge
	CompletedTasks        Counter
	TotalBlockedTasks     Counter
	CurrentlyBlockedTasks Counter
	MaxPoolSize           Gauge
}

// CompactionStats embeds stats for Compaction
type CompactionStats struct {
	BytesCompacted Counter
	PendingTasks   Gauge
	CompletedTasks Counter
}

// ClientRequestStats embeds stats for client requests
type ClientRequestStats struct {
	RequestType    string
	RequestLatency Latency
	Timeouts       Counter
	Failures       Counter
	Unavailables   Counter
}

// MemoryStats embeds stats about Java memory such as how much
// heap and off-heap memory is being utilised
type MemoryStats struct {
	HeapUsed    BytesGauge
	NonHeapUsed BytesGauge
}

// GCStats embeds information for each type of GC that occurs in the process
type GCStats struct {
	Name        string
	Count       Counter // How many collections have occurred
	LastGC      time.Duration
	Accumulated time.Duration
}
