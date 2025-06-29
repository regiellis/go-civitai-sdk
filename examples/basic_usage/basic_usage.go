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

// Package example demonstrates how to use the go-civitai-sdk
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/regiellis/go-civitai-sdk"
)

func main() {
	// Create a client without authentication for public endpoints
	client := civitai.NewClientWithoutAuth()

	// Or create a client with authentication
	// client := civitai.NewClient(os.Getenv("CIVITAI_TOKEN"))

	// You can also customize the client with options
	clientWithOptions := civitai.NewClient(
		os.Getenv("CIVITAI_TOKEN"),
		civitai.WithTimeout(60*time.Second),
		civitai.WithUserAgent("my-app/1.0.0"),
	)

	ctx := context.Background()

	// Example 1: Search for models
	fmt.Println("=== Searching for photo models ===")
	searchParams := civitai.SearchParams{
		// Using tag instead of query - API has issues with query parameter
		Tag:   "photo",
		Limit: 5,
	}

	models, metadata, err := client.SearchModels(ctx, searchParams)
	if err != nil {
		log.Fatalf("Failed to search models: %v", err)
	}

	fmt.Printf("Found %d models\n", len(models))
	if metadata.TotalItems > 0 {
		fmt.Printf("(total available: %d)\n", metadata.TotalItems)
	}
	for _, model := range models {
		fmt.Printf("- %s (%s) - %d downloads\n",
			model.Name, model.Type, model.Stats.DownloadCount)
	}

	// Example 2: Get a specific model
	fmt.Println("\n=== Getting specific model ===")
	if len(models) > 0 {
		modelDetail, err := client.GetModel(ctx, models[0].ID)
		if err != nil {
			log.Printf("Failed to get model details: %v", err)
		} else {
			fmt.Printf("Model: %s\n", modelDetail.Name)
			fmt.Printf("Description: %s\n", modelDetail.Description)
			fmt.Printf("Versions: %d\n", len(modelDetail.ModelVersions))
		}
	}

	// Example 3: Get model version by hash
	fmt.Println("\n=== Getting model version by hash ===")
	// This is an example hash - you would use a real hash from a downloaded file
	hash := "5493A0EC49E72336B89F7E0A0BF9B2B2E03F3E2E9E7A6F8B5F3C3E9A3C9E2F9"
	versionByHash, err := clientWithOptions.GetModelVersionByHash(ctx, hash)
	if err != nil {
		log.Printf("Failed to get model version by hash: %v", err)
	} else {
		fmt.Printf("Found model version: %s for model: %s\n",
			versionByHash.Name, versionByHash.Model.Name)
	}

	// Example 4: Search for creators
	fmt.Println("\n=== Searching for creators ===")
	creatorParams := civitai.CreatorParams{
		Query: "civitai",
		Limit: 3,
	}

	creators, _, err := client.GetCreators(ctx, creatorParams)
	if err != nil {
		log.Printf("Failed to search creators: %v", err)
	} else {
		for _, creator := range creators {
			fmt.Printf("- %s (%d models)\n", creator.Username, creator.ModelCount)
		}
	}

	// Example 5: Get images
	fmt.Println("\n=== Getting images ===")
	imageParams := civitai.ImageParams{
		Limit: 5,
		Sort:  string(civitai.ImageSortNewest),
		NSFW:  string(civitai.NSFWLevelNone), // Safe images only
	}

	images, _, err := client.GetImages(ctx, imageParams)
	if err != nil {
		log.Printf("Failed to get images: %v", err)
	} else {
		for _, image := range images {
			fmt.Printf("- Image %d: %dx%d (%s)\n",
				image.ID, image.Width, image.Height, image.Username)
		}
	}

	// Example 6: Get tags
	fmt.Println("\n=== Getting tags ===")
	tagParams := civitai.TagParams{
		Query: "style",
		Limit: 5,
	}

	tags, _, err := client.GetTags(ctx, tagParams)
	if err != nil {
		log.Printf("Failed to get tags: %v", err)
	} else {
		for _, tag := range tags {
			fmt.Printf("- %s (%d models)\n", tag.Name, tag.ModelCount)
		}
	}

	// Example 7: Health check
	fmt.Println("\n=== Health check ===")
	if err := client.Health(ctx); err != nil {
		log.Printf("API health check failed: %v", err)
	} else {
		fmt.Println("API is healthy!")
	}

	// Example 8: Working with model versions
	fmt.Println("\n=== Getting model versions ===")
	if len(models) > 0 && len(models[0].ModelVersions) > 0 {
		versionID := models[0].ModelVersions[0].ID
		version, err := client.GetModelVersion(ctx, versionID)
		if err != nil {
			log.Printf("Failed to get model version: %v", err)
		} else {
			fmt.Printf("Version: %s\n", version.Name)
			fmt.Printf("Files: %d\n", len(version.Files))
			fmt.Printf("Download URL: %s\n", version.DownloadURL)
		}

		// Get all versions for this model
		allVersions, err := client.GetModelVersionsByModelID(ctx, models[0].ID)
		if err != nil {
			log.Printf("Failed to get all model versions: %v", err)
		} else {
			fmt.Printf("Total versions for model: %d\n", len(allVersions))
		}
	}

	fmt.Println("\n=== SDK Demo Complete! ===")
}
