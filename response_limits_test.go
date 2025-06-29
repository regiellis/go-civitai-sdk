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

func TestResponseSizeLimits(t *testing.T) {
	// Create a large response that exceeds the limit
	largeResponse := `{"items": [` + strings.Repeat(`{"id": 1, "name": "test"},`, 10000) + `], "metadata": {"totalItems": 10000}}`

	t.Run("Response size limit exceeded", func(t *testing.T) {
		// Create a mock server that returns a large response
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(largeResponse))
		}))
		defer server.Close()

		// Create client with small response size limit (1KB)
		client := NewClientWithoutAuth(
			WithBaseURL(server.URL),
			WithMaxResponseSize(1024), // 1KB limit
		)
		ctx := context.Background()

		// Try to search models - should fail due to size limit
		_, _, err := client.SearchModels(ctx, SearchParams{Limit: 10})
		
		if err == nil {
			t.Error("Expected error due to response size limit, got nil")
		}
		
		if !strings.Contains(err.Error(), "response size exceeded") {
			t.Errorf("Expected 'response size exceeded' in error, got: %s", err.Error())
		}
	})

	t.Run("Response within size limit", func(t *testing.T) {
		smallResponse := `{"items": [{"id": 1, "name": "Test Model", "type": "Checkpoint", "createdAt": "2024-01-01T00:00:00Z", "updatedAt": "2024-01-01T00:00:00Z"}], "metadata": {"totalItems": 1}}`

		// Create a mock server that returns a small response
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(smallResponse))
		}))
		defer server.Close()

		// Create client with generous response size limit (10MB)
		client := NewClientWithoutAuth(
			WithBaseURL(server.URL),
			WithMaxResponseSize(10*1024*1024), // 10MB limit
		)
		ctx := context.Background()

		// Try to search models - should succeed
		models, metadata, err := client.SearchModels(ctx, SearchParams{Limit: 10})
		
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		
		if len(models) != 1 {
			t.Errorf("Expected 1 model, got %d", len(models))
		}
		
		if metadata.TotalItems != 1 {
			t.Errorf("Expected metadata total items 1, got %d", metadata.TotalItems)
		}
	})

	t.Run("Error response within size limit", func(t *testing.T) {
		errorResponse := `{"code": "INVALID_REQUEST", "message": "Invalid request parameters"}`

		// Create a mock server that returns an error
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errorResponse))
		}))
		defer server.Close()

		// Create client with small response size limit
		client := NewClientWithoutAuth(
			WithBaseURL(server.URL),
			WithMaxResponseSize(1024), // 1KB limit
		)
		ctx := context.Background()

		// Try to search models - should get API error, not size limit error
		_, _, err := client.SearchModels(ctx, SearchParams{Limit: 10})
		
		if err == nil {
			t.Error("Expected API error, got nil")
		}
		
		if !strings.Contains(err.Error(), "INVALID_REQUEST") {
			t.Errorf("Expected 'INVALID_REQUEST' in error, got: %s", err.Error())
		}
		
		// Should not contain size limit error
		if strings.Contains(err.Error(), "response size exceeded") {
			t.Errorf("Should not contain size limit error, got: %s", err.Error())
		}
	})

	t.Run("Default response size limit", func(t *testing.T) {
		client := NewClientWithoutAuth()
		
		// Check that default limit is set
		if client.maxResponseSize != DefaultMaxResponseSize {
			t.Errorf("Expected default max response size %d, got %d", DefaultMaxResponseSize, client.maxResponseSize)
		}
	})

	t.Run("Custom response size limit option", func(t *testing.T) {
		customLimit := int64(5 * 1024 * 1024) // 5MB
		client := NewClientWithoutAuth(WithMaxResponseSize(customLimit))
		
		// Check that custom limit is set
		if client.maxResponseSize != customLimit {
			t.Errorf("Expected custom max response size %d, got %d", customLimit, client.maxResponseSize)
		}
	})
}
