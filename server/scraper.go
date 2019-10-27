package server

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/suhailpatel/seastat/jolokia"
)

// tableScrapeInterval defines how often we will ask for a full table list
const tableScrapeInterval = 5 * time.Minute

type Scraper struct {
	client  jolokia.Client
	stopped chan struct{}

	// Keep track of all our tables and when we last scraped them
	tables          []jolokia.Table
	lastTableScrape time.Time

	metrics           ScrapedMetrics
	lastMetricsScrape time.Time
}

type ScrapedMetrics struct {
	TableStats       []jolokia.TableStats
	CQLStats         jolokia.CQLStats
	ThreadPoolStats  []jolokia.ThreadPoolStats
	ConnectedClients jolokia.Gauge
	MemoryStats      jolokia.MemoryStats
	GCStats          []jolokia.GCStats

	ScrapeDuration time.Duration
}

// NewScraper returns a new instance of a Scraper
func NewScraper(client jolokia.Client) Scraper {
	return Scraper{
		client:  client,
		stopped: make(chan struct{}),
	}
}

// Run blocks whilst attempting to scrape. It will stop scraping
func (s *Scraper) Run(interval time.Duration) error {
	// Run an initial scrape before we kick off the timer
	s.runScrape()

	t := time.Tick(interval)
	for {
		select {
		case <-t:
			s.runScrape()
		case <-s.stopped:
			return nil
		}
	}
}

// runScrape is the mammoth function which handles all the scraping via Jolokia
func (s *Scraper) runScrape() {
	start := time.Now()

	// Do a quick version sanity check, if this fails, we will bail out
	_, err := s.client.Version()
	if err != nil {
		// bail out, we don't want to continue if we couldn't even get
		// the version string
		logrus.Debugf("🦂 Could not fetch version, bailing out")
		return
	}

	// First check to see if our tables need a refresh
	if len(s.tables) == 0 || time.Now().Sub(s.lastTableScrape) > tableScrapeInterval {
		tables, err := s.client.Tables()
		if err != nil {
			// bail out, we don't want to continue if we don't have updated tables
			logrus.Debugf("🦂 Could not refresh tables, bailing out")
			return
		}
		s.tables = tables
		s.lastTableScrape = time.Now()
		logrus.Debugf("🐝 Refreshed table list, got %d tables (took %d ms)", len(s.tables), time.Since(start).Milliseconds())
	}

	s.metrics = s.scrapeAllMetrics()
	s.lastMetricsScrape = time.Now()

	logrus.Debugf("🕸️ Finished scrape for %d tables (took %d ms)", len(s.tables), time.Since(start).Milliseconds())
}

func (s *Scraper) scrapeAllMetrics() ScrapedMetrics {
	scrapeStart := time.Now()

	tableStats := s.scrapeTableMetrics()

	cqlStats, err := s.client.CQLStats()
	if err != nil {
		logrus.Debugf("🦂 Could not get CQL stats: %v", err)
	}

	tpStats, err := s.client.ThreadPoolStats()
	if err != nil {
		logrus.Debugf("🦂 Could not get ThreadPool stats: %v", err)
	}

	connectedClients, err := s.client.ConnectedClients()
	if err != nil {
		logrus.Debugf("🦂 Could not get Client stats: %v", err)
	}

	memoryStats, err := s.client.MemoryStats()
	if err != nil {
		logrus.Debugf("🦂 Could not get Memory stats: %v", err)
	}

	gcStats, err := s.client.GarbageCollectionStats()
	if err != nil {
		logrus.Debugf("🦂 Could not get GC stats: %v", err)
	}

	return ScrapedMetrics{
		TableStats:       tableStats,
		CQLStats:         cqlStats,
		ThreadPoolStats:  tpStats,
		ConnectedClients: connectedClients,
		MemoryStats:      memoryStats,
		GCStats:          gcStats,
		ScrapeDuration:   time.Since(scrapeStart),
	}
}

func (s *Scraper) scrapeTableMetrics() []jolokia.TableStats {
	// The goal of this function is to scrape the table metrics in parallel.
	// A few numbers were tried and 8 seemed to be the sweet spot with Jolokia.
	// Ramping this too high may lead to stuck conns
	const workers = 8

	type result struct {
		table      jolokia.Table
		tableStats jolokia.TableStats
		err        error
	}

	workerCh := make(chan jolokia.Table, workers)
	resultCh := make(chan result, len(s.tables))

	wg := sync.WaitGroup{}
	workerFunc := func() {
		for {
			select {
			case table := <-workerCh:
				if table.KeyspaceName == "" && table.TableName == "" {
					return // closed channel
				}

				stats, err := s.client.TableStats(table)
				resultCh <- result{
					table:      table,
					tableStats: stats,
					err:        err,
				}
				wg.Done()
			}
		}
	}

	for i := 0; i < workers; i++ {
		go workerFunc()
	}
	for _, table := range s.tables {
		wg.Add(1)
		workerCh <- table
	}
	wg.Wait()
	close(workerCh)
	close(resultCh)

	tableStats := make([]jolokia.TableStats, 0, len(s.tables))
	for res := range resultCh {
		// Occassionally, we might not be abkle to fetch table stats for a
		// table. This isn't the end of the world
		if res.err != nil {
			logrus.Debugf("🦂 Could not get table stats for %s.%s: %v", res.table.KeyspaceName,
				res.table.TableName, res.err)
			continue
		}
		tableStats = append(tableStats, res.tableStats)
	}
	return tableStats
}

// Stop informs the scraper to stop scraping any further
func (s *Scraper) Stop() {
	close(s.stopped)
}
