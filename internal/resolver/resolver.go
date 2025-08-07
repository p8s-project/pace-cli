package resolver

import (
	"fmt"
	"strings"

	"github.com/Vezia/vez-cli/internal/types"
)

// Resolver holds the state for the resolution process.
type Resolver struct {
	App     *types.AppManifest
	Catalog *types.Catalog
}

// New creates a new Resolver.
func New(app *types.AppManifest, catalog *types.Catalog) *Resolver {
	return &Resolver{
		App:     app,
		Catalog: catalog,
	}
}

// ResolveUses resolves a shorthand or fully qualified `uses` string to a fully qualified name.
func (r *Resolver) ResolveUses(uses string) (string, error) {
	parts := strings.Split(uses, "/")
	if len(parts) == 1 {
		// This is a shorthand. Search for the module in the current stack.
		for name := range r.Catalog.Resources {
			if strings.HasSuffix(name, "/"+uses) && strings.Contains(name, "/"+r.App.Stack+"/") {
				return name, nil
			}
		}
		return "", fmt.Errorf("module '%s' not found in stack '%s'", uses, r.App.Stack)
	}

	// This is a fully qualified name. Validate that it belongs to the current stack.
	if parts[1] != r.App.Stack {
		return "", fmt.Errorf("module '%s' does not belong to stack '%s'", uses, r.App.Stack)
	}
	return uses, nil
}

// ResolveDependencies is a recursive function that walks the dependency tree and populates the outputs map.
func (r *Resolver) ResolveDependencies(request types.ResourceRequest, outputs map[string]map[string]string) error {
	// Avoid processing the same resource twice.
	if _, ok := outputs[request.ID]; ok {
		return nil
	}

	// Resolve the resource request.
	resolvedUses, err := r.ResolveUses(request.Uses)
	if err != nil {
		return err
	}

	// Find the corresponding resource specification in the catalog.
	catalogEntry, ok := r.Catalog.Resources[resolvedUses]
	if !ok {
		return fmt.Errorf("resource type '%s' requested by '%s' not found in catalog", resolvedUses, request.ID)
	}

	// First, process all the dependencies of this resource.
	for _, dep := range catalogEntry.Dependencies {
		if err := r.ResolveDependencies(dep, outputs); err != nil {
			return err
		}
	}

	// Build the outputs map for this resource.
	resourceOutputs := make(map[string]string)
	for _, outputSpec := range catalogEntry.Outputs {
		resourceOutputs[outputSpec.To] = fmt.Sprintf("module.%s.%s", request.ID, outputSpec.From)
	}
	outputs[request.ID] = resourceOutputs

	return nil
}
