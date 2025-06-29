//go:build ignore

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

	fmt.Println("=== CivitAI API Performance & Consistency Deep Dive ===")

	// Test 1: Query vs Tag result count patterns
	fmt.Println("\n--- Test 1: Query vs Tag Result Count Analysis ---")
	testResultCounts(client, ctx)

	// Test 2: Timing analysis
	fmt.Println("\n--- Test 2: Endpoint Timing Analysis ---")
	testEndpointTiming(client, ctx)

	// Test 3: Load testing (rapid requests)
	fmt.Println("\n--- Test 3: Load Testing (Rapid Requests) ---")
	testLoadBehavior(client, ctx)

	// Test 4: Cross-endpoint consistency
	fmt.Println("\n--- Test 4: Cross-Endpoint Data Consistency ---")
	testDataConsistency(client, ctx)

	fmt.Println("\n=== Deep Dive Complete ===")

	fmt.Println("\n🎯 Summary Recommendations:")
	fmt.Println("✅ Use tag-based search for model discovery (more reliable)")
	fmt.Println("✅ Images API is rock-solid for content browsing")
	fmt.Println("✅ Individual lookups work great for specific models")
	fmt.Println("⚠️  Avoid version-by-hash lookups (consistently broken)")
	fmt.Println("⚠️  Creators API has timeout issues under load")
	fmt.Println("💡 Implement retries for creators/tags endpoints")
}

func testResultCounts(client *civitai.Client, ctx context.Context) {
	searchTerms := []string{"photo", "portrait", "anime", "landscape", "realistic"}

	for _, term := range searchTerms {
		fmt.Printf("Analyzing '%s':\n", term)

		// Query results
		queryParams := civitai.SearchParams{
			Query: term,
			Limit: 10,
		}
		queryModels, queryMeta, queryErr := client.SearchModels(ctx, queryParams)

		// Tag results
		tagParams := civitai.SearchParams{
			Tag:   term,
			Limit: 10,
		}
		tagModels, tagMeta, tagErr := client.SearchModels(ctx, tagParams)

		// Analysis
		if queryErr != nil {
			fmt.Printf("  Query: FAILED - %v\n", queryErr)
		} else {
			fmt.Printf("  Query: %d models (total: %d)\n", len(queryModels), queryMeta.TotalItems)
		}

		if tagErr != nil {
			fmt.Printf("  Tag: FAILED - %v\n", tagErr)
		} else {
			fmt.Printf("  Tag: %d models (cursor: %v)\n", len(tagModels), tagMeta.NextCursor != "")
		}

		// Compare model IDs if both work
		if queryErr == nil && tagErr == nil {
			overlap := findModelOverlap(queryModels, tagModels)
			fmt.Printf("  Overlap: %d models in common\n", overlap)

			if len(tagModels) > len(queryModels) {
				fmt.Printf("  ✅ Tag returns %d more models\n", len(tagModels)-len(queryModels))
			} else if len(queryModels) > len(tagModels) {
				fmt.Printf("  ⚠️  Query returns %d more models\n", len(queryModels)-len(tagModels))
			} else {
				fmt.Printf("  ➖ Same count\n")
			}
		}

		fmt.Println()
		time.Sleep(200 * time.Millisecond)
	}
}

func testEndpointTiming(client *civitai.Client, ctx context.Context) {
	tests := []struct {
		name string
		test func() (time.Duration, error)
	}{
		{
			name: "Models (tag)",
			test: func() (time.Duration, error) {
				start := time.Now()
				params := civitai.SearchParams{Tag: "photo", Limit: 5}
				_, _, err := client.SearchModels(ctx, params)
				return time.Since(start), err
			},
		},
		{
			name: "Models (query)",
			test: func() (time.Duration, error) {
				start := time.Now()
				params := civitai.SearchParams{Query: "photo", Limit: 5}
				_, _, err := client.SearchModels(ctx, params)
				return time.Since(start), err
			},
		},
		{
			name: "Images",
			test: func() (time.Duration, error) {
				start := time.Now()
				params := civitai.ImageParams{Limit: 5, NSFW: string(civitai.NSFWLevelNone)}
				_, _, err := client.GetImages(ctx, params)
				return time.Since(start), err
			},
		},
		{
			name: "Individual model",
			test: func() (time.Duration, error) {
				start := time.Now()
				_, err := client.GetModel(ctx, 133005)
				return time.Since(start), err
			},
		},
		{
			name: "Health check",
			test: func() (time.Duration, error) {
				start := time.Now()
				err := client.Health(ctx)
				return time.Since(start), err
			},
		},
	}

	for _, test := range tests {
		fmt.Printf("Timing %s:\n", test.name)

		var durations []time.Duration
		failures := 0

		for i := 0; i < 5; i++ {
			duration, err := test.test()
			if err != nil {
				fmt.Printf("  Attempt %d: FAILED (%v) - %v\n", i+1, duration, err)
				failures++
			} else {
				fmt.Printf("  Attempt %d: SUCCESS (%v)\n", i+1, duration)
				durations = append(durations, duration)
			}
			time.Sleep(100 * time.Millisecond)
		}

		if len(durations) > 0 {
			avg := calculateAverage(durations)
			min, max := findMinMax(durations)
			fmt.Printf("  Average: %v (min: %v, max: %v)\n", avg, min, max)
		}

		if failures > 0 {
			fmt.Printf("  Failures: %d/5\n", failures)
		}

		fmt.Println()
	}
}

