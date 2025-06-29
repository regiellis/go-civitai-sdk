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

// Package civitai provides a comprehensive Go SDK for the CivitAI API.
//
// The CivitAI API allows developers to search for AI models, browse images,
// discover creators, and access detailed model information programmatically.
// This SDK handles authentication, rate limiting, retries, and provides
// type-safe interfaces for all API endpoints.
//
// # Quick Start
//
// Create a client and start searching for models:
//
//	client := civitai.NewClientWithoutAuth()
//	models, metadata, err := client.SearchModels(context.Background(), civitai.SearchParams{
//		Tag:   "anime",
//		Limit: 10,
//	})
//
// # Authentication
//
// For endpoints requiring authentication, create a client with your API token:
//
//	client := civitai.NewClient("your-api-token")
//
// # Best Practices
//
// 1. Use tag-based search instead of query search for better results:
//
//	// Recommended - returns more results
//	params := civitai.SearchParams{Tag: "realistic"}
//
//	// Less reliable
//	params := civitai.SearchParams{Query: "realistic"}
//
// 2. Use cursor-based pagination for consistent results:
//
//	var allModels []civitai.Model
//	cursor := ""
//	for {
//		models, meta, err := client.SearchModels(ctx, civitai.SearchParams{
//			Tag:    "anime",
//			Cursor: cursor,
//			Limit:  50,
//		})
//		if err != nil || len(models) == 0 {
//			break
//		}
//		allModels = append(allModels, models...)
//		if meta.NextCursor == "" {
//			break
//		}
//		cursor = meta.NextCursor
//	}
//
// 3. Configure timeouts and retries for production use:
//
//	client := civitai.NewClientWithoutAuth(
//		civitai.WithTimeout(60*time.Second),
//		civitai.WithRetryConfig(3, 2*time.Second, 30*time.Second),
//	)
//
// # API Endpoints
//
// Models:
//   - SearchModels: Search for AI models by various criteria
//   - GetModel: Retrieve detailed information about a specific model
//   - GetModelVersion: Get details about a specific model version
//   - GetModelVersionsByModelID: List all versions of a model
//   - GetModelVersionByHash: Find a model version by file hash
//
// Images:
//   - GetImages: Browse AI-generated images with filtering options
//
// Creators:
//   - GetCreators: Discover model creators and their statistics
//
// Tags:
//   - GetTags: Explore available tags for categorizing models
//
// # Error Handling
//
// The SDK provides comprehensive error handling with typed errors:
//
//	models, _, err := client.SearchModels(ctx, params)
//	if err != nil {
//		var apiErr *civitai.APIError
//		if errors.As(err, &apiErr) {
//			log.Printf("API error: %s (code: %d)", apiErr.Message, apiErr.StatusCode)
//		} else {
//			log.Printf("Request error: %v", err)
//		}
//	}
//
// # Known API Limitations
//
// Based on extensive testing, be aware of these API behaviors:
//   - Tag-based search returns 2-5x more results than query-based search
//   - Creators endpoint has ~20% timeout rate under load
//   - Version-by-hash endpoint is currently non-functional
//   - Page-based pagination is unreliable; use cursor-based pagination
//
// For the latest API behavior analysis, see the comprehensive test suite
// and documentation in the examples directory.
package civitai

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	// DefaultBaseURL is the default CivitAI API base URL
	DefaultBaseURL = "https://civitai.com/api/v1"

	// DefaultTimeout is the default HTTP client timeout
	DefaultTimeout = 30 * time.Second

	// DefaultUserAgent is the default user agent string
	DefaultUserAgent = "go-civitai-sdk/1.0.0"

	// DefaultMaxResponseSize is the default maximum response size (10MB)
	DefaultMaxResponseSize = 10 * 1024 * 1024 // 10MB

	// DefaultMaxRetries is the default number of retry attempts
	DefaultMaxRetries = 3

	// DefaultRetryDelay is the base delay for exponential backoff
	DefaultRetryDelay = 1 * time.Second

	// DefaultMaxRetryDelay is the maximum delay between retries
	DefaultMaxRetryDelay = 30 * time.Second
)

