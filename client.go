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
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
)

// Client represents a CivitAI API client
type Client struct {
	baseURL    string
	apiToken   string
	httpClient *http.Client
	userAgent  string
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

// NewClient creates a new CivitAI API client
func NewClient(apiToken string, options ...ClientOption) *Client {
	client := &Client{
		baseURL:  DefaultBaseURL,
		apiToken: apiToken,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		userAgent: DefaultUserAgent,
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

// doRequest executes an HTTP request and returns the response
func (c *Client) doRequest(ctx context.Context, method, url string, body []byte) (*http.Response, error) {
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

	// Add authentication if token is provided
	if c.apiToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return resp, nil
}

// handleResponse processes the HTTP response and unmarshals JSON
func (c *Client) handleResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiErr APIError
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
			return fmt.Errorf("API request failed with status %d", resp.StatusCode)
		}
		return fmt.Errorf("API error [%s]: %s", apiErr.Code, apiErr.Message)
	}

	if target != nil {
		if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// SearchModels searches for models with the given parameters
func (c *Client) SearchModels(ctx context.Context, params SearchParams) ([]Model, *Metadata, error) {
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
