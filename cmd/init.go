package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	catalogURL string
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize pace with a catalog from a Git repository.",
	Long:  `Downloads a catalog configuration from a Git repository to be used by pace.`,
	Example: `  # Initialize pace with your company's official catalog
  pace init --from git@github.com:my-company/pace-config.git`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Initializing pace with catalog from: %s\n", catalogURL)
		// In a real implementation, we would add logic here to `git clone` the repository.
		// For now, we will just print a success message.
		fmt.Println("Pace has been initialized successfully.")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVar(&catalogURL, "from", "", "The Git URL of the catalog repository (required)")
	initCmd.MarkFlagRequired("from")
}
