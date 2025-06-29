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

	// Test 1: Query parameter (problematic)
	fmt.Println("=== Test 1: Using query parameter ===")
	queryParams := civitai.SearchParams{
		Query: "photo",
		Limit: 5,
	}
	models1, metadata1, err := client.SearchModels(ctx, queryParams)
	if err != nil {
		log.Printf("Query search failed: %v", err)
	} else {
		fmt.Printf("Query search found %d models (total: %d)\n", len(models1), metadata1.TotalItems)
	}

	// Test 2: Tag parameter (should work)
	fmt.Println("\n=== Test 2: Using tag parameter ===")
	tagParams := civitai.SearchParams{
		Tag:   "photo",
		Limit: 5,
	}
	models2, metadata2, err := client.SearchModels(ctx, tagParams)
	if err != nil {
		log.Printf("Tag search failed: %v", err)
	} else {
		fmt.Printf("Tag search found %d models (total: %d)\n", len(models2), metadata2.TotalItems)
	}

	// Test 3: No search parameters (should work)
	fmt.Println("\n=== Test 3: No search parameters ===")
	emptyParams := civitai.SearchParams{
		Limit: 5,
	}
	models3, metadata3, err := client.SearchModels(ctx, emptyParams)
	if err != nil {
		log.Printf("Empty search failed: %v", err)
	} else {
		fmt.Printf("Empty search found %d models (total: %d)\n", len(models3), metadata3.TotalItems)
	}

	// Test 4: Different tag
	fmt.Println("\n=== Test 4: Using different tag ===")
	tagParams2 := civitai.SearchParams{
		Tag:   "portrait",
		Limit: 5,
	}
	models4, metadata4, err := client.SearchModels(ctx, tagParams2)
	if err != nil {
		log.Printf("Portrait tag search failed: %v", err)
	} else {
		fmt.Printf("Portrait tag search found %d models (total: %d)\n", len(models4), metadata4.TotalItems)
	}
}
