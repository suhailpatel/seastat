package cmd

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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
		logrus.Fatalf("interval must be a minimum of 10 seconds")
	}
}
