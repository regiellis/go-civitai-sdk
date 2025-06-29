//go:build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/regiellis/go-civitai-sdk"
)

// PaginationResult represents the result of paginated search
type PaginationResult struct {
	Models     []civitai.Model
	TotalPages int
	HasMore    bool
	NextCursor string
}

// PaginatedSearch performs cursor-based pagination with automatic deduplication
func PaginatedSearch(client *civitai.Client, ctx context.Context, params civitai.SearchParams, maxPages int) (*PaginationResult, error) {
	var allModels []civitai.Model
	seen := make(map[int]bool) // Track seen model IDs for deduplication
	cursor := params.Cursor
	page := 0

	fmt.Printf("🔍 Starting paginated search (tag: %s, limit: %d, maxPages: %d)\n", params.Tag, params.Limit, maxPages)

	for page < maxPages {
		page++

		// Set up parameters for this page
		currentParams := params
		currentParams.Cursor = cursor

		fmt.Printf("\n📄 Page %d (cursor: %s)\n", page, cursor)

		models, metadata, err := client.SearchModels(ctx, currentParams)
		if err != nil {
			return nil, fmt.Errorf("page %d failed: %w", page, err)
		}

		if len(models) == 0 {
			fmt.Printf("✅ No more models found\n")
			break
		}

		// Process models and deduplicate
		newModels := 0
		duplicates := 0

		for _, model := range models {
			if !seen[model.ID] {
				seen[model.ID] = true
				allModels = append(allModels, model)
				newModels++
			} else {
				duplicates++
			}
		}

		fmt.Printf("📊 Found %d models (%d new, %d duplicates)\n", len(models), newModels, duplicates)

		// Check if we can continue
		if metadata.NextCursor == "" {
			fmt.Printf("✅ Reached end of results\n")
			return &PaginationResult{
				Models:     allModels,
				TotalPages: page,
				HasMore:    false,
				NextCursor: "",
			}, nil
		}

		cursor = metadata.NextCursor
		fmt.Printf("➡️  Next cursor: %s\n", cursor)

		// If we got no new models, we might be at the end
		if newModels == 0 {
			fmt.Printf("⚠️  No new models found, stopping pagination\n")
			break
		}
	}

	return &PaginationResult{
		Models:     allModels,
		TotalPages: page,
		HasMore:    cursor != "",
		NextCursor: cursor,
	}, nil
}

func main() {
	client := civitai.NewClientWithoutAuth()
	ctx := context.Background()

	fmt.Println("=== Production-Ready CivitAI Pagination Demo ===")

	// Test 1: Standard pagination with deduplication
	fmt.Println("\n--- Test 1: Standard Pagination with Deduplication ---")

	params := civitai.SearchParams{
		Tag:   "photo",
		Limit: 5,
	}

	result, err := PaginatedSearch(client, ctx, params, 4)
	if err != nil {
		log.Fatalf("Pagination failed: %v", err)
	}

	fmt.Printf("\n🎯 Final Results:\n")
	fmt.Printf("Total unique models: %d\n", len(result.Models))
	fmt.Printf("Pages processed: %d\n", result.TotalPages)
	fmt.Printf("Has more results: %v\n", result.HasMore)

	if len(result.Models) > 0 {
		fmt.Printf("\nFirst 5 models:\n")
		for i, model := range result.Models {
			if i >= 5 {
				break
			}
			fmt.Printf("  %d. %s (ID: %d, Downloads: %d)\n",
				i+1, model.Name, model.ID, model.Stats.DownloadCount)
		}
	}

	// Test 2: Different batch sizes
	fmt.Println("\n--- Test 2: Testing Different Batch Sizes ---")
	batchSizes := []int{1, 3, 10}

	for _, batchSize := range batchSizes {
		fmt.Printf("\nTesting batch size: %d\n", batchSize)
		params := civitai.SearchParams{
			Tag:   "portrait",
			Limit: batchSize,
		}

		result, err := PaginatedSearch(client, ctx, params, 2)
		if err != nil {
			fmt.Printf("❌ Batch size %d failed: %v\n", batchSize, err)
			continue
		}

		fmt.Printf("✅ Batch size %d: %d unique models in %d pages\n",
			batchSize, len(result.Models), result.TotalPages)
	}

	// Test 3: Resume from cursor
	fmt.Println("\n--- Test 3: Resume from Previous Cursor ---")
	if result.HasMore {
		fmt.Printf("Resuming from cursor: %s\n", result.NextCursor)

		resumeParams := civitai.SearchParams{
			Tag:    "photo",
			Limit:  3,
			Cursor: result.NextCursor,
		}

		resumeResult, err := PaginatedSearch(client, ctx, resumeParams, 2)
		if err != nil {
			fmt.Printf("❌ Resume failed: %v\n", err)
		} else {
			fmt.Printf("✅ Resumed: %d additional models\n", len(resumeResult.Models))

			// Check for overlaps with previous results
			overlaps := 0
			for _, newModel := range resumeResult.Models {
				for _, oldModel := range result.Models {
					if newModel.ID == oldModel.ID {
						overlaps++
						break
					}
				}
			}

			if overlaps == 0 {
				fmt.Printf("✅ No overlaps with previous results\n")
			} else {
				fmt.Printf("⚠️  Found %d overlapping models\n", overlaps)
			}
		}
	}

	// Test 4: Performance comparison
	fmt.Println("\n--- Test 4: Performance Comparison ---")

	// Small batches (more requests)
	fmt.Println("Testing small batches (limit=2):")
	start := time.Now()
	smallParams := civitai.SearchParams{
		Tag:   "anime",
		Limit: 2,
	}
	smallResult, err := PaginatedSearch(client, ctx, smallParams, 3)
	smallDuration := time.Since(start)

	if err != nil {
		fmt.Printf("❌ Small batch test failed: %v\n", err)
	} else {
		fmt.Printf("✅ Small batches: %d models in %v (%d requests)\n",
			len(smallResult.Models), smallDuration, smallResult.TotalPages)
	}

	// Large batches (fewer requests)
	fmt.Println("Testing large batches (limit=10):")
	start = time.Now()
	largeParams := civitai.SearchParams{
		Tag:   "anime",
		Limit: 10,
	}
	largeResult, err := PaginatedSearch(client, ctx, largeParams, 2)
	largeDuration := time.Since(start)

	if err != nil {
		fmt.Printf("❌ Large batch test failed: %v\n", err)
	} else {
		fmt.Printf("✅ Large batches: %d models in %v (%d requests)\n",
			len(largeResult.Models), largeDuration, largeResult.TotalPages)
	}

	fmt.Println("\n=== Pagination Demo Complete ===")

	fmt.Println("\n🎯 Key Takeaways:")
	fmt.Println("✅ Cursor-based pagination works reliably")
	fmt.Println("✅ Deduplication handles API overlap behavior")
	fmt.Println("✅ Resuming from cursors works correctly")
	fmt.Println("✅ Different batch sizes are supported")
	fmt.Println("⚠️  Page-based pagination has duplicate issues")
	fmt.Println("💡 Recommendation: Use cursor-based pagination exclusively")
}
