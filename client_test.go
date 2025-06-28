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
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient("test-token")

	if client.apiToken != "test-token" {
		t.Errorf("Expected token 'test-token', got '%s'", client.apiToken)
	}

	if client.baseURL != DefaultBaseURL {
		t.Errorf("Expected base URL '%s', got '%s'", DefaultBaseURL, client.baseURL)
	}
}

func TestNewClientWithoutAuth(t *testing.T) {
	client := NewClientWithoutAuth()

	if client.apiToken != "" {
		t.Errorf("Expected empty token, got '%s'", client.apiToken)
	}
}

func TestClientOptions(t *testing.T) {
	customTimeout := 60 * time.Second
	customUserAgent := "test-agent/1.0"
	customBaseURL := "https://test.api.com"

	client := NewClient("test-token",
		WithTimeout(customTimeout),
		WithUserAgent(customUserAgent),
		WithBaseURL(customBaseURL),
	)

	if client.httpClient.Timeout != customTimeout {
		t.Errorf("Expected timeout %v, got %v", customTimeout, client.httpClient.Timeout)
	}

	if client.userAgent != customUserAgent {
		t.Errorf("Expected user agent '%s', got '%s'", customUserAgent, client.userAgent)
	}

	if client.baseURL != customBaseURL {
		t.Errorf("Expected base URL '%s', got '%s'", customBaseURL, client.baseURL)
	}
}

func TestBuildURL(t *testing.T) {
	client := NewClient("test")

	url := client.buildURL("models/123")
	expected := DefaultBaseURL + "/models/123"

	if url != expected {
		t.Errorf("Expected URL '%s', got '%s'", expected, url)
	}
}

func TestHealth(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/models" {
			t.Errorf("Expected path '/models', got '%s'", r.URL.Path)
		}
		// Return a minimal valid response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"items":[],"metadata":{"totalItems":0,"currentPage":1,"pageSize":0,"totalPages":0}}`))
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewClient("test", WithBaseURL(server.URL))

	// Test health check
	err := client.Health(context.Background())
	if err != nil {
		t.Errorf("Health check failed: %v", err)
	}
}

func TestAPIError(t *testing.T) {
	err := APIError{
		Code:    "VALIDATION_ERROR",
		Message: "Invalid model ID",
		Details: "Model ID must be a positive integer",
	}

	expected := "CivitAI API error [VALIDATION_ERROR]: Invalid model ID - Model ID must be a positive integer"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}

	// Test without details
	err2 := APIError{
		Code:    "NOT_FOUND",
		Message: "Model not found",
	}

	expected2 := "CivitAI API error [NOT_FOUND]: Model not found"
	if err2.Error() != expected2 {
		t.Errorf("Expected error message '%s', got '%s'", expected2, err2.Error())
	}
}