// Client represents a CivitAI API client
type Client struct {
	baseURL         string
	apiToken        string
	httpClient      *http.Client
	userAgent       string
	maxResponseSize int64
	maxRetries      int
	retryDelay      time.Duration
	maxRetryDelay   time.Duration
}

// ClientOption represents a function that configures the client
type ClientOption func(*Client)

// WithBaseURL sets a custom base URL for the API
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = strings.TrimSuffix(baseURL, "/")
	}
}

// WithTimeout sets a custom timeout for HTTP requests
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// WithUserAgent sets a custom user agent string
func WithUserAgent(userAgent string) ClientOption {
	return func(c *Client) {
		c.userAgent = userAgent
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithMaxResponseSize sets the maximum allowed response size in bytes
func WithMaxResponseSize(size int64) ClientOption {
	return func(c *Client) {
		c.maxResponseSize = size
	}
}

// WithRetryConfig sets the retry configuration for failed requests
func WithRetryConfig(maxRetries int, baseDelay, maxDelay time.Duration) ClientOption {
	return func(c *Client) {
		c.maxRetries = maxRetries
		c.retryDelay = baseDelay
		c.maxRetryDelay = maxDelay
	}
}

// WithConnectionPooling configures the HTTP client for connection pooling and compression
func WithConnectionPooling(maxIdleConns, maxIdleConnsPerHost int) ClientOption {
	return func(c *Client) {
		transport := &http.Transport{
			MaxIdleConns:        maxIdleConns,
			MaxIdleConnsPerHost: maxIdleConnsPerHost,
			IdleConnTimeout:     90 * time.Second,
			DisableCompression:  false, // Enable compression
		}
		c.httpClient.Transport = transport
	}
}

// NewClient creates a new CivitAI API client
func NewClient(apiToken string, options ...ClientOption) *Client {
	client := &Client{
		baseURL:  DefaultBaseURL,
		apiToken: apiToken,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		userAgent:       DefaultUserAgent,
		maxResponseSize: DefaultMaxResponseSize,
		maxRetries:      DefaultMaxRetries,
		retryDelay:      DefaultRetryDelay,
		maxRetryDelay:   DefaultMaxRetryDelay,
	}

	// Apply options
	for _, option := range options {
		option(client)
	}

	return client
}

// NewClientWithoutAuth creates a new CivitAI API client without authentication
// This can be used for public endpoints that don't require an API token
func NewClientWithoutAuth(options ...ClientOption) *Client {
	return NewClient("", options...)
}

// buildURL constructs a full URL from the base URL and path
func (c *Client) buildURL(path string) string {
	return fmt.Sprintf("%s/%s", c.baseURL, strings.TrimPrefix(path, "/"))
}

// addQueryParams adds query parameters to a URL
func (c *Client) addQueryParams(baseURL string, params map[string]string) string {
	if len(params) == 0 {
		return baseURL
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		return baseURL
	}

	q := u.Query()
	for key, value := range params {
		if value != "" {
			q.Set(key, value)
		}
	}

	u.RawQuery = q.Encode()
	return u.String()
}

// Input validation functions

// validateModelID validates that a model ID is positive
func validateModelID(modelID int) error {
	if modelID <= 0 {
		return errors.New("model ID must be a positive integer")
	}
	return nil
}

// validateVersionID validates that a version ID is positive
func validateVersionID(versionID int) error {
	if versionID <= 0 {
		return errors.New("version ID must be a positive integer")
	}
	return nil
}

// validateHash validates that a hash string is not empty and contains only valid characters
func validateHash(hash string) error {
	if hash == "" {
		return errors.New("hash cannot be empty")
	}

	// Hash should only contain alphanumeric characters
	hashRegex := regexp.MustCompile(`^[a-fA-F0-9]+$`)
	if !hashRegex.MatchString(hash) {
		return errors.New("hash must contain only hexadecimal characters")
	}

	// Common hash lengths: MD5(32), SHA1(40), SHA256(64), etc.
	// Allow reasonable range
	if len(hash) < 8 || len(hash) > 128 {
		return errors.New("hash length must be between 8 and 128 characters")
	}

	return nil
}

// validateSearchParams validates search parameters for safety
func validateSearchParams(params SearchParams) error {
	// Validate page and limit bounds
	if params.Page < 0 {
		return errors.New("page cannot be negative")
	}
	if params.Limit < 0 {
		return errors.New("limit cannot be negative")
	}
	if params.Limit > 200 {
		return errors.New("limit cannot exceed 200")
	}
	if params.Rating < 0 || params.Rating > 5 {
		return errors.New("rating must be between 0 and 5")
	}

	// Validate string parameters for length to prevent abuse
	if len(params.Query) > 500 {
		return errors.New("query parameter too long (max 500 characters)")
	}
	if len(params.Tag) > 100 {
		return errors.New("tag parameter too long (max 100 characters)")
	}
	if len(params.Username) > 100 {
		return errors.New("username parameter too long (max 100 characters)")
	}

	return nil
}

// validateImageParams validates image search parameters
func (c *Client) validateImageParams(params ImageParams) error {
	if params.Limit < 0 || params.Limit > 200 {
		return errors.New("limit must be between 0 and 200")
	}
	if params.Page < 0 {
		return errors.New("page cannot be negative")
	}
	if params.PostID < 0 {
		return errors.New("post ID cannot be negative")
	}
	if params.ModelID < 0 {
		return errors.New("model ID cannot be negative")
	}
	if params.ModelVersionID < 0 {
		return errors.New("model version ID cannot be negative")
	}
	if len(params.Username) > 100 {
		return errors.New("username parameter too long (max 100 characters)")
	}
	return nil
}

// validateCreatorParams validates creator search parameters
func (c *Client) validateCreatorParams(params CreatorParams) error {
	if params.Limit < 0 || params.Limit > 200 {
		return errors.New("limit must be between 0 and 200")
	}
	if params.Page < 0 {
		return errors.New("page cannot be negative")
	}
	if len(params.Query) > 500 {
		return errors.New("query parameter too long (max 500 characters)")
	}
	return nil
}

// validateTagParams validates tag search parameters
func (c *Client) validateTagParams(params TagParams) error {
	if params.Limit < 0 || params.Limit > 200 {
		return errors.New("limit must be between 0 and 200")
	}
	if params.Page < 0 {
		return errors.New("page cannot be negative")
	}
	if len(params.Query) > 500 {
		return errors.New("query parameter too long (max 500 characters)")
	}
	return nil
}

// isRetryableError determines if an error is worth retrying
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// Check for specific context errors
	if err == context.DeadlineExceeded {
		return true
	}
	if err == context.Canceled {
		return false // Don't retry cancelled contexts
	}

	// Retry on network errors, timeouts, and temporary failures
	errStr := err.Error()
	return strings.Contains(errStr, "timeout") ||
		strings.Contains(errStr, "connection refused") ||
		strings.Contains(errStr, "no such host") ||
		strings.Contains(errStr, "temporary failure") ||
		strings.Contains(errStr, "network is unreachable") ||
		strings.Contains(errStr, "connection reset")
}

// isRetryableStatusCode determines if an HTTP status code is worth retrying
func isRetryableStatusCode(statusCode int) bool {
	// Retry on server errors and rate limiting
	return statusCode == http.StatusTooManyRequests ||
		statusCode == http.StatusInternalServerError ||
		statusCode == http.StatusBadGateway ||
		statusCode == http.StatusServiceUnavailable ||
		statusCode == http.StatusGatewayTimeout
}

// calculateBackoffDelay calculates the delay for exponential backoff with jitter
func (c *Client) calculateBackoffDelay(attempt int) time.Duration {
	// Exponential backoff: baseDelay * 2^attempt
	delay := time.Duration(float64(c.retryDelay) * math.Pow(2, float64(attempt)))

	// Add jitter (±25% random variation)
	jitter := time.Duration(float64(delay) * 0.25 * (2*rand.Float64() - 1))
	delay += jitter

	// Cap at maximum delay
	if delay > c.maxRetryDelay {
		delay = c.maxRetryDelay
	}

	return delay
}

// doRequest executes an HTTP request with retry logic and returns the response
func (c *Client) doRequest(ctx context.Context, method, url string, body []byte) (*http.Response, error) {
	var lastErr error

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		// Create request for this attempt
		var req *http.Request
		var err error

		if body != nil {
			req, err = http.NewRequestWithContext(ctx, method, url, strings.NewReader(string(body)))
		} else {
			req, err = http.NewRequestWithContext(ctx, method, url, nil)
		}

		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Set headers
		req.Header.Set("User-Agent", c.userAgent)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept-Encoding", "gzip, deflate") // Request compression

		// Add authentication if token is provided
		if c.apiToken != "" {
			req.Header.Set("Authorization", "Bearer "+c.apiToken)
		}

		resp, err := c.httpClient.Do(req)

		// If successful or non-retryable error, return immediately
		if err == nil {
			if !isRetryableStatusCode(resp.StatusCode) {
				return resp, nil
			}
			// Close response body for retryable status codes
			resp.Body.Close()
			lastErr = fmt.Errorf("HTTP %d: %s", resp.StatusCode, http.StatusText(resp.StatusCode))
		} else {
			lastErr = err
			if !isRetryableError(err) {
				return nil, fmt.Errorf("failed to execute request: %w", err)
			}
		}

		// Don't wait after the last attempt
		if attempt < c.maxRetries {
			delay := c.calculateBackoffDelay(attempt)

			// Create timer with context cancellation support
			timer := time.NewTimer(delay)
			select {
			case <-ctx.Done():
				timer.Stop()
				return nil, ctx.Err()
			case <-timer.C:
				// Continue to next attempt
			}
		}
	}

	return nil, fmt.Errorf("failed to execute request after %d attempts: %w", c.maxRetries+1, lastErr)
}

