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

// Package civitai - Tag Discovery and Categorization
//
// This file provides functionality for discovering and browsing tags
// used to categorize AI models on the CivitAI platform.
//
// # Basic Tag Browsing
//
// Discover available tags for model categorization:
//
//	client := civitai.NewClientWithoutAuth()
//	tags, metadata, err := client.GetTags(context.Background(), civitai.TagParams{
//		Limit: 50,
//	})
//
// # Searching Tags
//
// Search for specific tags by name:
//
//	params := civitai.TagParams{
//		Query: "anime",
//		Limit: 20,
//	}
//	tags, _, err := client.GetTags(ctx, params)
//
// # Tag Information
//
// Each tag contains usage statistics and metadata:
//
//	for _, tag := range tags {
//		fmt.Printf("Tag: %s\n", tag.Name)
//		fmt.Printf("Models using this tag: %d\n", tag.ModelCount)
//		fmt.Printf("Link: %s\n", tag.Link)
//
//		// Check tag color for UI display
//		if tag.Color != "" {
//			fmt.Printf("Display color: %s\n", tag.Color)
//		}
//	}
//
// # Using Tags for Model Search
//
// Tags are essential for effective model discovery:
//
//	// Get popular tags first
//	tags, _, err := client.GetTags(ctx, civitai.TagParams{Limit: 100})
//
//	// Use the most popular tags for model search
//	for _, tag := range tags {
//		if tag.ModelCount > 1000 { // Popular tags
//			models, _, err := client.SearchModels(ctx, civitai.SearchParams{
//				Tag: tag.Name,
//				Limit: 10,
//			})
//			// Process models...
//		}
//	}
//
// # Error Handling
//
// The Tags endpoint can experience timeout issues:
//
//	tags, _, err := client.GetTags(ctx, params)
//	if err != nil {
//		if strings.Contains(err.Error(), "timeout") {
//			// Retry without query parameter
//			params.Query = ""
//			tags, _, err = client.GetTags(ctx, params)
//		}
//	}
//
// # Performance Notes
//
// The Tags endpoint has mixed reliability:
//   - Generally stable for basic requests
//   - Occasional timeout issues with search queries
//   - Better performance with smaller page sizes
//   - Consider caching tag lists for better UX

package civitai

import (
	"context"
	"fmt"
	"strconv"
)

// GetTags retrieves a list of tags from the CivitAI API
// GET /api/v1/tags
func (c *Client) GetTags(ctx context.Context, params TagParams) ([]TagResponse, *Metadata, error) {
	if err := c.validateTagParams(params); err != nil {
		return nil, nil, fmt.Errorf("invalid tag parameters: %w", err)
	}

	queryParams := c.buildTagParams(params)
	url := c.addQueryParams(c.buildURL("tags"), queryParams)

	resp, err := c.doRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var apiResp struct {
		Items    []TagResponse `json:"items"`
		Metadata *Metadata     `json:"metadata"`
	}

	if err := c.handleResponse(resp, &apiResp); err != nil {
		return nil, nil, err
	}

	return apiResp.Items, apiResp.Metadata, nil
}

// buildTagParams converts TagParams to query parameters
func (c *Client) buildTagParams(params TagParams) map[string]string {
	queryParams := make(map[string]string)

	if params.Limit > 0 {
		queryParams["limit"] = strconv.Itoa(params.Limit)
	}
	if params.Page > 0 {
		queryParams["page"] = strconv.Itoa(params.Page)
	}
	if params.Query != "" {
		queryParams["query"] = params.Query
	}

	return queryParams
}
