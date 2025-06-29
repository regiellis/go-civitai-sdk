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

func TestFilterVersions(t *testing.T) {
	versions := []ModelVersion{
		{
			ID:           1,
			Name:         "Version 1.0",
			BaseModel:    BaseModelSD1_5,
			TrainedWords: []string{"character"},
			Files: []File{
				{SizeKB: 1000, Metadata: FileMetadata{Format: FileFormatSafeTensors}},
			},
		},
		{
			ID:        2,
			Name:      "Version 2.0",
			BaseModel: BaseModelSDXL,
			Files: []File{
				{SizeKB: 3000, Metadata: FileMetadata{Format: FileFormatPickleTensor}},
			},
		},
		{
			ID:           3,
			Name:         "Version 3.0",
			BaseModel:    BaseModelSD1_5,
			TrainedWords: []string{"style", "anime"},
			Files: []File{
				{SizeKB: 500, Metadata: FileMetadata{Format: FileFormatSafeTensors}},
			},
		},
	}

	t.Run("Filter by base model", func(t *testing.T) {
		filter := VersionFilter{BaseModels: []BaseModel{BaseModelSD1_5}}
		filtered := FilterVersions(versions, filter)
		
		if len(filtered) != 2 {
			t.Errorf("Expected 2 versions, got %d", len(filtered))
		}
		
		for _, version := range filtered {
			if version.BaseModel != BaseModelSD1_5 {
				t.Errorf("Expected SD 1.5 model, got %s", version.BaseModel)
			}
		}
	})

	t.Run("Filter by file format", func(t *testing.T) {
		filter := VersionFilter{FileFormats: []FileFormat{FileFormatSafeTensors}}
		filtered := FilterVersions(versions, filter)
		
		if len(filtered) != 2 {
			t.Errorf("Expected 2 versions, got %d", len(filtered))
		}
		
		for _, version := range filtered {
			if !version.HasFormat(FileFormatSafeTensors) {
				t.Error("Expected version to have SafeTensor format")
			}
		}
	})

	t.Run("Filter by size range", func(t *testing.T) {
		filter := VersionFilter{MinSize: 600, MaxSize: 2000}
		filtered := FilterVersions(versions, filter)
		
		if len(filtered) != 1 {
			t.Errorf("Expected 1 version, got %d", len(filtered))
		}
		
		if filtered[0].ID != 1 {
			t.Errorf("Expected version 1, got version %d", filtered[0].ID)
		}
	})

	t.Run("Filter by trained words", func(t *testing.T) {
		hasWords := true
		filter := VersionFilter{HasTrainedWords: &hasWords}
		filtered := FilterVersions(versions, filter)
		
		if len(filtered) != 2 {
			t.Errorf("Expected 2 versions, got %d", len(filtered))
		}
		
		for _, version := range filtered {
			if !version.HasTrainedWords() {
				t.Error("Expected version to have trained words")
			}
		}
		
		// Test filter for no trained words
		noWords := false
		filter = VersionFilter{HasTrainedWords: &noWords}
		filtered = FilterVersions(versions, filter)
		
		if len(filtered) != 1 {
			t.Errorf("Expected 1 version, got %d", len(filtered))
		}
		
		if filtered[0].HasTrainedWords() {
			t.Error("Expected version to have no trained words")
		}
	})

	t.Run("Empty versions slice", func(t *testing.T) {
		filter := VersionFilter{BaseModels: []BaseModel{BaseModelSD1_5}}
		filtered := FilterVersions([]ModelVersion{}, filter)
		
		if len(filtered) != 0 {
			t.Errorf("Expected 0 versions, got %d", len(filtered))
		}
	})
}

func TestSortVersions(t *testing.T) {
	now := time.Now()
	versions := []ModelVersion{
		{ID: 1, Name: "Version A", CreatedAt: now.Add(-time.Hour)},
		{ID: 2, Name: "Version B", CreatedAt: now},
		{ID: 3, Name: "Version C", CreatedAt: now.Add(-30 * time.Minute)},
	}

	t.Run("Sort newest first", func(t *testing.T) {
		sorted := SortVersions(versions, true)
		
		if sorted[0].ID != 2 {
			t.Errorf("Expected Version B first, got Version %d", sorted[0].ID)
		}
		if sorted[1].ID != 3 {
			t.Errorf("Expected Version C second, got Version %d", sorted[1].ID)
		}
		if sorted[2].ID != 1 {
			t.Errorf("Expected Version A third, got Version %d", sorted[2].ID)
		}
	})

	t.Run("Sort oldest first", func(t *testing.T) {
		sorted := SortVersions(versions, false)
		
		if sorted[0].ID != 1 {
			t.Errorf("Expected Version A first, got Version %d", sorted[0].ID)
		}
		if sorted[2].ID != 2 {
			t.Errorf("Expected Version B last, got Version %d", sorted[2].ID)
		}
	})

	t.Run("Empty versions slice", func(t *testing.T) {
		sorted := SortVersions([]ModelVersion{}, true)
		
		if len(sorted) != 0 {
			t.Errorf("Expected 0 versions, got %d", len(sorted))
		}
	})
}

