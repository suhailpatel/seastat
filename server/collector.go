package server

import (
	"github.com/prometheus/client_golang/prometheus"
)

// SeastatCollector is here to satisfy the Prometheus Collector interface
type SeastatCollector struct {
	scraper *Scraper
}

// NewSeastatCollector returns an initialized Prometheus collector
// ready for exporting Seastat metrics
func NewSeastatCollector(s *Scraper) prometheus.Collector {
	return &SeastatCollector{scraper: s}
}

// Describe sends the super-set of all possible descriptors of metrics
// collected by this Collector to the provided channel and returns once
// the last descriptor has been sent.
func (c *SeastatCollector) Describe(ch chan<- *prometheus.Desc) {
	descs := []*prometheus.Desc{
		// ScrapeStats
		PromScrapeTimestamp,
		PromScrapeDuration,

		// TableStats
		PromTableCoordinatorRead,
		PromTableCoordinatorWrite,
		PromTableCoordinatorRangeScan,
		PromTableRead,
		PromTableWrite,
		PromTableRangeScan,
		PromTableEstimatedPartitionCount,
		PromTablePendingCompactions,
		PromTableLiveDiskSpaceUsed,
		PromTableTotalDiskSpaceUsed,
		PromTableLiveSSTables,
		PromTableMaxPartitionSize,
		PromTableMeanPartitionSize,
		PromTableBloomFilterFalseRatio,
		PromTableTombstonesScanned,
		PromTableLiveCellsScanned,
		PromTableKeyCacheHitRate,
		PromTablePercentRepaired,
		PromTableSpeculativeRetries,
		PromTableSpeculativeFailedRetries,

		// CQLStats
		PromCQLPreparedStatementsCount,
		PromCQLPreparedStatementsEvicted,
		PromCQLPreparedStatementsExecuted,
		PromCQLRegularStatementsExecuted,
		PromCQLPreparedStatementsRatio,

		// ThreadPoolStats
		PromThreadPoolActiveTasks,
		PromThreadPoolPendingTasks,
		PromThreadPoolCompletedTasks,
		PromThreadPoolTotalBlockedTasks,
		PromThreadPoolCurrentlyBlockedTasks,
		PromThreadPoolMaxPoolSize,

		// CompactionStats
		PromCompactionBytesCompacted,
		PromCompactionPendingTasks,
		PromCompactionCompletedTasks,

		// ConnectedClients
		PromConnectedClients,

		// MemoryStats
		PromMemoryStatsHeapUsed,
		PromMemoryStatsNonHeapUsed,

		// GCStats
		PromGCStatsCountTotal,
		PromGCStatsLastGC,
		PromGCStatsAccumulatedGC,
	}

	for _, desc := range descs {
		ch <- desc
	}
}

// Collect is called by the Prometheus registry when collecting
// metrics. The implementation sends each collected metric via the
// provided channel and returns once the last metric has been sent.
func (c *SeastatCollector) Collect(ch chan<- prometheus.Metric) {
	metrics := c.scraper.Get()

	// ScrapeStats
	ch <- prometheus.MustNewConstMetric(PromScrapeTimestamp,
		prometheus.GaugeValue, float64(metrics.ScrapeTime.Unix()))
	ch <- prometheus.MustNewConstMetric(PromScrapeDuration,
		prometheus.GaugeValue, float64(metrics.ScrapeDuration.Seconds()))

	addTableStats(metrics, ch)
	addCQLStats(metrics, ch)
	addThreadPoolStats(metrics, ch)
	addCompactionStats(metrics, ch)
	addClientRequestStats(metrics, ch)
	addConnectedClientStats(metrics, ch)
	addMemoryStats(metrics, ch)
	addGCStats(metrics, ch)
}

