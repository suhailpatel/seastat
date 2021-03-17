package cmd

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/suhailpatel/seastat/jolokia"
	"github.com/suhailpatel/seastat/server"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Server kicks off the Cassandra Exporter",
	Long:  `Server kicks off the Cassandra Exporter`,
	Run: func(cmd *cobra.Command, args []string) {
		run(cmd)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.PersistentFlags().String("endpoint", "http://localhost:8778", "endpoint where Jolokia is running")
	serverCmd.PersistentFlags().Duration("interval", 30*time.Second, "how often we attempt to extract metrics (minimum 10s)")
	serverCmd.PersistentFlags().Int("port", 8080, "port to run the Seastat server on (for Prometheus to scrape)")
	serverCmd.PersistentFlags().Duration("timeout", 3*time.Second, "how long before we timeout a Jolokia request")
	serverCmd.PersistentFlags().Int("concurrency", 10, "maximum number of concurrent requests to Jolokia")
	viper.BindPFlag("endpoint", serverCmd.Flags().Lookup("endpoint"))
	viper.BindPFlag("interval", serverCmd.Flags().Lookup("interval"))
	viper.BindPFlag("port", serverCmd.Flags().Lookup("port"))
	viper.BindPFlag("timeout", serverCmd.Flags().Lookup("timeout"))
	viper.BindPFlag("concurrency", serverCmd.Flags().Lookup("concurrency"))
}

func run(cmd *cobra.Command) {
	endpoint := viper.GetString("endpoint")
	interval := viper.GetDuration("interval")
	port := viper.GetInt("port")
	timeout := viper.GetDuration("timeout")
	concurrency := viper.GetInt("concurrency")

	if endpoint == "" {
		logrus.Fatalf("'endpoint' can not be empty")
	}

	if interval < 10*time.Second {
		interval = 10 * time.Second
	}

	if port < 0 {
		port = 8000
	}

	client := jolokia.Init(endpoint, timeout)

	// Run a quick sanity check of the provided endpoint
	version, err := client.Version()
	if err != nil {
		logrus.Fatalf("could not connect to Jolokia: %v", err)
	}
	logrus.Infof("☕ Communicating with Jolokia %s (%s)", version, endpoint)
	server.Run(client, interval, port, concurrency)
}
