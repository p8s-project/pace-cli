package types

// AppManifest represents the structure of a developer's app.yaml file.
// It is the top-level object that defines an application's infrastructure needs.
type AppManifest struct {
	// Name is the unique identifier for the application.
	Name string `yaml:"name"`
	// Stack is the name of the cloud stack to use (e.g., "aws", "gcp").
	Stack string `yaml:"stack"`
	// Tap is the default tap to use for resolving modules.
	Tap string `yaml:"tap"`
	// Resources is a list of the infrastructure components required by the application.
	Resources []ResourceRequest `yaml:"resources"`
}

// ResourceRequest represents a single infrastructure component requested by a developer.
type ResourceRequest struct {
	// ID is the unique identifier for this specific resource instance (e.g., "primary-db").
	ID string `yaml:"id"`
	// Uses is the identifier for the resource type from the catalog (e.g., "postgres:v1").
	Uses string `yaml:"uses"`
	// With is a map of the user-provided inputs for this resource.
	With map[string]interface{} `yaml:"with"`
}

// Catalog represents the structure of a platform team's catalog.yaml file.
type Catalog struct {
	// Resources is a map of the available resource types that can be requested.
	Resources map[string]ResourceSpec `yaml:"resources"`
}

// ResourceSpec defines a single resource type available in the catalog.
type ResourceSpec struct {
	// Source is the URL or path to the underlying Terraform module.
	Source string `yaml:"source"`
	// Version is the specific version of the Terraform module to use.
	Version string `yaml:"version"`
	// Inputs defines the API contract for this resource type.
	Inputs []InputSpec `yaml:"inputs"`
	// Outputs defines the outputs that this resource exposes to other resources.
	Outputs []OutputSpec `yaml:"outputs"`
	// Dependencies is a list of other resources that this resource depends on.
	Dependencies []ResourceRequest `yaml:"dependencies"`
}

// InputSpec defines the mapping from a developer's input in the `with` block
// to a variable in the underlying Terraform module.
type InputSpec struct {
	// From is the key used in the `with` block of a developer's ResourceRequest.
	From string `yaml:"from"`
	// To is the variable name in the underlying Terraform module.
	// Dot notation can be used to access nested variables (e.g., "versioning.enabled").
	To string `yaml:"to"`
	// Required specifies whether the developer must provide this input.
	Required bool `yaml:"required"`
	// Default is the value to use if an input is not provided by the developer.
	Default interface{} `yaml:"default"`
}

// OutputSpec defines the mapping from a Terraform module's output to an output
// that can be consumed by other resources.
type OutputSpec struct {
	// From is the name of the output from the Terraform module.
	From string `yaml:"from"`
	// To is the name of the output to expose to other resources.
	To string `yaml:"to"`
}
