package jolokia

import "time"

// Client embeds all the methods which can be called by a Jolokia client
// running alongside Cassandra
type Client interface {
	// Version gives the running agent version of Jolokia
	Version() (string, error)

	// Tables returns the list of tables from Cassandra
	Tables() ([]Table, error)

	// CQLStats returns info about the kinds of CQL statements being processed
	// and how many were prepared vs non-prepared. It also gives some insight
	// into the Prepared Statement cache
	CQLStats() (CQLStats, error)

	// ConnectedClients returns the number of connected clients via the
	// Native Protocol in Cassandra
	ConnectedClients() (int64, error)

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

// CQLStats embeds information about Prepared and Regular CQL statements
// to give insight of the kinds of queries hitting the cluster
type CQLStats struct {
	PreparedStatementsCount    int64
	PreparedStatementsEvicted  int64
	PreparedStatementsExecuted int64
	RegularStatementsExecuted  int64
	PreparedStatementsRatio    float64
}

// MemoryStats embeds information from the Java memory stats about how much
// heap and off-heap memory is being utilised
type MemoryStats struct {
	HeapUsed    int64 // bytes
	NonHeapUsed int64 // bytes
}

// GCStats embeds information for each type of GC that occurs in the process
type GCStats struct {
	Name        string
	Count       int64 // How many collections have occurred
	LastGC      time.Duration
	Accumulated time.Duration //
}
