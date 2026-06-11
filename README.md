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
import (
    "context"
    "log"
    "time"

    "github.com/artyconst/fb-storefront-go"
    "github.com/artyconst/fb-storefront-go/pkg/config"
)

// Initialize SDK with basic configuration (replace with your actual API key)
client, err := storefront.NewStorefront("sk_test_your_storefront_key_here",
    storefront.WithUserAgent("MyApp/1.0"),                            // Required: identifies your application to the server
    // Uncomment and modify these for advanced configuration:
    // storefront.WithAPIHost("https://api.custom-domain.com"),      // Custom API host URL
    // storefront.WithTimeout(60 * time.Second),                     // HTTP timeout (time.Duration)
    // storefront.WithLogLevel(config.LevelDebug),                    // Set log level (Error, Warn, Info, Debug)
    // storefront.WithDebugMode(),                                    // Enable debug logging for API calls
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
    "time"

    "github.com/artyconst/fb-storefront-go"
    "github.com/artyconst/fb-storefront-go/pkg/product"
)

func main() {
    client, err := storefront.NewStorefront("sk_test_your_storefront_key_here",
        storefront.WithUserAgent("MyApp/1.0"),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Use context with timeout for cancellation
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    products, err := client.Products().List(ctx, product.WithLimit(20))
    if err != nil {
        log.Fatal(err)
    }

    // Monetary values are int64 (cents/smallest currency unit)
    for _, p := range products {
        fmt.Printf("%s: %d\n", p.Name, p.Price)
        if p.SalePrice != nil {
            fmt.Printf("  Sale price: %d\n", *p.SalePrice)
        }
    }
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
| `WithUserAgent(ua)` | **Required.** Custom User-Agent header identifying your application. The server rejects requests without one (403). | `WithUserAgent("MyApp/1.0")` |
| `WithAPIHost(host)` | Custom API host URL | `WithAPIHost("https://api.example.com")` |
| `WithTimeout(duration)` | HTTP timeout (time.Duration) | `WithTimeout(60 * time.Second)` |
| `WithLogLevel(level)` | Log level (Error, Warn, Info, Debug) | `WithLogLevel(config.LevelDebug)` |
| `WithLoggerOutput(w)` | Custom logger output writer | `WithLoggerOutput(os.Stdout)` |
| `WithDebugMode()` | Enable debug logging for API calls | `WithDebugMode()` |

## Core Concepts

The SDK uses a service-based architecture. Access services via method calls on the client:

```go
client := storefront.NewStorefront(YOUR_API_KEY)

// Available services
client.Store()       // Store configuration and information retrieval
client.Products()    // Product catalog
client.Cart()        // Shopping cart
client.Checkout()    // Checkout sessions
client.Customers()   // Customer accounts
client.Orders()      // Order management
client.Categories()  // Product categories
client.Reviews()     // Product review submission and listing
client.Locations()   // Store location and working hours retrieval
```

All operations accept `context.Context` for cancellation and timeouts.

## Product Type

The `Product` type represents a product in the store:

```go
type Product struct {
    ID              string                 `json:"id"`
    Name            string                 `json:"name"`
    Description     *string                `json:"description,omitempty"`
    Price           int64                  `json:"price"`
    SalePrice       *int64                 `json:"sale_price,omitempty"`
    Currency        string                 `json:"currency"`
    SKU             *string                `json:"sku,omitempty"`
    PrimaryImageURL *string                `json:"primary_image_url,omitempty"`
    Tags            []string               `json:"tags"`
    Status          string                 `json:"status"`
    Meta            interface{}            `json:"meta"`
    Slug            string                 `json:"slug"`
    IsOnSale        bool                   `json:"is_on_sale"`
    IsRecommended   bool                   `json:"is_recommended"`
    IsService       bool                   `json:"is_service"`
    IsBookable      bool                   `json:"is_bookable"`
    IsAvailable     bool                   `json:"is_available"`
    Translations    []product.Translation  `json:"translations"`
    AddonCategories []product.AddonCategory `json:"addon_categories"`
    Variants        []product.Variant      `json:"variants"`
    Images          []string               `json:"images"`
    Videos          []string               `json:"videos"`
    Hours           []product.Hour         `json:"hours"`
    YouTubeURLs     []string               `json:"youtube_urls"`
    CreatedAt       string                 `json:"created_at"`
    UpdatedAt       string                 `json:"updated_at"`
}
```

> **Note:** Monetary values (`Price`, `SalePrice`) are `int64` representing amounts in the smallest currency unit (e.g., cents). Use `%d` format verbs when printing — never `.String()`.

## Service Documentation

- [Products Guide](./user-guide-products.md) - Browse, search, manage products
- [Categories Guide](./user-guide-categories.md) - Organize hierarchies
- [Cart Guide](./user-guide-cart.md) - Shopping cart operations
- [Checkout Guide](./user-guide-checkout.md) - Checkout and payment flow
- [Customers Guide](./user-guide-customers.md) - Customer account management
- [Orders Guide](./user-guide-orders.md) - Order viewing, management
- [Store Guide](./user-guide-store.md) - Store configuration and payment gateways
- [Reviews Guide](./user-guide-reviews.md) - Product review submission and listing
- [Locations Guide](./user-guide-locations.md) - Store location and working hours

## Store Service

Retrieve store configuration and payment gateways:

```go
storeService := storeSDK.NewStoreService(client)

// Get store information
about, err := storeService.About(ctx)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Store: %s\n", about.Store.Name)

// List payment gateways
gateways, err := storeService.ListGateways(ctx)
if err != nil {
    log.Fatal(err)
}
for _, g := range gateways.Data {
    fmt.Printf("Gateway: %s\n", g.Name)
}

// Get specific gateway
gateway, err := storeService.GetGateway(ctx, "gw_xxxxx")
if err != nil {
    log.Fatal(err)
}
```

## Location Service

Retrieve store locations and working hours:

```go
locService := locationSDK.NewService(client)

// List all locations for a store
locations, err := locService.List(ctx, storeID)
if err != nil {
    log.Fatal(err)
}
for _, loc := range locations {
    fmt.Printf("Location: %s\n", loc.Name)
    if loc.Place != nil {
        fmt.Printf("  Address: %s\n", loc.Place.Address)
        fmt.Printf("  Lat/Lng: %f, %f\n", loc.Place.Latitude, loc.Place.Longitude)
    }
    for _, hour := range loc.Hours {
        fmt.Printf("  %s: %s - %s\n", hour.Day, *hour.Start, *hour.End)
    }
}

// Get a specific location
location, err := locService.Get(ctx, storeID, locationID)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Location: %s\n", location.Name)
```

## Review Service

Submit and list product reviews:

```go
reviewService := reviewSDK.NewReviewService(client)

// Count reviews for a store
count, err := reviewService.CountByStore(ctx, storeID)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Total reviews: %d\n", count)

// Count reviews by rating
fiveStarCount, err := reviewService.CountByRating(ctx, 5)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("5-star reviews: %d\n", fiveStarCount)
```

## Error Handling

```go
import (
    "errors"
    "fmt"
)

client, err := storefront.NewStorefront(YOUR_API_KEY)
if err != nil {
    if errors.Is(err, storefront.ErrInvalidAPIKey) {
        log.Fatal("Invalid or missing API key")
    }
}

products, err := client.Products().List(ctx, opts)
if err != nil {
    // Handle generic API errors
    log.Printf("API error: %v", err)

    // Check for specific API error codes
    var apiErr *storefront.APIError
    if errors.As(err, &apiErr) {
        log.Printf("API error %s: %s (status: %d)", apiErr.Code, apiErr.Message, apiErr.Status)
    }
}
```

## For Agents

When implementing features with this SDK, consult the comprehensive [SDK Usage Guide](./docs/skills/sdk-usage-guide/SKILL.md) for initialization patterns, service methods reference, and best practices.

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