func addTableStats(metrics ScrapedMetrics, ch chan<- prometheus.Metric) {
	// TableStats
	for _, stat := range metrics.TableStats {
		ch <- prometheus.MustNewConstSummary(PromTableCoordinatorRead,
			uint64(stat.CoordinatorRead.Count),
			float64(stat.CoordinatorRead.Count)*stat.CoordinatorRead.Mean.Seconds(),
			map[float64]float64{
				75.0: stat.CoordinatorRead.Percentile75.Seconds(),
				95.0: stat.CoordinatorRead.Percentile95.Seconds(),
				99.0: stat.CoordinatorRead.Percentile99.Seconds(),
				99.9: stat.CoordinatorRead.Percentile999.Seconds(),
			}, stat.Table.KeyspaceName, stat.Table.TableName)
		ch <- prometheus.MustNewConstSummary(PromTableCoordinatorWrite,
			uint64(stat.CoordinatorWrite.Count),
			float64(stat.CoordinatorWrite.Count)*stat.CoordinatorWrite.Mean.Seconds(),
			map[float64]float64{
				75.0: stat.CoordinatorWrite.Percentile75.Seconds(),
				95.0: stat.CoordinatorWrite.Percentile95.Seconds(),
				99.0: stat.CoordinatorWrite.Percentile99.Seconds(),
				99.9: stat.CoordinatorWrite.Percentile999.Seconds(),
			}, stat.Table.KeyspaceName, stat.Table.TableName)
		ch <- prometheus.MustNewConstSummary(PromTableCoordinatorRangeScan,
			uint64(stat.CoordinatorScan.Count),
			float64(stat.CoordinatorScan.Count)*stat.CoordinatorScan.Mean.Seconds(),
			map[float64]float64{
				75.0: stat.CoordinatorScan.Percentile75.Seconds(),
				95.0: stat.CoordinatorScan.Percentile95.Seconds(),
				99.0: stat.CoordinatorScan.Percentile99.Seconds(),
				99.9: stat.CoordinatorScan.Percentile999.Seconds(),
			}, stat.Table.KeyspaceName, stat.Table.TableName)
		ch <- prometheus.MustNewConstSummary(PromTableRead,
			uint64(stat.ReadLatency.Count),
			float64(stat.ReadLatency.Count)*stat.ReadLatency.Mean.Seconds(),
			map[float64]float64{
				75.0: stat.ReadLatency.Percentile75.Seconds(),
				95.0: stat.ReadLatency.Percentile95.Seconds(),
				99.0: stat.ReadLatency.Percentile99.Seconds(),
				99.9: stat.ReadLatency.Percentile999.Seconds(),
			}, stat.Table.KeyspaceName, stat.Table.TableName)
		ch <- prometheus.MustNewConstSummary(PromTableWrite,
			uint64(stat.WriteLatency.Count),
			float64(stat.WriteLatency.Count)*stat.WriteLatency.Mean.Seconds(),
			map[float64]float64{
				75.0: stat.WriteLatency.Percentile75.Seconds(),
				95.0: stat.WriteLatency.Percentile95.Seconds(),
				99.0: stat.WriteLatency.Percentile99.Seconds(),
				99.9: stat.WriteLatency.Percentile999.Seconds(),
			}, stat.Table.KeyspaceName, stat.Table.TableName)
		ch <- prometheus.MustNewConstSummary(PromTableRangeScan,
			uint64(stat.RangeLatency.Count),
			float64(stat.RangeLatency.Count)*stat.RangeLatency.Mean.Seconds(),
			map[float64]float64{
				75.0: stat.RangeLatency.Percentile75.Seconds(),
				95.0: stat.RangeLatency.Percentile95.Seconds(),
				99.0: stat.RangeLatency.Percentile99.Seconds(),
				99.9: stat.RangeLatency.Percentile999.Seconds(),
			}, stat.Table.KeyspaceName, stat.Table.TableName)

		ch <- prometheus.MustNewConstMetric(PromTableEstimatedPartitionCount,
			prometheus.GaugeValue, float64(stat.EstimatedPartitionCount),
			stat.Table.KeyspaceName, stat.Table.TableName)
		ch <- prometheus.MustNewConstMetric(PromTablePendingCompactions,
			prometheus.GaugeValue, float64(stat.PendingCompactions),
			stat.Table.KeyspaceName, stat.Table.TableName)
		ch <- prometheus.MustNewConstMetric(PromTableLiveDiskSpaceUsed,
			prometheus.GaugeValue, float64(stat.LiveDiskSpaceUsed),
			stat.Table.KeyspaceName, stat.Table.TableName)
		ch <- prometheus.MustNewConstMetric(PromTableTotalDiskSpaceUsed,
			prometheus.GaugeValue, float64(stat.TotalDiskSpaceUsed),
			stat.Table.KeyspaceName, stat.Table.TableName)
		ch <- prometheus.MustNewConstMetric(PromTableLiveSSTables,
			prometheus.GaugeValue, float64(stat.LiveSSTables),
			stat.Table.KeyspaceName, stat.Table.TableName)
		ch <- prometheus.MustNewConstSummary(PromTableSSTablesPerRead,
			uint64(stat.SSTablesPerRead.Count),
			float64(stat.SSTablesPerRead.Count)*float64(stat.SSTablesPerRead.Mean),
			map[float64]float64{
				75.0: float64(stat.SSTablesPerRead.Percentile75),
				95.0: float64(stat.SSTablesPerRead.Percentile95),
				99.0: float64(stat.SSTablesPerRead.Percentile99),
				99.9: float64(stat.SSTablesPerRead.Percentile999),
			}, stat.Table.KeyspaceName, stat.Table.TableName)
		ch <- prometheus.MustNewConstMetric(PromTableMaxPartitionSize,
			prometheus.GaugeValue, float64(stat.MaxPartitionSize),
			stat.Table.KeyspaceName, stat.Table.TableName)
		ch <- prometheus.MustNewConstMetric(PromTableMeanPartitionSize,
			prometheus.GaugeValue, float64(stat.MeanPartitionSize),
			stat.Table.KeyspaceName, stat.Table.TableName)
		ch <- prometheus.MustNewConstMetric(PromTableBloomFilterFalseRatio,
			prometheus.GaugeValue, float64(stat.BloomFilterFalseRatio),
			stat.Table.KeyspaceName, stat.Table.TableName)
		ch <- prometheus.MustNewConstSummary(PromTableTombstonesScanned,
			uint64(stat.TombstonesScanned.Count),
			float64(stat.TombstonesScanned.Count)*float64(stat.TombstonesScanned.Mean),
			map[float64]float64{
				75.0: float64(stat.TombstonesScanned.Percentile75),
				95.0: float64(stat.TombstonesScanned.Percentile95),
				99.0: float64(stat.TombstonesScanned.Percentile99),
				99.9: float64(stat.TombstonesScanned.Percentile999),
			}, stat.Table.KeyspaceName, stat.Table.TableName)
		ch <- prometheus.MustNewConstSummary(PromTableLiveCellsScanned,
			uint64(stat.LiveCellsScanned.Count),
			float64(stat.LiveCellsScanned.Count)*float64(stat.LiveCellsScanned.Mean),
			map[float64]float64{
				75.0: float64(stat.LiveCellsScanned.Percentile75),
				95.0: float64(stat.LiveCellsScanned.Percentile95),
				99.0: float64(stat.LiveCellsScanned.Percentile99),
				99.9: float64(stat.LiveCellsScanned.Percentile999),
			}, stat.Table.KeyspaceName, stat.Table.TableName)
		ch <- prometheus.MustNewConstMetric(PromTableKeyCacheHitRate,
			prometheus.GaugeValue, float64(stat.KeyCacheHitRate),
			stat.Table.KeyspaceName, stat.Table.TableName)
		ch <- prometheus.MustNewConstMetric(PromTablePercentRepaired,
			prometheus.GaugeValue, float64(stat.PercentRepaired),
			stat.Table.KeyspaceName, stat.Table.TableName)
		ch <- prometheus.MustNewConstMetric(PromTableSpeculativeRetries,
			prometheus.GaugeValue, float64(stat.SpeculativeRetries),
			stat.Table.KeyspaceName, stat.Table.TableName)
		ch <- prometheus.MustNewConstMetric(PromTableSpeculativeFailedRetries,
			prometheus.GaugeValue, float64(stat.SpeculativeFailedRetries),
			stat.Table.KeyspaceName, stat.Table.TableName)
	}
}

