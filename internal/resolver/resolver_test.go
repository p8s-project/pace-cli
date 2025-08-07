package resolver

import (
	"testing"

	"github.com/Vezia/vez-cli/internal/types"
)

func TestResolveUses(t *testing.T) {
	catalog := &types.Catalog{
		Resources: map[string]types.ResourceSpec{
			"community/aws/data-management/s3-bucket": {},
		},
	}
	app := &types.AppManifest{
		Stack: "aws",
		Tap:   "community",
	}

	t.Run("shorthand", func(t *testing.T) {
		r := New(app, catalog)
		resolved, err := r.ResolveUses("s3-bucket")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if resolved != "community/aws/data-management/s3-bucket" {
			t.Errorf("expected resolved to be 'community/aws/data-management/s3-bucket', got '%s'", resolved)
		}
	})

	t.Run("fully qualified", func(t *testing.T) {
		r := New(app, catalog)
		resolved, err := r.ResolveUses("community/aws/data-management/s3-bucket")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if resolved != "community/aws/data-management/s3-bucket" {
			t.Errorf("expected resolved to be 'community/aws/data-management/s3-bucket', got '%s'", resolved)
		}
	})

	t.Run("wrong stack", func(t *testing.T) {
		r := New(app, catalog)
		_, err := r.ResolveUses("community/gcp/data-management/gcs-bucket")
		if err == nil {
			t.Error("expected an error, got nil")
		}
	})
}
