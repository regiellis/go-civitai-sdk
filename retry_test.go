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
	"sync/atomic"
	"testing"
	"time"
)

func TestRetryLogic(t *testing.T) {
	t.Run("Successful request on first attempt", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"items": [], "metadata": {"totalItems": 0}}`))
		}))
		defer server.Close()

		client := NewClientWithoutAuth(
			WithBaseURL(server.URL),
			WithRetryConfig(3, 100*time.Millisecond, 1*time.Second),
		)

		ctx := context.Background()
		_, _, err := client.SearchModels(ctx, SearchParams{Limit: 10})
		if err != nil {
			t.Errorf("Expected successful request, got error: %v", err)
		}
	})

	t.Run("Retry on server error", func(t *testing.T) {
		var attempts int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attempt := atomic.AddInt32(&attempts, 1)
			if attempt <= 2 {
				// Fail first two attempts
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			// Success on third attempt
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"items": [], "metadata": {"totalItems": 0}}`))
		}))
		defer server.Close()

		client := NewClientWithoutAuth(
			WithBaseURL(server.URL),
			WithRetryConfig(3, 100*time.Millisecond, 1*time.Second),
		)

		ctx := context.Background()
		start := time.Now()
		_, _, err := client.SearchModels(ctx, SearchParams{Limit: 10})
		duration := time.Since(start)

		if err != nil {
			t.Errorf("Expected successful request after retries, got error: %v", err)
		}

		if attempts != 3 {
			t.Errorf("Expected 3 attempts, got %d", attempts)
		}

		// Should have taken some time due to backoff delays
		if duration < 200*time.Millisecond {
			t.Errorf("Expected duration >= 200ms due to retries, got %v", duration)
		}
	})

	t.Run("Retry on rate limiting", func(t *testing.T) {
		var attempts int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attempt := atomic.AddInt32(&attempts, 1)
			if attempt == 1 {
				// Rate limited on first attempt
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
			// Success on second attempt
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"items": [], "metadata": {"totalItems": 0}}`))
		}))
		defer server.Close()

		client := NewClientWithoutAuth(
			WithBaseURL(server.URL),
			WithRetryConfig(3, 100*time.Millisecond, 1*time.Second),
		)

		ctx := context.Background()
		_, _, err := client.SearchModels(ctx, SearchParams{Limit: 10})

		if err != nil {
			t.Errorf("Expected successful request after retry, got error: %v", err)
		}

		if attempts != 2 {
			t.Errorf("Expected 2 attempts, got %d", attempts)
		}
	})

	t.Run("No retry on client error", func(t *testing.T) {
		var attempts int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			atomic.AddInt32(&attempts, 1)
			w.WriteHeader(http.StatusBadRequest) // Client error, should not retry
		}))
		defer server.Close()

		client := NewClientWithoutAuth(
			WithBaseURL(server.URL),
			WithRetryConfig(3, 100*time.Millisecond, 1*time.Second),
		)

		ctx := context.Background()
		_, _, err := client.SearchModels(ctx, SearchParams{Limit: 10})

		if err == nil {
			t.Error("Expected error for bad request")
		}

		if attempts != 1 {
			t.Errorf("Expected 1 attempt (no retry), got %d", attempts)
		}
	})

	t.Run("Exhaust all retries", func(t *testing.T) {
		var attempts int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			atomic.AddInt32(&attempts, 1)
			w.WriteHeader(http.StatusInternalServerError) // Always fail
		}))
		defer server.Close()

		client := NewClientWithoutAuth(
			WithBaseURL(server.URL),
			WithRetryConfig(2, 50*time.Millisecond, 500*time.Millisecond),
		)

		ctx := context.Background()
		_, _, err := client.SearchModels(ctx, SearchParams{Limit: 10})

		if err == nil {
			t.Error("Expected error after exhausting retries")
		}

		if !strings.Contains(err.Error(), "after 3 attempts") {
			t.Errorf("Expected error message about attempts, got: %v", err)
		}

		if attempts != 3 { // 2 retries + 1 initial attempt
			t.Errorf("Expected 3 total attempts, got %d", attempts)
		}
	})

	t.Run("Context cancellation during retry", func(t *testing.T) {
		var attempts int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			atomic.AddInt32(&attempts, 1)
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		client := NewClientWithoutAuth(
			WithBaseURL(server.URL),
			WithRetryConfig(5, 1*time.Second, 5*time.Second), // Long delays
		)

		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		_, _, err := client.SearchModels(ctx, SearchParams{Limit: 10})

		if err == nil {
			t.Error("Expected error due to context cancellation")
		}

		if err != context.DeadlineExceeded {
			t.Errorf("Expected context deadline exceeded, got: %v", err)
		}

		// Should have made at least one attempt but not all
		if attempts == 0 {
			t.Error("Expected at least one attempt")
		}
		if attempts > 3 {
			t.Errorf("Expected context cancellation to prevent too many attempts, got %d", attempts)
		}
	})
}

func TestRetryHelperFunctions(t *testing.T) {
	t.Run("isRetryableError", func(t *testing.T) {
		testCases := []struct {
			err      error
			expected bool
		}{
			{nil, false},
			{context.DeadlineExceeded, true},
			{context.Canceled, false},
		}

		for _, tc := range testCases {
			result := isRetryableError(tc.err)
			if result != tc.expected {
				t.Errorf("isRetryableError(%v) = %v, expected %v", tc.err, result, tc.expected)
			}
		}
	})

	t.Run("isRetryableStatusCode", func(t *testing.T) {
		testCases := []struct {
			code     int
			expected bool
		}{
			{http.StatusOK, false},
			{http.StatusBadRequest, false},
			{http.StatusUnauthorized, false},
			{http.StatusForbidden, false},
			{http.StatusNotFound, false},
			{http.StatusTooManyRequests, true},
			{http.StatusInternalServerError, true},
			{http.StatusBadGateway, true},
			{http.StatusServiceUnavailable, true},
			{http.StatusGatewayTimeout, true},
		}

		for _, tc := range testCases {
			result := isRetryableStatusCode(tc.code)
			if result != tc.expected {
				t.Errorf("isRetryableStatusCode(%d) = %v, expected %v", tc.code, result, tc.expected)
			}
		}
	})

	t.Run("calculateBackoffDelay", func(t *testing.T) {
		client := NewClientWithoutAuth(
			WithRetryConfig(3, 100*time.Millisecond, 1*time.Second),
		)

		// Test exponential backoff
		delay1 := client.calculateBackoffDelay(0)
		delay2 := client.calculateBackoffDelay(1)
		delay3 := client.calculateBackoffDelay(2)

		// Base delay should be around 100ms (with jitter)
		if delay1 < 75*time.Millisecond || delay1 > 125*time.Millisecond {
			t.Errorf("First delay %v should be around 100ms ±25%%", delay1)
		}

		// Second delay should be roughly double (with jitter)
		if delay2 < 150*time.Millisecond || delay2 > 250*time.Millisecond {
			t.Errorf("Second delay %v should be around 200ms ±25%%", delay2)
		}

		// Third delay should be roughly quadruple (with jitter)
		if delay3 < 300*time.Millisecond || delay3 > 500*time.Millisecond {
			t.Errorf("Third delay %v should be around 400ms ±25%%", delay3)
		}

		// Test maximum delay cap
		delay10 := client.calculateBackoffDelay(10)
		if delay10 > client.maxRetryDelay {
			t.Errorf("Delay %v should not exceed max delay %v", delay10, client.maxRetryDelay)
		}
	})
}

func TestRetryConfiguration(t *testing.T) {
	t.Run("Default retry configuration", func(t *testing.T) {
		client := NewClientWithoutAuth()

		if client.maxRetries != DefaultMaxRetries {
			t.Errorf("Expected default max retries %d, got %d", DefaultMaxRetries, client.maxRetries)
		}
		if client.retryDelay != DefaultRetryDelay {
			t.Errorf("Expected default retry delay %v, got %v", DefaultRetryDelay, client.retryDelay)
		}
		if client.maxRetryDelay != DefaultMaxRetryDelay {
			t.Errorf("Expected default max retry delay %v, got %v", DefaultMaxRetryDelay, client.maxRetryDelay)
		}
	})

	t.Run("Custom retry configuration", func(t *testing.T) {
		maxRetries := 5
		baseDelay := 200 * time.Millisecond
		maxDelay := 10 * time.Second

		client := NewClientWithoutAuth(
			WithRetryConfig(maxRetries, baseDelay, maxDelay),
		)

		if client.maxRetries != maxRetries {
			t.Errorf("Expected max retries %d, got %d", maxRetries, client.maxRetries)
		}
		if client.retryDelay != baseDelay {
			t.Errorf("Expected retry delay %v, got %v", baseDelay, client.retryDelay)
		}
		if client.maxRetryDelay != maxDelay {
			t.Errorf("Expected max retry delay %v, got %v", maxDelay, client.maxRetryDelay)
		}
	})

	t.Run("Zero retries configuration", func(t *testing.T) {
		var attempts int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			atomic.AddInt32(&attempts, 1)
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		client := NewClientWithoutAuth(
			WithBaseURL(server.URL),
			WithRetryConfig(0, 100*time.Millisecond, 1*time.Second), // No retries
		)

		ctx := context.Background()
		_, _, err := client.SearchModels(ctx, SearchParams{Limit: 10})

		if err == nil {
			t.Error("Expected error with no retries")
		}

		if attempts != 1 {
			t.Errorf("Expected exactly 1 attempt with no retries, got %d", attempts)
		}
	})
}
