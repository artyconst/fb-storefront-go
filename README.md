# Fleetbase Storefront Go SDK

A powerful Go SDK for building custom shopping experiences with Fleetbase's headless commerce platform.

[![Go Reference](https://pkg.go.dev/badge/github.com/artyconst/fb-storefront-go.svg)](https://pkg.go.dev/github.com/artyconst/fb-storefront-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Installation

```bash
go get github.com/artyconst/fb-storefront-go
```

### Prerequisites

Get your Storefront API key from [console.fleetbase.io](https://console.fleetbase.io) or use a self-hosted Fleetbase instance.

## Quick Start

### Basic Initialization

Initialize the SDK with your Storefront API key:

```go
// Initialize SDK with basic configuration (replace with your actual API key)
sf, err := storefront.NewStorefront("sk_test_your_storefront_key_here",
    // Uncomment and modify these for advanced configuration:
    // storefront.WithAPIHost("https://api.custom-domain.com"),  // Custom API host URL
    // storefront.WithTimeout(60),                              // HTTP timeout in seconds
    // storefront.WithLogLevel(sf.LevelDebug),                  // Set log level (Error, Warn, Info, Debug)
    // storefront.WithDebugMode(),                              // Enable debug logging for API calls
)
if err != nil {
    log.Fatal(err)
}
```

### Example Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/artyconst/fb-storefront-go"
    "github.com/artyconst/fb-storefront-go/pkg/product"
)

func main() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    products, err := sf.Products().List(context.Background(), &product.ListOptions{Limit: 20})
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Found %d products\n", len(products))
}
```

### Advanced Examples

See the example scripts in the `examples` directory for complete demonstrations:

- **[`examples/demo/main.go`](./examples/demo/main.go)** - Complete demonstration of all SDK features including cart, orders, and checkout operations
- **[`examples/customer_auth/main.go`](./examples/customer_auth/main.go)** - Customer authentication flows

Note: These examples require setting the `STOREFRONT_KEY` environment variable with a valid API key. Unlike tests which use mocked responses, these examples make real API calls to your Fleetbase instance.

## Configuration Options

| Option | Description | Example |
|--------|-------------|---------|
| `WithAPIHost(host)` | Custom API host URL | `WithAPIHost("https://api.example.com")` |
| `WithTimeout(seconds)` | HTTP timeout in seconds | `WithTimeout(60)` |
| `WithLogLevel(level)` | Log level (Error, Warn, Info, Debug) | `WithLogLevel(sf.LevelDebug)` |
| `WithLoggerOutput(w)` | Custom logger output writer | `WithLoggerOutput(os.Stdout)` |
| `WithDebugMode()` | Enable debug logging for API calls | `WithDebugMode()` |

## Core Concepts

The SDK uses a service-based architecture. Access services via method calls:

```go
sf := storefront.NewStorefront(YOUR_API_KEY)

// Available services
sf.Products()    // Product catalog
sf.Cart()        // Shopping cart
sf.Checkout()    // Checkout sessions
sf.Customers()   // Customer accounts
sf.Orders()      // Order management
sf.Categories()  // Product categories
```

All operations accept `context.Context` for cancellation and timeouts.

## Service Documentation

- [Products Guide](./user-guide-products.md) - Browse, search, manage products
- [Categories Guide](./user-guide-categories.md) - Organize hierarchies  
- [Cart Guide](./user-guide-cart.md) - Shopping cart operations
- [Checkout Guide](./user-guide-checkout.md) - Checkout and payment flow
- [Customers Guide](./user-guide-customers.md) - Customer account management
- [Orders Guide](./user-guide-orders.md) - Order viewing, management

## Error Handling

```go
import "errors"

sf, err := storefront.NewStorefront(YOUR_API_KEY)
if err != nil {
    if errors.Is(err, storefront.ErrInvalidAPIKey) {
        log.Fatal("Invalid or missing API key")
    }
}

products, err := sf.Products().List(ctx, opts)
if err != nil {
    // Handle generic API errors
    log.Printf("API error: %v", err)
}
```

## Testing

Tests use mocked API responses and require no environment setup:
```bash
go test -v ./...
go test -cover ./...
```

## License

MIT License. See [LICENSE](./LICENSE) file. For support, visit our [GitHub Issues](https://github.com/artyconst/fb-storefront-go/issues).

## Found This Useful?

If you're building with Fleetbase Storefront and find this SDK valuable, please consider starring the repository. Your support helps others discover it and encourages continued development.

[![GitHub stars](https://img.shields.io/github/stars/artyconst/fb-storefront-go?style=social)](https://github.com/artyconst/fb-storefront-go/stargazers)

---

*This Go SDK is an independent implementation using the official Fleetbase Storefront OpenAPI specification.*
