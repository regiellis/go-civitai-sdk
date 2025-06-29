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

// Package civitai - Model Utilities and Helper Functions
//
// This file provides utility functions for working with AI models,
// including filtering, sorting, and data extraction helpers.
//
// # Filtering Models
//
// Filter models by various criteria:
//
//	// Filter by model type
//	checkpoints := civitai.FilterModels(models, func(m civitai.Model) bool {
//		return m.Type == civitai.ModelTypeCheckpoint
//	})
//
//	// Filter by rating
//	highRated := civitai.FilterModels(models, func(m civitai.Model) bool {
//		return m.Stats.Rating >= 4.5
//	})
//
//	// Filter by tag
//	animeModels := civitai.FilterModels(models, func(m civitai.Model) bool {
//		return m.HasTag("anime")
//	})
//
// # Sorting Models
//
// Sort models by different criteria:
//
//	// Sort by highest rated
//	civitai.SortModels(models, civitai.ModelSortByRating)
//
//	// Sort by most downloads
//	civitai.SortModels(models, civitai.ModelSortByDownloads)
//
//	// Sort by newest first
//	civitai.SortModels(models, civitai.ModelSortByNewest)
//
// # Model Methods
//
// Models provide convenient methods for common operations:
//
//	model := models[0]
//
//	// Check commercial use permission
//	if model.IsCommercialUseAllowed() {
//		fmt.Println("Commercial use allowed")
//	}
//
//	// Get latest version
//	latest := model.GetLatestVersion()
//	if latest != nil {
//		fmt.Printf("Latest version: %s\n", latest.Name)
//	}
//
//	// Check for specific tags
//	if model.HasTag("realistic") {
//		fmt.Println("This is a realistic model")
//	}
//
//	// Get model summary
//	summary := model.GetModelSummary()
//	fmt.Printf("Model: %s (%d downloads)\n", summary.Name, summary.Downloads)

package civitai

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// ModelFilter provides filtering options for model collections
type ModelFilter struct {
	Types     []ModelType
	NSFW      *bool
	MinRating float64
	Tags      []string
}

// FilterModels filters a slice of models based on the given criteria
func FilterModels(models []Model, filter ModelFilter) []Model {
	if len(models) == 0 {
		return models
	}

	var filtered []Model
	for _, model := range models {
		if shouldIncludeModel(model, filter) {
			filtered = append(filtered, model)
		}
	}

	return filtered
}

// shouldIncludeModel checks if a model matches the filter criteria
func shouldIncludeModel(model Model, filter ModelFilter) bool {
	// Filter by model type
	if len(filter.Types) > 0 {
		typeMatch := false
		for _, t := range filter.Types {
			if model.Type == t {
				typeMatch = true
				break
			}
		}
		if !typeMatch {
			return false
		}
	}

	// Filter by NSFW setting
	if filter.NSFW != nil && model.NSFW != *filter.NSFW {
		return false
	}

	// Filter by minimum rating
	if filter.MinRating > 0 && model.Stats.Rating < filter.MinRating {
		return false
	}

	// Filter by tags (if model has at least one matching tag)
	if len(filter.Tags) > 0 {
		tagMatch := false
		for _, filterTag := range filter.Tags {
			for _, modelTag := range model.Tags {
				if strings.EqualFold(modelTag, filterTag) {
					tagMatch = true
					break
				}
			}
			if tagMatch {
				break
			}
		}
		if !tagMatch {
			return false
		}
	}

	return true
}

