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

// Package civitai - Model Version Management and Analysis
//
// This file provides comprehensive utilities for working with AI model versions,
// including filtering, sorting, file analysis, and metadata extraction.
//
// # Version Filtering
//
// Filter model versions by various criteria:
//
//	// Filter by base model compatibility
//	sdxlVersions := civitai.FilterVersions(versions, func(v civitai.ModelVersion) bool {
//		return v.BaseModel == civitai.BaseModelSDXL10
//	})
//
//	// Filter by file format
//	safetensorVersions := civitai.FilterVersions(versions, func(v civitai.ModelVersion) bool {
//		return v.HasFormat("SafeTensor")
//	})
//
//	// Filter by size range (in bytes)
//	compactVersions := civitai.FilterVersions(versions, func(v civitai.ModelVersion) bool {
//		return v.GetDownloadSize() < 4*1024*1024*1024 // Less than 4GB
//	})
//
// # Version Sorting
//
// Sort versions by different criteria:
//
//	// Sort by newest first
//	civitai.SortVersions(versions, civitai.VersionSortByNewest)
//
//	// Sort by oldest first
//	civitai.SortVersions(versions, civitai.VersionSortByOldest)
//
// # File Analysis
//
// Analyze model files and formats:
//
//	version := versions[0]
//
//	// Get primary download file
//	primaryFile := version.GetPrimaryFile()
//	if primaryFile != nil {
//		fmt.Printf("Primary file: %s (%s)\n", primaryFile.Name, primaryFile.Type)
//		fmt.Printf("Size: %d bytes\n", primaryFile.SizeKB*1024)
//	}
//
//	// Get files by format
//	safetensorFiles := version.GetFilesByFormat("SafeTensor")
//	pickleFiles := version.GetFilesByFormat("PickleTensor")
//
//	// Check for clean/safe files
//	cleanFiles := version.GetCleanFiles()
//	fmt.Printf("Found %d verified clean files\n", len(cleanFiles))
//
// # Version Metadata
//
// Extract version information and statistics:
//
//	// Get training information
//	if len(version.TrainedWords) > 0 {
//		fmt.Printf("Trained words: %s\n", strings.Join(version.TrainedWords, ", "))
//	}
//
//	// Check version age
//	age := version.GetVersionAge()
//	fmt.Printf("Version age: %s\n", version.GetVersionAgeString())
//
//	// Get file statistics
//	stats := version.GetFileStats()
//	fmt.Printf("Total files: %d, Total size: %d bytes\n", stats.Count, stats.TotalSize)
//
//	// Check early access status
//	if version.IsEarlyAccess() {
//		fmt.Println("This is an early access version")
//	}
//
// # Version Utilities
//
// Additional helper functions for version management:
//
//	// Find specific version by ID
//	targetVersion := civitai.FindVersionByID(versions, 12345)
//	if targetVersion != nil {
//		fmt.Printf("Found version: %s\n", targetVersion.Name)
//	}
//
//	// Group versions by base model
//	grouped := civitai.GroupVersionsByBaseModel(versions)
//	for baseModel, versionList := range grouped {
//		fmt.Printf("%s: %d versions\n", baseModel, len(versionList))
//	}
//
// # Security and Safety
//
// Check file safety and scan results:
//
//	for _, file := range version.Files {
//		if civitai.IsFileClean(file) {
//			fmt.Printf("✅ %s is verified clean\n", file.Name)
//		} else {
//			fmt.Printf("⚠️ %s has scan issues\n", file.Name)
//		}
//	}

package civitai

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// VersionFilter provides filtering options for model version collections
type VersionFilter struct {
	BaseModels         []BaseModel
	FileFormats        []FileFormat
	MinSize            float64 // in KB
	MaxSize            float64 // in KB
	HasTrainedWords    *bool
	ExcludeEarlyAccess bool
}

// FilterVersions filters a slice of model versions based on the given criteria
func FilterVersions(versions []ModelVersion, filter VersionFilter) []ModelVersion {
	if len(versions) == 0 {
		return versions
	}

	var filtered []ModelVersion
	for _, version := range versions {
		if shouldIncludeVersion(version, filter) {
			filtered = append(filtered, version)
		}
	}

	return filtered
}

