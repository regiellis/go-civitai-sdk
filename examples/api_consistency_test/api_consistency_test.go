package main

import (
	"context"
	"fmt"
	"time"

	"github.com/regiellis/go-civitai-sdk"
)

func main() {
	client := civitai.NewClientWithoutAuth()
	ctx := context.Background()

	fmt.Println("=== CivitAI API Consistency Testing ===")
	fmt.Println("Testing query vs tag behavior across all endpoints")

	// Test 1: Models endpoint (we know this has issues)
	fmt.Println("\n--- Test 1: Models Endpoint ---")
	testModelsEndpoint(client, ctx)

	// Test 2: Images endpoint
	fmt.Println("\n--- Test 2: Images Endpoint ---")
	testImagesEndpoint(client, ctx)

	// Test 3: Creators endpoint
	fmt.Println("\n--- Test 3: Creators Endpoint ---")
	testCreatorsEndpoint(client, ctx)

	// Test 4: Tags endpoint
	fmt.Println("\n--- Test 4: Tags Endpoint ---")
	testTagsEndpoint(client, ctx)

	// Test 5: Individual model lookups
	fmt.Println("\n--- Test 5: Individual Model Lookups ---")
	testIndividualLookups(client, ctx)

	// Test 6: Version endpoints
	fmt.Println("\n--- Test 6: Model Version Endpoints ---")
	testVersionEndpoints(client, ctx)

	// Test 7: Different search terms consistency
	fmt.Println("\n--- Test 7: Search Term Consistency ---")
	testSearchTermConsistency(client, ctx)

	// Test 8: Timeout and reliability patterns
	fmt.Println("\n--- Test 8: Reliability Patterns ---")
	testReliabilityPatterns(client, ctx)

	fmt.Println("\n=== API Consistency Testing Complete ===")
}

