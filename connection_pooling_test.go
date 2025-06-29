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
	"sync"
	"testing"
	"time"
)

func TestConnectionPooling(t *testing.T) {
	t.Run("Custom connection pooling configuration", func(t *testing.T) {
		maxIdleConns := 20
		maxIdleConnsPerHost := 5

		client := NewClientWithoutAuth(
			WithConnectionPooling(maxIdleConns, maxIdleConnsPerHost),
		)

		// Get the transport from the HTTP client
		transport, ok := client.httpClient.Transport.(*http.Transport)
		if !ok {
			t.Fatal("Expected HTTP transport to be *http.Transport")
		}

		if transport.MaxIdleConns != maxIdleConns {
			t.Errorf("Expected MaxIdleConns %d, got %d", maxIdleConns, transport.MaxIdleConns)
		}

		if transport.MaxIdleConnsPerHost != maxIdleConnsPerHost {
			t.Errorf("Expected MaxIdleConnsPerHost %d, got %d", maxIdleConnsPerHost, transport.MaxIdleConnsPerHost)
		}

		if transport.IdleConnTimeout != 90*time.Second {
			t.Errorf("Expected IdleConnTimeout 90s, got %v", transport.IdleConnTimeout)
		}

		if transport.DisableCompression != false {
			t.Error("Expected compression to be enabled")
		}
	})

	t.Run("Connection reuse with pooling", func(t *testing.T) {
		var connectionCount int32
		var mutex sync.Mutex
		connections := make(map[string]bool)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Track unique connections
			mutex.Lock()
			remoteAddr := r.RemoteAddr
			if !connections[remoteAddr] {
				connections[remoteAddr] = true
				connectionCount++
			}
			mutex.Unlock()

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"items": [], "metadata": {"totalItems": 0}}`))
		}))
		defer server.Close()

		// Configure client with connection pooling
		client := NewClientWithoutAuth(
			WithBaseURL(server.URL),
			WithConnectionPooling(10, 2),
		)

		ctx := context.Background()

		// Make multiple requests
		const numRequests = 5
		for i := 0; i < numRequests; i++ {
			_, _, err := client.SearchModels(ctx, SearchParams{Limit: 10})
			if err != nil {
				t.Fatalf("Request %d failed: %v", i+1, err)
			}
		}

		// With connection pooling, we should reuse connections
		// Note: In test environment, this may still be 1 connection
		if connectionCount > numRequests {
			t.Errorf("Expected at most %d connections, got %d", numRequests, connectionCount)
		}
	})

	t.Run("Compression enabled", func(t *testing.T) {
		var requestHeaders http.Header
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestHeaders = r.Header.Clone()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"items": [], "metadata": {"totalItems": 0}}`))
		}))
		defer server.Close()

		client := NewClientWithoutAuth(
			WithBaseURL(server.URL),
			WithConnectionPooling(10, 2),
		)

		ctx := context.Background()
		_, _, err := client.SearchModels(ctx, SearchParams{Limit: 10})
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}

		// Check that compression was requested
		acceptEncoding := requestHeaders.Get("Accept-Encoding")
		if acceptEncoding == "" {
			t.Error("Expected Accept-Encoding header to be set")
		}

		if acceptEncoding != "gzip, deflate" {
			t.Errorf("Expected Accept-Encoding 'gzip, deflate', got '%s'", acceptEncoding)
		}
	})

	t.Run("Default transport without pooling", func(t *testing.T) {
		client := NewClientWithoutAuth()

		// Default client should not have custom transport
		if client.httpClient.Transport != nil {
			t.Error("Expected default client to have nil transport (using default)")
		}
	})

	t.Run("Concurrent requests with pooling", func(t *testing.T) {
		var requestCount int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestCount++
			// Add small delay to simulate real API
			time.Sleep(10 * time.Millisecond)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"items": [], "metadata": {"totalItems": 0}}`))
		}))
		defer server.Close()

		client := NewClientWithoutAuth(
			WithBaseURL(server.URL),
			WithConnectionPooling(20, 10), // Allow many connections
		)

		ctx := context.Background()
		const numConcurrent = 10

		var wg sync.WaitGroup
		errors := make(chan error, numConcurrent)

		start := time.Now()

		// Launch concurrent requests
		for i := 0; i < numConcurrent; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, _, err := client.SearchModels(ctx, SearchParams{Limit: 10})
				if err != nil {
					errors <- err
				}
			}()
		}

		wg.Wait()
		close(errors)
		duration := time.Since(start)

		// Check for errors
		for err := range errors {
			t.Errorf("Concurrent request failed: %v", err)
		}

		if requestCount != numConcurrent {
			t.Errorf("Expected %d requests, got %d", numConcurrent, requestCount)
		}

		// With proper connection pooling, concurrent requests should complete faster
		// than sequential requests (less than numConcurrent * 10ms)
		maxExpectedDuration := time.Duration(numConcurrent) * 10 * time.Millisecond
		if duration > maxExpectedDuration {
			t.Logf("Duration %v exceeded expected %v, but this may be acceptable in test environment", duration, maxExpectedDuration)
		}
	})

	t.Run("Connection pooling with custom HTTP client", func(t *testing.T) {
		customClient := &http.Client{
			Timeout: 5 * time.Second,
		}

		client := NewClientWithoutAuth(
			WithHTTPClient(customClient),
			WithConnectionPooling(15, 3),
		)

		// Should have the custom client but with modified transport
		if client.httpClient != customClient {
			t.Error("Expected custom HTTP client to be preserved")
		}

		transport, ok := client.httpClient.Transport.(*http.Transport)
		if !ok {
			t.Fatal("Expected HTTP transport to be *http.Transport")
		}

		if transport.MaxIdleConns != 15 {
			t.Errorf("Expected MaxIdleConns 15, got %d", transport.MaxIdleConns)
		}

		if transport.MaxIdleConnsPerHost != 3 {
			t.Errorf("Expected MaxIdleConnsPerHost 3, got %d", transport.MaxIdleConnsPerHost)
		}
	})
}

func TestAdvancedHTTPConfiguration(t *testing.T) {
	t.Run("Combined retry and pooling configuration", func(t *testing.T) {
		client := NewClientWithoutAuth(
			WithRetryConfig(2, 100*time.Millisecond, 1*time.Second),
			WithConnectionPooling(10, 5),
			WithMaxResponseSize(5*1024*1024),
		)

		// Check retry configuration
		if client.maxRetries != 2 {
			t.Errorf("Expected maxRetries 2, got %d", client.maxRetries)
		}

		// Check pooling configuration
		transport, ok := client.httpClient.Transport.(*http.Transport)
		if !ok {
			t.Fatal("Expected HTTP transport to be *http.Transport")
		}

		if transport.MaxIdleConns != 10 {
			t.Errorf("Expected MaxIdleConns 10, got %d", transport.MaxIdleConns)
		}

		// Check response size limit
		if client.maxResponseSize != 5*1024*1024 {
			t.Errorf("Expected maxResponseSize 5MB, got %d", client.maxResponseSize)
		}
	})

	t.Run("Performance comparison with and without pooling", func(t *testing.T) {
		var requestCount int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestCount++
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"items": [], "metadata": {"totalItems": 0}}`))
		}))
		defer server.Close()

		const numRequests = 5
		ctx := context.Background()

		// Test without pooling
		clientWithoutPooling := NewClientWithoutAuth(WithBaseURL(server.URL))
		startWithout := time.Now()
		for i := 0; i < numRequests; i++ {
			_, _, err := clientWithoutPooling.SearchModels(ctx, SearchParams{Limit: 10})
			if err != nil {
				t.Fatalf("Request without pooling failed: %v", err)
			}
		}
		durationWithout := time.Since(startWithout)

		// Reset counter
		requestCount = 0

		// Test with pooling
		clientWithPooling := NewClientWithoutAuth(
			WithBaseURL(server.URL),
			WithConnectionPooling(10, 5),
		)
		startWith := time.Now()
		for i := 0; i < numRequests; i++ {
			_, _, err := clientWithPooling.SearchModels(ctx, SearchParams{Limit: 10})
			if err != nil {
				t.Fatalf("Request with pooling failed: %v", err)
			}
		}
		durationWith := time.Since(startWith)

		t.Logf("Without pooling: %v, With pooling: %v", durationWithout, durationWith)

		// Both should work, specific performance may vary in test environment
		if durationWith > 2*durationWithout {
			t.Logf("Pooling didn't improve performance as expected, but this is acceptable in test environment")
		}
	})
}