// shouldIncludeVersion checks if a version matches the filter criteria
func shouldIncludeVersion(version ModelVersion, filter VersionFilter) bool {
	// Filter by base model
	if len(filter.BaseModels) > 0 {
		baseMatch := false
		for _, baseModel := range filter.BaseModels {
			if version.BaseModel == baseModel {
				baseMatch = true
				break
			}
		}
		if !baseMatch {
			return false
		}
	}

	// Filter by file format
	if len(filter.FileFormats) > 0 {
		formatMatch := false
		for _, file := range version.Files {
			for _, format := range filter.FileFormats {
				if file.Metadata.Format == format {
					formatMatch = true
					break
				}
			}
			if formatMatch {
				break
			}
		}
		if !formatMatch {
			return false
		}
	}

	// Filter by size
	totalSize := version.GetDownloadSize()
	if filter.MinSize > 0 && totalSize < filter.MinSize {
		return false
	}
	if filter.MaxSize > 0 && totalSize > filter.MaxSize {
		return false
	}

	// Filter by trained words
	if filter.HasTrainedWords != nil {
		hasWords := len(version.TrainedWords) > 0
		if *filter.HasTrainedWords != hasWords {
			return false
		}
	}

	// Filter out early access if requested
	if filter.ExcludeEarlyAccess && version.IsEarlyAccess() {
		return false
	}

	return true
}

// SortVersions sorts a slice of model versions by creation date (newest first by default)
func SortVersions(versions []ModelVersion, newestFirst bool) []ModelVersion {
	if len(versions) == 0 {
		return versions
	}

	// Make a copy to avoid modifying the original slice
	sorted := make([]ModelVersion, len(versions))
	copy(sorted, versions)

	sort.Slice(sorted, func(i, j int) bool {
		if newestFirst {
			return sorted[i].CreatedAt.After(sorted[j].CreatedAt)
		}
		return sorted[i].CreatedAt.Before(sorted[j].CreatedAt)
	})

	return sorted
}

// GetFilesByFormat returns all files matching the specified format
func (mv *ModelVersion) GetFilesByFormat(format FileFormat) []File {
	var files []File
	for _, file := range mv.Files {
		if file.Metadata.Format == format {
			files = append(files, file)
		}
	}
	return files
}

// GetSafeTensorFiles returns all SafeTensor format files
func (mv *ModelVersion) GetSafeTensorFiles() []File {
	return mv.GetFilesByFormat(FileFormatSafeTensors)
}

// GetPickleFiles returns all Pickle format files
func (mv *ModelVersion) GetPickleFiles() []File {
	return mv.GetFilesByFormat(FileFormatPickleTensor)
}

// HasFormat checks if the version has files in the specified format
func (mv *ModelVersion) HasFormat(format FileFormat) bool {
	for _, file := range mv.Files {
		if file.Metadata.Format == format {
			return true
		}
	}
	return false
}

// GetCleanFiles returns files that have passed security scans
func (mv *ModelVersion) GetCleanFiles() []File {
	var cleanFiles []File
	for _, file := range mv.Files {
		if isFileClean(file) {
			cleanFiles = append(cleanFiles, file)
		}
	}
	return cleanFiles
}

// isFileClean checks if a file has passed security scans
func isFileClean(file File) bool {
	// Check pickle scan result
	if file.PickleScanResult != "" && !strings.EqualFold(file.PickleScanResult, "success") {
		return false
	}

	// Check virus scan result
	if file.VirusScanResult != "" && !strings.EqualFold(file.VirusScanResult, "success") {
		return false
	}

	return true
}

