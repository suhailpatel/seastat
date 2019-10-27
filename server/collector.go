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