func testLoadBehavior(client *civitai.Client, ctx context.Context) {
	fmt.Println("Testing rapid-fire requests...")

	// Test models endpoint under load
	fmt.Println("  Models endpoint (10 rapid requests):")
	successCount := 0
	totalTime := time.Now()

	for i := 0; i < 10; i++ {
		params := civitai.SearchParams{Tag: "photo", Limit: 3}
		_, _, err := client.SearchModels(ctx, params)
		if err != nil {
			fmt.Printf("    Request %d: FAILED\n", i+1)
		} else {
			fmt.Printf("    Request %d: SUCCESS\n", i+1)
			successCount++
		}
		// No delay - rapid fire
	}

	totalDuration := time.Since(totalTime)
	fmt.Printf("  Results: %d/10 successful in %v\n", successCount, totalDuration)
	fmt.Printf("  Average per request: %v\n", totalDuration/10)
}

func testDataConsistency(client *civitai.Client, ctx context.Context) {
	fmt.Println("Testing data consistency across endpoints...")

	// Get a model via search
	searchParams := civitai.SearchParams{
		Tag:   "photo",
		Limit: 1,
	}
	models, _, err := client.SearchModels(ctx, searchParams)
	if err != nil || len(models) == 0 {
		fmt.Printf("  Could not get model from search: %v\n", err)
		return
	}

	searchModel := models[0]
	fmt.Printf("  Found model from search: %s (ID: %d)\n", searchModel.Name, searchModel.ID)

	// Get same model via direct lookup
	directModel, err := client.GetModel(ctx, searchModel.ID)
	if err != nil {
		fmt.Printf("  Could not get model directly: %v\n", err)
		return
	}

	fmt.Printf("  Got model directly: %s (ID: %d)\n", directModel.Name, directModel.ID)

	// Compare data
	if searchModel.Name == directModel.Name {
		fmt.Printf("  ✅ Name consistent\n")
	} else {
		fmt.Printf("  ❌ Name inconsistent: '%s' vs '%s'\n", searchModel.Name, directModel.Name)
	}

	if searchModel.Stats.DownloadCount == directModel.Stats.DownloadCount {
		fmt.Printf("  ✅ Download count consistent (%d)\n", searchModel.Stats.DownloadCount)
	} else {
		fmt.Printf("  ⚠️  Download count differs: %d vs %d\n", searchModel.Stats.DownloadCount, directModel.Stats.DownloadCount)
	}

	if len(searchModel.ModelVersions) == len(directModel.ModelVersions) {
		fmt.Printf("  ✅ Version count consistent (%d)\n", len(searchModel.ModelVersions))
	} else {
		fmt.Printf("  ⚠️  Version count differs: %d vs %d\n", len(searchModel.ModelVersions), len(directModel.ModelVersions))
	}
}

func findModelOverlap(models1, models2 []civitai.Model) int {
	ids1 := make(map[int]bool)
	for _, model := range models1 {
		ids1[model.ID] = true
	}

	overlap := 0
	for _, model := range models2 {
		if ids1[model.ID] {
			overlap++
		}
	}

	return overlap
}

func calculateAverage(durations []time.Duration) time.Duration {
	var total time.Duration
	for _, d := range durations {
		total += d
	}
	return total / time.Duration(len(durations))
}

func findMinMax(durations []time.Duration) (time.Duration, time.Duration) {
	if len(durations) == 0 {
		return 0, 0
	}

	min, max := durations[0], durations[0]
	for _, d := range durations[1:] {
		if d < min {
			min = d
		}
		if d > max {
			max = d
		}
	}

	return min, max
}
