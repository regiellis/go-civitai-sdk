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

	// Test 1: No query search
	fmt.Println("=== Test 1: No query search ===")
	models1, _, err := client.SearchModels(ctx, civitai.SearchParams{
		Limit: 1,
	})
	if err != nil {
		log.Fatalf("Test 1 failed: %v", err)
	}
	fmt.Printf("✅ Test 1 passed: Found %d models\n", len(models1))

	// Test 2: Simple search with debug
	fmt.Println("\n=== Test 2: Simple search with debug ===")

	// First try direct curl to see what should work
	fmt.Println("Direct curl test:")
	// Let's manually build the URL to see what we get
	models2, metadata2, err := client.SearchModels(ctx, civitai.SearchParams{
		Query: "anime",
		Limit: 1,
	})
	if err != nil {
		fmt.Printf("❌ Test 2 failed: %v\n", err)
		// Try different query
		fmt.Println("Trying different query...")
		models2, metadata2, err = client.SearchModels(ctx, civitai.SearchParams{
			Query: "realistic",
			Limit: 1,
		})
		if err != nil {
			log.Fatalf("Both queries failed: %v", err)
		}
	}
	fmt.Printf("✅ Test 2 passed: Found %d models\n", len(models2))
	if metadata2 != nil {
		fmt.Printf("   Metadata: totalItems=%d\n", metadata2.TotalItems)
	}

	// Test 3: Add Sort
	fmt.Println("\n=== Test 3: With Sort ===")
	models3, _, err := client.SearchModels(ctx, civitai.SearchParams{
		Query: "anime",
		Sort:  civitai.SortMostDownload,
		Limit: 1,
	})
	if err != nil {
		log.Fatalf("Test 3 failed: %v", err)
	}
	fmt.Printf("✅ Test 3 passed: Found %d models\n", len(models3))

	fmt.Println("\n🎉 All tests passed!")
}
