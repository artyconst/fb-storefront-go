---
name: sdk-usage-guide
description: Comprehensive guide for using the Fleetbase Storefront Go SDK including initialization, service patterns, and best practices
license: MIT
compatibility: opencode
metadata:
  domain: go-sdk-development
  language: go
---

## What I do

I provide comprehensive guidance for working with the Fleetbase Storefront Go SDK. I help agents with:

- **SDK Initialization** - Proper client setup and configuration options
- **Service Access Patterns** - How to access Products, Cart, Checkout, Customers, Orders, Categories, Reviews, and Store services
- **API Usage** - Correct method signatures and parameter patterns for all service operations
- **Error Handling** - Best practices for typed errors, error wrapping, and checking specific error types
- **Context Management** - Proper use of context.Context for cancellation and timeouts
- **Testing Patterns** - How to write tests using mocked HTTP clients

## When to use me

Use this skill when you need to:

1. Initialize the Storefront client with proper configuration
2. Access and use any of the SDK's service methods (Products, Cart, Checkout, etc.)
3. Understand the correct parameter types for API operations
4. Implement proper error handling using typed errors
5. Write tests that mock HTTP responses without network calls
6. Follow Go idioms specific to this SDK (pointer usage, Decimal types, etc.)

This skill is your primary reference when implementing features with the Fleetbase Storefront Go SDK. Consult me before writing any SDK integration code.

## Core Usage Patterns

### Client Initialization

Always initialize the client with proper error handling:

```go
import sf "github.com/artyconst/fb-storefront-go"

sfClient, err := sf.NewStorefront("sk_test_your_api_key_here",
    sf.WithAPIHost("https://api.storefront.fleetbase.io/v1"),  // Optional custom host
    sf.WithTimeout(60),                                        // Optional timeout in seconds
    sf.WithLogLevel(sf.LevelDebug),                           // Optional log level
)
if err != nil {
    return fmt.Errorf("failed to initialize client: %w", err)
}
```

### Service Access Pattern

Access services through the main client instance:

```go
// All services are accessed via methods on the client
sf.Products()    // Product catalog operations
sf.Cart()        // Shopping cart operations  
sf.Checkout()    // Checkout session operations
sf.Customers()   // Customer account operations
sf.Orders()      // Order management operations
sf.Categories()  // Category browsing operations
sf.Reviews()     // Product review operations
sf.Store()       // Store configuration operations
```

### Context Usage

All service methods accept `context.Context` as the first parameter:

```go
import "time"

ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

products, err := sfClient.Products().List(ctx, product.ListOptions{})
if ctx.Err() == context.DeadlineExceeded {
    return fmt.Errorf("request timed out")
}
```

### Type Conventions

- **Currency Values**: Use `int64` for price fields (store values in smallest currency unit)
- **Optional Fields**: Use pointer types (`*string`, `*int64`) in request/response structs
- **Service Types**: Each domain has its own service type with dedicated methods
- **Method Naming**: Methods follow Go conventions, not OpenAPI operation IDs

## Service Methods Reference

### Products Service

```go
// List products with filtering and pagination
products, err := sf.Products().List(ctx, 
    product.WithLimit(20),
    product.WithOffset(0),
    product.WithSort("created_at"),
)

// Search products by query
products, err := sf.Products().Search(ctx, "running shoes", 
    product.WithPriceRange(1000, 50000), // prices in cents
)

// Get single product by ID
product, err := sf.Products().Get(ctx, "prod_abc123")

// List reviews for a product
reviews, err := sf.Products().ListReviews(ctx, "prod_abc123", 
    review.WithLimit(10),
)
```

### Cart Service

```go
// Get current cart (creates if doesn't exist)
cart, err := sf.Cart().GetOrCreateCart(ctx)

// Add item to cart
cart, err = sf.Cart().AddItem(ctx, "product_id_123", 2)

// Update cart item quantity
cart, err = sf.Cart().UpdateItem(ctx, "item_id", 5)

// Remove item from cart
err = sf.Cart().RemoveItem(ctx, "item_id")

// Clear entire cart
err = sf.Cart().ClearCart(ctx)

// Get cart totals
total := cart.GetTotal() // returns int64 in cents
```

### Checkout Service

```go
// Create checkout session from cart
session, err := sf.Checkout().CreateSession(ctx, 
    checkout.WithPaymentMethod("stripe"),
    checkout.WithMetadata(map[string]string{"order_type": "express"}),
)

// Complete payment for checkout session
result, err := sf.Checkout().CompleteSession(ctx, "session_id", 
    checkout.WithPaymentIntentData("pi_abc123"),
)

// Retrieve checkout session
session, err := sf.Checkout().GetSession(ctx, "session_id")
```

### Customers Service

