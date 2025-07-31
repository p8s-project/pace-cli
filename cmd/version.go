package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version is set at build time
var Version = "v0.1.0-dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of pace",
	Long:  `All software has versions. This is pace's.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("pace version %s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