// handleResponse processes the HTTP response and unmarshals JSON
func (c *Client) handleResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()

	// Handle gzip compression
	var reader io.Reader = resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer gzipReader.Close()
		reader = gzipReader
	}

	// Apply response size limit to prevent DoS attacks
	limitedReader := io.LimitReader(reader, c.maxResponseSize)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiErr APIError
		if err := json.NewDecoder(limitedReader).Decode(&apiErr); err != nil {
			return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, resp.Status)
		}
		return fmt.Errorf("API error [%s]: %s", apiErr.Code, apiErr.Message)
	}

	if target != nil {
		decoder := json.NewDecoder(limitedReader)
		if err := decoder.Decode(target); err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				return fmt.Errorf("response size exceeded maximum allowed size of %d bytes", c.maxResponseSize)
			}
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// SearchModels searches for models with the given parameters
func (c *Client) SearchModels(ctx context.Context, params SearchParams) ([]Model, *Metadata, error) {
	if err := validateSearchParams(params); err != nil {
		return nil, nil, fmt.Errorf("invalid search parameters: %w", err)
	}

	queryParams := c.buildSearchParams(params)
	url := c.addQueryParams(c.buildURL("models"), queryParams)

	resp, err := c.doRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var apiResp struct {
		Items    []Model   `json:"items"`
		Metadata *Metadata `json:"metadata"`
	}

	if err := c.handleResponse(resp, &apiResp); err != nil {
		return nil, nil, err
	}

	return apiResp.Items, apiResp.Metadata, nil
}