func addCQLStats(metrics ScrapedMetrics, ch chan<- prometheus.Metric) {
	if metrics.CQLStats == nil {
		return
	}

	// CQLStats
	ch <- prometheus.MustNewConstMetric(PromCQLPreparedStatementsCount,
		prometheus.GaugeValue, float64(metrics.CQLStats.PreparedStatementsCount))
	ch <- prometheus.MustNewConstMetric(PromCQLPreparedStatementsEvicted,
		prometheus.CounterValue, float64(metrics.CQLStats.PreparedStatementsEvicted))
	ch <- prometheus.MustNewConstMetric(PromCQLPreparedStatementsExecuted,
		prometheus.CounterValue, float64(metrics.CQLStats.PreparedStatementsExecuted))
	ch <- prometheus.MustNewConstMetric(PromCQLRegularStatementsExecuted,
		prometheus.CounterValue, float64(metrics.CQLStats.RegularStatementsExecuted))
	ch <- prometheus.MustNewConstMetric(PromCQLPreparedStatementsRatio,
		prometheus.GaugeValue, float64(metrics.CQLStats.PreparedStatementsRatio))
}

func addThreadPoolStats(metrics ScrapedMetrics, ch chan<- prometheus.Metric) {
	if metrics.ThreadPoolStats == nil {
		return
	}

	// ThreadPoolStats
	for _, pool := range metrics.ThreadPoolStats {
		ch <- prometheus.MustNewConstMetric(PromThreadPoolActiveTasks,
			prometheus.GaugeValue, float64(pool.ActiveTasks), pool.PoolName)
		ch <- prometheus.MustNewConstMetric(PromThreadPoolPendingTasks,
			prometheus.GaugeValue, float64(pool.PendingTasks), pool.PoolName)
		ch <- prometheus.MustNewConstMetric(PromThreadPoolCompletedTasks,
			prometheus.CounterValue, float64(pool.CompletedTasks), pool.PoolName)
		ch <- prometheus.MustNewConstMetric(PromThreadPoolTotalBlockedTasks,
			prometheus.CounterValue, float64(pool.TotalBlockedTasks), pool.PoolName)
		ch <- prometheus.MustNewConstMetric(PromThreadPoolCurrentlyBlockedTasks,
			prometheus.GaugeValue, float64(pool.CurrentlyBlockedTasks), pool.PoolName)
		ch <- prometheus.MustNewConstMetric(PromThreadPoolMaxPoolSize,
			prometheus.GaugeValue, float64(pool.MaxPoolSize), pool.PoolName)
	}
}