// SortModels sorts a slice of models by the specified criteria
func SortModels(models []Model, sortBy SortType) []Model {
	if len(models) == 0 {
		return models
	}

	// Make a copy to avoid modifying the original slice
	sorted := make([]Model, len(models))
	copy(sorted, models)

	sort.Slice(sorted, func(i, j int) bool {
		switch sortBy {
		case SortHighestRated:
			return sorted[i].Stats.Rating > sorted[j].Stats.Rating
		case SortMostDownload:
			return sorted[i].Stats.DownloadCount > sorted[j].Stats.DownloadCount
		case SortMostLiked:
			return sorted[i].Stats.ThumbsUpCount > sorted[j].Stats.ThumbsUpCount
		case SortNewest:
			return sorted[i].CreatedAt.After(sorted[j].CreatedAt)
		case SortOldest:
			return sorted[i].CreatedAt.Before(sorted[j].CreatedAt)
		default:
			return sorted[i].Stats.DownloadCount > sorted[j].Stats.DownloadCount
		}
	})

	return sorted
}

// GetLatestVersion returns the most recently created model version
func (m *Model) GetLatestVersion() *ModelVersion {
	if len(m.ModelVersions) == 0 {
		return nil
	}

	latest := &m.ModelVersions[0]
	for i := 1; i < len(m.ModelVersions); i++ {
		if m.ModelVersions[i].CreatedAt.After(latest.CreatedAt) {
			latest = &m.ModelVersions[i]
		}
	}

	return latest
}

// GetPrimaryFile returns the primary file from the model version
func (mv *ModelVersion) GetPrimaryFile() *File {
	for i := range mv.Files {
		if mv.Files[i].Primary {
			return &mv.Files[i]
		}
	}

	// If no primary file is marked, return the first file
	if len(mv.Files) > 0 {
		return &mv.Files[0]
	}

	return nil
}

// GetFileByFormat returns the first file matching the specified format
func (mv *ModelVersion) GetFileByFormat(format FileFormat) *File {
	for i := range mv.Files {
		if mv.Files[i].Metadata.Format == format {
			return &mv.Files[i]
		}
	}
	return nil
}

// HasTag checks if the model has a specific tag (case-insensitive)
func (m *Model) HasTag(tag string) bool {
	for _, modelTag := range m.Tags {
		if strings.EqualFold(modelTag, tag) {
			return true
		}
	}
	return false
}

// IsCommercialUseAllowed checks if the model allows commercial use
func (m *Model) IsCommercialUseAllowed() bool {
	for _, use := range m.AllowCommercialUse {
		if use != string(CommercialUseNone) {
			return true
		}
	}
	return false
}

// GetDownloadSize returns the total download size in KB for all files in the version
func (mv *ModelVersion) GetDownloadSize() float64 {
	var totalSize float64
	for _, file := range mv.Files {
		totalSize += file.SizeKB
	}
	return totalSize
}

// GetTrainedWordsString returns trained words as a comma-separated string
func (mv *ModelVersion) GetTrainedWordsString() string {
	return strings.Join(mv.TrainedWords, ", ")
}

// IsEarlyAccess checks if the model version is still in early access
func (mv *ModelVersion) IsEarlyAccess() bool {
	if mv.EarlyAccessTimeFrame <= 0 {
		return false
	}

	if mv.PublishedAt == nil {
		return false
	}

	earlyAccessEnd := mv.PublishedAt.Add(time.Duration(mv.EarlyAccessTimeFrame) * time.Hour)
	return time.Now().Before(earlyAccessEnd)
}

// GetModelSummary returns a formatted summary string for the model
func (m *Model) GetModelSummary() string {
	return fmt.Sprintf("%s (%s) - %d downloads, %.1f rating, %d versions",
		m.Name,
		m.Type,
		m.Stats.DownloadCount,
		m.Stats.Rating,
		len(m.ModelVersions),
	)
}

// GetVersionSummary returns a formatted summary string for the model version
func (mv *ModelVersion) GetVersionSummary() string {
	primaryFile := mv.GetPrimaryFile()
	fileInfo := "no files"
	if primaryFile != nil {
		fileInfo = fmt.Sprintf("%.1f MB", primaryFile.SizeKB/1024)
	}

	return fmt.Sprintf("%s (%s) - %s, %d images",
		mv.Name,
		mv.BaseModel,
		fileInfo,
		len(mv.Images),
	)
}