// GetModel retrieves a specific model by ID
func (c *Client) GetModel(ctx context.Context, modelID int) (*Model, error) {
	if err := validateModelID(modelID); err != nil {
		return nil, fmt.Errorf("invalid model ID: %w", err)
	}

	url := c.buildURL(fmt.Sprintf("models/%d", modelID))

	resp, err := c.doRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	var model Model
	if err := c.handleResponse(resp, &model); err != nil {
		return nil, err
	}

	return &model, nil
}

// GetModelVersion retrieves a specific model version by ID
func (c *Client) GetModelVersion(ctx context.Context, versionID int) (*ModelVersion, error) {
	if err := validateVersionID(versionID); err != nil {
		return nil, fmt.Errorf("invalid version ID: %w", err)
	}

	url := c.buildURL(fmt.Sprintf("model-versions/%d", versionID))

	resp, err := c.doRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	var version ModelVersion
	if err := c.handleResponse(resp, &version); err != nil {
		return nil, err
	}

	return &version, nil
}

// GetModelVersionsByModelID retrieves all versions for a specific model
func (c *Client) GetModelVersionsByModelID(ctx context.Context, modelID int) ([]ModelVersion, error) {
	if err := validateModelID(modelID); err != nil {
		return nil, fmt.Errorf("invalid model ID: %w", err)
	}

	url := c.buildURL(fmt.Sprintf("models/%d/versions", modelID))

	resp, err := c.doRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	var versions []ModelVersion
	if err := c.handleResponse(resp, &versions); err != nil {
		return nil, err
	}

	return versions, nil
}