```go
// Register new customer
customer, err := sf.Customers().Register(ctx, 
    customer.RegisterRequest{
        Email: "user@example.com",
        Password: "secure_password_123",
        FirstName: "John",
        LastName: "Doe",
    },
)

// Login existing customer
customer, err := sf.Customers().Login(ctx, 
    customer.LoginRequest{
        Email: "user@example.com",
        Password: "secure_password_123",
    },
)

// Get current authenticated customer
me, err := sf.Customers().GetMe(ctx)

// Update customer profile
customer, err = sf.Customers().UpdateProfile(ctx, 
    customer.UpdateRequest{
        FirstName: pt.Str("Jonathan"),
        Phone:     pt.Str("+1234567890"),
    },
)

// Add address to customer
address, err := sf.Customers().AddAddress(ctx, 
    customer.AddressRequest{
        Type:      "shipping",
        Address1:  "123 Main St",
        City:      "San Francisco",
        State:     "CA",
        ZipCode:   "94105",
        Country:   "US",
    },
)

// Get customer orders
orders, err := sf.Customers().GetOrders(ctx)
```

### Orders Service

```go
// List all orders for current customer
orders, err := sf.Orders().List(ctx, 
    order.WithLimit(20),
    order.WithStatus("delivered"),
)

// Get single order by ID
order, err := sf.Orders().Get(ctx, "order_abc123")

// Cancel an order (if allowed)
err = sf.Orders().Cancel(ctx, "order_abc123")

// List orders for specific customer
orders, err := sf.Orders().ListForCustomer(ctx, "customer_id", 
    order.WithStatus("processing"),
)

// Get order items
items := order.GetItems() // returns []*OrderItem
```

### Categories Service

```go
// List top-level categories
categories, err := sf.Categories().List(ctx)

// Get category by ID with children
category, err := sf.Categories().Get(ctx, "cat_abc123")

// Search categories by name
categories, err := sf.Categories().Search(ctx, "electronics")

// Get products in a category
products, err := sf.Categories().ListProducts(ctx, "cat_abc123", 
    product.WithLimit(50),
)
```

### Reviews Service

```go
// List reviews for a product
reviews, err := sf.Reviews().List(ctx, review.ListOptions{
    ProductID: "prod_abc123",
    Limit:     10,
})

// Submit a new review
review, err := sf.Reviews().Submit(ctx, 
    review.SubmitRequest{
        ProductID: "prod_abc123",
        Rating:    5,
        Title:     "Excellent product!",
        Body:      "Very satisfied with this purchase.",
    },
)

// Count reviews by rating
counts, err := sf.Reviews().CountByRating(ctx, "prod_abc123")
```

### Store Service

```go
// Get store configuration and information
storeInfo, err := sf.Store().Get(ctx)

// List stores (for multi-store setups)
stores, err := sf.Store().List(ctx)

// Get store by ID
store, err := sf.Store().GetByID(ctx, "store_abc123")
```

## Configuration Options

All configuration uses the functional options pattern:

| Option | Description | Example |
|--------|-------------|---------|
| `WithAPIHost(host)` | Custom API host URL | `WithAPIHost("https://api.example.com")` |
| `WithAPIPath(path)` | Custom API path suffix | `WithAPIPath("/v1/custom")` |
| `WithTimeout(seconds)` | HTTP timeout in seconds | `WithTimeout(60)` |
| `WithLogLevel(level)` | Log level setting | `WithLogLevel(sf.LevelDebug)` |
| `WithLoggerOutput(w)` | Custom output writer | `WithLoggerOutput(os.Stdout)` |
| `WithDebugMode()` | Enable debug logging | `WithDebugMode()` |

### Debug Mode Usage

Enable detailed request/response logging:

```go
sfClient, err := sf.NewStorefront("sk_test_key",
    sf.WithDebugMode(),  // Logs all HTTP requests/responses at debug level
)
```

## Error Handling Best Practices

### Typed Errors

The SDK defines specific error types for programmatic handling:

```go
import (
    "errors"
    sf "github.com/artyconst/fb-storefront-go"
)

// Check for API key errors
sf, err := sf.NewStorefront("invalid_key")
if err != nil {
    if errors.Is(err, sf.ErrInvalidAPIKey) {
        log.Fatal("Please check your API key configuration")
    }
}

// Check for resource not found
products, err := sf.Products().List(ctx, opts)
if errors.Is(err, sf.ErrResourceNotFound) {
    log.Println("No products found matching criteria")
}

// Handle cart-specific errors
err = sf.Cart().AddItem(ctx, "product_123", 0)
if errors.Is(err, sf.ErrCartEmpty) {
    // Cart is empty
}
```

### APIError Type

Handle specific API error responses:

