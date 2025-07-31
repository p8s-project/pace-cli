package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pace",
	Short: "pace is a command-line tool for generating infrastructure from a simple manifest.",
	Long: `pace is the CLI for the p8s platform.
It generates standardized, production-ready infrastructure code
from a high-level declarative manifest.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Default action when no subcommand is given
		fmt.Println("Welcome to pace! Use 'pace --help' for more information.")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