// GetModelVersionByHash retrieves a model version by file hash
// GET /api/v1/model-versions/by-hash/:hash
// Supports AutoV1, AutoV2, SHA256, CRC32, and Blake3 hash algorithms
func (c *Client) GetModelVersionByHash(ctx context.Context, hash string) (*ModelVersionByHashResponse, error) {
	if err := validateHash(hash); err != nil {
		return nil, fmt.Errorf("invalid hash: %w", err)
	}

	url := c.buildURL(fmt.Sprintf("model-versions/by-hash/%s", hash))

	resp, err := c.doRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	var version ModelVersionByHashResponse
	if err := c.handleResponse(resp, &version); err != nil {
		return nil, err
	}

	return &version, nil
}

// buildSearchParams converts SearchParams to query parameters
func (c *Client) buildSearchParams(params SearchParams) map[string]string {
	queryParams := make(map[string]string)

	if params.Query != "" {
		queryParams["query"] = params.Query
	}
	if len(params.Types) > 0 {
		var types []string
		for _, t := range params.Types {
			types = append(types, string(t))
		}
		queryParams["types"] = strings.Join(types, ",")
	}
	if params.Sort != "" {
		queryParams["sort"] = string(params.Sort)
	}
	if params.Period != "" {
		queryParams["period"] = string(params.Period)
	}
	if params.Rating > 0 {
		queryParams["rating"] = strconv.Itoa(params.Rating)
	}
	if params.Page > 0 {
		queryParams["page"] = strconv.Itoa(params.Page)
	}
	if params.Limit > 0 {
		queryParams["limit"] = strconv.Itoa(params.Limit)
	}
	if params.Cursor != "" {
		queryParams["cursor"] = params.Cursor
	}
	if params.Tag != "" {
		queryParams["tag"] = params.Tag
	}
	if params.Username != "" {
		queryParams["username"] = params.Username
	}
	if params.Favorites {
		queryParams["favorites"] = "true"
	}
	if params.Hidden {
		queryParams["hidden"] = "true"
	}
	if params.PrimaryFileOnly {
		queryParams["primaryFileOnly"] = "true"
	}
	if params.AllowNoCredit {
		queryParams["allowNoCredit"] = "true"
	}
	if params.AllowDerivatives {
		queryParams["allowDerivatives"] = "true"
	}
	if params.AllowDifferentLicense {
		queryParams["allowDifferentLicense"] = "true"
	}
	if len(params.AllowCommercialUse) > 0 {
		queryParams["allowCommercialUse"] = strings.Join(params.AllowCommercialUse, ",")
	}
	if params.NSFW != nil {
		if *params.NSFW {
			queryParams["nsfw"] = "true"
		} else {
			queryParams["nsfw"] = "false"
		}
	}
	if params.SupportsGeneration != nil {
		if *params.SupportsGeneration {
			queryParams["supportsGeneration"] = "true"
		} else {
			queryParams["supportsGeneration"] = "false"
		}
	}

	return queryParams
}

