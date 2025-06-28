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

// Package main is a simple test program for the CivitAI SDK
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/regiellis/go-civitai-sdk"
)

func main() {
	// Create a client without authentication
	client := civitai.NewClientWithoutAuth()
	ctx := context.Background()

	fmt.Println("🤖 CivitAI SDK Test")
	fmt.Println("==================")

	// Test 1: Health check
	fmt.Print("⚡ Testing API health... ")
	if err := client.Health(ctx); err != nil {
		log.Fatalf("❌ Health check failed: %v", err)
	}
	fmt.Println("✅ API is healthy!")

	// Test 2: Search for models
	fmt.Print("🔍 Searching for models... ")
	searchParams := civitai.SearchParams{
		Query: "anime",
		Types: []civitai.ModelType{civitai.ModelTypeCheckpoint},
		Limit: 3,
	}

	models, metadata, err := client.SearchModels(ctx, searchParams)
	if err != nil {
		log.Fatalf("❌ Model search failed: %v", err)
	}
	fmt.Printf("✅ Found %d models (total: %d)\n", len(models), metadata.TotalItems)

	// Display model details
	for i, model := range models {
		fmt.Printf("   %d. %s (%s) - %d downloads\n", 
			i+1, model.Name, model.Type, model.Stats.DownloadCount)
	}

	// Test 3: Get model details
	if len(models) > 0 {
		fmt.Print("📋 Getting model details... ")
		modelDetail, err := client.GetModel(ctx, models[0].ID)
		if err != nil {
			log.Printf("❌ Failed to get model details: %v", err)
		} else {
			fmt.Printf("✅ Model: %s has %d versions\n", 
				modelDetail.Name, len(modelDetail.ModelVersions))
		}
	}

	// Test 4: Get images
	fmt.Print("🖼️  Getting images... ")
	imageParams := civitai.ImageParams{
		Sort:  string(civitai.ImageSortNewest),
		NSFW:  string(civitai.NSFWLevelNone),
		Limit: 3,
	}

	images, _, err := client.GetImages(ctx, imageParams)
	if err != nil {
		log.Printf("❌ Failed to get images: %v", err)
	} else {
		fmt.Printf("✅ Found %d safe images\n", len(images))
		for i, image := range images {
			fmt.Printf("   %d. Image %d (%dx%d) by %s\n", 
				i+1, image.ID, image.Width, image.Height, image.Username)
		}
	}

	// Test 5: Get creators
	fmt.Print("👥 Getting creators... ")
	creatorParams := civitai.CreatorParams{
		Limit: 3,
	}

	creators, _, err := client.GetCreators(ctx, creatorParams)
	if err != nil {
		log.Printf("❌ Failed to get creators: %v", err)
	} else {
		fmt.Printf("✅ Found %d creators\n", len(creators))
		for i, creator := range creators {
			fmt.Printf("   %d. %s (%d models)\n", 
				i+1, creator.Username, creator.ModelCount)
		}
	}

	// Test 6: Get tags
	fmt.Print("🏷️  Getting tags... ")
	tagParams := civitai.TagParams{
		Query: "style",
		Limit: 3,
	}

	tags, _, err := client.GetTags(ctx, tagParams)
	if err != nil {
		log.Printf("❌ Failed to get tags: %v", err)
	} else {
		fmt.Printf("✅ Found %d style tags\n", len(tags))
		for i, tag := range tags {
			fmt.Printf("   %d. %s (%d models)\n", 
				i+1, tag.Name, tag.ModelCount)
		}
	}

	fmt.Println("\n🎉 All SDK tests completed successfully!")
	fmt.Println("📖 Check the examples/ directory for more usage examples.")
}
