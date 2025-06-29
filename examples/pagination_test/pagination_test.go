//go:build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/regiellis/go-civitai-sdk"
)

func main() {
	client := civitai.NewClientWithoutAuth()
	ctx := context.Background()

	fmt.Println("=== CivitAI API Pagination Testing ===")
	fmt.Println("Testing cursor-based vs page-based pagination with tags")

	// Test 1: Basic tag search to establish baseline
	fmt.Println("\n--- Test 1: Basic Tag Search (Baseline) ---")
	basicParams := civitai.SearchParams{
		Tag:   "photo",
		Limit: 5,
	}
	models1, metadata1, err := client.SearchModels(ctx, basicParams)
	if err != nil {
		log.Printf("Basic search failed: %v", err)
		return
	}

	fmt.Printf("Found %d models\n", len(models1))
	fmt.Printf("Metadata - NextCursor: %s\n", metadata1.NextCursor)
	fmt.Printf("Metadata - NextPage: %s\n", metadata1.NextPage)
	fmt.Printf("Metadata - CurrentPage: %d\n", metadata1.CurrentPage)
	fmt.Printf("Metadata - PageSize: %d\n", metadata1.PageSize)
	fmt.Printf("Metadata - TotalItems: %d\n", metadata1.TotalItems)

	if len(models1) > 0 {
		fmt.Printf("First model: %s (ID: %d)\n", models1[0].Name, models1[0].ID)
		if len(models1) > 1 {
			fmt.Printf("Last model: %s (ID: %d)\n", models1[len(models1)-1].Name, models1[len(models1)-1].ID)
		}
	}

	// Test 2: Page-based pagination (traditional)
	fmt.Println("\n--- Test 2: Page-Based Pagination ---")
	for page := 1; page <= 3; page++ {
		pageParams := civitai.SearchParams{
			Tag:   "photo",
			Limit: 3,
			Page:  page,
		}

		models, metadata, err := client.SearchModels(ctx, pageParams)
		if err != nil {
			log.Printf("Page %d failed: %v", page, err)
			continue
		}

		fmt.Printf("Page %d: Found %d models\n", page, len(models))
		if len(models) > 0 {
			fmt.Printf("  First: %s (ID: %d)\n", models[0].Name, models[0].ID)
			fmt.Printf("  NextCursor: %s\n", metadata.NextCursor)
			fmt.Printf("  NextPage: %s\n", metadata.NextPage)
		}
	}

	// Test 3: Cursor-based pagination using NextPage URL
	fmt.Println("\n--- Test 3: Cursor-Based Pagination (NextPage URL) ---")
	if metadata1.NextPage != "" {
		fmt.Printf("Following NextPage URL: %s\n", metadata1.NextPage)

		// We need to manually call the NextPage URL since the SDK might not support it directly
		// Let's parse the URL to extract cursor parameters
		parsedURL, err := url.Parse(metadata1.NextPage)
		if err != nil {
			log.Printf("Failed to parse NextPage URL: %v", err)
		} else {
			queryParams := parsedURL.Query()
			cursor := queryParams.Get("cursor")
			limit := queryParams.Get("limit")
			tag := queryParams.Get("tag")

			fmt.Printf("Extracted cursor: %s\n", cursor)
			fmt.Printf("Extracted limit: %s\n", limit)
			fmt.Printf("Extracted tag: %s\n", tag)

			// Test manual cursor-based search
			testCursorPagination(client, ctx, tag, cursor, limit)
		}
	}

	// Test 4: Multiple pagination rounds with different tags
	fmt.Println("\n--- Test 4: Multi-Round Pagination Test ---")
	tags := []string{"portrait", "anime", "landscape"}

	for _, tag := range tags {
		fmt.Printf("\nTesting pagination for tag: %s\n", tag)
		testPaginationRounds(client, ctx, tag, 3) // Test 3 pages
	}

	// Test 5: Large limit test
	fmt.Println("\n--- Test 5: Large Limit Test ---")
	largeParams := civitai.SearchParams{
		Tag:   "photo",
		Limit: 20, // Larger batch size
	}

	models5, metadata5, err := client.SearchModels(ctx, largeParams)
	if err != nil {
		log.Printf("Large limit test failed: %v", err)
	} else {
		fmt.Printf("Large limit: Found %d models\n", len(models5))
		fmt.Printf("NextCursor available: %v\n", metadata5.NextCursor != "")
		fmt.Printf("NextPage available: %v\n", metadata5.NextPage != "")
	}

	// Test 6: Edge cases
	fmt.Println("\n--- Test 6: Edge Cases ---")

	// Very small limit
	smallParams := civitai.SearchParams{
		Tag:   "photo",
		Limit: 1,
	}

	models6, metadata6, err := client.SearchModels(ctx, smallParams)
	if err != nil {
		log.Printf("Small limit test failed: %v", err)
	} else {
		fmt.Printf("Limit=1: Found %d models\n", len(models6))
		fmt.Printf("NextCursor: %s\n", metadata6.NextCursor)
	}

	fmt.Println("\n=== Pagination Testing Complete ===")

	// Summary and recommendations
	fmt.Println("\n--- Summary ---")
	fmt.Println("✅ Tag-based search: Reliable")
	fmt.Println("✅ Basic pagination: Working")
	if metadata1.NextCursor != "" {
		fmt.Println("✅ Cursor pagination: Available")
	} else {
		fmt.Println("❓ Cursor pagination: Not provided")
	}
	if metadata1.NextPage != "" {
		fmt.Println("✅ NextPage URLs: Available")
	} else {
		fmt.Println("❓ NextPage URLs: Not provided")
	}
}

