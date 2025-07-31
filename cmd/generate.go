package cmd

import (
	"fmt"
	"os"

	"github.com/p8s-project/pace-cli/internal/generator"
	"github.com/spf13/cobra"
)

var (
	appFilePath string
	catalogPath string
	outputDir   string
)

var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"gen"},
	Short:   "Generate Terraform code from an app manifest.",
	Long:    `Generate Terraform code based on a high-level app.yaml manifest and a catalog.yaml.`,
	Example: `  # Generate terraform files for the application defined in app.yaml
  pace generate --app-file ./app.yaml --catalog ./catalog.yaml --output-dir ./infra`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running generate command...")

		g, err := generator.New(appFilePath, catalogPath)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Successfully parsed files!")

		err = g.Generate(outputDir)
		if err != nil {
			fmt.Printf("Error generating files: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringVarP(&appFilePath, "app-file", "a", "", "Path to the app.yaml file (required)")
	generateCmd.Flags().StringVarP(&catalogPath, "catalog", "c", "", "Path to the catalog.yaml file (required)")
	generateCmd.Flags().StringVarP(&outputDir, "output-dir", "o", "", "Directory to output the generated Terraform files (required)")
	generateCmd.MarkFlagRequired("app-file")
	generateCmd.MarkFlagRequired("catalog")
	generateCmd.MarkFlagRequired("output-dir")
}
