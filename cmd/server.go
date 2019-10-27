package cmd

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/suhailpatel/seastat/jolokia"
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
}

func run(cmd *cobra.Command) {
	endpoint, _ := cmd.Flags().GetString("endpoint")
	interval, _ := cmd.Flags().GetDuration("interval")

	if endpoint == "" {
		logrus.Fatalf("'endpoint' can not be empty")
	}

	if interval < 1*time.Second {
		interval = 10 * time.Second
	}

	// Test our connection to Jolokia to make sure everything is good!
	client := jolokia.Init(endpoint)
	version, err := client.Version()
	if err != nil {
		logrus.Fatalf("could not get version from Jolokia: %v", err)
	}

	logrus.Debugf("Running with Jolokia version: %v", version)
}
