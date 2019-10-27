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
		PromTableMaxPartitionSize,
		PromTableMeanPartitionSize,
		PromTableBloomFilterFalseRatio,
		PromTableKeyCacheHitRate,
		PromTablePercentRepaired,

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
		ch <- prometheus.MustNewConstMetric(PromTableMaxPartitionSize,
			prometheus.GaugeValue, float64(stat.MaxPartitionSize),
			stat.Table.KeyspaceName, stat.Table.TableName)
		ch <- prometheus.MustNewConstMetric(PromTableMeanPartitionSize,
			prometheus.GaugeValue, float64(stat.MeanPartitionSize),
			stat.Table.KeyspaceName, stat.Table.TableName)
		ch <- prometheus.MustNewConstMetric(PromTableBloomFilterFalseRatio,
			prometheus.GaugeValue, float64(stat.BloomFilterFalseRatio),
			stat.Table.KeyspaceName, stat.Table.TableName)
		ch <- prometheus.MustNewConstMetric(PromTableKeyCacheHitRate,
			prometheus.GaugeValue, float64(stat.KeyCacheHitRate),
			stat.Table.KeyspaceName, stat.Table.TableName)
		ch <- prometheus.MustNewConstMetric(PromTablePercentRepaired,
			prometheus.GaugeValue, float64(stat.PercentRepaired),
			stat.Table.KeyspaceName, stat.Table.TableName)
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

	// MemoryStats
	ch <- prometheus.MustNewConstMetric(PromMemoryStatsHeapUsed,
		prometheus.GaugeValue, float64(metrics.MemoryStats.HeapUsed))
	ch <- prometheus.MustNewConstMetric(PromMemoryStatsNonHeapUsed,
		prometheus.GaugeValue, float64(metrics.MemoryStats.NonHeapUsed))

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
