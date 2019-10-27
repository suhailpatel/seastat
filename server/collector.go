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
