/*
Copyright (c) 2025 Regi Ellis

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package civitai

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAPIMethodsWithMockServer(t *testing.T) {
	// Create a mock server for testing
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "/models") && r.Method == "GET":
			if strings.Contains(r.URL.Path, "/versions") {
				// GetModelVersionsByModelID
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`[{"id": 1, "name": "Version 1.0", "createdAt": "2024-01-01T00:00:00Z", "updatedAt": "2024-01-01T00:00:00Z"}]`))
			} else if len(r.URL.Path) > 8 { // Specific model ID
				// GetModel
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"id": 123, "name": "Test Model", "type": "Checkpoint", "createdAt": "2024-01-01T00:00:00Z", "updatedAt": "2024-01-01T00:00:00Z"}`))
			} else {
				// SearchModels
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"items": [{"id": 1, "name": "Test Model", "type": "Checkpoint", "createdAt": "2024-01-01T00:00:00Z", "updatedAt": "2024-01-01T00:00:00Z"}], "metadata": {"totalItems": 1}}`))
			}
		case strings.Contains(r.URL.Path, "/model-versions"):
			if strings.Contains(r.URL.Path, "/by-hash/") {
				// GetModelVersionByHash
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"id": 456, "name": "Version by hash", "createdAt": "2024-01-01T00:00:00Z", "updatedAt": "2024-01-01T00:00:00Z", "model": {"name": "Test Model", "type": "Checkpoint"}}`))
			} else {
				// GetModelVersion
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"id": 456, "name": "Test Version", "createdAt": "2024-01-01T00:00:00Z", "updatedAt": "2024-01-01T00:00:00Z"}`))
			}
		case strings.Contains(r.URL.Path, "/images"):
			// GetImages
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"items": [{"id": 1, "url": "https://example.com/image.jpg", "width": 512, "height": 512, "createdAt": "2024-01-01T00:00:00Z", "username": "testuser"}], "metadata": {"totalItems": 1}}`))
		case strings.Contains(r.URL.Path, "/creators"):
			// GetCreators
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"items": [{"username": "testcreator", "modelCount": 5, "link": "https://civitai.com/user/testcreator"}], "metadata": {"totalItems": 1}}`))
		case strings.Contains(r.URL.Path, "/tags"):
			// GetTags
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"items": [{"name": "anime", "modelCount": 100, "link": "https://civitai.com/tag/anime"}], "metadata": {"totalItems": 1}}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	// Create client with mock server URL
	client := NewClientWithoutAuth(WithBaseURL(server.URL))
	ctx := context.Background()

	t.Run("GetModel", func(t *testing.T) {
		model, err := client.GetModel(ctx, 123)
		if err != nil {
			t.Fatalf("GetModel failed: %v", err)
		}

		if model.ID != 123 {
			t.Errorf("Expected model ID 123, got %d", model.ID)
		}
		if model.Name != "Test Model" {
			t.Errorf("Expected model name 'Test Model', got %s", model.Name)
		}
	})

	t.Run("GetModelVersion", func(t *testing.T) {
		version, err := client.GetModelVersion(ctx, 456)
		if err != nil {
			t.Fatalf("GetModelVersion failed: %v", err)
		}

		if version.ID != 456 {
			t.Errorf("Expected version ID 456, got %d", version.ID)
		}
		if version.Name != "Test Version" {
			t.Errorf("Expected version name 'Test Version', got %s", version.Name)
		}
	})

	t.Run("GetModelVersionsByModelID", func(t *testing.T) {
		versions, err := client.GetModelVersionsByModelID(ctx, 123)
		if err != nil {
			t.Fatalf("GetModelVersionsByModelID failed: %v", err)
		}

		if len(versions) != 1 {
			t.Errorf("Expected 1 version, got %d", len(versions))
		}
		if versions[0].ID != 1 {
			t.Errorf("Expected version ID 1, got %d", versions[0].ID)
		}
	})

	t.Run("GetModelVersionByHash", func(t *testing.T) {
		version, err := client.GetModelVersionByHash(ctx, "abcdef1234567890")
		if err != nil {
			t.Fatalf("GetModelVersionByHash failed: %v", err)
		}

		if version.ID != 456 {
			t.Errorf("Expected version ID 456, got %d", version.ID)
		}
		if version.Name != "Version by hash" {
			t.Errorf("Expected version name 'Version by hash', got %s", version.Name)
		}
	})

	t.Run("GetImages", func(t *testing.T) {
		images, metadata, err := client.GetImages(ctx, ImageParams{Limit: 10})
		if err != nil {
			t.Fatalf("GetImages failed: %v", err)
		}

		if len(images) != 1 {
			t.Errorf("Expected 1 image, got %d", len(images))
		}
		if images[0].ID != 1 {
			t.Errorf("Expected image ID 1, got %d", images[0].ID)
		}
		if metadata.TotalItems != 1 {
			t.Errorf("Expected metadata total items 1, got %d", metadata.TotalItems)
		}
	})

	t.Run("GetCreators", func(t *testing.T) {
		creators, metadata, err := client.GetCreators(ctx, CreatorParams{Limit: 10})
		if err != nil {
			t.Fatalf("GetCreators failed: %v", err)
		}

		if len(creators) != 1 {
			t.Errorf("Expected 1 creator, got %d", len(creators))
		}
		if creators[0].Username != "testcreator" {
			t.Errorf("Expected creator username 'testcreator', got %s", creators[0].Username)
		}
		if metadata.TotalItems != 1 {
			t.Errorf("Expected metadata total items 1, got %d", metadata.TotalItems)
		}
	})

	t.Run("GetTags", func(t *testing.T) {
		tags, metadata, err := client.GetTags(ctx, TagParams{Limit: 10})
		if err != nil {
			t.Fatalf("GetTags failed: %v", err)
		}

		if len(tags) != 1 {
			t.Errorf("Expected 1 tag, got %d", len(tags))
		}
		if tags[0].Name != "anime" {
			t.Errorf("Expected tag name 'anime', got %s", tags[0].Name)
		}
		if metadata.TotalItems != 1 {
			t.Errorf("Expected metadata total items 1, got %d", metadata.TotalItems)
		}
	})
}

func TestParameterValidation(t *testing.T) {
	client := NewClientWithoutAuth()

	t.Run("ValidateImageParams", func(t *testing.T) {
		// Test valid params
		validParams := ImageParams{Limit: 10, Page: 1}
		err := client.validateImageParams(validParams)
		if err != nil {
			t.Errorf("Expected valid params to pass, got error: %v", err)
		}

		// Test invalid limit
		invalidParams := ImageParams{Limit: 300}
		err = client.validateImageParams(invalidParams)
		if err == nil {
			t.Error("Expected error for limit > 200")
		}

		// Test negative page
		invalidParams = ImageParams{Page: -1}
		err = client.validateImageParams(invalidParams)
		if err == nil {
			t.Error("Expected error for negative page")
		}

		// Test long username
		invalidParams = ImageParams{Username: strings.Repeat("a", 101)}
		err = client.validateImageParams(invalidParams)
		if err == nil {
			t.Error("Expected error for username too long")
		}
	})

	t.Run("ValidateCreatorParams", func(t *testing.T) {
		// Test valid params
		validParams := CreatorParams{Limit: 10, Page: 1, Query: "test"}
		err := client.validateCreatorParams(validParams)
		if err != nil {
			t.Errorf("Expected valid params to pass, got error: %v", err)
		}

		// Test invalid limit
		invalidParams := CreatorParams{Limit: -1}
		err = client.validateCreatorParams(invalidParams)
		if err == nil {
			t.Error("Expected error for negative limit")
		}

		// Test long query
		invalidParams = CreatorParams{Query: strings.Repeat("a", 501)}
		err = client.validateCreatorParams(invalidParams)
		if err == nil {
			t.Error("Expected error for query too long")
		}
	})

	t.Run("ValidateTagParams", func(t *testing.T) {
		// Test valid params
		validParams := TagParams{Limit: 10, Page: 1, Query: "test"}
		err := client.validateTagParams(validParams)
		if err != nil {
			t.Errorf("Expected valid params to pass, got error: %v", err)
		}

		// Test invalid limit
		invalidParams := TagParams{Limit: 201}
		err = client.validateTagParams(invalidParams)
		if err == nil {
			t.Error("Expected error for limit > 200")
		}

		// Test negative page
		invalidParams = TagParams{Page: -1}
		err = client.validateTagParams(invalidParams)
		if err == nil {
			t.Error("Expected error for negative page")
		}
	})
}

func TestBuildParams(t *testing.T) {
	client := NewClientWithoutAuth()

	t.Run("buildImageParams", func(t *testing.T) {
		params := ImageParams{
			Limit:          10,
			PostID:         123,
			ModelID:        456,
			ModelVersionID: 789,
			Username:       "testuser",
			NSFW:           "None",
			Sort:           "Newest",
			Period:         PeriodWeek,
			Page:           2,
		}

		queryParams := client.buildImageParams(params)

		if queryParams["limit"] != "10" {
			t.Errorf("Expected limit '10', got '%s'", queryParams["limit"])
		}
		if queryParams["postId"] != "123" {
			t.Errorf("Expected postId '123', got '%s'", queryParams["postId"])
		}
		if queryParams["modelId"] != "456" {
			t.Errorf("Expected modelId '456', got '%s'", queryParams["modelId"])
		}
		if queryParams["username"] != "testuser" {
			t.Errorf("Expected username 'testuser', got '%s'", queryParams["username"])
		}
	})

	t.Run("buildCreatorParams", func(t *testing.T) {
		params := CreatorParams{
			Limit: 20,
			Page:  3,
			Query: "anime",
		}

		queryParams := client.buildCreatorParams(params)

		if queryParams["limit"] != "20" {
			t.Errorf("Expected limit '20', got '%s'", queryParams["limit"])
		}
		if queryParams["page"] != "3" {
			t.Errorf("Expected page '3', got '%s'", queryParams["page"])
		}
		if queryParams["query"] != "anime" {
			t.Errorf("Expected query 'anime', got '%s'", queryParams["query"])
		}
	})

	t.Run("buildTagParams", func(t *testing.T) {
		params := TagParams{
			Limit: 15,
			Page:  1,
			Query: "style",
		}

		queryParams := client.buildTagParams(params)

		if queryParams["limit"] != "15" {
			t.Errorf("Expected limit '15', got '%s'", queryParams["limit"])
		}
		if queryParams["page"] != "1" {
			t.Errorf("Expected page '1', got '%s'", queryParams["page"])
		}
		if queryParams["query"] != "style" {
			t.Errorf("Expected query 'style', got '%s'", queryParams["query"])
		}
	})
}

func TestAdditionalClientOptions(t *testing.T) {
	t.Run("WithHTTPClient", func(t *testing.T) {
		customClient := &http.Client{}
		client := NewClientWithoutAuth(WithHTTPClient(customClient))

		if client.httpClient != customClient {
			t.Error("Expected custom HTTP client to be set")
		}
	})

	t.Run("GetAPIToken", func(t *testing.T) {
		token := "test-token-123"
		client := NewClient(token)

		if client.GetAPIToken() != token {
			t.Errorf("Expected token '%s', got '%s'", token, client.GetAPIToken())
		}
	})
}

func TestAPIErrorHandling(t *testing.T) {
	// Create mock server that returns errors
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"code": "INVALID_REQUEST", "message": "Invalid request parameters"}`))
	}))
	defer server.Close()

	client := NewClientWithoutAuth(WithBaseURL(server.URL))
	ctx := context.Background()

	t.Run("API Error Response", func(t *testing.T) {
		_, err := client.GetModel(ctx, 123)
		if err == nil {
			t.Error("Expected error from API, got nil")
		}

		if !strings.Contains(err.Error(), "INVALID_REQUEST") {
			t.Errorf("Expected error to contain 'INVALID_REQUEST', got: %s", err.Error())
		}
	})
}

func TestExceptionTypes(t *testing.T) {
	t.Run("NewAPIError", func(t *testing.T) {
		err := NewAPIError(nil, "TEST_CODE", "Test message")

		if err.Code != "TEST_CODE" {
			t.Errorf("Expected code 'TEST_CODE', got '%s'", err.Code)
		}
		if err.Message != "Test message" {
			t.Errorf("Expected message 'Test message', got '%s'", err.Message)
		}
	})

	t.Run("NewAPIErrorWithDetails", func(t *testing.T) {
		err := NewAPIErrorWithDetails("TEST_CODE", "Test message", "Test details")

		if err.Code != "TEST_CODE" {
			t.Errorf("Expected code 'TEST_CODE', got '%s'", err.Code)
		}
		if err.Message != "Test message" {
			t.Errorf("Expected message 'Test message', got '%s'", err.Message)
		}
		if err.Details != "Test details" {
			t.Errorf("Expected details 'Test details', got '%s'", err.Details)
		}
	})

	t.Run("APIError Error method", func(t *testing.T) {
		err := APIError{
			Code:    "TEST_CODE",
			Message: "Test message",
		}

		expected := "CivitAI API error [TEST_CODE]: Test message"
		if err.Error() != expected {
			t.Errorf("Expected '%s', got '%s'", expected, err.Error())
		}
	})
}
