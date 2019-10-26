package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version is the current version of the app
	Version string
	// GitCommitHash contains the commit hash used to build
	// this version of the app
	GitCommitHash string
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version of Seastat",
	Long:  `Prints the version of Seastat`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %v\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
