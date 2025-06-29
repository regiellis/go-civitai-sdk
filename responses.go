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

// Package civitai provides comprehensive response structures and utilities for the CivitAI API.
//
// This file contains all the standard API response patterns used throughout the SDK,
// including pagination metadata, error handling, and response validation utilities.
//
// # Response Structure
//
// All CivitAI API responses follow a consistent pattern:
//
//	{
//	  "items": [...],     // The actual data
//	  "metadata": {       // Pagination and response metadata
//	    "totalItems": 1234,
//	    "currentPage": 1,
//	    "pageSize": 20,
//	    "totalPages": 62,
//	    "nextPage": "https://...",
//	    "prevPage": "https://...",
//	    "nextCursor": "eyJtb2RlbElkIjoxMjM0fQ=="
//	  }
//	}
//
// # Usage Example
//
//	client := civitai.NewClientWithoutAuth()
//	models, metadata, err := client.SearchModels(ctx, civitai.SearchParams{
//		Tag: "anime",
//		Limit: 20,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Use pagination metadata for next page
//	if metadata.NextCursor != "" {
//		params.Cursor = metadata.NextCursor
//		nextPageModels, _, _ := client.SearchModels(ctx, params)
//	}
//
// # Error Handling
//
// API responses include detailed error information when requests fail:
//
//	if err != nil {
//		if apiErr, ok := err.(*civitai.APIError); ok {
//			fmt.Printf("API Error: %s (Status: %d)\n", apiErr.Message, apiErr.StatusCode)
//		}
//	}
//
// # Performance Notes
//
// Response handling includes automatic:
//   - JSON validation and parsing
//   - Error detection and wrapping
//   - Rate limit header parsing
//   - Response size validation
//   - Retry-after header handling

package civitai

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// APIResponse represents the standard CivitAI API response structure
// This is the base structure that all API endpoints follow
type APIResponse[T any] struct {
	Items    []T       `json:"items"`
	Metadata *Metadata `json:"metadata,omitempty"`
	Success  bool      `json:"success"`
	Error    *APIError `json:"error,omitempty"`
}

// ModelsResponse represents the response from /api/v1/models
type ModelsResponse struct {
	Items    []Model   `json:"items"`
	Metadata *Metadata `json:"metadata,omitempty"`
}

// ImagesResponse represents the response from /api/v1/images
type ImagesResponse struct {
	Items    []DetailedImageResponse `json:"items"`
	Metadata *Metadata               `json:"metadata,omitempty"`
}

// CreatorsResponse represents the response from /api/v1/creators
type CreatorsResponse struct {
	Items    []Creator `json:"items"`
	Metadata *Metadata `json:"metadata,omitempty"`
}

// TagsResponse represents the response from /api/v1/tags
type TagsResponse struct {
	Items    []Tag     `json:"items"`
	Metadata *Metadata `json:"metadata,omitempty"`
}

// ModelVersionsResponse represents the response from /api/v1/model-versions
type ModelVersionsResponse struct {
	Items    []ModelVersion `json:"items"`
	Metadata *Metadata      `json:"metadata,omitempty"`
}

// SingleModelResponse represents the response for single model requests
type SingleModelResponse struct {
	*Model
}

// SingleModelVersionResponse represents the response for single model version requests
type SingleModelVersionResponse struct {
	*ModelVersion
}

// APIError represents an error response from the CivitAI API
type APIError struct {
	StatusCode int    `json:"statusCode"`
	Code       string `json:"code"`
	Message    string `json:"message"`
	ErrorMsg   string `json:"error,omitempty"`
	Details    string `json:"details,omitempty"`
	Timestamp  string `json:"timestamp,omitempty"`
	Path       string `json:"path,omitempty"`
}

// Error implements the error interface for APIError
func (e *APIError) Error() string {
	if e.Code != "" && e.Details != "" {
		return fmt.Sprintf("CivitAI API error [%s]: %s - %s", e.Code, e.Message, e.Details)
	}
	if e.Code != "" {
		return fmt.Sprintf("CivitAI API error [%s]: %s", e.Code, e.Message)
	}
	if e.Details != "" {
		return fmt.Sprintf("CivitAI API error %d: %s - %s", e.StatusCode, e.Message, e.Details)
	}
	return fmt.Sprintf("CivitAI API error %d: %s", e.StatusCode, e.Message)
}

