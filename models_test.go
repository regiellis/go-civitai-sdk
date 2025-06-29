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
	"testing"
	"time"
)

func TestFilterModels(t *testing.T) {
	models := []Model{
		{
			ID:    1,
			Name:  "Test Model 1",
			Type:  ModelTypeCheckpoint,
			NSFW:  false,
			Stats: Stats{Rating: 4.5},
			Tags:  []string{"anime", "character"},
		},
		{
			ID:    2,
			Name:  "Test Model 2",
			Type:  ModelTypeLORA,
			NSFW:  true,
			Stats: Stats{Rating: 3.8},
			Tags:  []string{"realistic", "portrait"},
		},
		{
			ID:    3,
			Name:  "Test Model 3",
			Type:  ModelTypeCheckpoint,
			NSFW:  false,
			Stats: Stats{Rating: 4.9},
			Tags:  []string{"anime", "style"},
		},
	}

	t.Run("Filter by type", func(t *testing.T) {
		filter := ModelFilter{Types: []ModelType{ModelTypeCheckpoint}}
		filtered := FilterModels(models, filter)

		if len(filtered) != 2 {
			t.Errorf("Expected 2 models, got %d", len(filtered))
		}

		for _, model := range filtered {
			if model.Type != ModelTypeCheckpoint {
				t.Errorf("Expected checkpoint model, got %s", model.Type)
			}
		}
	})

	t.Run("Filter by NSFW", func(t *testing.T) {
		nsfw := false
		filter := ModelFilter{NSFW: &nsfw}
		filtered := FilterModels(models, filter)

		if len(filtered) != 2 {
			t.Errorf("Expected 2 models, got %d", len(filtered))
		}

		for _, model := range filtered {
			if model.NSFW {
				t.Errorf("Expected non-NSFW model, got NSFW model")
			}
		}
	})

	t.Run("Filter by rating", func(t *testing.T) {
		filter := ModelFilter{MinRating: 4.0}
		filtered := FilterModels(models, filter)

		if len(filtered) != 2 {
			t.Errorf("Expected 2 models, got %d", len(filtered))
		}

		for _, model := range filtered {
			if model.Stats.Rating < 4.0 {
				t.Errorf("Expected rating >= 4.0, got %.1f", model.Stats.Rating)
			}
		}
	})

	t.Run("Filter by tag", func(t *testing.T) {
		filter := ModelFilter{Tags: []string{"anime"}}
		filtered := FilterModels(models, filter)

		if len(filtered) != 2 {
			t.Errorf("Expected 2 models, got %d", len(filtered))
		}

		for _, model := range filtered {
			if !model.HasTag("anime") {
				t.Errorf("Expected model with 'anime' tag")
			}
		}
	})

	t.Run("Empty models slice", func(t *testing.T) {
		filter := ModelFilter{Types: []ModelType{ModelTypeCheckpoint}}
		filtered := FilterModels([]Model{}, filter)

		if len(filtered) != 0 {
			t.Errorf("Expected 0 models, got %d", len(filtered))
		}
	})
}

func TestSortModels(t *testing.T) {
	now := time.Now()
	models := []Model{
		{
			ID:        1,
			Name:      "Model A",
			Stats:     Stats{Rating: 3.5, DownloadCount: 100, ThumbsUpCount: 50},
			CreatedAt: now.Add(-time.Hour),
		},
		{
			ID:        2,
			Name:      "Model B",
			Stats:     Stats{Rating: 4.5, DownloadCount: 200, ThumbsUpCount: 30},
			CreatedAt: now,
		},
		{
			ID:        3,
			Name:      "Model C",
			Stats:     Stats{Rating: 4.0, DownloadCount: 150, ThumbsUpCount: 80},
			CreatedAt: now.Add(-30 * time.Minute),
		},
	}

	t.Run("Sort by highest rated", func(t *testing.T) {
		sorted := SortModels(models, SortHighestRated)

		if sorted[0].ID != 2 {
			t.Errorf("Expected Model B first, got Model %d", sorted[0].ID)
		}
		if sorted[1].ID != 3 {
			t.Errorf("Expected Model C second, got Model %d", sorted[1].ID)
		}
		if sorted[2].ID != 1 {
			t.Errorf("Expected Model A third, got Model %d", sorted[2].ID)
		}
	})

	t.Run("Sort by most downloads", func(t *testing.T) {
		sorted := SortModels(models, SortMostDownload)

		if sorted[0].Stats.DownloadCount != 200 {
			t.Errorf("Expected 200 downloads first, got %d", sorted[0].Stats.DownloadCount)
		}
	})

	t.Run("Sort by newest", func(t *testing.T) {
		sorted := SortModels(models, SortNewest)

		if sorted[0].ID != 2 {
			t.Errorf("Expected newest model first, got Model %d", sorted[0].ID)
		}
	})

	t.Run("Empty models slice", func(t *testing.T) {
		sorted := SortModels([]Model{}, SortHighestRated)

		if len(sorted) != 0 {
			t.Errorf("Expected 0 models, got %d", len(sorted))
		}
	})
}

