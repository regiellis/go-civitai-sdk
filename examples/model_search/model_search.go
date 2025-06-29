//go:build ignore

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

// Package main demonstrates advanced model searching capabilities
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

	// Example 1: Basic search
	fmt.Println("=== Basic Model Search ===")
	basicParams := civitai.SearchParams{
		Tag:   "portrait", // Using tag instead of query for reliability
		Limit: 5,
	}

	models, metadata, err := client.SearchModels(ctx, basicParams)
	if err != nil {
		log.Fatalf("Basic search failed: %v", err)
	}

	fmt.Printf("Found %d models (page %d of %d)\n",
		len(models), metadata.CurrentPage, metadata.TotalPages)
	for _, model := range models {
		fmt.Printf("- %s (%s)\n", model.Name, model.Type)
	}

	// Example 2: Advanced search with filters
	fmt.Println("\n=== Advanced Model Search ===")
	advancedParams := civitai.SearchParams{
		// Using only tag for reliability, removed problematic query parameter
		Tag:   "photorealistic",
		Types: []civitai.ModelType{civitai.ModelTypeCheckpoint},
		Sort:  civitai.SortHighestRated,
		Limit: 10,
	}

	models, metadata, err = client.SearchModels(ctx, advancedParams)
	if err != nil {
		log.Fatalf("Advanced search failed: %v", err)
	}

	fmt.Printf("Found %d filtered models\n", len(models))
	for _, model := range models {
		fmt.Printf("- %s (Rating: %.1f, Downloads: %d)\n",
			model.Name, model.Stats.Rating, model.Stats.DownloadCount)
	}

	// Example 3: Search by specific creator
	fmt.Println("\n=== Search by Creator ===")
	creatorParams := civitai.SearchParams{
		Username: "RunDiffusion", // Known creator with popular models
		Types:    []civitai.ModelType{civitai.ModelTypeCheckpoint},
		Sort:     civitai.SortMostDownload,
		Limit:    3,
	}

	models, _, err = client.SearchModels(ctx, creatorParams)
	if err != nil {
		log.Fatalf("Creator search failed: %v", err)
	}

	fmt.Printf("Found %d models by specific creator\n", len(models))
	for _, model := range models {
		fmt.Printf("- %s by %s\n", model.Name, model.Creator.Username)
	}

	// Example 4: Search different model types
	fmt.Println("\n=== Search by Model Type ===")
	types := []civitai.ModelType{
		civitai.ModelTypeCheckpoint,
		civitai.ModelTypeLORA,
		civitai.ModelTypeEmbedding,
		civitai.ModelTypeControlNet,
	}

	for _, modelType := range types {
		typeParams := civitai.SearchParams{
			Types: []civitai.ModelType{modelType},
			Sort:  civitai.SortMostDownload,
			Limit: 3,
		}

		typeModels, _, err := client.SearchModels(ctx, typeParams)
		if err != nil {
			log.Printf("Search for %s failed: %v", modelType, err)
			continue
		}

		fmt.Printf("\nTop %s models:\n", modelType)
		for _, model := range typeModels {
			fmt.Printf("- %s (%d downloads)\n", model.Name, model.Stats.DownloadCount)
		}
	}

	// Example 5: Pagination example
	fmt.Println("\n=== Pagination Example ===")
	page := 1
	totalFound := 0

	for page <= 3 { // Get first 3 pages
		pageParams := civitai.SearchParams{
			Query: "girl",
			Page:  page,
			Limit: 5,
		}

		pageModels, pageMetadata, err := client.SearchModels(ctx, pageParams)
		if err != nil {
			log.Printf("Page %d search failed: %v", page, err)
			break
		}

		fmt.Printf("Page %d: %d models\n", page, len(pageModels))
		totalFound += len(pageModels)

		if page >= pageMetadata.TotalPages {
			break
		}
		page++
	}

	fmt.Printf("Total models found across pages: %d\n", totalFound)
}