```go
var apiErr *sf.APIError
if errors.As(err, &apiErr) {
    log.Printf("API Error %s: %s (status: %d)", 
        apiErr.Code, 
        apiErr.Message, 
        apiErr.Status,
    )
}
```

### Error Wrapping

Always wrap errors with context using `%w`:

```go
func getListProducts(ctx context.Context) ([]*sf.Product, error) {
    products, err := sfClient.Products().List(ctx, opts)
    if err != nil {
        return nil, fmt.Errorf("failed to list products: %w", err)
    }
    return products, nil
}

// Wrapped errors can still be checked
err := getListProducts(ctx)
if errors.Is(err, sf.ErrInvalidAPIKey) {
    // This works even with wrapped errors
}
```

## Testing Patterns

### Mocked HTTP Client

Tests use mocked API responses via `setupTestClient()`:

```go
import (
    "net/http"
    "testing"
    sf "github.com/artyconst/fb-storefront-go"
)

func setupTestClient(t *testing.T, handler http.Handler) *sf.StorefrontClient {
    server := httptest.NewServer(handler)
    t.Cleanup(server.Close)

    client, err := sf.NewStorefront("sk_test_key", 
        sf.WithAPIHost(server.URL),
    )
    if err != nil {
        t.Fatalf("Failed to create test client: %v", err)
    }
    return client
}

func TestProductList(t *testing.T) {
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        fmt.Fprintln(w, `{"data": [{"id": "1", "name": "Test"}]}`)
    })

    client := setupTestClient(t, handler)
    
    products, err := client.Products().List(context.Background(), product.WithLimit(10))
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }
    
    if len(products) != 1 {
        t.Errorf("Expected 1 product, got %d", len(products))
    }
}
```

### Test Utilities

Use the `internal/utils` package for common test patterns:
- `setupTestClient()` - Creates client with mock handler
- `HTTPClient` interface - Abstraction for mocking HTTP operations
- Table-driven tests using `t.Run()` for multiple scenarios

### Testing Patterns Summary

1. **Always use `httptest.NewServer()`** to create a test server
2. **Define handlers** that return appropriate JSON responses
3. **Use t.Cleanup()** to ensure servers are closed after tests
4. **Validate both success and error cases** for each operation
5. **Check specific error types** using `errors.Is()` and `errors.As()`

## Common Pitfalls

### 1. Forgetting Context

Always pass context as the first parameter:

```go
// ❌ Wrong - no context
products, err := sf.Products().List(product.WithLimit(20))

// ✅ Correct
products, err := sf.Products().List(context.Background(), product.WithLimit(20))
```

### 2. Using String for Prices

Never use string types for currency values:

```go
// ❌ Wrong - loses precision
type Product struct {
    Price string `json:"price"`
}

// ✅ Correct - maintain precision
type Product struct {
    Price int64 `json:"price"` // value in cents/smallest unit
}
```

### 3. Not Checking Optional Fields

Optional fields are pointers and may be nil:

```go
product, err := sf.Products().Get(ctx, "id")
if product.Description != nil {
    fmt.Println(*product.Description)
} else {
    fmt.Println("No description available")
}
```

### 4. Ignoring Context Cancellation

Always check for context cancellation:

```go
products, err := sf.Products().List(ctx, opts)
if ctx.Err() != nil {
    return fmt.Errorf("operation cancelled: %w", ctx.Err())
}
```

### 5. Not Using Pointer Types for Optional Request Fields

When sending requests, use pointers for optional fields:

```go
// ❌ Wrong - always sends field even if empty
sf.Customers().UpdateProfile(ctx, customer.UpdateRequest{
    FirstName: "John",
})

// ✅ Correct - only updates if pointer is non-nil
firstName := "John"
sf.Customers().UpdateProfile(ctx, customer.UpdateRequest{
    FirstName: &firstName,
})
```

## Quick Reference: Helper Functions

Available helper functions for creating pointers:

```go
// String helper
name := sf.Str("value")

// Int64 helper  
qty := sf.Int64(10)

// Bool helper
active := sf.Bool(true)
```

Use these helpers when you need to send optional fields in requests.

## Integration Checklist

Before implementing SDK features, verify:

- [ ] Client initialized with proper error handling
- [ ] All service methods called with context.Context as first parameter
- [ ] Errors checked and wrapped appropriately using `%w`
- [ ] Specific errors validated using `errors.Is()` and `errors.As()`
- [ ] Optional fields use pointer types in structs
- [ ] Currency values stored as int64 (not strings)
- [ ] Tests use mocked HTTP clients via setupTestClient()
- [ ] Context timeouts are set for all operations

---

*This skill is your primary reference when working with the Fleetbase Storefront Go SDK. Consult it whenever you need guidance on initialization, service usage, error handling, or testing patterns.*