func TestVersionMethods(t *testing.T) {
	version := ModelVersion{
		ID:           1,
		Name:         "Test Version",
		BaseModel:    BaseModelSD1_5,
		CreatedAt:    time.Now().Add(-time.Hour),
		TrainedWords: []string{"character", "anime"},
		Files: []File{
			{
				ID:                1,
				Primary:           true,
				SizeKB:            1024,
				Metadata:          FileMetadata{Format: FileFormatSafeTensors},
				PickleScanResult:  "Success",
				VirusScanResult:   "Success",
			},
			{
				ID:                2,
				Primary:           false,
				SizeKB:            512,
				Metadata:          FileMetadata{Format: FileFormatPickleTensor},
				PickleScanResult:  "Success",
				VirusScanResult:   "Success",
			},
			{
				ID:                3,
				Primary:           false,
				SizeKB:            256,
				Metadata:          FileMetadata{Format: FileFormatSafeTensors},
				PickleScanResult:  "Failed",
				VirusScanResult:   "Success",
			},
		},
	}

	t.Run("GetFilesByFormat", func(t *testing.T) {
		safeTensorFiles := version.GetFilesByFormat(FileFormatSafeTensors)
		
		if len(safeTensorFiles) != 2 {
			t.Errorf("Expected 2 SafeTensor files, got %d", len(safeTensorFiles))
		}
		
		pickleFiles := version.GetFilesByFormat(FileFormatPickleTensor)
		if len(pickleFiles) != 1 {
			t.Errorf("Expected 1 Pickle file, got %d", len(pickleFiles))
		}
	})

	t.Run("GetSafeTensorFiles", func(t *testing.T) {
		files := version.GetSafeTensorFiles()
		
		if len(files) != 2 {
			t.Errorf("Expected 2 SafeTensor files, got %d", len(files))
		}
	})

	t.Run("GetPickleFiles", func(t *testing.T) {
		files := version.GetPickleFiles()
		
		if len(files) != 1 {
			t.Errorf("Expected 1 Pickle file, got %d", len(files))
		}
	})

	t.Run("HasFormat", func(t *testing.T) {
		if !version.HasFormat(FileFormatSafeTensors) {
			t.Error("Expected version to have SafeTensor format")
		}
		
		if !version.HasFormat(FileFormatPickleTensor) {
			t.Error("Expected version to have Pickle format")
		}
		
		if version.HasFormat(FileFormatCKPT) {
			t.Error("Expected version not to have CKPT format")
		}
	})

	t.Run("GetCleanFiles", func(t *testing.T) {
		cleanFiles := version.GetCleanFiles()
		
		if len(cleanFiles) != 2 {
			t.Errorf("Expected 2 clean files, got %d", len(cleanFiles))
		}
		
		// Verify the failed file is not included
		for _, file := range cleanFiles {
			if file.ID == 3 {
				t.Error("Expected failed file to be excluded from clean files")
			}
		}
	})

	t.Run("GetCompatibleBaseModels", func(t *testing.T) {
		compatible := version.GetCompatibleBaseModels()
		
		if len(compatible) != 2 {
			t.Errorf("Expected 2 compatible models, got %d", len(compatible))
		}
		
		// Should include SD 1.5 and SD 2.0
		hasSD15 := false
		hasSD20 := false
		for _, model := range compatible {
			if model == BaseModelSD1_5 {
				hasSD15 = true
			}
			if model == BaseModelSD2_0 {
				hasSD20 = true
			}
		}
		
		if !hasSD15 {
			t.Error("Expected SD 1.5 to be included in compatible models")
		}
		if !hasSD20 {
			t.Error("Expected SD 2.0 to be included in compatible models")
		}
	})

	t.Run("GetRecommendedFile", func(t *testing.T) {
		recommended := version.GetRecommendedFile()
		
		if recommended == nil {
			t.Error("Expected recommended file, got nil")
		}
		
		// Should prefer clean SafeTensor files
		if recommended.ID != 1 {
			t.Errorf("Expected file ID 1 (clean SafeTensor), got %d", recommended.ID)
		}
	})

	t.Run("GetVersionAge", func(t *testing.T) {
		age := version.GetVersionAge()
		
		if age < time.Hour {
			t.Error("Expected age to be at least 1 hour")
		}
		if age > 2*time.Hour {
			t.Error("Expected age to be less than 2 hours")
		}
	})

	t.Run("GetVersionAgeString", func(t *testing.T) {
		ageString := version.GetVersionAgeString()
		
		if ageString == "" {
			t.Error("Expected non-empty age string")
		}
		
		// Test with different ages
		testCases := []struct {
			age      time.Duration
			expected string
		}{
			{30 * time.Minute, "30 minutes ago"},
			{2 * time.Hour, "2 hours ago"},
			{3 * 24 * time.Hour, "3 days ago"},
		}
		
		for _, tc := range testCases {
			testVersion := ModelVersion{CreatedAt: time.Now().Add(-tc.age)}
			result := testVersion.GetVersionAgeString()
			if result != tc.expected {
				t.Errorf("Expected '%s', got '%s'", tc.expected, result)
			}
		}
	})

	t.Run("GetFileStats", func(t *testing.T) {
		stats := version.GetFileStats()
		
		if stats["total_files"] != 3 {
			t.Errorf("Expected 3 total files, got %v", stats["total_files"])
		}
		
		expectedSize := 1024.0 + 512.0 + 256.0
		if stats["total_size_kb"] != expectedSize {
			t.Errorf("Expected total size %.1f KB, got %v", expectedSize, stats["total_size_kb"])
		}
		
		if stats["clean_files"] != 2 {
			t.Errorf("Expected 2 clean files, got %v", stats["clean_files"])
		}
		
		expectedRate := 2.0 / 3.0
		if stats["scan_pass_rate"] != expectedRate {
			t.Errorf("Expected scan pass rate %.2f, got %v", expectedRate, stats["scan_pass_rate"])
		}
	})

	t.Run("HasTrainedWords", func(t *testing.T) {
		if !version.HasTrainedWords() {
			t.Error("Expected version to have trained words")
		}
		
		emptyVersion := ModelVersion{}
		if emptyVersion.HasTrainedWords() {
			t.Error("Expected empty version to have no trained words")
		}
	})

	t.Run("GetTrainedWordsCount", func(t *testing.T) {
		count := version.GetTrainedWordsCount()
		
		if count != 2 {
			t.Errorf("Expected 2 trained words, got %d", count)
		}
	})
}

