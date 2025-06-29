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


//go:build ignore

// Package main demonstrates image browsing and discovery capabilities
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

	// Example 1: Get recent safe images
	fmt.Println("=== Recent Safe Images ===")
	safeParams := civitai.ImageParams{
		Sort:  string(civitai.ImageSortNewest),
		NSFW:  string(civitai.NSFWLevelNone),
		Limit: 10,
	}

	images, _, err := client.GetImages(ctx, safeParams)
	if err != nil {
		log.Fatalf("Failed to get safe images: %v", err)
	}

	fmt.Printf("Found %d safe images\n", len(images))
	for _, image := range images {
		fmt.Printf("- Image %d: %dx%d by %s\n",
			image.ID, image.Width, image.Height, image.Username)
		if prompt, ok := image.Meta["prompt"].(string); ok && prompt != "" {
			// Truncate long prompts
			if len(prompt) > 100 {
				prompt = prompt[:97] + "..."
			}
			fmt.Printf("  Prompt: %s\n", prompt)
		}
	}

	// Example 2: Images from a specific model
	fmt.Println("\n=== Images from Specific Model ===")
	// First, find a popular model
	searchParams := civitai.SearchParams{
		Query: "realistic",
		Types: []civitai.ModelType{civitai.ModelTypeCheckpoint},
		Sort:  civitai.SortMostDownload,
		Limit: 1,
	}

	models, _, err := client.SearchModels(ctx, searchParams)
	if err != nil {
		log.Fatalf("Failed to search models: %v", err)
	}

	if len(models) > 0 {
		modelID := models[0].ID
		modelParams := civitai.ImageParams{
			ModelID: modelID,
			Sort:    string(civitai.ImageSortMostReactions),
			NSFW:    string(civitai.NSFWLevelNone),
			Limit:   5,
		}

		modelImages, _, err := client.GetImages(ctx, modelParams)
		if err != nil {
			log.Printf("Failed to get model images: %v", err)
		} else {
			fmt.Printf("Top images from model '%s':\n", models[0].Name)
			for _, image := range modelImages {
				fmt.Printf("- Image %d: %d reactions, %d comments\n",
					image.ID, image.Stats.LikeCount, image.Stats.CommentCount)
			}
		}
	}

	// Example 3: Images by username
	fmt.Println("\n=== Images by Creator ===")
	creatorParams := civitai.ImageParams{
		Username: "civitai",
		Sort:     string(civitai.ImageSortMostReactions),
		NSFW:     string(civitai.NSFWLevelNone),
		Limit:    5,
	}

	creatorImages, _, err := client.GetImages(ctx, creatorParams)
	if err != nil {
		log.Printf("Failed to get creator images: %v", err)
	} else {
		fmt.Printf("Top images by creator:\n")
		for _, image := range creatorImages {
			fmt.Printf("- Image %d: %d reactions\n",
				image.ID, image.Stats.LikeCount+image.Stats.HeartCount)
		}
	}

	// Example 4: Browse images by different time periods
	fmt.Println("\n=== Trending Images by Period ===")
	periods := []civitai.Period{
		civitai.PeriodDay,
		civitai.PeriodWeek,
		civitai.PeriodMonth,
	}

	for _, period := range periods {
		periodParams := civitai.ImageParams{
			Sort:   string(civitai.ImageSortMostReactions),
			Period: period,
			NSFW:   string(civitai.NSFWLevelNone),
			Limit:  3,
		}

		periodImages, _, err := client.GetImages(ctx, periodParams)
		if err != nil {
			log.Printf("Failed to get %s images: %v", period, err)
			continue
		}

		fmt.Printf("\nTop images from past %s:\n", period)
		for _, image := range periodImages {
			fmt.Printf("- Image %d: %d total reactions (by %s)\n",
				image.ID, image.Stats.LikeCount+image.Stats.HeartCount+image.Stats.CryCount+image.Stats.LaughCount, image.Username)
		}
	}

	// Example 5: Pagination through images
	fmt.Println("\n=== Image Pagination ===")
	page := 1
	totalImages := 0
	var nextCursor string

	paginationParams := civitai.ImageParams{
		Sort:  string(civitai.ImageSortNewest),
		NSFW:  string(civitai.NSFWLevelNone),
		Limit: 20,
	}

	// Get first few pages using cursor-based pagination
	for page <= 3 {
		if nextCursor != "" {
			// For subsequent requests, you would typically use the nextPage URL
			// This is just a demonstration of the concept
			fmt.Printf("Would fetch next page using cursor: %s\n", nextCursor)
			break
		}

		pageImages, pageMetadata, err := client.GetImages(ctx, paginationParams)
		if err != nil {
			log.Printf("Page %d failed: %v", page, err)
			break
		}

		fmt.Printf("Page %d: %d images\n", page, len(pageImages))
		totalImages += len(pageImages)

		nextCursor = pageMetadata.NextPage
		if nextCursor == "" {
			fmt.Println("No more pages available")
			break
		}

		page++
	}

	fmt.Printf("Total images browsed: %d\n", totalImages)

	// Example 6: Image metadata analysis
	fmt.Println("\n=== Image Metadata Analysis ===")
	metaParams := civitai.ImageParams{
		Sort:  string(civitai.ImageSortNewest),
		NSFW:  string(civitai.NSFWLevelNone),
		Limit: 5,
	}

	metaImages, _, err := client.GetImages(ctx, metaParams)
	if err != nil {
		log.Printf("Failed to get images for metadata analysis: %v", err)
	} else {
		for _, image := range metaImages {
			fmt.Printf("\nImage %d Analysis:\n", image.ID)
			fmt.Printf("- Dimensions: %dx%d\n", image.Width, image.Height)
			fmt.Printf("- NSFW Level: %s\n", image.NSFWLevel)
			fmt.Printf("- Reactions: %d likes, %d hearts\n",
				image.Stats.LikeCount, image.Stats.HeartCount)

			if model, ok := image.Meta["model"].(string); ok && model != "" {
				fmt.Printf("- Model: %s\n", model)
			}
			if sampler, ok := image.Meta["sampler"].(string); ok && sampler != "" {
				fmt.Printf("- Sampler: %s\n", sampler)
			}
			if steps, ok := image.Meta["steps"]; ok {
				if stepsFloat, ok := steps.(float64); ok && stepsFloat > 0 {
					fmt.Printf("- Steps: %.0f\n", stepsFloat)
				}
			}
			if cfgScale, ok := image.Meta["cfgScale"]; ok {
				if cfgFloat, ok := cfgScale.(float64); ok && cfgFloat > 0 {
					fmt.Printf("- CFG Scale: %.1f\n", cfgFloat)
				}
			}
		}
	}
}
