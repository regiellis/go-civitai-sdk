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

// Package main demonstrates creator and tag discovery capabilities
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

	// Example 1: Search for popular creators
	fmt.Println("=== Popular Creators ===")
	popularParams := civitai.CreatorParams{
		Limit: 10,
	}

	creators, metadata, err := client.GetCreators(ctx, popularParams)
	if err != nil {
		log.Fatalf("Failed to get creators: %v", err)
	}

	fmt.Printf("Found %d creators (total: %d)\n", len(creators), metadata.TotalItems)
	for i, creator := range creators {
		fmt.Printf("%d. %s - %d models\n", i+1, creator.Username, creator.ModelCount)
	}

	// Example 2: Search for specific creators
	fmt.Println("\n=== Search Creators by Name ===")
	searchParams := civitai.CreatorParams{
		Query: "anime",
		Limit: 5,
	}

	searchCreators, _, err := client.GetCreators(ctx, searchParams)
	if err != nil {
		log.Printf("Failed to search creators: %v", err)
	} else {
		fmt.Printf("Creators with 'anime' in their profile:\n")
		for _, creator := range searchCreators {
			fmt.Printf("- %s (%d models)\n", creator.Username, creator.ModelCount)
		}
	}

	// Example 3: Explore different tag categories
	fmt.Println("\n=== Popular Tags ===")
	tagParams := civitai.TagParams{
		Limit: 15,
	}

	tags, _, err := client.GetTags(ctx, tagParams)
	if err != nil {
		log.Printf("Failed to get tags: %v", err)
	} else {
		fmt.Printf("Most popular tags:\n")
		for i, tag := range tags {
			fmt.Printf("%d. %s (%d models)\n", i+1, tag.Name, tag.ModelCount)
		}
	}

	// Example 4: Search for specific tag categories
	fmt.Println("\n=== Style Tags ===")
	styleParams := civitai.TagParams{
		Query: "style",
		Limit: 10,
	}

	styleTags, _, err := client.GetTags(ctx, styleParams)
	if err != nil {
		log.Printf("Failed to get style tags: %v", err)
	} else {
		for _, tag := range styleTags {
			fmt.Printf("- %s (%d models)\n", tag.Name, tag.ModelCount)
		}
	}

	// Example 5: Character tags
	fmt.Println("\n=== Character Tags ===")
	characterParams := civitai.TagParams{
		Query: "character",
		Limit: 8,
	}

	characterTags, _, err := client.GetTags(ctx, characterParams)
	if err != nil {
		log.Printf("Failed to get character tags: %v", err)
	} else {
		for _, tag := range characterTags {
			fmt.Printf("- %s (%d models)\n", tag.Name, tag.ModelCount)
		}
	}

	// Example 6: Concept tags
	fmt.Println("\n=== Concept Tags ===")
	conceptParams := civitai.TagParams{
		Query: "concept",
		Limit: 8,
	}

	conceptTags, _, err := client.GetTags(ctx, conceptParams)
	if err != nil {
		log.Printf("Failed to get concept tags: %v", err)
	} else {
		for _, tag := range conceptTags {
			fmt.Printf("- %s (%d models)\n", tag.Name, tag.ModelCount)
		}
	}

	// Example 7: Use tags to find related models
	fmt.Println("\n=== Models with Popular Tags ===")
	if len(tags) > 0 {
		// Use the most popular tag to find models
		popularTag := tags[0].Name

		tagModelParams := civitai.SearchParams{
			Tag:   popularTag,
			Sort:  civitai.SortMostDownload,
			Limit: 5,
		}

		tagModels, _, err := client.SearchModels(ctx, tagModelParams)
		if err != nil {
			log.Printf("Failed to search models by tag: %v", err)
		} else {
			fmt.Printf("Top models tagged with '%s':\n", popularTag)
			for _, model := range tagModels {
				fmt.Printf("- %s (%d downloads)\n", model.Name, model.Stats.DownloadCount)
			}
		}
	}

	// Example 8: Discover trending creators and their content
	fmt.Println("\n=== Creator Deep Dive ===")
	if len(creators) > 0 {
		// Pick a creator with multiple models
		var selectedCreator *civitai.Creator
		for _, creator := range creators {
			if creator.ModelCount >= 3 {
				selectedCreator = &creator
				break
			}
		}

		if selectedCreator != nil {
			fmt.Printf("Exploring creator: %s\n", selectedCreator.Username)

			// Find their models
			creatorModelParams := civitai.SearchParams{
				Username: selectedCreator.Username,
				Sort:     civitai.SortMostDownload,
				Limit:    5,
			}

			creatorModels, _, err := client.SearchModels(ctx, creatorModelParams)
			if err != nil {
				log.Printf("Failed to get creator models: %v", err)
			} else {
				fmt.Printf("Top models by %s:\n", selectedCreator.Username)
				for _, model := range creatorModels {
					fmt.Printf("- %s (%s) - %d downloads\n",
						model.Name, model.Type, model.Stats.DownloadCount)
				}
			}

			// Find images from this creator
			creatorImageParams := civitai.ImageParams{
				Username: selectedCreator.Username,
				Sort:     string(civitai.ImageSortMostReactions),
				NSFW:     string(civitai.NSFWLevelNone),
				Limit:    3,
			}

			creatorImages, _, err := client.GetImages(ctx, creatorImageParams)
			if err != nil {
				log.Printf("Failed to get creator images: %v", err)
			} else {
				fmt.Printf("Popular images by %s:\n", selectedCreator.Username)
				for _, image := range creatorImages {
					fmt.Printf("- Image %d: %dx%d (%d reactions)\n",
						image.ID, image.Width, image.Height,
						image.Stats.LikeCount+image.Stats.HeartCount)
				}
			}
		}
	}

	// Example 9: Cross-reference tags and creators
	fmt.Println("\n=== Tag and Creator Analysis ===")
	if len(styleTags) > 0 && len(creators) > 0 {
		styleTag := styleTags[0].Name

		// Find creators who work with this style
		styleCreatorParams := civitai.SearchParams{
			Tag:   styleTag,
			Sort:  civitai.SortMostDownload,
			Limit: 3,
		}

		styleModels, _, err := client.SearchModels(ctx, styleCreatorParams)
		if err != nil {
			log.Printf("Failed to search models by style tag: %v", err)
		} else {
			fmt.Printf("Creators working with '%s' style:\n", styleTag)
			creatorSet := make(map[string]bool)
			for _, model := range styleModels {
				if !creatorSet[model.Creator.Username] {
					fmt.Printf("- %s (model: %s)\n", model.Creator.Username, model.Name)
					creatorSet[model.Creator.Username] = true
				}
			}
		}
	}

	fmt.Println("\n=== Creator and Tag Discovery Complete! ===")
}
