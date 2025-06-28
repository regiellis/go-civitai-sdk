# Go CivitAI SDK - Project Structure

This document outlines the complete structure of the Go CivitAI SDK library.

## 📁 Directory Structure

```
go-civitai-sdk/
├── README.md                 # Main documentation and usage guide
├── LICENSE                   # Restricted Use License - Non-Commercial Only
├── .gitignore               # Git ignore rules for Go library
├── go.mod                   # Go module definition
├── verify-structure.sh      # Library structure verification script
├── STRUCTURE.md             # This documentation file
│
├── 📚 Core Library Files (package civitai)
├── client.go                # Main SDK client implementation
├── types.go                 # Type definitions and constants
├── exceptions.go            # Error handling and custom exceptions
├── models.go               # Model-related API methods
├── model-versions.go       # Model version API methods
├── images.go               # Image API methods
├── creators.go             # Creator API methods
├── tags.go                 # Tag API methods
├── responses.go            # API response structures
├── utils.go                # Utility functions
│
├── 🧪 Testing (package civitai - same package as library)
├── client_test.go          # Unit tests for client functionality
├── types_test.go           # Unit tests for types and validation
├── integration_test.go     # Integration tests (real API calls)
│
├── 📖 Examples
├── examples/
│   ├── basic_usage.go       # Complete SDK demonstration
│   ├── model_search.go      # Advanced model searching examples
│   ├── image_browsing.go    # Image browsing and discovery
│   └── creator_discovery.go # Creator and tag exploration
│
├── 🔧 Command Line Tools
└── cmd/
    └── test/
        └── main.go          # SDK test and validation program
```

## 📦 Package Structure

### Main Package: `civitai`

The SDK uses a single main package `civitai` for all public APIs:

```go
import "github.com/regiellis/go-civitai-sdk"

client := civitai.NewClientWithoutAuth()
```

### Core Components

1. **Client (`client.go`)**
   - HTTP client configuration
   - Authentication handling
   - Request/response processing
   - Error handling

2. **Types (`types.go`)**
   - API request/response structures
   - Enums and constants
   - Validation logic

3. **API Methods**
   - `models.go` - Model search and retrieval
   - `model-versions.go` - Model version operations
   - `images.go` - Image browsing and search
   - `creators.go` - Creator discovery
   - `tags.go` - Tag management

## 🚀 Usage Patterns

### Basic Usage
```go
client := civitai.NewClientWithoutAuth()
models, _, err := client.SearchModels(ctx, civitai.SearchParams{
    Query: "anime",
    Limit: 10,
})
```

### Authenticated Usage
```go
client := civitai.NewClient("your-api-token")
// Access to additional endpoints requiring authentication
```

### Configuration
```go
client := civitai.NewClient("token",
    civitai.WithTimeout(60*time.Second),
    civitai.WithUserAgent("my-app/1.0.0"),
)
```

## 🧪 Testing Strategy

### Unit Tests
- Mock HTTP responses
- Test all public APIs
- Validate type marshaling/unmarshaling
- Error handling scenarios

### Integration Tests
- Real API calls (when network available)
- End-to-end functionality
- Performance benchmarks

### Examples as Tests
- All examples must compile
- Examples serve as documentation
- Demonstrate real-world usage

## 📋 Development Guidelines

### Code Style
- Follow Go best practices
- Use `gofmt` for formatting
- Clear and concise naming
- Comprehensive documentation

### Error Handling
- Structured error types
- Context-aware errors
- Graceful degradation

### Performance
- Efficient JSON parsing
- HTTP connection reuse
- Proper context handling
- Memory-conscious design

## 🔗 Dependencies

The SDK uses **only Go standard library**:
- `net/http` - HTTP client
- `encoding/json` - JSON parsing
- `context` - Request contexts
- `time` - Timeouts and timestamps
- `fmt`, `strings`, `url` - Utilities

No external dependencies ensures:
- ✅ Minimal footprint
- ✅ Easy integration
- ✅ Reduced security surface
- ✅ Long-term stability

## 📚 Documentation

1. **README.md** - Main documentation with examples
2. **GoDoc** - Inline documentation for all public APIs
3. **Examples** - Working code demonstrations
4. **Integration Tests** - Real usage scenarios

## 🚀 Distribution

Ready for:
- ✅ Go modules (`go get`)
- ✅ GitHub repository hosting
- ✅ Go package registry
- ✅ Private repositories
- ✅ Vendoring

## 🔄 Version Management

- Semantic versioning (v1.0.0+)
- Go module compatibility
- Backward compatibility commitment
- Clear upgrade paths

---

**Note**: This SDK follows Go library best practices and is ready for standalone repository hosting.
