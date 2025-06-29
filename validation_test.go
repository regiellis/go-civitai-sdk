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
	"strings"
	"testing"
)

func TestValidateModelID(t *testing.T) {
	tests := []struct {
		name    string
		modelID int
		wantErr bool
	}{
		{"valid positive ID", 123, false},
		{"zero ID", 0, true},
		{"negative ID", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateModelID(tt.modelID)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateModelID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateVersionID(t *testing.T) {
	tests := []struct {
		name      string
		versionID int
		wantErr   bool
	}{
		{"valid positive ID", 456, false},
		{"zero ID", 0, true},
		{"negative ID", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateVersionID(tt.versionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateVersionID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateHash(t *testing.T) {
	tests := []struct {
		name    string
		hash    string
		wantErr bool
	}{
		{"valid SHA256 hash", "5493A0EC49E72336B89F7E0A0BF9B2B2E03F3E2E9E7A6F8B5F3C3E9A3C9E2F9", false},
		{"valid MD5 hash", "5d41402abc4b2a76b9719d911017c592", false},
		{"empty hash", "", true},
		{"invalid characters", "invalid_hash_123", true},
		{"too short", "abc123", true},
		{"too long", strings.Repeat("a", 129), true},
		{"valid short hash", "abcdef12", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateHash(tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateSearchParams(t *testing.T) {
	tests := []struct {
		name    string
		params  SearchParams
		wantErr bool
	}{
		{"valid params", SearchParams{Query: "test", Limit: 10, Page: 1, Rating: 4}, false},
		{"negative page", SearchParams{Page: -1}, true},
		{"negative limit", SearchParams{Limit: -1}, true},
		{"limit too high", SearchParams{Limit: 300}, true},
		{"invalid rating", SearchParams{Rating: 6}, true},
		{"query too long", SearchParams{Query: strings.Repeat("a", 501)}, true},
		{"tag too long", SearchParams{Tag: strings.Repeat("b", 101)}, true},
		{"username too long", SearchParams{Username: strings.Repeat("c", 101)}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSearchParams(tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateSearchParams() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidationInAPIMethods(t *testing.T) {
	client := NewClientWithoutAuth()
	ctx := context.Background()

	t.Run("GetModel with invalid ID", func(t *testing.T) {
		_, err := client.GetModel(ctx, -1)
		if err == nil {
			t.Error("expected error for invalid model ID")
		}
		if !strings.Contains(err.Error(), "invalid model ID") {
			t.Errorf("expected 'invalid model ID' in error, got %v", err)
		}
	})

	t.Run("GetModelVersion with invalid ID", func(t *testing.T) {
		_, err := client.GetModelVersion(ctx, 0)
		if err == nil {
			t.Error("expected error for invalid version ID")
		}
		if !strings.Contains(err.Error(), "invalid version ID") {
			t.Errorf("expected 'invalid version ID' in error, got %v", err)
		}
	})

	t.Run("GetModelVersionByHash with invalid hash", func(t *testing.T) {
		_, err := client.GetModelVersionByHash(ctx, "")
		if err == nil {
			t.Error("expected error for invalid hash")
		}
		if !strings.Contains(err.Error(), "invalid hash") {
			t.Errorf("expected 'invalid hash' in error, got %v", err)
		}
	})

	t.Run("SearchModels with invalid params", func(t *testing.T) {
		_, _, err := client.SearchModels(ctx, SearchParams{Limit: -1})
		if err == nil {
			t.Error("expected error for invalid search params")
		}
		if !strings.Contains(err.Error(), "invalid search parameters") {
			t.Errorf("expected 'invalid search parameters' in error, got %v", err)
		}
	})
}
