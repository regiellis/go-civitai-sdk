# Go CivitAI SDK 🎨

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go&logoColor=white)
![Zero Dependencies](https://img.shields.io/badge/Dependencies-Zero-00D084?logo=go&logoColor=white)
![CivitAI API](https://img.shields.io/badge/CivitAI-API%20Ready-FF6B35?logo=artifact&logoColor=white)
![License](https://img.shields.io/badge/License-Restricted-red?logo=law&logoColor=white)

---

>[!TIP]
> **Zero Dependencies Magic!** ✨ This SDK uses **only Go's standard library** - no external dependencies, no version conflicts, no headaches. Just pure Go goodness that compiles to a single binary and works everywhere Go works!

>[!IMPORTANT]
> **Ready for Production!** 🚀 This SDK provides complete coverage of the CivitAI REST API with proper error handling, context support, and type safety. It's battle-tested and ready for both hobby projects and production applications.
>
> **Perfect for AI Developers:** Whether you're building model discovery tools, content management systems, or AI workflow automation - this SDK has you covered with a clean, idiomatic Go API.

>[!NOTE]
> **Authentication Optional:** Most endpoints work without authentication! Get started immediately with public data, then add your API token later for authenticated features like favorites and private models.

---

## ✨ Why Go CivitAI SDK?

- 🚀 **Zero Dependencies:** Only uses Go standard library - no external packages, no conflicts!
- 🎯 **Complete API Coverage:** All CivitAI REST API endpoints supported with full type safety
- ⚡ **Lightning Fast:** Efficient HTTP client with connection reuse and smart JSON parsing  
- 🛡️ **Battle-Tested:** Comprehensive test suite with real API integration tests
- 🔧 **Developer Friendly:** Context support, timeouts, retries, and structured errors
- 📖 **Extensive Examples:** Working code for every use case - copy, paste, ship!
- 🌍 **Cross-Platform:** Works on Linux, macOS, Windows - anywhere Go compiles
- 🎨 **Type Safe:** Rich Go types for all API responses with proper validation

### Quick Demo

![SDK Demo](https://img.shields.io/badge/Demo-Coming%20Soon-blue)

## 🚀 Quick Start

### Install the SDK

```bash
go get github.com/regiellis/go-civitai-sdk
```

### Hello, CivitAI! 👋

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/regiellis/go-civitai-sdk"
)

func main() {
    // Create client (no auth needed for this!)
    client := civitai.NewClientWithoutAuth()
    
    // Search for anime models
    models, _, err := client.SearchModels(context.Background(), civitai.SearchParams{
        Query: "anime",
        Types: []civitai.ModelType{civitai.ModelTypeCheckpoint},
        Limit: 5,
    })
    
    if err != nil {
        log.Fatal("Oops:", err)
    }
    
    fmt.Printf("🎉 Found %d amazing anime models!\n", len(models))
    for i, model := range models {
        fmt.Printf("%d. %s (%d downloads) 🔥\n", 
            i+1, model.Name, model.Stats.DownloadCount)
    }
}
```

**Run it:**
```bash
go run main.go
# 🎉 Found 5 amazing anime models!
# 1. Amazing Anime Model (50000 downloads) 🔥
# 2. Cute Character LoRA (30000 downloads) 🔥
# ...
```

## 🔑 Authentication (Optional!)

>[!TIP]
> **Start Without Auth!** Most CivitAI endpoints work without authentication. Jump right in and add your API token later when you need authenticated features!

Get your API token from [CivitAI User Account Settings](https://civitai.com/user/account) when you're ready.

```go
// 🌟 For public endpoints (start here!)
client := civitai.NewClientWithoutAuth()

// 🔐 For authenticated requests (when you need favorites, etc.)
client := civitai.NewClient("your-api-token-here")

// ⚙️ Power user configuration
client := civitai.NewClient("your-token",
    civitai.WithTimeout(60*time.Second),           // Custom timeout
    civitai.WithUserAgent("MyApp/1.0.0"),         // Custom user agent  
    civitai.WithBaseURL("https://custom-api.com"), // Custom API URL
)
```

## 🎯 More Examples

>[!NOTE]
> **Learning by Example:** Check out our comprehensive examples directory! Each file is a complete, runnable program that demonstrates different aspects of the SDK.

### 🔍 Advanced Model Search

```go
// Find high-quality realistic models
params := civitai.SearchParams{
    Query:                 "realistic portrait",
    Types:                 []civitai.ModelType{civitai.ModelTypeCheckpoint},
    Sort:                  civitai.SortHighestRated,
    Rating:                4, // Minimum 4 stars
    AllowCommercialUse:    []string{string(civitai.CommercialUseSell)},
    NSFW:                  &[]bool{false}[0], // Safe content only
    SupportsGeneration:    &[]bool{true}[0],  // Works with generation
    Limit:                 20,
}

models, metadata, err := client.SearchModels(ctx, params)
```

### 🖼️ Browse Beautiful Images

```go
// Get the latest safe, high-quality images
images, _, err := client.GetImages(ctx, civitai.ImageParams{
    Sort:  string(civitai.ImageSortNewest),
    NSFW:  string(civitai.NSFWLevelNone),
    Limit: 50,
})

for _, image := range images {
    fmt.Printf("✨ %dx%d by %s\n", image.Width, image.Height, image.Username)
    if prompt, ok := image.Meta["prompt"].(string); ok {
        fmt.Printf("   Prompt: %s\n", prompt)
    }
}
```

### 👥 Discover Amazing Creators

```go
// Find top creators in your favorite style
creators, _, err := client.GetCreators(ctx, civitai.CreatorParams{
    Query: "anime",
    Limit: 10,
})

for i, creator := range creators {
    fmt.Printf("%d. %s (%d models) 🎨\n", 
        i+1, creator.Username, creator.ModelCount)
}
```

### 🏷️ Explore Tags & Trends

```go
// See what's trending
tags, _, err := client.GetTags(ctx, civitai.TagParams{
    Query: "style",
    Limit: 20,
})

fmt.Println("🔥 Trending Style Tags:")
for _, tag := range tags {
    fmt.Printf("  #%s (%d models)\n", tag.Name, tag.ModelCount)
}
```

## �️ Available Endpoints

| Feature | Endpoint | Authentication | Description |
|---------|----------|----------------|-------------|
| 🔍 **Model Search** | `SearchModels()` | Optional | Find models by query, type, rating, etc. |
| 📋 **Model Details** | `GetModel()` | Optional | Get complete model information |
| 🔄 **Model Versions** | `GetModelVersion()` | Optional | Access specific model versions |
| 🖼️ **Image Gallery** | `GetImages()` | Optional | Browse community images |
| 👥 **Creator Discovery** | `GetCreators()` | Optional | Find talented creators |
| 🏷️ **Tag Explorer** | `GetTags()` | Optional | Explore tags and categories |
| ⚡ **Health Check** | `Health()` | None | Test API connectivity |

### 🎨 Model Types Supported

```go
civitai.ModelTypeCheckpoint    // 🎯 Full Stable Diffusion models
civitai.ModelTypeLORA          // 🎨 LoRA adaptations  
civitai.ModelTypeEmbedding     // 📝 Textual inversions
civitai.ModelTypeHypernetwork  // 🧠 Hypernetworks
civitai.ModelTypeControlNet    // 🎮 ControlNet models
civitai.ModelTypePose          // 🤸 Pose models
// ... and more!
```

### 📊 Smart Sorting Options

```go
civitai.SortHighestRated  // ⭐ Best rated content
civitai.SortMostDownload  // 📈 Most popular downloads  
civitai.SortMostLiked     // ❤️ Community favorites
civitai.SortNewest        // 🆕 Latest uploads
```

## � Production Ready Features

### 🛡️ Error Handling That Actually Helps

```go
models, _, err := client.SearchModels(ctx, params)
if err != nil {
    // Check if it's a structured API error
    if apiErr, ok := err.(*civitai.APIError); ok {
        fmt.Printf("API Error [%s]: %s 🚨\n", apiErr.Code, apiErr.Message)
        // Handle specific error codes
    } else {
        fmt.Printf("Network/Request Error: %v 📡\n", err)
        // Handle network issues
    }
}
```

### 📄 Smart Pagination

```go
models, metadata, err := client.SearchModels(ctx, params)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("📊 Page %d of %d (%d total models)\n", 
    metadata.CurrentPage, metadata.TotalPages, metadata.TotalItems)

// For images, use cursor-based pagination
if metadata.NextPage != "" {
    fmt.Printf("🔗 Next page: %s\n", metadata.NextPage)
}
```

### ⚡ Context & Timeout Support

```go
// Set timeouts for operations
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

// Graceful cancellation
ctx, cancel := context.WithCancel(context.Background())
go func() {
    time.Sleep(5*time.Second)
    cancel() // Cancel after 5 seconds
}()

models, _, err := client.SearchModels(ctx, params)
// Will respect cancellation and timeouts!
```

### 🔄 Retry Logic & Resilience

>[!TIP]
> **Built-in Resilience:** The SDK automatically handles connection pooling, HTTP/2 support, and proper resource cleanup. No need to worry about connection leaks or performance issues!

```go
// SDK handles retries and connection reuse automatically
client := civitai.NewClient("token", 
    civitai.WithTimeout(60*time.Second), // Overall timeout
)

// All requests benefit from:
// ✅ HTTP/2 connection reuse
// ✅ Automatic decompression  
// ✅ Proper resource cleanup
// ✅ Context cancellation support
```

## 🧪 Testing & Development

>[!NOTE]
> **Comprehensive Test Suite:** We've got unit tests, integration tests, and example validation. Every feature is tested with both mocked and real API responses!

### Run the Tests

```bash
# Quick unit tests (no network required)
go test -v -short

# Full integration tests (requires network)
go test -v

# Benchmark performance
go test -bench=. -benchmem

# Test with race detection
go test -race -v
```

### Try the Examples

```bash
# Run any example
go run examples/basic_usage.go
go run examples/model_search.go
go run examples/image_browsing.go
go run examples/creator_discovery.go

# Or build them
go build examples/basic_usage.go
./basic_usage
```

### Test Your Integration

```bash
# Use our validation tool
go run cmd/test/main.go
# 🤖 CivitAI SDK Test
# ==================
# ⚡ Testing API health... ✅ API is healthy!
# 🔍 Searching for models... ✅ Found 10 models
# 🖼️ Getting images... ✅ Found 20 safe images
# 👥 Getting creators... ✅ Found 15 creators
# 🏷️ Getting tags... ✅ Found 25 style tags
# 🎉 All SDK tests completed successfully!
```

## 📚 Complete Examples

| Example | Description | Features Shown |
|---------|-------------|----------------|
| [**basic_usage.go**](./examples/basic_usage.go) | 🌟 Complete SDK tour | All major endpoints, error handling |
| [**model_search.go**](./examples/model_search.go) | 🔍 Advanced searching | Filters, pagination, type safety |
| [**image_browsing.go**](./examples/image_browsing.go) | 🖼️ Image discovery | Metadata parsing, safe browsing |
| [**creator_discovery.go**](./examples/creator_discovery.go) | 👥 Creator exploration | Cross-referencing, tag analysis |

## 🤝 Contributing

>[!TIP]
> **Want to Contribute?** We'd love your help! This SDK is part of a larger ecosystem and welcomes improvements, bug fixes, and new features.

### Development Setup

```bash
# Clone and setup
git clone https://github.com/regiellis/go-civitai-sdk.git
cd go-civitai-sdk

# Verify everything works
./verify-structure.sh

# Run tests
go test -v

# Make your changes, then:
go fmt ./...
go vet ./...
go test -v
```

### What We'd Love Help With

- 🐛 **Bug Reports:** Found an issue? Please report it!
- ✨ **New Features:** Missing an endpoint? Want better docs?
- 🧪 **More Tests:** Help us achieve even better coverage
- 📖 **Documentation:** Improve examples and guides
- 🎨 **Type Improvements:** Better Go types for API responses

## 📄 License & Legal

>[!IMPORTANT]
> **Restricted Use License:** This SDK is licensed under a **Restricted Use License - Non-Commercial Only**. Please read the full license terms before using in commercial projects.

Licensed under the Restricted Use License - Non-Commercial Only.
See [LICENSE](./LICENSE) for complete details.

## 🔗 Links & Resources

- 🌐 **[CivitAI Platform](https://civitai.com/)** - The amazing AI art community
- 📚 **[CivitAI API Docs](https://developer.civitai.com/docs/api/public-rest)** - Official API documentation  
- 🔑 **[Get API Token](https://civitai.com/user/account)** - Your account settings
- 🐙 **[GitHub Repository](https://github.com/regiellis/go-civitai-sdk)** - Source code and issues
- 💬 **[Discussions](https://github.com/regiellis/go-civitai-sdk/discussions)** - Questions and community

## 🎉 What's Next?

>[!TIP]
> **Ready to Build Something Amazing?** 
> 
> - 🤖 **AI Tools:** Build model discovery and management tools
> - 🎨 **Content Apps:** Create galleries and showcases  
> - 🔧 **Workflow Automation:** Integrate with your AI pipelines
> - 📊 **Analytics:** Track trends and discover new content
> - 🎮 **Games & Apps:** Add AI art integration to your projects

**Start exploring the amazing world of AI art with Go! 🚀**

---

**🎨 Original work by [Regi Ellis](https://github.com/regiellis)** - Built with ❤️ for the AI art community

![Go Gopher](https://img.shields.io/badge/Made%20with-Go%20Gopher%20Power-00ADD8?logo=go&logoColor=white) ![AI Art](https://img.shields.io/badge/Powered%20by-AI%20Art%20Community-FF6B35)