// Health checks the API health status
func (c *Client) Health(ctx context.Context) error {
	// CivitAI doesn't have a dedicated health endpoint, so we'll use a simple model request
	url := c.buildURL("models")
	queryParams := map[string]string{"limit": "1"}
	url = c.addQueryParams(url, queryParams)

	resp, err := c.doRequest(ctx, "GET", url, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API health check failed with status %d", resp.StatusCode)
	}

	return nil
}

// GetAPIToken returns the API token used by this client
// WARNING: This method exposes sensitive credentials and should be used with caution.
// Consider using HasAPIToken() instead to check if a token is configured.
// Deprecated: This method will be removed in a future version for security reasons.
func (c *Client) GetAPIToken() string {
	return c.apiToken
}

// HasAPIToken returns true if an API token is configured for this client
func (c *Client) HasAPIToken() bool {
	return c.apiToken != ""
}

// GetMaskedAPIToken returns a masked version of the API token for logging/debugging purposes
// Returns the first 8 characters followed by asterisks, or "none" if no token is set
func (c *Client) GetMaskedAPIToken() string {
	if c.apiToken == "" {
		return "none"
	}

	if len(c.apiToken) <= 8 {
		return strings.Repeat("*", len(c.apiToken))
	}

	return c.apiToken[:8] + strings.Repeat("*", len(c.apiToken)-8)
}

// IsAuthenticated returns true if the client has an API token configured
func (c *Client) IsAuthenticated() bool {
	return c.HasAPIToken()
}

// GetModelByAIR retrieves a model using an AIR identifier
func (c *Client) GetModelByAIR(ctx context.Context, air *AIR) (*Model, error) {
	if air == nil {
		return nil, errors.New("AIR cannot be nil")
	}

	if !air.IsCivitAI() {
		return nil, fmt.Errorf("AIR source '%s' is not supported by CivitAI client", air.Source)
	}

	modelID, err := air.GetModelID()
	if err != nil {
		return nil, fmt.Errorf("failed to extract model ID from AIR: %w", err)
	}

	return c.GetModel(ctx, modelID)
}

// GetModelVersionByAIR retrieves a model version using an AIR identifier
func (c *Client) GetModelVersionByAIR(ctx context.Context, air *AIR) (*ModelVersion, error) {
	if air == nil {
		return nil, errors.New("AIR cannot be nil")
	}

	if !air.IsCivitAI() {
		return nil, fmt.Errorf("AIR source '%s' is not supported by CivitAI client", air.Source)
	}

	if !air.IsVersionSpecific() {
		return nil, errors.New("AIR must specify a version to retrieve model version")
	}

	versionID, err := air.GetVersionID()
	if err != nil {
		return nil, fmt.Errorf("failed to extract version ID from AIR: %w", err)
	}

	return c.GetModelVersion(ctx, versionID)
}

// SearchModelsByAIRType searches for models by AIR type
func (c *Client) SearchModelsByAIRType(ctx context.Context, airType AIRType, params SearchParams) ([]Model, *Metadata, error) {
	// Convert AIR type to CivitAI model type
	air := &AIR{Type: string(airType)}
	modelType := air.ToModelType()

	// Add type filter to search params
	if params.Types == nil {
		params.Types = []ModelType{modelType}
	} else {
		// Check if type is already in the list
		found := false
		for _, t := range params.Types {
			if t == modelType {
				found = true
				break
			}
		}
		if !found {
			params.Types = append(params.Types, modelType)
		}
	}

	return c.SearchModels(ctx, params)
}