func TestVersionUtilityFunctions(t *testing.T) {
	versions := []ModelVersion{
		{ID: 1, Name: "Version 1", BaseModel: BaseModelSD1_5},
		{ID: 2, Name: "Version 2", BaseModel: BaseModelSDXL},
		{ID: 3, Name: "Version 3", BaseModel: BaseModelSD1_5},
	}

	t.Run("FindVersionByID", func(t *testing.T) {
		found := FindVersionByID(versions, 2)
		
		if found == nil {
			t.Error("Expected to find version 2, got nil")
		}
		if found.ID != 2 {
			t.Errorf("Expected version ID 2, got %d", found.ID)
		}
		
		notFound := FindVersionByID(versions, 99)
		if notFound != nil {
			t.Error("Expected nil for non-existent version, got version")
		}
	})

	t.Run("GroupVersionsByBaseModel", func(t *testing.T) {
		groups := GroupVersionsByBaseModel(versions)
		
		if len(groups) != 2 {
			t.Errorf("Expected 2 groups, got %d", len(groups))
		}
		
		sd15Group := groups[BaseModelSD1_5]
		if len(sd15Group) != 2 {
			t.Errorf("Expected 2 SD 1.5 versions, got %d", len(sd15Group))
		}
		
		sdxlGroup := groups[BaseModelSDXL]
		if len(sdxlGroup) != 1 {
			t.Errorf("Expected 1 SDXL version, got %d", len(sdxlGroup))
		}
	})
}

func TestIsFileClean(t *testing.T) {
	t.Run("Clean file", func(t *testing.T) {
		cleanFile := File{
			PickleScanResult: "Success",
			VirusScanResult:  "Success",
		}
		
		if !isFileClean(cleanFile) {
			t.Error("Expected file to be clean")
		}
	})

	t.Run("Empty scan results", func(t *testing.T) {
		emptyFile := File{}
		
		if !isFileClean(emptyFile) {
			t.Error("Expected file with empty scan results to be considered clean")
		}
	})

	t.Run("Failed pickle scan", func(t *testing.T) {
		failedFile := File{
			PickleScanResult: "Failed",
			VirusScanResult:  "Success",
		}
		
		if isFileClean(failedFile) {
			t.Error("Expected file with failed pickle scan to be unclean")
		}
	})

	t.Run("Failed virus scan", func(t *testing.T) {
		failedFile := File{
			PickleScanResult: "Success",
			VirusScanResult:  "Failed",
		}
		
		if isFileClean(failedFile) {
			t.Error("Expected file with failed virus scan to be unclean")
		}
	})
}
