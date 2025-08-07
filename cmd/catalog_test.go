package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCatalogGenerateCmd(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		// Create a temporary directory for the test module.
		moduleDir, err := os.MkdirTemp("", "test-module")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(moduleDir)

		// Create a simple variables.tf file.
		variablesTf := `
variable "name" {
  type = string
}
`
		if err := os.WriteFile(filepath.Join(moduleDir, "variables.tf"), []byte(variablesTf), 0644); err != nil {
			t.Fatalf("failed to write variables.tf: %v", err)
		}

		// Create a temporary directory for the output.
		outputDir, err := os.MkdirTemp("", "test-catalog")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(outputDir)

		// Run the catalog generate command.
		rootCmd.SetArgs([]string{"catalog", "generate", "--from", moduleDir, "--output", outputDir})
		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("command execution failed: %v", err)
		}

		// Read the generated file.
		files, err := filepath.Glob(filepath.Join(outputDir, "*.yaml"))
		if err != nil {
			t.Fatalf("failed to glob generated files: %v", err)
		}
		if len(files) != 1 {
			t.Fatalf("expected 1 generated file, got %d", len(files))
		}
		generatedBytes, err := os.ReadFile(files[0])
		if err != nil {
			t.Fatalf("failed to read generated file: %v", err)
		}

		// Define the golden file content.
		goldenContent := `resources:
    test-module:
        source: git::` + moduleDir + `
        version: 0.0.0
        inputs:
            - from: name
              to: name
              required: true
              default: null
        outputs: []
        dependencies: []
`

		// Compare the generated file to the golden file.
		generatedContent := strings.ReplaceAll(string(generatedBytes), "\r\n", "\n")
		goldenContent = strings.ReplaceAll(goldenContent, "\r\n", "\n")

		if generatedContent != goldenContent {
			t.Errorf("generated content does not match golden file.\n\nGOT:\n%s\n\nWANT:\n%s", generatedContent, goldenContent)
		}
	})
}
