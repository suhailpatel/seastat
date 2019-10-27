package cmd

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

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
}

func run(cmd *cobra.Command) {
	endpoint, _ := cmd.Flags().GetString("endpoint")
	interval, _ := cmd.Flags().GetDuration("interval")
	port, _ := cmd.Flags().GetInt("port")

	if endpoint == "" {
		logrus.Fatalf("'endpoint' can not be empty")
	}

	if interval < 1*time.Second {
		interval = 30 * time.Second
	}

	if port < 0 {
		port = 8000
	}

	client := jolokia.Init(endpoint)

	// Run a quick sanity check of the provided endpoint
	version, err := client.Version()
	if err != nil {
		logrus.Fatalf("could not connect to Jolokia: %v", err)
	}
	logrus.Infof("☕ Communicating with Jolokia %s (%s)", version, endpoint)
	server.Run(client, interval, port)
}
