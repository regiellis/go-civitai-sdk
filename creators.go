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
	"strconv"
)

// GetCreators retrieves a list of creators from the CivitAI API
// GET /api/v1/creators
func (c *Client) GetCreators(ctx context.Context, params CreatorParams) ([]Creator, *Metadata, error) {
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