// ConvertModelToAIR converts a CivitAI model to an AIR identifier
func ConvertModelToAIR(model *Model, ecosystem string, versionID ...int) *AIR {
	if model == nil {
		return nil
	}

	// Determine ecosystem if not provided
	if ecosystem == "" {
		// Try to infer from model tags or default to sdxl
		ecosystem = string(AIREcosystemSDXL)
		if model.Tags != nil {
			for _, tag := range model.Tags {
				switch strings.ToLower(tag) {
				case "sd 1.5", "stable diffusion 1.5":
					ecosystem = string(AIREcosystemSD1)
				case "sd 2.0", "sd 2.1", "stable diffusion 2":
					ecosystem = string(AIREcosystemSD2)
				case "flux", "flux.1":
					ecosystem = string(AIREcosystemFlux)
				}
			}
		}
	}

	// Determine AIR type from model type
	var airType string
	switch model.Type {
	case ModelTypeCheckpoint:
		airType = string(AIRTypeModel)
	case ModelTypeLORA:
		airType = string(AIRTypeLora)
	case ModelTypeTextualInversion:
		airType = string(AIRTypeEmbedding)
	case ModelTypeVAE:
		airType = string(AIRTypeVAE)
	case ModelTypeControlNet:
		airType = string(AIRTypeControl)
	default:
		airType = string(AIRTypeModel)
	}

	air := NewCivitAIModelAIR(ecosystem, model.ID)
	air.Type = airType

	// Add version if provided
	if len(versionID) > 0 && versionID[0] > 0 {
		air.Version = strconv.Itoa(versionID[0])
	}

	return air
}

// ConvertVersionToAIR converts a CivitAI model version to an AIR identifier
func ConvertVersionToAIR(version *ModelVersion, ecosystem string) *AIR {
	if version == nil {
		return nil
	}

	// Determine ecosystem if not provided
	if ecosystem == "" {
		ecosystem = string(AIREcosystemSDXL) // Default
	}

	air := NewCivitAIModelAIR(ecosystem, version.ModelID, version.ID)

	// Try to determine type from version files
	if len(version.Files) > 0 {
		primaryFile := version.Files[0]
		switch {
		case strings.Contains(strings.ToLower(primaryFile.Name), "lora"):
			air.Type = string(AIRTypeLora)
		case strings.Contains(strings.ToLower(primaryFile.Name), "vae"):
			air.Type = string(AIRTypeVAE)
		case strings.Contains(strings.ToLower(primaryFile.Name), "embedding"):
			air.Type = string(AIRTypeEmbedding)
		default:
			air.Type = string(AIRTypeModel)
		}

		// Add format if determinable
		if strings.HasSuffix(strings.ToLower(primaryFile.Name), ".safetensors") {
			air.Format = "safetensors"
		} else if strings.HasSuffix(strings.ToLower(primaryFile.Name), ".ckpt") {
			air.Format = "ckpt"
		}
	}

	return air
}

// Helper methods for common use cases (KISS principle)

// QuickSearch performs a simple text search for models
func (c *Client) QuickSearch(ctx context.Context, query string, limit int) ([]Model, error) {
	models, _, err := c.SearchModels(ctx, SearchParams{
		Query: query,
		Limit: limit,
	})
	return models, err
}

// GetPopularModels returns the most downloaded models
func (c *Client) GetPopularModels(ctx context.Context, limit int) ([]Model, error) {
	models, _, err := c.SearchModels(ctx, SearchParams{
		Sort:  SortMostDownload,
		Limit: limit,
	})
	return models, err
}

// GetNewestModels returns the latest uploaded models
func (c *Client) GetNewestModels(ctx context.Context, limit int) ([]Model, error) {
	models, _, err := c.SearchModels(ctx, SearchParams{
		Sort:  SortNewest,
		Limit: limit,
	})
	return models, err
}

// GetSafeImages returns safe-for-work images
func (c *Client) GetSafeImages(ctx context.Context, limit int) ([]DetailedImageResponse, error) {
	images, _, err := c.GetImages(ctx, ImageParams{
		NSFW:  string(NSFWLevelNone),
		Limit: limit,
	})
	return images, err
}

// IsWorking performs a simple health check to see if the API is accessible
func (c *Client) IsWorking(ctx context.Context) bool {
	return c.Health(ctx) == nil
}
