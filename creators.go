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

// Package civitai - Creator Discovery and Statistics
//
// This file provides functionality for discovering AI model creators
// and accessing their statistics and model portfolios.
//
// # Basic Creator Search
//
// Browse and search for model creators:
//
//	client := civitai.NewClientWithoutAuth()
//	creators, metadata, err := client.GetCreators(context.Background(), civitai.CreatorParams{
//		Limit: 20,
//	})
//
// # Searching Specific Creators
//
// Find creators by username or query:
//
//	params := civitai.CreatorParams{
//		Query: "artist-name",
//		Limit: 10,
//	}
//	creators, _, err := client.GetCreators(ctx, params)
//
// # Creator Information
//
// Each creator object contains comprehensive statistics:
//
//	for _, creator := range creators {
//		fmt.Printf("Creator: %s\n", creator.Username)
//		fmt.Printf("Models: %d\n", creator.ModelCount)
//		fmt.Printf("Followers: %d\n", creator.FollowerCount)
//		fmt.Printf("Upload count: %d\n", creator.UploadCount)
//
//		// Check if creator has profile link
//		if creator.Link != "" {
//			fmt.Printf("Profile: %s\n", creator.Link)
//		}
//	}
//
// # Error Handling and Reliability
//
// Important: The Creators endpoint has known reliability issues:
//
//	creators, _, err := client.GetCreators(ctx, params)
//	if err != nil {
//		// Timeout errors are common (~20% failure rate)
//		if strings.Contains(err.Error(), "timeout") {
//			log.Println("Creators endpoint timeout - retrying...")
//			// Implement retry logic
//		}
//	}
//
// # Best Practices
//
// 1. Implement retry logic with exponential backoff
// 2. Use larger timeouts (60+ seconds) for this endpoint
// 3. Consider fallback strategies for timeout scenarios
// 4. Monitor success rates in production
//
// # Performance Characteristics
//
// Based on extensive testing:
//   - Average response time: 2+ seconds (when successful)
//   - Success rate: ~80% (20% timeout rate)
//   - Timeout issues increase under load
//   - Most reliable with smaller page sizes (limit ≤ 10)

package civitai

import (
	"context"
	"fmt"
	"strconv"
)

// GetCreators retrieves a list of creators from the CivitAI API
// GET /api/v1/creators
func (c *Client) GetCreators(ctx context.Context, params CreatorParams) ([]Creator, *Metadata, error) {
	if err := c.validateCreatorParams(params); err != nil {
		return nil, nil, fmt.Errorf("invalid creator parameters: %w", err)
	}

	queryParams := c.buildCreatorParams(params)
	url := c.addQueryParams(c.buildURL("creators"), queryParams)

	resp, err := c.doRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var apiResp struct {
		Items    []Creator `json:"items"`
		Metadata *Metadata `json:"metadata"`
	}

	if err := c.handleResponse(resp, &apiResp); err != nil {
		return nil, nil, err
	}

	return apiResp.Items, apiResp.Metadata, nil
}

// buildCreatorParams converts CreatorParams to query parameters
func (c *Client) buildCreatorParams(params CreatorParams) map[string]string {
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
