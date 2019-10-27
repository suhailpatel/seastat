package jolokia

// Client embeds all the methods which can be called by a Jolokia client
// running alongside Cassandra
type Client interface {
	// Version gives the running agent version of Jolokia
	Version() (string, error)

	// Tables returns the list of tables from Cassandra
	Tables() ([]Table, error)

	// MemoryStats returns memory information about the Java process
	MemoryStats() (*MemoryStats, error)
}

// MemoryStats embeds information from the Java memory stats about how much
// heap and off-heap memory is being utilised
type MemoryStats struct {
	HeapUsed    int64
	NonHeapUsed int64
}

// Table embeds information about a Keyspace and Table that exists in
// Cassandra
type Table struct {
	KeyspaceName string
	TableName    string
}
