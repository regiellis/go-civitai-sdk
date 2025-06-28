/*
Copyright (c) 2025 Regi Ellis

This file is part of Go CivitAI SDK.

Licensed under the Restricted Use License - Non-Commercial Only.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://github.com/regiellis/go-civitai-sdk/blob/main/LICENSE

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Original work by Regi Ellis (https://github.com/regiellis)
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
		Query: "portrait",
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
		Query:                 "realistic woman",
		Types:                 []civitai.ModelType{civitai.ModelTypeCheckpoint},
		Sort:                  civitai.SortHighestRated,
		Period:                civitai.PeriodMonth,
		Rating:                4, // Minimum rating of 4
		Tag:                   "photorealistic",
		AllowNoCredit:         false,
		AllowDerivatives:      true,
		AllowDifferentLicense: false,
		AllowCommercialUse:    []string{string(civitai.CommercialUseSell)},
		NSFW:                  &[]bool{false}[0], // Safe content only
		SupportsGeneration:    &[]bool{true}[0],  // Models that support generation
		Limit:                 10,
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
		Username: "civitai",
		Types:    []civitai.ModelType{civitai.ModelTypeCheckpoint, civitai.ModelTypeLORA},
		Sort:     civitai.SortMostDownload,
		Limit:    5,
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
			Query: "anime",
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
