package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	tomb "gopkg.in/tomb.v2"

	"github.com/suhailpatel/seastat/flags"
	"github.com/suhailpatel/seastat/jolokia"
)

// Run takes in the Jolokia client and some options and does everything needed
// to start scraping and serving metrics
func Run(client jolokia.Client, interval time.Duration, port int) {
	// Parent context to track all our child goroutines
	ctx, cancel := context.WithCancel(context.Background())

	// This tomb will take care of all our goroutines such as the scraper and
	// the webserver. If something unexpected happens or we need to gracefully
	// terminate, it'll keep track of everything pending
	t := tomb.Tomb{}

	// Start up our scraper
	scraper := NewScraper(client)
	t.Go(func() error {
		// Set up our scraper for shutdown when our context terminates
		t.Go(func() error {
			<-ctx.Done()
			scraper.Stop()
			return nil
		})

		logrus.Infof("🕷️ Starting scraper (interval: %v)", interval)
		if err := scraper.Run(interval); err != nil {
			logrus.Errorf("error whilst scraping: %v", err)
			t.Kill(fmt.Errorf("error whilst scraping: %v", err))
		}
		logrus.Infof("🦠 Stopping scraper")
		return nil
	})

	// Set up the Prometheus collector
	collector := NewSeastatCollector(scraper)
	prometheus.MustRegister(collector)

	// Set up our webserver
	addr := fmt.Sprintf(":%d", port)
	srv := &http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			w.Header().Set("Seastat-Version", flags.Version)
			http.DefaultServeMux.ServeHTTP(w, r)
			logrus.Infof("%s %s %.2fms", r.Method, r.URL, time.Since(start).Seconds()*1000.0)
		}),
	}
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/healthz", handleHealthz(client))
	http.HandleFunc("/", handleRoot())

	// Start up our webserver
	t.Go(func() error {
		// Set up our server for graceful shutdown when our context terminates
		t.Go(func() error {
			<-ctx.Done()
			srv.Shutdown(ctx)
			return nil
		})

		logrus.Infof("👂 Listening on %s", addr)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			logrus.Errorf("error whilst serving: %v", err)
			t.Kill(fmt.Errorf("error whilst serving: %v", err))
		}
		logrus.Infof("😴 Server has shut down")
		return nil
	})

	// Handle signal termination by cancelling our context which should
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sig := <-sigs:
		logrus.Infof("🏁 Received OS signal %v, shutting down", sig)
	case <-t.Dying():
		logrus.Infof("⚰️ Tomb is dying, shutting down")
	}
	cancel() // cancel our context to kick off the shutdown chain

	// Wait a maximum of 10 seconds for everything to cleanly shut down
	select {
	case <-t.Dead():
		logrus.Infof("👋 Goodbye!")
	case <-time.After(10 * time.Second):
		logrus.Errorf("🔴 Did not gracefully terminate in time, force exiting")
		os.Exit(128)
	}
}

func handleRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "🌊 Seastat Cassandra Exporter %v (Commit: %v)", flags.Version, flags.GitCommitHash)
	}
}

func handleHealthz(client jolokia.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jolokiaVersion, err := client.Version()
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			v, _ := json.Marshal(map[string]string{"error": fmt.Sprintf("%v", err)}) // not much we can do if this errors
			w.Write(v)
			return
		}

		w.WriteHeader(http.StatusOK)
		v, _ := json.Marshal(map[string]string{"jolokia": jolokiaVersion, "seastat": flags.Version})
		w.Write(v)
	}
}
