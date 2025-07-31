package cmd

import (
	"fmt"
	"os"

	"github.com/p8s-project/pace-cli/internal/generator"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate an app.yaml against a catalog.yaml",
	Long:  `Parses and validates an app.yaml against a catalog.yaml to ensure all resource requests are valid.`,
	Example: `  # Validate that my-app.yaml is compatible with the production catalog
  pace validate --app-file ./my-app.yaml --catalog ./prod-catalog.yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running validation...")
		_, err := generator.New(appFilePath, catalogPath)
		if err != nil {
			fmt.Printf("Validation failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Validation successful: app.yaml is compatible with the catalog.")
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
	validateCmd.Flags().StringVarP(&appFilePath, "app-file", "a", "", "Path to the app.yaml file (required)")
	validateCmd.Flags().StringVarP(&catalogPath, "catalog", "c", "", "Path to the catalog.yaml file (required)")
	validateCmd.MarkFlagRequired("app-file")
	validateCmd.MarkFlagRequired("catalog")
}
