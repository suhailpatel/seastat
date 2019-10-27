package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// config options
var (
	cfgFile      string
	logVerbosity string
)

var rootCmd = &cobra.Command{
	Use:   "seastat",
	Short: "Seastat is a Cassandra Prometheus Exporter",
	Long: `Seastat is a speedy Cassandra Prometheus Exporter 🏎️

You point it at an instance running Cassandra 3.0+ and it'll go and
poll for metrics periodically in a speedy manner via Jolokia`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "seastat.yaml", "Config file (default is seastat.yaml)")
	rootCmd.PersistentFlags().StringVarP(&logVerbosity, "verbosity", "v", logrus.InfoLevel.String(), "Log level (debug, info, warn, error, fatal, panic")

	// pre-start hooks
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		fmt.Printf("🌊 Seastat Cassandra Exporter %v\n", Version)

		lvl, err := logrus.ParseLevel(logVerbosity)
		if err != nil {
			return err
		}
		logrus.SetLevel(lvl)
		return nil
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		logrus.Info("Using config file:", viper.ConfigFileUsed())
	}
}
