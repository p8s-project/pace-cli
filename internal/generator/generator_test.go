package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Vezia/vez-cli/internal/loader"
)

func TestLoadAppManifest(t *testing.T) {
	t.Run("valid manifest", func(t *testing.T) {
		tmpfile, err := os.CreateTemp("", "app.yaml")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpfile.Name())

		content := []byte(`
name: my-dynamic-app
resources:
  - id: test-bucket
    uses: s3-bucket:v1
`)
		if _, err := tmpfile.Write(content); err != nil {
			t.Fatal(err)
		}
		if err := tmpfile.Close(); err != nil {
			t.Fatal(err)
		}

		manifest, err := loader.LoadAppManifest(tmpfile.Name())
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if manifest.Name != "my-dynamic-app" {
			t.Errorf("expected name to be 'my-dynamic-app', got '%s'", manifest.Name)
		}
		if len(manifest.Resources) != 1 {
			t.Fatal("expected 1 resource")
		}
		if manifest.Resources[0].ID != "test-bucket" {
			t.Errorf("expected resource ID to be 'test-bucket', got '%s'", manifest.Resources[0].ID)
		}
	})
}

func TestGenerate_GoldenFile(t *testing.T) {
	t.Run("happy path s3", func(t *testing.T) {
		// Define paths relative to the test file.
		testDataBasePath := "../../cmd/testdata/happy_path_s3"
		appFilePath := filepath.Join(testDataBasePath, "app.yaml")
		catalogPath := filepath.Join(testDataBasePath, "catalog.yaml")
		goldenFilePath := filepath.Join(testDataBasePath, "golden/test-bucket.tf")

		// Create a temporary directory for the output.
		outputDir, err := os.MkdirTemp("", "pace-test-output")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(outputDir)

		// Create and run the generator.
		generator, err := New(appFilePath, catalogPath)
		if err != nil {
			t.Fatalf("failed to create generator: %v", err)
		}
		opts := &Options{Verbose: false} // We don't need verbose output in tests.
		if err := generator.Generate(outputDir, opts); err != nil {
			t.Fatalf("generate failed: %v", err)
		}

		// Read the generated file.
		generatedFilePath := filepath.Join(outputDir, "test-bucket.tf")
		generatedBytes, err := os.ReadFile(generatedFilePath)
		if err != nil {
			t.Fatalf("failed to read generated file: %v", err)
		}

		// Read the golden file.
		goldenBytes, err := os.ReadFile(goldenFilePath)
		if err != nil {
			t.Fatalf("failed to read golden file: %v", err)
		}

		// Compare the generated file to the golden file.
		// We normalize line endings to avoid issues between different OSes.
		generatedContent := strings.ReplaceAll(string(generatedBytes), "\r\n", "\n")
		goldenContent := strings.ReplaceAll(string(goldenBytes), "\r\n", "\n")

		if generatedContent != goldenContent {
			t.Errorf("generated content does not match golden file.\n\nGOT:\n%s\n\nWANT:\n%s", generatedContent, goldenContent)
		}
	})
}
