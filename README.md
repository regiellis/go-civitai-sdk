# Go CivitAI SDK

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/License-MIT-green)](./LICENSE)
[![Zero Dependencies](https://img.shields.io/badge/Dependencies-Zero-00D084)](https://pkg.go.dev/github.com/regiellis/go-civitai-sdk)

A robust, production-ready Go client for the CivitAI API. Search models, browse images, discover creators, and work with AIR (AI Resource Identifiers) — all with zero dependencies and comprehensive GoDoc documentation.

---

## What's New

- **Robust error handling:** Unified `APIError` type with helpers for rate limits, auth, and server errors. All error handling is now consolidated in `responses.go`.
- **API quirks handled:** Integration tests and SDK logic are resilient to CivitAI API timeouts, unreliable endpoints, and tag vs query search inconsistencies.
- **Comprehensive GoDoc:** All main SDK files include package-level GoDoc and usage examples.
- **Advanced pagination:** Production-ready cursor-based pagination with deduplication and resume support. See `examples/production_pagination_demo/`.
- **AIR support:** Full support for AI Resource Identifiers (AIR) for future-proof model referencing. See `examples/air_integration/`.
- **No duplicate error types:** All error handling uses the canonical `APIError` from `responses.go`.

---

## Install

```bash
go get github.com/regiellis/go-civitai-sdk
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "github.com/regiellis/go-civitai-sdk"
)

func main() {
    // No auth needed for basic usage
    client := civitai.NewClientWithoutAuth()
    
    // Search for models (use Tag for best reliability)
    models, _, err := client.SearchModels(context.Background(), civitai.SearchParams{
        Tag: "anime",
        Limit: 5,
    })
    if err != nil {
        log.Fatal(err)
    }
    for _, model := range models {
        fmt.Printf("%s - %d downloads\n", model.Name, model.Stats.DownloadCount)
    }
}
```

## Pagination Best Practices

- **Use cursor-based pagination** for reliable, duplicate-free results. See `examples/production_pagination_demo/` for a full demo.
- **Deduplicate results** when paginating, as the API may return overlapping items.
- **Resume from cursor** for robust, resumable data fetching.

## AIR (AI Resource Identifier) Support

- **Parse, create, and use AIRS** for models and versions.
- **Future-proof**: AIRs enable standardized resource identification across AI platforms.
- See `examples/air_integration/air_integration.go` for advanced usage.

## Error Handling

All error handling uses the canonical `APIError` from `responses.go`:

```go
models, _, err := client.SearchModels(ctx, params)
if err != nil {
    if apiErr, ok := err.(*civitai.APIError); ok {
        if apiErr.IsRateLimitError() {
            // Handle rate limit
        }
        fmt.Printf("API Error: %s (Status: %d)\n", apiErr.Message, apiErr.StatusCode)
    } else {
        fmt.Printf("Network Error: %v\n", err)
    }
}
```

## API Quirks & Best Practices

- **Prefer `Tag` over `Query`** for model search — more reliable and less likely to fail.
- **Handle timeouts and unreliable endpoints**: The SDK and tests skip or retry known-problem endpoints (creators, tags, health).
- **All public endpoints** work without authentication. Use your API token for private data.

## Comprehensive GoDoc & Examples

- All main files include GoDoc package comments and usage examples.
- See the [examples/](./examples/) directory for:
  - `hello_world.go` — Simple validation
  - `basic_usage.go` — Core features
  - `model_search.go` — Advanced filtering
  - `image_browsing.go` — Image discovery
  - `air_integration.go` — AIR support
  - `production_pagination_demo.go` — Production-ready pagination

## Advanced Configuration

```go
client := civitai.NewClient("token",
    civitai.WithTimeout(60*time.Second),
    civitai.WithRetryConfig(3, time.Second, 30*time.Second),
    civitai.WithConnectionPooling(10, 2),
)
```

## Testing & Verification

```bash
# Run all tests (including integration, skips unreliable endpoints gracefully)
go test -v ./...

# Run production pagination demo
go run examples/production_pagination_demo/production_pagination_demo.go

# Run AIR integration example
go run examples/air_integration/air_integration.go
```

## Contributing

1. Fork the repo
2. Make your changes
3. Run `go test`
4. Submit a PR

## License

MIT License - see [LICENSE](./LICENSE) for details.

## Links

- [CivitAI Platform](https://civitai.com/)
- [API Documentation](https://developer.civitai.com/docs/api/public-rest)
- [Get API Token](https://civitai.com/user/account)