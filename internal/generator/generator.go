package generator

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/flosch/pongo2/v6"
	"github.com/Vezia/vez-cli/internal/builder"
	"github.com/Vezia/vez-cli/internal/loader"
	"github.com/Vezia/vez-cli/internal/resolver"
	"github.com/Vezia/vez-cli/internal/types"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

// Options defines configuration for the generation process.
type Options struct {
	Verbose bool
}

// Generator holds the state for the generation process.
type Generator struct {
	App     *types.AppManifest
	Catalog *types.Catalog
}

// New creates a new Generator.
func New(appFilePath, catalogPath string) (*Generator, error) {
	app, err := loader.LoadApp(appFilePath)
	if err != nil {
		return nil, err
	}

	catalog, err := loader.LoadCatalog(catalogPath)
	if err != nil {
		return nil, err
	}

	return &Generator{
		App:     app,
		Catalog: catalog,
	}, nil
}

// Generate processes the AppManifest and Catalog, executing templates to create Terraform files.
func (g *Generator) Generate(outputDir string, opts *Options) error {
	// Ensure the output directory exists.
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Validate the app manifest.
	if g.App.Stack == "" {
		return fmt.Errorf("missing required 'stack' field in app.yaml")
	}

	// Create a new resolver.
	r := resolver.New(g.App, g.Catalog)

	// Pass 1: Resolve dependencies and build the outputs map.
	outputs := make(map[string]map[string]string)
	for _, request := range g.App.Resources {
		if err := r.ResolveDependencies(request, outputs); err != nil {
			return err
		}
	}

	// Create a new builder.
	b := builder.New(g.App, g.Catalog)

	// Pass 2: Generate the Terraform files.
	for _, request := range g.App.Resources {
		if err := g.generateResource(request, outputDir, "templates/module.tf.tmpl", outputs, b, r, opts); err != nil {
			return err
		}
	}
	return nil
}

// generateResource generates the Terraform file for a single resource.
func (g *Generator) generateResource(request types.ResourceRequest, outputDir string, tmplPath string, outputs map[string]map[string]string, b *builder.Builder, r *resolver.Resolver, opts *Options) error {
	// Resolve the resource request.
	resolvedUses, err := r.ResolveUses(request.Uses)
	if err != nil {
		return err
	}

	// Find the corresponding resource specification in the catalog.
	catalogEntry, ok := g.Catalog.Resources[resolvedUses]
	if !ok {
		return fmt.Errorf("resource type '%s' requested by '%s' not found in catalog", resolvedUses, request.ID)
	}

	// Build the map of inputs to pass to the Terraform module.
	inputs, err := b.BuildInputs(request, catalogEntry, outputs)
	if err != nil {
		return err
	}

	// Prepare the data structure for template execution.
	templateData := pongo2.Context{
		"Request":      request,
		"CatalogEntry": catalogEntry,
		"Inputs":       inputs,
	}

	// Read the template from the embedded filesystem.
	tplBytes, err := templateFS.ReadFile(tmplPath)
	if err != nil {
		return fmt.Errorf("failed to read template file %s: %w", tmplPath, err)
	}

	// Register custom filters.
	pongo2.RegisterFilter("is_string", isString)

	// Execute the template to generate the HCL code.
	tpl, err := pongo2.FromString(string(tplBytes))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", tmplPath, err)
	}
	output, err := tpl.Execute(templateData)
	if err != nil {
		return fmt.Errorf("failed to execute template for resource %s: %w", request.ID, err)
	}

	// Write the generated HCL to a file.
	fileName := fmt.Sprintf("%s.tf", request.ID)
	outputPath := filepath.Join(outputDir, fileName)
	if err := os.WriteFile(outputPath, []byte(output), 0644); err != nil {
		return fmt.Errorf("failed to write output file %s: %w", outputPath, err)
	}
	fmt.Printf("Successfully generated %s\n", outputPath)

	return nil
}

func isString(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	return pongo2.AsValue(in.IsString()), nil
}
