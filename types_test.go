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
	"testing"
	"time"
)

func TestModelType(t *testing.T) {
	tests := []struct {
		modelType ModelType
		expected  string
	}{
		{ModelTypeCheckpoint, "Checkpoint"},
		{ModelTypeLORA, "LORA"},
		{ModelTypeEmbedding, "TextualInversion"},
		{ModelTypeHypernetwork, "Hypernetwork"},
	}

	for _, test := range tests {
		if string(test.modelType) != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, string(test.modelType))
		}
	}
}

func TestSortType(t *testing.T) {
	tests := []struct {
		sortType SortType
		expected string
	}{
		{SortHighestRated, "Highest Rated"},
		{SortMostDownload, "Most Downloaded"},
		{SortNewest, "Newest"},
		{SortOldest, "Oldest"},
	}

	for _, test := range tests {
		if string(test.sortType) != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, string(test.sortType))
		}
	}
}

func TestSearchParams(t *testing.T) {
	params := SearchParams{
		Query:              "anime",
		Types:              []ModelType{ModelTypeCheckpoint, ModelTypeLORA},
		Sort:               SortMostDownload,
		Period:             PeriodWeek,
		Rating:             4,
		Page:               1,
		Limit:              50,
		AllowCommercialUse: []string{"Sell", "RentCivit"},
	}

	if params.Query != "anime" {
		t.Errorf("Expected query 'anime', got '%s'", params.Query)
	}

	if len(params.Types) != 2 {
		t.Errorf("Expected 2 types, got %d", len(params.Types))
	}

	if params.Types[0] != ModelTypeCheckpoint {
		t.Errorf("Expected first type to be Checkpoint, got %s", params.Types[0])
	}

	if len(params.AllowCommercialUse) != 2 {
		t.Errorf("Expected 2 commercial use types, got %d", len(params.AllowCommercialUse))
	}
}

func TestModel(t *testing.T) {
	model := Model{
		ID:          12345,
		Name:        "Test Model",
		Description: "A test model",
		Type:        ModelTypeCheckpoint,
		NSFW:        false,
		Tags:        []string{"anime", "character"},
		Stats: Stats{
			DownloadCount: 1000,
			FavoriteCount: 50,
			CommentCount:  25,
			Rating:        4.5,
			RatingCount:   100,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if model.ID != 12345 {
		t.Errorf("Expected ID 12345, got %d", model.ID)
	}

	if model.Name != "Test Model" {
		t.Errorf("Expected name 'Test Model', got '%s'", model.Name)
	}

	if model.Type != ModelTypeCheckpoint {
		t.Errorf("Expected type Checkpoint, got %s", model.Type)
	}

	if len(model.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(model.Tags))
	}

	if model.Stats.DownloadCount != 1000 {
		t.Errorf("Expected 1000 downloads, got %d", model.Stats.DownloadCount)
	}
}

func TestMetadata(t *testing.T) {
	metadata := Metadata{
		CurrentPage: 1,
		PageSize:    20,
		TotalPages:  5,
		TotalItems:  100,
		NextPage:    "https://api.civitai.com/v1/models?page=2",
	}

	if metadata.CurrentPage != 1 {
		t.Errorf("Expected current page 1, got %d", metadata.CurrentPage)
	}

	if metadata.TotalItems != 100 {
		t.Errorf("Expected 100 total items, got %d", metadata.TotalItems)
	}
}

func TestImageParams(t *testing.T) {
	params := ImageParams{
		ModelID:        12345,
		ModelVersionID: 67890,
		Username:       "artist",
		NSFW:           "Soft",
		Sort:           string(ImageSortNewest),
		Period:         PeriodWeek,
		Page:           1,
		Limit:          100,
	}

	if params.ModelID != 12345 {
		t.Errorf("Expected model ID 12345, got %d", params.ModelID)
	}

	if params.NSFW != "Soft" {
		t.Errorf("Expected NSFW to be 'Soft', got '%s'", params.NSFW)
	}

	if params.Sort != string(ImageSortNewest) {
		t.Errorf("Expected sort 'newest', got '%s'", params.Sort)
	}
}
