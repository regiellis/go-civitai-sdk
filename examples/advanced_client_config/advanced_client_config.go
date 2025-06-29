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

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	civitai "github.com/regiellis/go-civitai-sdk"
)

func main() {
	// Example 1: Basic client with retry configuration
	fmt.Println("=== Example 1: Client with Retry Logic ===")
	clientWithRetry := civitai.NewClientWithoutAuth(
		civitai.WithRetryConfig(
			5,                    // Max 5 retries
			500*time.Millisecond, // Base delay of 500ms
			10*time.Second,       // Max delay of 10s
		),
	)

	ctx := context.Background()
	models, metadata, err := clientWithRetry.SearchModels(ctx, civitai.SearchParams{
		Limit: 5,
		Query: "girl",
	})

	if err != nil {
		log.Printf("Search failed even with retries: %v", err)
	} else {
		fmt.Printf("Found %d models (total: %d)\n", len(models), metadata.TotalItems)
		for _, model := range models {
			fmt.Printf("- %s (ID: %d)\n", model.Name, model.ID)
		}
	}

	// Example 2: Client with connection pooling and compression
	fmt.Println("\n=== Example 2: Client with Connection Pooling ===")
	clientWithPooling := civitai.NewClientWithoutAuth(
		civitai.WithConnectionPooling(
			50, // Max 50 idle connections total
			10, // Max 10 idle connections per host
		),
	)

	// Make multiple concurrent requests to demonstrate pooling
	const numConcurrent = 3
	results := make(chan string, numConcurrent)

	for i := 0; i < numConcurrent; i++ {
		go func(id int) {
			models, _, err := clientWithPooling.SearchModels(ctx, civitai.SearchParams{
				Limit: 2,
				Types: []civitai.ModelType{civitai.ModelTypeCheckpoint},
			})

			if err != nil {
				results <- fmt.Sprintf("Request %d failed: %v", id, err)
			} else {
				results <- fmt.Sprintf("Request %d: Found %d models", id, len(models))
			}
		}(i + 1)
	}

	// Collect results
	for i := 0; i < numConcurrent; i++ {
		fmt.Println(<-results)
	}

	// Example 3: Production-ready client with all optimizations
	fmt.Println("\n=== Example 3: Production-Ready Client ===")
	productionClient := civitai.NewClient(
		"your-api-token-here", // Replace with actual token
		civitai.WithRetryConfig(3, 1*time.Second, 30*time.Second),
		civitai.WithConnectionPooling(100, 20),
		civitai.WithMaxResponseSize(50*1024*1024), // 50MB limit
		civitai.WithTimeout(60*time.Second),
	)

	// Example of using secure token methods
	if productionClient.HasAPIToken() {
		fmt.Printf("Client authenticated with token: %s\n", productionClient.GetMaskedAPIToken())
	} else {
		fmt.Println("Client not authenticated")
	}

	// Example 4: Error handling with retries
	fmt.Println("\n=== Example 4: Error Handling with Retries ===")

	// This will demonstrate retry logic (will likely fail in examples)
	models, _, err = productionClient.SearchModels(ctx, civitai.SearchParams{
		Limit: 10,
		Query: "landscape",
	})

	if err != nil {
		fmt.Printf("Request failed after retries: %v\n", err)
		// In production, you might want to implement circuit breaker logic here
	} else {
		fmt.Printf("Successfully found %d models\n", len(models))
	}

	// Example 5: Performance monitoring
	fmt.Println("\n=== Example 5: Performance Monitoring ===")

	start := time.Now()
	_, _, err = clientWithPooling.SearchModels(ctx, civitai.SearchParams{
		Limit: 1,
	})
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("Performance test failed: %v\n", err)
	} else {
		fmt.Printf("Request completed in %v\n", duration)
	}

	// Example 6: Custom HTTP client with advanced features
	fmt.Println("\n=== Example 6: Custom HTTP Client Configuration ===")

	customClient := civitai.NewClientWithoutAuth(
		civitai.WithRetryConfig(2, 200*time.Millisecond, 5*time.Second),
		civitai.WithConnectionPooling(25, 5),
		civitai.WithMaxResponseSize(10*1024*1024), // 10MB
		civitai.WithUserAgent("MyApp/1.0 (custom-client)"),
	)

	// Test with a simple health check equivalent
	_, _, err = customClient.SearchModels(ctx, civitai.SearchParams{Limit: 1})
	if err != nil {
		fmt.Printf("Health check failed: %v\n", err)
	} else {
		fmt.Println("API is responsive")
	}

	fmt.Println("\n=== Examples Complete ===")
	fmt.Println("Key benefits of advanced configuration:")
	fmt.Println("- Retry logic improves reliability")
	fmt.Println("- Connection pooling improves performance")
	fmt.Println("- Response size limits prevent DoS attacks")
	fmt.Println("- Secure token handling prevents credential leaks")
	fmt.Println("- Compression reduces bandwidth usage")
}
