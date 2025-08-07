package builder

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Vezia/vez-cli/internal/types"
)

// Builder holds the state for the building process.
type Builder struct {
	App     *types.AppManifest
	Catalog *types.Catalog
}

// New creates a new Builder.
func New(app *types.AppManifest, catalog *types.Catalog) *Builder {
	return &Builder{
		App:     app,
		Catalog: catalog,
	}
}

// BuildInputs constructs the final map of variables to be passed to a Terraform module.
func (b *Builder) BuildInputs(request types.ResourceRequest, catalogEntry types.ResourceSpec, outputs map[string]map[string]string) (map[string]interface{}, error) {
	inputs := make(map[string]interface{})

	// Always map the resource ID to the 'name' input if it's defined in the catalog.
	// This is a common convention.
	for _, inputSpec := range catalogEntry.Inputs {
		if inputSpec.From == "id" {
			inputs[inputSpec.To] = request.ID
		}
	}

	// Process the `with` block from the developer's request.
	for _, inputSpec := range catalogEntry.Inputs {
		// Skip the special 'id' mapping as it's already handled.
		if inputSpec.From == "id" {
			continue
		}

		var val interface{}
		var ok bool

		// Safely check if the 'With' map exists before accessing it.
		if request.With != nil {
			val, ok = request.With[inputSpec.From]
		}

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

		// Check for dependency injection.
		if s, ok := val.(string); ok && strings.HasPrefix(s, "{{") && strings.HasSuffix(s, "}}") {
			parts := strings.Split(strings.Trim(s, "{}"), ".")
			if len(parts) != 3 {
				return nil, fmt.Errorf("invalid dependency reference: %s", s)
			}
			depID := parts[1]
			outputName := parts[2]
			if _, ok := outputs[depID]; !ok {
				return nil, fmt.Errorf("dependency '%s' not found", depID)
			}
			if _, ok := outputs[depID][outputName]; !ok {
				return nil, fmt.Errorf("output '%s' not found on dependency '%s'", outputName, depID)
			}
			val = outputs[depID][outputName]
		}

		inputs[inputSpec.To] = val
	}

	// Sort the keys of the map to ensure a consistent order.
	sortedInputs := make(map[string]interface{})
	keys := make([]string, 0, len(inputs))
	for k := range inputs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		sortedInputs[k] = inputs[k]
	}

	return sortedInputs, nil
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
