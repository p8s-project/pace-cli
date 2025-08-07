package loader

import (
	"os"
	"testing"
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
stack: aws
tap: community
resources:
  - id: test-bucket
    uses: s3-bucket
`)
		if _, err := tmpfile.Write(content); err != nil {
			t.Fatal(err)
		}
		if err := tmpfile.Close(); err != nil {
			t.Fatal(err)
		}

		manifest, err := LoadAppManifest(tmpfile.Name())
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

func TestLoadCatalog(t *testing.T) {
	t.Run("valid catalog", func(t *testing.T) {
		tmpfile, err := os.CreateTemp("", "catalog.yaml")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpfile.Name())

		content := []byte(`
resources:
  s3-bucket:
    source: "some-source"
    version: "1.0.0"
    inputs:
      - from: "name"
        to: "bucket"
`)
		if _, err := tmpfile.Write(content); err != nil {
			t.Fatal(err)
		}
		if err := tmpfile.Close(); err != nil {
			t.Fatal(err)
		}

		catalog, err := LoadCatalog(tmpfile.Name())
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		spec, ok := catalog.Resources["s3-bucket"]
		if !ok {
			t.Fatal("expected 's3-bucket' resource to be in catalog")
		}
		if len(spec.Inputs) != 1 {
			t.Error("expected 1 input spec")
		}
	})
}
