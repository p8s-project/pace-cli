package cmd

import (
	"fmt"
	"os"

	"github.com/Vezia/vez-cli/internal/generator"
	"github.com/spf13/cobra"
)

var (
	appFilePath string
	catalogPath string
	outputDir   string
	verbose     bool
)

var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"gen"},
	Short:   "Generate Terraform code from an app manifest.",
	Long:    `Generate Terraform code based on a high-level app.yaml manifest and a catalog.yaml.`,
	Example: `  # Generate terraform files for the application defined in app.yaml
  pace generate --app-file ./app.yaml --catalog ./catalog.yaml --output-dir ./infra`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if verbose {
			fmt.Println("Running generate command in verbose mode...")
		}

		// If no catalog path is provided, use the default public starter kit.
		if catalogPath == "" {
			if verbose {
				fmt.Println("No catalog specified, using default public starter kit...")
			}
			// In a real implementation, we would fetch this from a remote URL.
			// For now, we will simulate this by creating a temporary catalog file.
			// TODO: Replace this with a real git fetch in the future.
			tempCatalog, err := createTempDefaultCatalog()
			if err != nil {
				return fmt.Errorf("failed to create temporary default catalog: %w", err)
			}
			defer os.Remove(tempCatalog.Name())
			catalogPath = tempCatalog.Name()
		}

		g, err := generator.New(appFilePath, catalogPath)
		if err != nil {
			return fmt.Errorf("failed to create generator: %w", err)
		}

		if verbose {
			fmt.Println("Successfully parsed input files.")
		}

		opts := &generator.Options{
			Verbose: verbose,
		}

		if err := g.Generate(outputDir, opts); err != nil {
			return fmt.Errorf("failed to generate files: %w", err)
		}

		fmt.Println("Successfully generated all files.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringVarP(&appFilePath, "app-file", "a", "app.yaml", "Path to the app.yaml file")
	generateCmd.Flags().StringVarP(&catalogPath, "catalog", "c", "catalogs", "Path to a catalog.yaml file or a directory of .yaml files (optional)")
	generateCmd.Flags().StringVarP(&outputDir, "output-dir", "o", "infra", "Directory to output the generated Terraform files")
	generateCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output for debugging")
}

// createTempDefaultCatalog simulates fetching the public starter kit.
// TODO: Replace with a real implementation that fetches from a Git URL.
func createTempDefaultCatalog() (*os.File, error) {
	tmpfile, err := os.CreateTemp("", "default-catalog.yaml")
	if err != nil {
		return nil, err
	}

	content := []byte(`
resources:
  s3-bucket:
    source: "terraform-aws-modules/s3-bucket/aws"
    version: "3.15.1"
    inputs:
      - from: "name"
        to: "bucket"
        required: true
      - from: "versioning"
        to: "versioning.enabled"
        required: false
        default: false
`)
	if _, err := tmpfile.Write(content); err != nil {
		return nil, err
	}
	if err := tmpfile.Close(); err != nil {
		return nil, err
	}
	return tmpfile, nil
}
