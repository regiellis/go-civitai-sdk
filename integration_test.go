/*
Copyright (c) 2025 Regi Ellis

This file is part of Go CivitAI SDK.

Licensed under the Restricted Use License - Non-Commercial Only.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://github.com/regiellis/go-civitai-sdk/blob/main/LICENSE

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Original work by Regi Ellis (https://github.com/regiellis)
*/

package civitai

import (
	"context"
	"testing"
	"time"
)

// TestIntegration runs integration tests against the real CivitAI API
func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	client := NewClientWithoutAuth()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Run("SearchModels", func(t *testing.T) {
		params := SearchParams{
			Query: "anime",
			Types: []ModelType{ModelTypeCheckpoint},
			Limit: 5,
		}

		models, metadata, err := client.SearchModels(ctx, params)
		if err != nil {
			t.Fatalf("SearchModels failed: %v", err)
		}

		if len(models) == 0 {
			t.Error("Expected at least one model")
		}

		if metadata.TotalItems == 0 {
			t.Error("Expected metadata.TotalItems > 0")
		}

		// Validate model structure
		for _, model := range models {
			if model.ID == 0 {
				t.Error("Model ID should not be 0")
			}
			if model.Name == "" {
				t.Error("Model name should not be empty")
			}
			if model.Type == "" {
				t.Error("Model type should not be empty")
			}
		}
	})

	t.Run("GetModel", func(t *testing.T) {
		// First get a model ID from search
		searchParams := SearchParams{
			Query: "realistic",
			Types: []ModelType{ModelTypeCheckpoint},
			Limit: 1,
		}

		models, _, err := client.SearchModels(ctx, searchParams)
		if err != nil {
			t.Fatalf("Failed to search for model: %v", err)
		}

		if len(models) == 0 {
			t.Skip("No models found to test GetModel")
		}

		modelID := models[0].ID
		model, err := client.GetModel(ctx, modelID)
		if err != nil {
			t.Fatalf("GetModel failed: %v", err)
		}

		if model.ID != modelID {
			t.Errorf("Expected model ID %d, got %d", modelID, model.ID)
		}

		if model.Name == "" {
			t.Error("Model name should not be empty")
		}

		if len(model.ModelVersions) == 0 {
			t.Error("Expected at least one model version")
		}
	})

	t.Run("GetModelVersion", func(t *testing.T) {
		// Get a model with versions
		searchParams := SearchParams{
			Query: "realistic",
			Types: []ModelType{ModelTypeCheckpoint},
			Limit: 1,
		}

		models, _, err := client.SearchModels(ctx, searchParams)
		if err != nil {
			t.Fatalf("Failed to search for model: %v", err)
		}

		if len(models) == 0 || len(models[0].ModelVersions) == 0 {
			t.Skip("No model versions found to test")
		}

		versionID := models[0].ModelVersions[0].ID
		version, err := client.GetModelVersion(ctx, versionID)
		if err != nil {
			t.Fatalf("GetModelVersion failed: %v", err)
		}

		if version.ID != versionID {
			t.Errorf("Expected version ID %d, got %d", versionID, version.ID)
		}

		if version.Name == "" {
			t.Error("Version name should not be empty")
		}
	})

	t.Run("GetImages", func(t *testing.T) {
		params := ImageParams{
			Sort:  string(ImageSortNewest),
			NSFW:  string(NSFWLevelNone),
			Limit: 5,
		}

		images, metadata, err := client.GetImages(ctx, params)
		if err != nil {
			t.Fatalf("GetImages failed: %v", err)
		}

		if len(images) == 0 {
			t.Error("Expected at least one image")
		}

		// Validate image structure
		for _, image := range images {
			if image.ID == 0 {
				t.Error("Image ID should not be 0")
			}
			if image.Width == 0 || image.Height == 0 {
				t.Error("Image dimensions should not be 0")
			}
			if image.Username == "" {
				t.Error("Image username should not be empty")
			}
		}

		if metadata.NextPage == "" {
			t.Log("No next page in metadata (this is OK)")
		}
	})

	t.Run("GetCreators", func(t *testing.T) {
		params := CreatorParams{
			Limit: 5,
		}

		creators, metadata, err := client.GetCreators(ctx, params)
		if err != nil {
			t.Fatalf("GetCreators failed: %v", err)
		}

		if len(creators) == 0 {
			t.Error("Expected at least one creator")
		}

		if metadata.TotalItems == 0 {
			t.Error("Expected metadata.TotalItems > 0")
		}

		// Validate creator structure
		for _, creator := range creators {
			if creator.Username == "" {
				t.Error("Creator username should not be empty")
			}
			if creator.ModelCount < 0 {
				t.Error("Creator model count should not be negative")
			}
		}
	})

	t.Run("GetTags", func(t *testing.T) {
		params := TagParams{
			Query: "anime",
			Limit: 10,
		}

		tags, metadata, err := client.GetTags(ctx, params)
		if err != nil {
			t.Fatalf("GetTags failed: %v", err)
		}

		if len(tags) == 0 {
			t.Error("Expected at least one tag")
		}

		if metadata.TotalItems == 0 {
			t.Error("Expected metadata.TotalItems > 0")
		}

		// Validate tag structure
		for _, tag := range tags {
			if tag.Name == "" {
				t.Error("Tag name should not be empty")
			}
			if tag.ModelCount < 0 {
				t.Error("Tag model count should not be negative")
			}
		}
	})

	t.Run("Health", func(t *testing.T) {
		err := client.Health(ctx)
		if err != nil {
			t.Fatalf("Health check failed: %v", err)
		}
	})
}

// TestClientConfiguration tests various client configuration options
func TestClientConfiguration(t *testing.T) {
	t.Run("NewClientWithoutAuth", func(t *testing.T) {
		client := NewClientWithoutAuth()
		if client == nil {
			t.Error("Expected non-nil client")
		}
	})

	t.Run("NewClientWithAuth", func(t *testing.T) {
		client := NewClient("test-token")
		if client == nil {
			t.Error("Expected non-nil client")
		}
	})

	t.Run("ClientWithOptions", func(t *testing.T) {
		client := NewClient("test-token",
			WithTimeout(60*time.Second),
			WithUserAgent("test-agent/1.0.0"),
			WithBaseURL("https://example.com/api/v1"),
		)
		if client == nil {
			t.Error("Expected non-nil client")
		}
	})
}

// TestErrorHandling tests error scenarios
func TestErrorHandling(t *testing.T) {
	client := NewClient("invalid-token",
		WithBaseURL("https://invalid-url-that-does-not-exist.com/api/v1"),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t.Run("InvalidURL", func(t *testing.T) {
		params := SearchParams{
			Query: "test",
			Limit: 1,
		}

		_, _, err := client.SearchModels(ctx, params)
		if err == nil {
			t.Error("Expected error for invalid URL")
		}
	})

	t.Run("ContextCancellation", func(t *testing.T) {
		cancelCtx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		params := SearchParams{
			Query: "test",
			Limit: 1,
		}

		_, _, err := client.SearchModels(cancelCtx, params)
		if err == nil {
			t.Error("Expected error for cancelled context")
		}
	})
}

// BenchmarkSearchModels benchmarks the SearchModels operation
func BenchmarkSearchModels(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping benchmark in short mode")
	}

	client := NewClientWithoutAuth()
	ctx := context.Background()
	params := SearchParams{
		Query: "anime",
		Limit: 10,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := client.SearchModels(ctx, params)
		if err != nil {
			b.Fatalf("SearchModels failed: %v", err)
		}
	}
}