// GetCompatibleBaseModels returns a list of base models this version is compatible with
func (mv *ModelVersion) GetCompatibleBaseModels() []BaseModel {
	var models []BaseModel

	// Add the primary base model
	if mv.BaseModel != "" {
		models = append(models, mv.BaseModel)
	}

	// For certain model types, add compatible variants
	switch mv.BaseModel {
	case BaseModelSD1_5:
		// SD 1.5 models might work with SD 2.0 with some compatibility
		models = append(models, BaseModelSD2_0)
	case BaseModelSDXL:
		// SDXL is generally standalone
	case BaseModelSD2_0, BaseModelSD2_1:
		// SD 2.x models are generally compatible with each other
		if mv.BaseModel == BaseModelSD2_0 {
			models = append(models, BaseModelSD2_1)
		} else {
			models = append(models, BaseModelSD2_0)
		}
	}

	return models
}

// GetRecommendedFile returns the recommended file for download based on preferences
func (mv *ModelVersion) GetRecommendedFile() *File {
	// First preference: clean SafeTensor files
	safeTensorFiles := mv.GetSafeTensorFiles()
	for _, file := range safeTensorFiles {
		if isFileClean(file) {
			return &file
		}
	}

	// Second preference: primary file if clean
	primary := mv.GetPrimaryFile()
	if primary != nil && isFileClean(*primary) {
		return primary
	}

	// Third preference: any clean file
	cleanFiles := mv.GetCleanFiles()
	if len(cleanFiles) > 0 {
		return &cleanFiles[0]
	}

	// Last resort: any file
	if len(mv.Files) > 0 {
		return &mv.Files[0]
	}

	return nil
}

// GetVersionAge returns how long ago the version was created
func (mv *ModelVersion) GetVersionAge() time.Duration {
	return time.Since(mv.CreatedAt)
}

// GetVersionAgeString returns a human-readable age string
func (mv *ModelVersion) GetVersionAgeString() string {
	age := mv.GetVersionAge()

	switch {
	case age < time.Hour:
		return fmt.Sprintf("%d minutes ago", int(age.Minutes()))
	case age < 24*time.Hour:
		return fmt.Sprintf("%d hours ago", int(age.Hours()))
	case age < 30*24*time.Hour:
		return fmt.Sprintf("%d days ago", int(age.Hours()/24))
	case age < 365*24*time.Hour:
		return fmt.Sprintf("%d months ago", int(age.Hours()/(24*30)))
	default:
		return fmt.Sprintf("%d years ago", int(age.Hours()/(24*365)))
	}
}

// GetFileStats returns statistics about the files in this version
func (mv *ModelVersion) GetFileStats() map[string]interface{} {
	stats := make(map[string]interface{})

	stats["total_files"] = len(mv.Files)
	stats["total_size_kb"] = mv.GetDownloadSize()
	stats["total_size_mb"] = mv.GetDownloadSize() / 1024

	// Count by format
	formatCounts := make(map[FileFormat]int)
	for _, file := range mv.Files {
		formatCounts[file.Metadata.Format]++
	}
	stats["format_counts"] = formatCounts

	// Security scan status
	cleanFiles := len(mv.GetCleanFiles())
	stats["clean_files"] = cleanFiles
	stats["scan_pass_rate"] = float64(cleanFiles) / float64(len(mv.Files))

	return stats
}

// HasTrainedWords checks if the version has any trained words
func (mv *ModelVersion) HasTrainedWords() bool {
	return len(mv.TrainedWords) > 0
}

// GetTrainedWordsCount returns the number of trained words
func (mv *ModelVersion) GetTrainedWordsCount() int {
	return len(mv.TrainedWords)
}

// FindVersionByID finds a version with the specified ID from a slice
func FindVersionByID(versions []ModelVersion, id int) *ModelVersion {
	for i := range versions {
		if versions[i].ID == id {
			return &versions[i]
		}
	}
	return nil
}

// GroupVersionsByBaseModel groups versions by their base model
func GroupVersionsByBaseModel(versions []ModelVersion) map[BaseModel][]ModelVersion {
	groups := make(map[BaseModel][]ModelVersion)

	for _, version := range versions {
		baseModel := version.BaseModel
		if baseModel == "" {
			baseModel = BaseModelOther
		}
		groups[baseModel] = append(groups[baseModel], version)
	}

	return groups
}
