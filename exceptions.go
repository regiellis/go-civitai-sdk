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
	"fmt"
	"net/http"
)

// Error implements the error interface for APIError
func (e APIError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("CivitAI API error [%s]: %s - %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("CivitAI API error [%s]: %s", e.Code, e.Message)
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