func addCompactionStats(metrics ScrapedMetrics, ch chan<- prometheus.Metric) {
	if metrics.CompactionStats == nil {
		return
	}

	// CompactionStats
	ch <- prometheus.MustNewConstMetric(PromCompactionBytesCompacted,
		prometheus.CounterValue, float64(metrics.CompactionStats.BytesCompacted))
	ch <- prometheus.MustNewConstMetric(PromCompactionPendingTasks,
		prometheus.GaugeValue, float64(metrics.CompactionStats.PendingTasks))
	ch <- prometheus.MustNewConstMetric(PromCompactionCompletedTasks,
		prometheus.CounterValue, float64(metrics.CompactionStats.CompletedTasks))
}

func addClientRequestStats(metrics ScrapedMetrics, ch chan<- prometheus.Metric) {
	// ClientRequestStats
	for _, stat := range metrics.ClientRequestStats {
		ch <- prometheus.MustNewConstSummary(PromClientRequestLatency,
			uint64(stat.RequestLatency.Count),
			float64(stat.RequestLatency.Count)*stat.RequestLatency.Mean.Seconds(),
			map[float64]float64{
				75.0: stat.RequestLatency.Percentile75.Seconds(),
				95.0: stat.RequestLatency.Percentile95.Seconds(),
				99.0: stat.RequestLatency.Percentile99.Seconds(),
				99.9: stat.RequestLatency.Percentile999.Seconds(),
			}, stat.RequestType)
		ch <- prometheus.MustNewConstMetric(PromClientRequestTimeouts,
			prometheus.CounterValue, float64(stat.Timeouts), stat.RequestType)
		ch <- prometheus.MustNewConstMetric(PromClientRequestFailures,
			prometheus.CounterValue, float64(stat.Failures), stat.RequestType)
		ch <- prometheus.MustNewConstMetric(PromClientRequestUnavailable,
			prometheus.CounterValue, float64(stat.Unavailables), stat.RequestType)
	}
}

func addConnectedClientStats(metrics ScrapedMetrics, ch chan<- prometheus.Metric) {
	if metrics.ConnectedClients == nil {
		return
	}

	ch <- prometheus.MustNewConstMetric(PromConnectedClients,
		prometheus.GaugeValue, float64(*metrics.ConnectedClients))
}

func addMemoryStats(metrics ScrapedMetrics, ch chan<- prometheus.Metric) {
	if metrics.MemoryStats == nil {
		return
	}

	// MemoryStats
	ch <- prometheus.MustNewConstMetric(PromMemoryStatsHeapUsed,
		prometheus.GaugeValue, float64(metrics.MemoryStats.HeapUsed))
	ch <- prometheus.MustNewConstMetric(PromMemoryStatsNonHeapUsed,
		prometheus.GaugeValue, float64(metrics.MemoryStats.NonHeapUsed))
}

func addGCStats(metrics ScrapedMetrics, ch chan<- prometheus.Metric) {
	if metrics.GCStats == nil {
		return
	}

	// GCStats
	for _, stat := range metrics.GCStats {
		ch <- prometheus.MustNewConstMetric(PromGCStatsCountTotal,
			prometheus.CounterValue, float64(stat.Count), stat.Name)
		ch <- prometheus.MustNewConstMetric(PromGCStatsLastGC,
			prometheus.GaugeValue, stat.LastGC.Seconds(), stat.Name)
		ch <- prometheus.MustNewConstMetric(PromGCStatsAccumulatedGC,
			prometheus.CounterValue, stat.Accumulated.Seconds(), stat.Name)
	}
}
