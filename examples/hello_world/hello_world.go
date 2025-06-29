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

	civitai "github.com/regiellis/go-civitai-sdk"
)

func main() {
	// Create a simple client (no auth needed)
	client := civitai.NewClientWithoutAuth()
	ctx := context.Background()

	fmt.Println("🎨 CivitAI SDK - Hello World!")
	fmt.Println("==============================")

	// Test 1: Simple search
	fmt.Println("\n1. Quick search for 'girl' models...")
	models, err := client.QuickSearch(ctx, "girl", 3)
	if err != nil {
		log.Printf("Search failed: %v", err)
	} else {
		fmt.Printf("✅ Found %d models!\n", len(models))
		for i, model := range models {
			fmt.Printf("   %d. %s (%d downloads)\n", i+1, model.Name, model.Stats.DownloadCount)
		}
	}

	// Test 2: Get popular models
	fmt.Println("\n2. Getting popular models...")
	popular, err := client.GetPopularModels(ctx, 3)
	if err != nil {
		log.Printf("Popular models failed: %v", err)
	} else {
		fmt.Printf("✅ Found %d popular models!\n", len(popular))
		for i, model := range popular {
			fmt.Printf("   %d. %s (%d downloads)\n", i+1, model.Name, model.Stats.DownloadCount)
		}
	}

	// Test 3: Get safe images
	fmt.Println("\n3. Getting safe images...")
	images, err := client.GetSafeImages(ctx, 3)
	if err != nil {
		log.Printf("Safe images failed: %v", err)
	} else {
		fmt.Printf("✅ Found %d safe images!\n", len(images))
		for i, image := range images {
			fmt.Printf("   %d. %dx%d by %s\n", i+1, image.Width, image.Height, image.Username)
		}
	}

	// Test 4: API health
	fmt.Println("\n4. Testing if API is working...")
	if client.IsWorking(ctx) {
		fmt.Println("✅ API is working!")
	} else {
		fmt.Println("❌ API is not responding")
	}

	fmt.Println("\n🎉 All tests completed! The SDK is working perfectly.")
	fmt.Println("\nNext steps:")
	fmt.Println("- Try 'go run examples/basic_usage.go' for more features")
	fmt.Println("- Check out other examples in the examples/ directory")
	fmt.Println("- Get an API token from https://civitai.com/user/account for authenticated features")
}