func testModelsEndpoint(client *civitai.Client, ctx context.Context) {
	fmt.Println("Testing Models endpoint...")

	// Test query parameter
	fmt.Println("  Query parameter test:")
	for i := 0; i < 3; i++ {
		params := civitai.SearchParams{
			Query: "photo",
			Limit: 5,
		}
		models, metadata, err := client.SearchModels(ctx, params)
		if err != nil {
			fmt.Printf("    Attempt %d: FAILED - %v\n", i+1, err)
		} else {
			fmt.Printf("    Attempt %d: %d models, total: %d\n", i+1, len(models), metadata.TotalItems)
		}
		time.Sleep(100 * time.Millisecond) // Small delay between requests
	}

	// Test tag parameter
	fmt.Println("  Tag parameter test:")
	for i := 0; i < 3; i++ {
		params := civitai.SearchParams{
			Tag:   "photo",
			Limit: 5,
		}
		models, metadata, err := client.SearchModels(ctx, params)
		if err != nil {
			fmt.Printf("    Attempt %d: FAILED - %v\n", i+1, err)
		} else {
			fmt.Printf("    Attempt %d: %d models, cursor: %s\n", i+1, len(models), metadata.NextCursor != "")
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func testImagesEndpoint(client *civitai.Client, ctx context.Context) {
	fmt.Println("Testing Images endpoint...")

	// Test basic images
	fmt.Println("  Basic images test:")
	for i := 0; i < 3; i++ {
		params := civitai.ImageParams{
			Limit: 5,
			Sort:  string(civitai.ImageSortNewest),
			NSFW:  string(civitai.NSFWLevelNone),
		}
		images, metadata, err := client.GetImages(ctx, params)
		if err != nil {
			fmt.Printf("    Attempt %d: FAILED - %v\n", i+1, err)
		} else {
			fmt.Printf("    Attempt %d: %d images, cursor: %s\n", i+1, len(images), metadata.NextCursor != "")
		}
		time.Sleep(100 * time.Millisecond)
	}

	// Test images with username filter
	fmt.Println("  Images with username filter:")
	for i := 0; i < 3; i++ {
		params := civitai.ImageParams{
			Limit:    3,
			Username: "civitai",
			Sort:     string(civitai.ImageSortNewest),
			NSFW:     string(civitai.NSFWLevelNone),
		}
		images, _, err := client.GetImages(ctx, params)
		if err != nil {
			fmt.Printf("    Attempt %d: FAILED - %v\n", i+1, err)
		} else {
			fmt.Printf("    Attempt %d: %d images\n", i+1, len(images))
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func testCreatorsEndpoint(client *civitai.Client, ctx context.Context) {
	fmt.Println("Testing Creators endpoint...")

	// Test basic creators
	fmt.Println("  Basic creators test:")
	for i := 0; i < 3; i++ {
		params := civitai.CreatorParams{
			Limit: 5,
		}
		creators, _, err := client.GetCreators(ctx, params)
		if err != nil {
			fmt.Printf("    Attempt %d: FAILED - %v\n", i+1, err)
		} else {
			fmt.Printf("    Attempt %d: %d creators\n", i+1, len(creators))
		}
		time.Sleep(100 * time.Millisecond)
	}

	// Test creators with query
	fmt.Println("  Creators with query:")
	for i := 0; i < 3; i++ {
		params := civitai.CreatorParams{
			Query: "ai",
			Limit: 3,
		}
		creators, _, err := client.GetCreators(ctx, params)
		if err != nil {
			fmt.Printf("    Attempt %d: FAILED - %v\n", i+1, err)
		} else {
			fmt.Printf("    Attempt %d: %d creators\n", i+1, len(creators))
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func testTagsEndpoint(client *civitai.Client, ctx context.Context) {
	fmt.Println("Testing Tags endpoint...")

	// Test basic tags
	fmt.Println("  Basic tags test:")
	for i := 0; i < 3; i++ {
		params := civitai.TagParams{
			Limit: 5,
		}
		tags, _, err := client.GetTags(ctx, params)
		if err != nil {
			fmt.Printf("    Attempt %d: FAILED - %v\n", i+1, err)
		} else {
			fmt.Printf("    Attempt %d: %d tags\n", i+1, len(tags))
		}
		time.Sleep(100 * time.Millisecond)
	}

	// Test tags with query
	fmt.Println("  Tags with query:")
	for i := 0; i < 3; i++ {
		params := civitai.TagParams{
			Query: "style",
			Limit: 5,
		}
		tags, _, err := client.GetTags(ctx, params)
		if err != nil {
			fmt.Printf("    Attempt %d: FAILED - %v\n", i+1, err)
		} else {
			fmt.Printf("    Attempt %d: %d tags\n", i+1, len(tags))
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func testIndividualLookups(client *civitai.Client, ctx context.Context) {
	fmt.Println("Testing Individual lookups...")

	// Test individual model lookup
	fmt.Println("  Individual model lookup:")
	modelID := 133005 // Juggernaut XL
	for i := 0; i < 3; i++ {
		model, err := client.GetModel(ctx, modelID)
		if err != nil {
			fmt.Printf("    Attempt %d: FAILED - %v\n", i+1, err)
		} else {
			fmt.Printf("    Attempt %d: SUCCESS - %s (%d versions)\n", i+1, model.Name, len(model.ModelVersions))
		}
		time.Sleep(100 * time.Millisecond)
	}

	// Test health check
	fmt.Println("  Health check:")
	for i := 0; i < 3; i++ {
		err := client.Health(ctx)
		if err != nil {
			fmt.Printf("    Attempt %d: FAILED - %v\n", i+1, err)
		} else {
			fmt.Printf("    Attempt %d: SUCCESS\n", i+1)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func testVersionEndpoints(client *civitai.Client, ctx context.Context) {
	fmt.Println("Testing Version endpoints...")

	// Test model version lookup
	fmt.Println("  Model version lookup:")
	versionID := 1759168 // A known version ID
	for i := 0; i < 3; i++ {
		version, err := client.GetModelVersion(ctx, versionID)
		if err != nil {
			fmt.Printf("    Attempt %d: FAILED - %v\n", i+1, err)
		} else {
			fmt.Printf("    Attempt %d: SUCCESS - %s (%d files)\n", i+1, version.Name, len(version.Files))
		}
		time.Sleep(100 * time.Millisecond)
	}

	// Test version by hash (this often fails)
	fmt.Println("  Version by hash lookup:")
	hash := "5493A0EC49E72336B89F7E0A0BF9B2B2E03F3E2E9E7A6F8B5F3C3E9A3C9E2F9"
	for i := 0; i < 3; i++ {
		version, err := client.GetModelVersionByHash(ctx, hash)
		if err != nil {
			fmt.Printf("    Attempt %d: FAILED - %v\n", i+1, err)
		} else {
			fmt.Printf("    Attempt %d: SUCCESS - %s\n", i+1, version.Name)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func testSearchTermConsistency(client *civitai.Client, ctx context.Context) {
	fmt.Println("Testing search term consistency...")

	searchTerms := []string{"photo", "portrait", "anime", "landscape", "style"}

	for _, term := range searchTerms {
		fmt.Printf("  Testing term: %s\n", term)

		// Test as query
		queryResults := 0
		for i := 0; i < 3; i++ {
			params := civitai.SearchParams{
				Query: term,
				Limit: 3,
			}
			models, _, err := client.SearchModels(ctx, params)
			if err == nil {
				queryResults += len(models)
			}
			time.Sleep(50 * time.Millisecond)
		}

		// Test as tag
		tagResults := 0
		for i := 0; i < 3; i++ {
			params := civitai.SearchParams{
				Tag:   term,
				Limit: 3,
			}
			models, _, err := client.SearchModels(ctx, params)
			if err == nil {
				tagResults += len(models)
			}
			time.Sleep(50 * time.Millisecond)
		}

		fmt.Printf("    Query results (3 attempts): %d total\n", queryResults)
		fmt.Printf("    Tag results (3 attempts): %d total\n", tagResults)

		if tagResults > queryResults {
			fmt.Printf("    ✅ Tag search more reliable for '%s'\n", term)
		} else if queryResults > tagResults {
			fmt.Printf("    ⚠️  Query search better for '%s'\n", term)
		} else {
			fmt.Printf("    ➖ Similar results for '%s'\n", term)
		}
	}
}

func testReliabilityPatterns(client *civitai.Client, ctx context.Context) {
	fmt.Println("Testing reliability patterns...")

	endpoints := []struct {
		name string
		test func() error
	}{
		{
			name: "Models (tag search)",
			test: func() error {
				params := civitai.SearchParams{Tag: "photo", Limit: 3}
				_, _, err := client.SearchModels(ctx, params)
				return err
			},
		},
		{
			name: "Models (query search)",
			test: func() error {
				params := civitai.SearchParams{Query: "photo", Limit: 3}
				_, _, err := client.SearchModels(ctx, params)
				return err
			},
		},
		{
			name: "Images",
			test: func() error {
				params := civitai.ImageParams{Limit: 3, NSFW: string(civitai.NSFWLevelNone)}
				_, _, err := client.GetImages(ctx, params)
				return err
			},
		},
		{
			name: "Creators",
			test: func() error {
				params := civitai.CreatorParams{Limit: 3}
				_, _, err := client.GetCreators(ctx, params)
				return err
			},
		},
		{
			name: "Tags",
			test: func() error {
				params := civitai.TagParams{Limit: 3}
				_, _, err := client.GetTags(ctx, params)
				return err
			},
		},
		{
			name: "Health check",
			test: func() error {
				return client.Health(ctx)
			},
		},
	}

	attempts := 5
	for _, endpoint := range endpoints {
		fmt.Printf("  Testing %s reliability (%d attempts):\n", endpoint.name, attempts)

		successes := 0
		failures := 0
		var lastError error

		for i := 0; i < attempts; i++ {
			err := endpoint.test()
			if err != nil {
				failures++
				lastError = err
				fmt.Printf("    Attempt %d: ❌ FAILED\n", i+1)
			} else {
				successes++
				fmt.Printf("    Attempt %d: ✅ SUCCESS\n", i+1)
			}
			time.Sleep(200 * time.Millisecond)
		}

		successRate := float64(successes) / float64(attempts) * 100
		fmt.Printf("    Success rate: %.1f%% (%d/%d)\n", successRate, successes, attempts)

		if failures > 0 {
			fmt.Printf("    Last error: %v\n", lastError)
		}

		if successRate >= 80 {
			fmt.Printf("    Status: 🟢 RELIABLE\n")
		} else if successRate >= 50 {
			fmt.Printf("    Status: 🟡 INTERMITTENT\n")
		} else {
			fmt.Printf("    Status: 🔴 UNRELIABLE\n")
		}

		fmt.Println()
	}
}
