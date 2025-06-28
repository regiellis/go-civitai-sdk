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

// GetImages retrieves a list of images from the CivitAI API
// GET /api/v1/images
func (c *Client) GetImages(ctx context.Context, params ImageParams) ([]DetailedImageResponse, *Metadata, error) {
	queryParams := c.buildImageParams(params)
	url := c.addQueryParams(c.buildURL("images"), queryParams)

	resp, err := c.doRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var apiResp struct {
		Items    []DetailedImageResponse `json:"items"`
		Metadata *Metadata              `json:"metadata"`
	}

	if err := c.handleResponse(resp, &apiResp); err != nil {
		return nil, nil, err
	}

	return apiResp.Items, apiResp.Metadata, nil
}

// buildImageParams converts ImageParams to query parameters
func (c *Client) buildImageParams(params ImageParams) map[string]string {
	queryParams := make(map[string]string)

	if params.Limit > 0 {
		queryParams["limit"] = strconv.Itoa(params.Limit)
	}
	if params.PostID > 0 {
		queryParams["postId"] = strconv.Itoa(params.PostID)
	}
	if params.ModelID > 0 {
		queryParams["modelId"] = strconv.Itoa(params.ModelID)
	}
	if params.ModelVersionID > 0 {
		queryParams["modelVersionId"] = strconv.Itoa(params.ModelVersionID)
	}
	if params.Username != "" {
		queryParams["username"] = params.Username
	}
	if params.NSFW != "" {
		queryParams["nsfw"] = params.NSFW
	}
	if params.Sort != "" {
		queryParams["sort"] = params.Sort
	}
	if params.Period != "" {
		queryParams["period"] = string(params.Period)
	}
	if params.Page > 0 {
		queryParams["page"] = strconv.Itoa(params.Page)
	}

	return queryParams
}