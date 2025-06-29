//go:build ignore

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/regiellis/go-civitai-sdk"
)

func main() {
	client := civitai.NewClientWithoutAuth()
	ctx := context.Background()

	fmt.Println("=== CivitAI Cursor-Based Pagination Test ===")

	// Test proper cursor-based pagination
	fmt.Println("\n--- Cursor-Based Pagination (Proper Implementation) ---")

	var allModels []civitai.Model
	cursor := ""
	page := 1
	limit := 5
	maxPages := 3

	for page <= maxPages {
		fmt.Printf("\nPage %d (cursor: %s):\n", page, cursor)

		params := civitai.SearchParams{
			Tag:    "photo",
			Limit:  limit,
			Cursor: cursor,
		}

		models, metadata, err := client.SearchModels(ctx, params)
		if err != nil {
			log.Printf("Page %d failed: %v", page, err)
			break
		}

		if len(models) == 0 {
			fmt.Printf("No more models found\n")
			break
		}

		fmt.Printf("Found %d models:\n", len(models))
		for i, model := range models {
			fmt.Printf("  %d. %s (ID: %d)\n", i+1, model.Name, model.ID)
		}

		allModels = append(allModels, models...)

		// Check for duplicates across pages
		duplicates := findDuplicateIDs(allModels)
		if len(duplicates) > 0 {
			fmt.Printf("⚠️  Warning: Found %d duplicate model IDs: %v\n", len(duplicates), duplicates)
		}

		// Get next cursor
		if metadata.NextCursor == "" {
			fmt.Printf("No more pages available\n")
			break
		}

		cursor = metadata.NextCursor
		fmt.Printf("Next cursor: %s\n", cursor)
		fmt.Printf("NextPage URL: %s\n", metadata.NextPage)

		page++
	}

	fmt.Printf("\n--- Results Summary ---\n")
	fmt.Printf("Total models collected: %d\n", len(allModels))
	unique := removeDuplicates(allModels)
	fmt.Printf("Unique models: %d\n", len(unique))

	if len(allModels) > len(unique) {
		fmt.Printf("⚠️  Found %d duplicate entries\n", len(allModels)-len(unique))
	} else {
		fmt.Printf("✅ No duplicate entries found\n")
	}

	// Test different tags with cursor pagination
	fmt.Println("\n--- Testing Different Tags ---")
	tags := []string{"portrait", "anime", "landscape"}

	for _, tag := range tags {
		fmt.Printf("\nTesting tag: %s\n", tag)
		testTagCursorPagination(client, ctx, tag, 2) // 2 pages max
	}

	// Test cursor edge cases
	fmt.Println("\n--- Testing Edge Cases ---")

	// Test with invalid cursor
	fmt.Println("\nTest 1: Invalid cursor")
	invalidParams := civitai.SearchParams{
		Tag:    "photo",
		Limit:  3,
		Cursor: "invalid-cursor",
	}
	_, _, err := client.SearchModels(ctx, invalidParams)
	if err != nil {
		fmt.Printf("✅ Invalid cursor properly rejected: %v\n", err)
	} else {
		fmt.Printf("⚠️  Invalid cursor was accepted\n")
	}

	// Test with very small limit
	fmt.Println("\nTest 2: Very small limit (1)")
	smallParams := civitai.SearchParams{
		Tag:   "photo",
		Limit: 1,
	}
	models, metadata, err := client.SearchModels(ctx, smallParams)
	if err != nil {
		fmt.Printf("❌ Small limit failed: %v\n", err)
	} else {
		fmt.Printf("✅ Small limit works: %d models, cursor: %s\n", len(models), metadata.NextCursor)
	}

	fmt.Println("\n=== Test Complete ===")
}

func testTagCursorPagination(client *civitai.Client, ctx context.Context, tag string, maxPages int) {
	cursor := ""
	var allModels []civitai.Model

	for page := 1; page <= maxPages; page++ {
		params := civitai.SearchParams{
			Tag:    tag,
			Limit:  3,
			Cursor: cursor,
		}

		models, metadata, err := client.SearchModels(ctx, params)
		if err != nil {
			fmt.Printf("  Page %d failed: %v\n", page, err)
			return
		}

		if len(models) == 0 {
			fmt.Printf("  No more models for tag '%s'\n", tag)
			break
		}

		allModels = append(allModels, models...)
		fmt.Printf("  Page %d: %d models (total: %d)\n", page, len(models), len(allModels))

		if metadata.NextCursor == "" {
			fmt.Printf("  No more pages available\n")
			break
		}

		cursor = metadata.NextCursor
	}

	unique := removeDuplicates(allModels)
	if len(allModels) == len(unique) {
		fmt.Printf("  ✅ No duplicates for tag '%s'\n", tag)
	} else {
		fmt.Printf("  ⚠️  Found %d duplicates for tag '%s'\n", len(allModels)-len(unique), tag)
	}
}

func findDuplicateIDs(models []civitai.Model) []int {
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

func removeDuplicates(models []civitai.Model) []civitai.Model {
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
