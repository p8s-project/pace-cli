package cmd

import (
	"fmt"
	"net/url"
	"os"

	"github.com/Vezia/vez-cli/internal/loader"
	"github.com/Vezia/vez-cli/internal/types"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate an app.yaml against a catalog.yaml",
	Long:  `Parses and validates an app.yaml against a catalog.yaml to ensure all resource requests are valid.`,
	Example: `  # Validate that my-app.yaml is compatible with the production catalog
  pace validate --app-file ./my-app.yaml --catalog ./prod-catalog.yaml`,
	RunE: func(cmd *cobra.Command, args []string) error {
		appFilePath, _ := cmd.Flags().GetString("app-file")
		catalogPath, _ := cmd.Flags().GetString("catalog")

		fmt.Println("Running validation...")
		if err := validate(appFilePath, catalogPath); err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}
		fmt.Println("Validation successful: app.yaml is compatible with the catalog.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
	validateCmd.Flags().String("app-file", "", "Path to the app.yaml file (required)")
	validateCmd.Flags().String("catalog", "", "Path to the catalog.yaml file (required)")
	validateCmd.MarkFlagRequired("app-file")
	validateCmd.MarkFlagRequired("catalog")
}

func validate(appFilePath, catalogPath string) error {
	app, err := loader.LoadApp(appFilePath)
	if err != nil {
		return err
	}

	catalog, err := loader.LoadCatalog(catalogPath)
	if err != nil {
		return err
	}

	if err := validateApp(app, catalog); err != nil {
		return err
	}

	if err := validateCatalog(catalog); err != nil {
		return err
	}

	return nil
}

func validateApp(app *types.AppManifest, catalog *types.Catalog) error {
	for _, resource := range app.Resources {
		spec, ok := catalog.Resources[resource.Uses]
		if !ok {
			return fmt.Errorf("resource '%s' not found in catalog", resource.Uses)
		}

		for _, input := range spec.Inputs {
			if input.Required {
				if _, ok := resource.With[input.From]; !ok {
					return fmt.Errorf("required input '%s' for resource '%s' not found", input.From, resource.Uses)
				}
			}
		}
	}
	return nil
}

func validateCatalog(catalog *types.Catalog) error {
	for name, resource := range catalog.Resources {
		if resource.Source == "" {
			return fmt.Errorf("source for resource '%s' is empty", name)
		}

		if _, err := url.ParseRequestURI(resource.Source); err != nil {
			if _, err := os.Stat(resource.Source); err != nil {
				return fmt.Errorf("source for resource '%s' is not a valid URL or file path", name)
			}
		}
	}
	return nil
}