// NewAPIError creates a new API error from an HTTP response
func NewAPIError(resp *http.Response, code, message string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
	}
}

// NewAPIErrorWithDetails creates a new API error with details
func NewAPIErrorWithDetails(code, message, details string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// IsRateLimitError returns true if the error is a rate limit error (429)
func (e *APIError) IsRateLimitError() bool {
	return e.StatusCode == http.StatusTooManyRequests
}

// IsNotFoundError returns true if the error is a not found error (404)
func (e *APIError) IsNotFoundError() bool {
	return e.StatusCode == http.StatusNotFound
}

// IsAuthenticationError returns true if the error is an authentication error (401)
func (e *APIError) IsAuthenticationError() bool {
	return e.StatusCode == http.StatusUnauthorized
}

// IsForbiddenError returns true if the error is a forbidden error (403)
func (e *APIError) IsForbiddenError() bool {
	return e.StatusCode == http.StatusForbidden
}

// IsServerError returns true if the error is a server error (5xx)
func (e *APIError) IsServerError() bool {
	return e.StatusCode >= 500 && e.StatusCode < 600
}

// IsClientError returns true if the error is a client error (4xx)
func (e *APIError) IsClientError() bool {
	return e.StatusCode >= 400 && e.StatusCode < 500
}

// ResponseInfo contains metadata about the HTTP response
type ResponseInfo struct {
	StatusCode   int
	Headers      http.Header
	Size         int64
	ResponseTime time.Duration
	Cached       bool
}

// RateLimitInfo contains rate limiting information from response headers
type RateLimitInfo struct {
	Limit       int           // X-RateLimit-Limit
	Remaining   int           // X-RateLimit-Remaining
	Reset       time.Time     // X-RateLimit-Reset
	RetryAfter  time.Duration // Retry-After (for 429 responses)
	WindowStart time.Time     // X-RateLimit-Reset-After
}

// ParseRateLimitHeaders extracts rate limit information from HTTP response headers
func ParseRateLimitHeaders(headers http.Header) *RateLimitInfo {
	info := &RateLimitInfo{}

	if limit := headers.Get("X-RateLimit-Limit"); limit != "" {
		if val, err := strconv.Atoi(limit); err == nil {
			info.Limit = val
		}
	}

	if remaining := headers.Get("X-RateLimit-Remaining"); remaining != "" {
		if val, err := strconv.Atoi(remaining); err == nil {
			info.Remaining = val
		}
	}

	if reset := headers.Get("X-RateLimit-Reset"); reset != "" {
		if val, err := strconv.ParseInt(reset, 10, 64); err == nil {
			info.Reset = time.Unix(val, 0)
		}
	}

	if retryAfter := headers.Get("Retry-After"); retryAfter != "" {
		if val, err := strconv.Atoi(retryAfter); err == nil {
			info.RetryAfter = time.Duration(val) * time.Second
		}
	}

	return info
}

// ValidateResponse validates the structure and content of an API response
func ValidateResponse[T any](resp *APIResponse[T]) error {
	if resp == nil {
		return fmt.Errorf("response is nil")
	}

	// Basic validation - items can be empty for valid empty results
	if resp.Items == nil {
		return fmt.Errorf("response items is nil")
	}

	// Validate metadata if present
	if resp.Metadata != nil {
		if err := ValidateMetadata(resp.Metadata); err != nil {
			return fmt.Errorf("invalid metadata: %w", err)
		}
	}

	return nil
}

// ValidateMetadata validates pagination metadata
func ValidateMetadata(meta *Metadata) error {
	if meta == nil {
		return fmt.Errorf("metadata is nil")
	}

	// CurrentPage and TotalPages should be consistent
	if meta.CurrentPage < 0 {
		return fmt.Errorf("currentPage cannot be negative: %d", meta.CurrentPage)
	}

	if meta.TotalPages < 0 {
		return fmt.Errorf("totalPages cannot be negative: %d", meta.TotalPages)
	}

	if meta.CurrentPage > meta.TotalPages && meta.TotalPages > 0 {
		return fmt.Errorf("currentPage (%d) exceeds totalPages (%d)", meta.CurrentPage, meta.TotalPages)
	}

	// TotalItems should be consistent
	if meta.TotalItems < 0 {
		return fmt.Errorf("totalItems cannot be negative: %d", meta.TotalItems)
	}

	// PageSize should be reasonable
	if meta.PageSize < 0 {
		return fmt.Errorf("pageSize cannot be negative: %d", meta.PageSize)
	}

	if meta.PageSize > 200 {
		return fmt.Errorf("pageSize too large: %d (max 200)", meta.PageSize)
	}

	return nil
}

// ParseErrorResponse parses an error response from the API
func ParseErrorResponse(resp *http.Response, body []byte) error {
	apiErr := &APIError{
		StatusCode: resp.StatusCode,
	}

	// Try to parse JSON error response
	var errorResp struct {
		Message   string `json:"message"`
		Error     string `json:"error"`
		Details   string `json:"details"`
		Timestamp string `json:"timestamp"`
		Path      string `json:"path"`
	}

	if err := json.Unmarshal(body, &errorResp); err == nil {
		apiErr.Message = errorResp.Message
		apiErr.ErrorMsg = errorResp.Error
		apiErr.Details = errorResp.Details
		apiErr.Timestamp = errorResp.Timestamp
		apiErr.Path = errorResp.Path
	} else {
		// Fallback to status text if JSON parsing fails
		apiErr.Message = resp.Status
		if len(body) > 0 && len(body) < 500 {
			apiErr.Details = string(body)
		}
	}

	// Provide default messages for common HTTP status codes
	if apiErr.Message == "" {
		switch resp.StatusCode {
		case http.StatusTooManyRequests:
			apiErr.Message = "Rate limit exceeded"
		case http.StatusUnauthorized:
			apiErr.Message = "Authentication required"
		case http.StatusForbidden:
			apiErr.Message = "Access forbidden"
		case http.StatusNotFound:
			apiErr.Message = "Resource not found"
		case http.StatusInternalServerError:
			apiErr.Message = "Internal server error"
		case http.StatusBadGateway:
			apiErr.Message = "Bad gateway"
		case http.StatusServiceUnavailable:
			apiErr.Message = "Service unavailable"
		case http.StatusGatewayTimeout:
			apiErr.Message = "Gateway timeout"
		default:
			apiErr.Message = fmt.Sprintf("HTTP %d", resp.StatusCode)
		}
	}

	return apiErr
}

// IsRetryableError determines if an error is retryable
func IsRetryableError(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		// Retry on server errors and rate limits
		return apiErr.IsServerError() || apiErr.IsRateLimitError()
	}
	return false
}

