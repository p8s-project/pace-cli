package loader

import (
	"fmt"
	"os"

	"github.com/Vezia/vez-cli/internal/types"
	"gopkg.in/yaml.v3"
)

func LoadApp(path string) (*types.AppManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read app file: %w", err)
	}

	var app types.AppManifest
	if err := yaml.Unmarshal(data, &app); err != nil {
		return nil, fmt.Errorf("failed to unmarshal app file: %w", err)
	}

	return &app, nil
}

func LoadCatalog(path string) (*types.Catalog, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read catalog file: %w", err)
	}

	var catalog types.Catalog
	if err := yaml.Unmarshal(data, &catalog); err != nil {
		return nil, fmt.Errorf("failed to unmarshal catalog file: %w", err)
	}

	return &catalog, nil
}
