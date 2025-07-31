package generator

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/p8s-project/pace-cli/internal/types"
	"gopkg.in/yaml.v3"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

// Generator holds the state for the generation process.
type Generator struct {
	App     *types.AppManifest
	Catalog *types.Catalog
}

// New creates a new Generator.
func New(appFilePath, catalogPath string) (*Generator, error) {
	app, err := loadAppManifest(appFilePath)
	if err != nil {
		return nil, err
	}

	catalog, err := loadCatalog(catalogPath)
	if err != nil {
		return nil, err
	}

	return &Generator{
		App:     app,
		Catalog: catalog,
	}, nil
}

func loadAppManifest(path string) (*types.AppManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read app manifest file: %w", err)
	}

	var manifest types.AppManifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse app manifest: %w", err)
	}

	return &manifest, nil
}

func loadCatalog(path string) (*types.Catalog, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read catalog file: %w", err)
	}

	var catalog types.Catalog
	if err := yaml.Unmarshal(data, &catalog); err != nil {
		return nil, fmt.Errorf("failed to parse catalog: %w", err)
	}

	return &catalog, nil
}

// Generate processes the AppManifest and Catalog, executing templates to create Terraform files.
func (g *Generator) Generate(outputDir string) error {
	// Ensure the output directory exists.
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Parse the generic Terraform module template.
	tmpl, err := template.ParseFS(templateFS, "templates/module.tf.tmpl")
	if err != nil {
		return fmt.Errorf("failed to parse generic module template: %w", err)
	}

	// Iterate over each resource requested in the app manifest.
	for _, request := range g.App.Resources {
		// Find the corresponding resource specification in the catalog.
		catalogEntry, ok := g.Catalog.Resources[request.Uses]
		if !ok {
			return fmt.Errorf("resource type '%s' requested by '%s' not found in catalog", request.Uses, request.ID)
		}

		// Build the map of inputs to pass to the Terraform module.
		inputs, err := g.buildInputs(request, catalogEntry)
		if err != nil {
			return err
		}

		// Prepare the data structure for template execution.
		templateData := struct {
			Request      types.ResourceRequest
			CatalogEntry types.ResourceSpec
			Inputs       map[string]interface{}
		}{
			Request:      request,
			CatalogEntry: catalogEntry,
			Inputs:       inputs,
		}

		// Execute the template to generate the HCL code.
		var output bytes.Buffer
		if err := tmpl.Execute(&output, templateData); err != nil {
			return fmt.Errorf("failed to execute template for resource %s: %w", request.ID, err)
		}

		// Write the generated HCL to a file.
		fileName := fmt.Sprintf("%s.tf", request.ID)
		outputPath := filepath.Join(outputDir, fileName)
		if err := os.WriteFile(outputPath, output.Bytes(), 0644); err != nil {
			return fmt.Errorf("failed to write output file %s: %w", outputPath, err)
		}
		fmt.Printf("Successfully generated %s\n", outputPath)
	}
	return nil
}

// buildInputs constructs the final map of variables to be passed to a Terraform module.
// It maps the inputs from the developer's request (`with` block) to the variable names
// expected by the module, handling defaults and required fields as defined in the catalog.
func (g *Generator) buildInputs(request types.ResourceRequest, catalogEntry types.ResourceSpec) (map[string]interface{}, error) {
	inputs := make(map[string]interface{})

	// Always map the resource ID to the 'name' input if it's defined in the catalog.
	// This is a common convention.
	for _, inputSpec := range catalogEntry.Inputs {
		if inputSpec.From == "id" {
			inputs[inputSpec.To] = fmt.Sprintf(`"%s"`, request.ID)
		}
	}

	// Process the `with` block from the developer's request.
	for _, inputSpec := range catalogEntry.Inputs {
		// Skip the special 'id' mapping as it's already handled.
		if inputSpec.From == "id" {
			continue
		}

		val, ok := request.With[inputSpec.From]
		if !ok {
			if inputSpec.Required {
				return nil, fmt.Errorf("missing required input '%s' for resource '%s'", inputSpec.From, request.ID)
			}
			// Use the default value from the catalog if one is defined.
			if inputSpec.Default != nil {
				val = inputSpec.Default
			} else {
				// If no default is specified, skip this input.
				continue
			}
		}

		// Perform value mapping for specific known inputs.
		if inputSpec.From == "size" {
			val = mapSizeToInstanceClass(val.(string))
		}

		// Format the value for HCL. Strings are quoted, other types are passed through.
		if s, ok := val.(string); ok {
			inputs[inputSpec.To] = fmt.Sprintf(`"%s"`, s)
		} else {
			inputs[inputSpec.To] = val
		}
	}
	return inputs, nil
}

// mapSizeToInstanceClass is a simple example of a value mapper.
func mapSizeToInstanceClass(size string) string {
	switch size {
	case "small":
		return "db.t3.small"
	case "medium":
		return "db.t3.medium"
	case "large":
		return "db.t3.large"
	default:
		return "db.t3.micro" // A safe default.
	}
}
