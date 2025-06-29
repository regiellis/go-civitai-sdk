# Civitai API Tester

⚠️ **SECURITY WARNING: DEVELOPMENT/TESTING TOOL ONLY** ⚠️

A comprehensive testing tool for the Civitai API using the go-civitai-sdk. This standalone application tests all major API endpoints and displays results in a beautiful status page-style web interface.

**🚨 IMPORTANT: This tool is NOT intended for production use. Do not expose this server to public networks or untrusted environments.**

## Features

- **Comprehensive API Testing**: Tests all major Civitai API endpoints
- **Real-time Web Dashboard**: Beautiful status page UI with live updates
- **Cross-platform**: Runs on Windows, macOS, and Linux
- **Configurable**: Customizable test parameters and API settings
- **Zero Dependencies**: Uses only Go standard library

## Quick Start

### Download Pre-built Binary

Download the appropriate binary for your platform from the releases section:

- **Windows**: `civitai-tester-windows-amd64.exe`
- **macOS**: `civitai-tester-darwin-amd64` (Intel) or `civitai-tester-darwin-arm64` (Apple Silicon)
- **Linux**: `civitai-tester-linux-amd64` or `civitai-tester-linux-arm64`

### Run the Tester

```bash
# Make executable (macOS/Linux)
chmod +x civitai-tester-*

# Run the tester (server only)
./civitai-tester-linux-amd64

# Or run with automatic test execution
./civitai-tester-linux-amd64 --run

# Or on Windows
civitai-tester-windows-amd64.exe
civitai-tester-windows-amd64.exe --run
```

The application will:
1. Start the web server on port 9999
2. Wait for manual test execution (click "Start Tests" or use `--run` flag)
3. Display results in real-time at http://localhost:9999

### Command Line Options

- `--run`: Automatically run tests on startup
- `--help`: Show help information

## Configuration

### Environment Variables

```bash
export CIVITAI_API_KEY="your_api_key_here"    # Optional: Your Civitai API key
export TESTER_PORT="9999"                     # Optional: Web server port (default: 9999)
export TEST_TIMEOUT="30"                      # Optional: Test timeout in seconds
```

### Configuration File

Create a `config.json` file in the same directory as the executable:

```json
{
  "api_key": "your_api_key_here",
  "server_port": 9999,
  "test_timeout_seconds": 30,
  "test_limits": {
    "models_limit": 5,
    "images_limit": 5,
    "creators_limit": 5,
    "tags_limit": 10
  },
  "custom_tests": {
    "skip_rate_limit_test": false,
    "only_tests": [],
    "skip_tests": ["Test Rate Limiting"]
  }
}
```

Copy `config.example.json` to `config.json` and customize as needed.

## Tests Performed

The tester runs the following comprehensive tests:

1. **API Health Check** - Verifies API is responding
2. **Get Models** - Tests model listing endpoint
3. **Get Model Details** - Tests individual model retrieval
4. **Get Model Versions** - Tests model version listing
5. **Get Images** - Tests image browsing functionality
6. **Get Creators** - Tests creator discovery
7. **Get Tags** - Tests tag listing
8. **Search Models by Query** - Tests model search functionality
9. **Test Pagination** - Verifies pagination works correctly
10. **Test Rate Limiting** - Tests rate limit handling

## Web Dashboard

The web dashboard provides:

- **Modern Dark Theme**: Professional dark UI with TailwindCSS styling
- **On-Demand Testing**: Click "Start Tests" to begin API testing
- **Real-time Updates**: Test results update every 2 seconds during tests, 5 minutes otherwise
- **Progress Indicators**: Animated progress bars for running tests
- **Manual Refresh**: Click the refresh button to re-run tests immediately
- **Status Overview**: Summary of operational, failed, and checking services
- **Detailed Dropdowns**: Click any service to see what was tested
- **Smooth Animations**: AnimeJS-powered transitions and effects
- **Interactive UI**: AlpineJS for reactive components
- **Responsive**: Works on desktop and mobile devices
- **Debug Information**: Built-in debugging and error tracking

## Project Structure

```
cmd/civitai-tester/
├── main.go              # Main application server
├── config.go            # Configuration management
├── static/              # Web assets (served at runtime)
│   ├── index.html       # Main web interface
│   ├── css/
│   │   └── style.css    # Custom dark theme styles
│   └── js/
│       └── app.js       # Alpine.js application logic
├── config.example.json  # Example configuration
├── build.sh            # Cross-platform build script
├── build.bat           # Windows build script
├── run.sh              # Quick run script
└── README.md           # This file
```

### Technologies Used

- **Backend**: Go with standard library HTTP server
- **Frontend**: Modern web stack
  - **TailwindCSS**: Utility-first CSS framework
  - **Alpine.js**: Lightweight reactive framework
  - **AnimeJS**: Animation library for smooth transitions
- **Real-time**: Polling-based updates with enhanced error handling

## Building from Source

### Prerequisites

- Go 1.21 or later

### Build for Current Platform

```bash
cd cmd/civitai-tester
go build -o civitai-tester .
```

### Cross-platform Build

Use the provided build scripts:

```bash
# Unix/Linux/macOS
./build.sh

# Windows
build.bat
```

This creates binaries for all supported platforms in the `builds/` directory.

## API Key Usage

While an API key is not required for most endpoints, providing one allows:

- Higher rate limits
- Access to authenticated endpoints
- More comprehensive testing

Get your free API key from [Civitai](https://civitai.com/user/account).

## 🔒 Security & Production Warnings

### ⚠️ **CRITICAL SECURITY NOTICE**

**This application is a development and testing tool ONLY. It is NOT designed for production environments.**

### 🚫 **DO NOT:**
- Expose this server to the public internet
- Run this on production servers
- Use this in untrusted network environments
- Deploy this to cloud instances accessible from the web
- Share the server URL publicly
- Run this with elevated privileges in production systems

### ✅ **SAFE USAGE:**
- Run locally on your development machine (localhost only)
- Use for API testing and development purposes
- Keep it behind firewalls and private networks
- Use for internal team testing (secure networks only)
- Monitor API functionality during development

### 🛡️ **Security Considerations:**
- The server binds to all interfaces (0.0.0.0) by default
- No authentication or authorization mechanisms
- No rate limiting on the web interface
- No input validation for malicious requests
- No HTTPS/TLS encryption
- No protection against common web attacks

### 📝 **Production Alternative:**
For production API monitoring, consider:
- Professional monitoring services (Datadog, New Relic, etc.)
- Internal monitoring systems with proper security
- Custom monitoring solutions with authentication
- Cloud provider monitoring tools

## Troubleshooting

### Common Issues

**Port already in use**: Change the port using `TESTER_PORT` environment variable or `config.json` (default is 9999)

**API rate limiting**: The tester handles rate limits gracefully. You can skip rate limit tests by setting `skip_rate_limit_test: true`

**Network issues**: Check your internet connection and firewall settings

### Verbose Output

The application provides detailed console output showing:
- Configuration settings
- Test progress
- Individual test results
- Web server status

## License

This tool is part of the go-civitai-sdk project and follows the same license terms.

## Contributing

Issues and pull requests are welcome. Please ensure all tests pass before submitting PRs.