package builder

import (
	"testing"

	"github.com/Vezia/vez-cli/internal/types"
)

func TestBuildInputs(t *testing.T) {
	catalog := &types.Catalog{
		Resources: map[string]types.ResourceSpec{
			"community/aws/data-management/s3-bucket": {
				Inputs: []types.InputSpec{
					{
						From: "id",
						To:   "bucket",
					},
					{
						From: "versioning",
						To:   "versioning.enabled",
					},
				},
			},
		},
	}
	app := &types.AppManifest{
		Stack: "aws",
		Tap:   "community",
		Resources: []types.ResourceRequest{
			{
				ID:   "my-bucket",
				Uses: "s3-bucket",
				With: map[string]interface{}{
					"versioning": true,
				},
			},
		},
	}
	outputs := make(map[string]map[string]string)

	t.Run("happy path", func(t *testing.T) {
		b := New(app, catalog)
		inputs, err := b.BuildInputs(app.Resources[0], catalog.Resources["community/aws/data-management/s3-bucket"], outputs)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if inputs["bucket"] != "my-bucket" {
			t.Errorf("expected bucket to be 'my-bucket', got '%s'", inputs["bucket"])
		}
		if inputs["versioning.enabled"] != true {
			t.Errorf("expected versioning.enabled to be true, got %v", inputs["versioning.enabled"])
		}
	})
}