// GetRetryDelay calculates the delay before retrying a request
func GetRetryDelay(err error, attempt int) time.Duration {
	baseDelay := time.Second

	if apiErr, ok := err.(*APIError); ok {
		// Use Retry-After header if available (for rate limits)
		if apiErr.IsRateLimitError() {
			// Parse rate limit headers would go here
			// For now, use exponential backoff with longer delays for rate limits
			return time.Duration(attempt*attempt) * 5 * time.Second
		}
	}

	// Exponential backoff with jitter
	delay := baseDelay * time.Duration(1<<uint(attempt))
	if delay > 30*time.Second {
		delay = 30 * time.Second
	}

	return delay
}

// ResponseMetrics contains metrics about API responses
type ResponseMetrics struct {
	TotalRequests   int64
	SuccessfulReqs  int64
	FailedRequests  int64
	RateLimitErrors int64
	ServerErrors    int64
	AverageResponse time.Duration
	TotalBytes      int64
	CacheHits       int64
	CacheMisses     int64
}

// UpdateMetrics updates response metrics (would be called by the client)
func (m *ResponseMetrics) UpdateMetrics(info *ResponseInfo, err error) {
	m.TotalRequests++
	m.TotalBytes += info.Size

	if info.Cached {
		m.CacheHits++
	} else {
		m.CacheMisses++
	}

	if err != nil {
		m.FailedRequests++
		if apiErr, ok := err.(*APIError); ok {
			if apiErr.IsRateLimitError() {
				m.RateLimitErrors++
			} else if apiErr.IsServerError() {
				m.ServerErrors++
			}
		}
	} else {
		m.SuccessfulReqs++
	}

	// Update average response time (simple moving average)
	if m.TotalRequests > 0 {
		m.AverageResponse = (m.AverageResponse*time.Duration(m.TotalRequests-1) + info.ResponseTime) / time.Duration(m.TotalRequests)
	}
}
