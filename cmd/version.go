package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/suhailpatel/seastat/flags"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version of Seastat",
	Long:  `Prints the version of Seastat`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %v (Commit: %v)\n", flags.Version, flags.GitCommitHash)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