func TestModelMethods(t *testing.T) {
	model := Model{
		ID:   1,
		Name: "Test Model",
		Type: ModelTypeCheckpoint,
		Stats: Stats{
			Rating:        4.5,
			DownloadCount: 1000,
		},
		Tags:               []string{"anime", "character", "style"},
		AllowCommercialUse: []string{string(CommercialUseSell)},
		ModelVersions: []ModelVersion{
			{
				ID:        1,
				Name:      "Version 1.0",
				CreatedAt: time.Now().Add(-time.Hour),
				Files: []File{
					{ID: 1, Primary: true, SizeKB: 1000},
					{ID: 2, Primary: false, SizeKB: 500},
				},
			},
			{
				ID:        2,
				Name:      "Version 2.0",
				CreatedAt: time.Now(),
				Files: []File{
					{ID: 3, Primary: true, SizeKB: 1200},
				},
			},
		},
	}

	t.Run("GetLatestVersion", func(t *testing.T) {
		latest := model.GetLatestVersion()

		if latest == nil {
			t.Error("Expected latest version, got nil")
		}
		if latest.ID != 2 {
			t.Errorf("Expected version 2, got version %d", latest.ID)
		}
	})

	t.Run("GetLatestVersion empty", func(t *testing.T) {
		emptyModel := Model{}
		latest := emptyModel.GetLatestVersion()

		if latest != nil {
			t.Error("Expected nil for empty model, got version")
		}
	})

	t.Run("HasTag", func(t *testing.T) {
		if !model.HasTag("anime") {
			t.Error("Expected model to have 'anime' tag")
		}
		if !model.HasTag("ANIME") { // Case insensitive
			t.Error("Expected case insensitive tag matching")
		}
		if model.HasTag("nonexistent") {
			t.Error("Expected model not to have 'nonexistent' tag")
		}
	})

	t.Run("IsCommercialUseAllowed", func(t *testing.T) {
		if !model.IsCommercialUseAllowed() {
			t.Error("Expected commercial use to be allowed")
		}

		noCommercialModel := Model{
			AllowCommercialUse: []string{string(CommercialUseNone)},
		}
		if noCommercialModel.IsCommercialUseAllowed() {
			t.Error("Expected commercial use not to be allowed")
		}
	})

	t.Run("GetModelSummary", func(t *testing.T) {
		summary := model.GetModelSummary()
		expected := "Test Model (Checkpoint) - 1000 downloads, 4.5 rating, 2 versions"

		if summary != expected {
			t.Errorf("Expected '%s', got '%s'", expected, summary)
		}
	})
}

func TestModelVersionMethods(t *testing.T) {
	version := ModelVersion{
		ID:           1,
		Name:         "Test Version",
		BaseModel:    BaseModelSD1_5,
		TrainedWords: []string{"character", "anime"},
		Files: []File{
			{
				ID:       1,
				Primary:  true,
				SizeKB:   1024,
				Metadata: FileMetadata{Format: FileFormatSafeTensors},
			},
			{
				ID:       2,
				Primary:  false,
				SizeKB:   512,
				Metadata: FileMetadata{Format: FileFormatPickleTensor},
			},
		},
		Images: []Image{
			{ID: 1},
			{ID: 2},
		},
	}

	t.Run("GetPrimaryFile", func(t *testing.T) {
		primary := version.GetPrimaryFile()

		if primary == nil {
			t.Error("Expected primary file, got nil")
		}
		if primary.ID != 1 {
			t.Errorf("Expected file ID 1, got %d", primary.ID)
		}
	})

	t.Run("GetPrimaryFile no primary", func(t *testing.T) {
		noPrimaryVersion := ModelVersion{
			Files: []File{
				{ID: 1, Primary: false},
				{ID: 2, Primary: false},
			},
		}
		primary := noPrimaryVersion.GetPrimaryFile()

		if primary == nil {
			t.Error("Expected first file when no primary, got nil")
		}
		if primary.ID != 1 {
			t.Errorf("Expected first file ID 1, got %d", primary.ID)
		}
	})

	t.Run("GetPrimaryFile empty", func(t *testing.T) {
		emptyVersion := ModelVersion{}
		primary := emptyVersion.GetPrimaryFile()

		if primary != nil {
			t.Error("Expected nil for empty version, got file")
		}
	})

	t.Run("GetFileByFormat", func(t *testing.T) {
		file := version.GetFileByFormat(FileFormatSafeTensors)

		if file == nil {
			t.Error("Expected SafeTensor file, got nil")
		}
		if file.ID != 1 {
			t.Errorf("Expected file ID 1, got %d", file.ID)
		}

		noFile := version.GetFileByFormat(FileFormatCKPT)
		if noFile != nil {
			t.Error("Expected nil for non-existent format, got file")
		}
	})

	t.Run("GetDownloadSize", func(t *testing.T) {
		size := version.GetDownloadSize()
		expected := 1024.0 + 512.0

		if size != expected {
			t.Errorf("Expected size %.1f, got %.1f", expected, size)
		}
	})

	t.Run("GetTrainedWordsString", func(t *testing.T) {
		words := version.GetTrainedWordsString()
		expected := "character, anime"

		if words != expected {
			t.Errorf("Expected '%s', got '%s'", expected, words)
		}
	})

	t.Run("IsEarlyAccess", func(t *testing.T) {
		// Test non-early access (no timeframe)
		if version.IsEarlyAccess() {
			t.Error("Expected not early access")
		}

		// Test early access
		now := time.Now()
		earlyVersion := ModelVersion{
			EarlyAccessTimeFrame: 24, // 24 hours
			PublishedAt:          &now,
		}
		if !earlyVersion.IsEarlyAccess() {
			t.Error("Expected early access")
		}

		// Test expired early access
		pastTime := now.Add(-48 * time.Hour)
		expiredVersion := ModelVersion{
			EarlyAccessTimeFrame: 24,
			PublishedAt:          &pastTime,
		}
		if expiredVersion.IsEarlyAccess() {
			t.Error("Expected early access to be expired")
		}
	})

	t.Run("GetVersionSummary", func(t *testing.T) {
		summary := version.GetVersionSummary()
		expected := "Test Version (SD 1.5) - 1.0 MB, 2 images"

		if summary != expected {
			t.Errorf("Expected '%s', got '%s'", expected, summary)
		}
	})
}