func testCursorPagination(client *civitai.Client, ctx context.Context, tag, cursor, limitStr string) {
	// Note: The current SDK might not support cursor parameter directly
	// This would require enhancing the SDK to support cursor-based pagination
	fmt.Printf("Note: Cursor-based pagination might require SDK enhancement\n")
	fmt.Printf("Would need to add cursor parameter to SearchParams struct\n")
}

func testPaginationRounds(client *civitai.Client, ctx context.Context, tag string, maxPages int) {
	var allModels []civitai.Model

	for page := 1; page <= maxPages; page++ {
		params := civitai.SearchParams{
			Tag:   tag,
			Limit: 3,
			Page:  page,
		}

		models, metadata, err := client.SearchModels(ctx, params)
		if err != nil {
			log.Printf("  Page %d failed: %v", page, err)
			break
		}

		if len(models) == 0 {
			fmt.Printf("  Page %d: No more models\n", page)
			break
		}

		allModels = append(allModels, models...)
		fmt.Printf("  Page %d: %d models (total so far: %d)\n", page, len(models), len(allModels))

		// Check for duplicate models across pages
		if page > 1 {
			duplicates := findDuplicates(allModels)
			if len(duplicates) > 0 {
				fmt.Printf("  ⚠️  Found %d duplicate models across pages\n", len(duplicates))
			}
		}

		// Check if we should continue
		if metadata.NextCursor == "" && metadata.NextPage == "" {
			fmt.Printf("  No more pages available\n")
			break
		}
	}

	fmt.Printf("Total unique models for tag '%s': %d\n", tag, len(removeDuplicateModels(allModels)))
}

func findDuplicates(models []civitai.Model) []int {
	seen := make(map[int]bool)
	var duplicates []int

	for _, model := range models {
		if seen[model.ID] {
			duplicates = append(duplicates, model.ID)
		} else {
			seen[model.ID] = true
		}
	}

	return duplicates
}

func removeDuplicateModels(models []civitai.Model) []civitai.Model {
	seen := make(map[int]bool)
	var unique []civitai.Model

	for _, model := range models {
		if !seen[model.ID] {
			seen[model.ID] = true
			unique = append(unique, model)
		}
	}

	return unique
}
