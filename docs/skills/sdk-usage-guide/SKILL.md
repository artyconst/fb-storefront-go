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
- **Service Access Patterns** - How to construct and use Products, Cart, Checkout, Customers, Orders, Categories, Reviews, Store, FoodTruck, and Location services
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
    sf.WithUserAgent("MyApp/1.0"),                            // Required: identifies your application to the server
    sf.WithAPIHost("https://api.storefront.fleetbase.io/v1"),  // Optional custom host
    sf.WithTimeout(60),                                        // Optional timeout in seconds
    sf.WithLogLevel(sf.LevelDebug),                           // Optional log level
)
if err != nil {
    return fmt.Errorf("failed to initialize client: %w", err)
}
```

### Service Access Pattern

Services are constructed manually using `NewXxxService(client)` constructors — NOT via accessor methods on the client. Each service is a standalone struct that holds a reference to the `StorefrontClient`.

```go
import (
    "github.com/artyconst/fb-storefront-go/pkg/resources/cart"
    "github.com/artyconst/fb-storefront-go/pkg/resources/category"
    "github.com/artyconst/fb-storefront-go/pkg/resources/checkout"
    "github.com/artyconst/fb-storefront-go/pkg/resources/customer"
    "github.com/artyconst/fb-storefront-go/pkg/resources/foodtruck"
    "github.com/artyconst/fb-storefront-go/pkg/resources/location"
    "github.com/artyconst/fb-storefront-go/pkg/resources/order"
    "github.com/artyconst/fb-storefront-go/pkg/resources/product"
    "github.com/artyconst/fb-storefront-go/pkg/resources/review"
    "github.com/artyconst/fb-storefront-go/pkg/resources/store"
)

// Construct each service independently
productSDK := product.NewProductService(sfClient)
cartSDK := cart.NewCartService(sfClient)
checkoutSDK := checkout.NewCheckoutService(sfClient)
customerSDK := customer.NewCustomerService(sfClient)
orderSDK := order.NewOrderService(sfClient)
categorySDK := category.NewCategoryService(sfClient)
reviewSDK := review.NewReviewService(sfClient)
storeSDK := store.NewStoreService(sfClient)
foodtruckSDK := foodtruck.NewService(sfClient)
locationSDK := location.NewService(sfClient)
```

### Context Usage

All service methods accept `context.Context` as the first parameter:

```go
import "time"

ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

products, err := productSDK.List(ctx, product.WithLimit(20))
if ctx.Err() == context.DeadlineExceeded {
    return fmt.Errorf("request timed out")
}
```

### Type Conventions

- **Currency Values**: Use `int64` for price fields (store values in smallest currency unit)
- **Optional Fields**: Use pointer types (`*string`, `*int64`) in request/response structs
- **Service Types**: Each domain has its own service type with dedicated methods
- **Method Naming**: Methods follow Go conventions, not OpenAPI operation IDs

## Functional Options Pattern

The SDK uses the functional options pattern extensively for optional parameters. Instead of passing large structs or many positional arguments, services accept variadic `...Option` parameters that modify an internal options struct before making the API call.

### How It Works

Each service defines option types and corresponding functions:

```go
// Product ListOptions with functional options
products, err := productSDK.List(ctx,
    product.WithLimit(20),
    product.WithOffset(100),
    product.WithCategory("cat_abc"),
    product.WithSortBy("created_at"),
    product.WithOrder("desc"),
)

// Cart AddProduct with functional options
cart, err := cartSDK.AddProduct(ctx, "cart_123", "prod_456",
    cart.WithQuantity(2),
    cart.WithVariant("size", "L"),
    cart.WithScheduledAt("2026-07-01T10:00:00Z"),
    cart.WithStoreLocation("loc_main"),
)

// Cart UpdateLineItem with functional options
cart, err = cartSDK.UpdateLineItem(ctx, "cart_123", "item_xyz",
    cart.WithQuantityForUpdate(5),
    cart.WithVariantForUpdate("color", "red"),
)

// Checkout Initialize with functional options
preview, err := checkoutSDK.Initialize(ctx,
    checkout.WithGateway("stripe"),
    checkout.WithCustomer("cust_123"),
    checkout.WithCart("cart_456"),
    checkout.WithCash(true),
    checkout.WithTip(500),
)

// Checkout CaptureCheckout with functional options
checkoutResult, err := checkoutSDK.CaptureCheckout(ctx, "token_abc",
    checkout.WithTransactionDetails(details),
    checkout.WithNotes("Rush order"),
)

// Order GenerateReceipt with functional options
receipt, err := orderSDK.GenerateReceipt(ctx, "order_123",
    order.WithEbarimtReceiverType("individual"),
    order.WithEbarimtReceiver(receiverInfo),
)
```

### Benefits of This Pattern

- **Optional parameters are truly optional** — only pass what you need
- **Composable** — chain multiple options together naturally
- **Type-safe** — the compiler catches mistakes at compile time
- **Extensible** — new options can be added without breaking existing code

## Service Methods Reference

### Products Service

```go
productSDK := product.NewProductService(sfClient)

// List products with filtering and pagination
products, err := productSDK.List(ctx,
    product.WithLimit(20),
    product.WithOffset(0),
    product.WithCategory("cat_abc"),
    product.WithSortBy("created_at"),
    product.WithOrder("desc"),
)

// Get single product by ID
product, err := productSDK.Get(ctx, "prod_abc123")

// Create a new product (functional options for fields like name, price, etc.)
newProduct, err := productSDK.Create(ctx,
    product.WithName("Running Shoes"),
    product.WithPrice(9999), // 99.99 in cents
)

// Update an existing product
updatedProduct, err := productSDK.Update(ctx, "prod_abc123",
    product.WithName("Updated Running Shoes"),
)
```

### Cart Service

Carts are retrieved via `Get()` — if no cart exists for the given ID, the server creates one automatically. There is no explicit Create method.

```go
cartSDK := cart.NewCartService(sfClient)

// Get or create cart (creates implicitly if it doesn't exist)
cart, err := cartSDK.Get(ctx, "cart_abc123")
if err != nil {
    return fmt.Errorf("failed to get cart: %w", err)
}

// Add item to cart (product ID is in URL path, quantity as positional arg)
cart, err = cartSDK.AddItem(ctx, cart.ID, "prod_456", 2,
    cart.WithVariant("size", "L"),
    cart.WithScheduledAt("2026-07-01T10:00:00Z"),
)

// Add product to cart (functional options for quantity and variants)
cart, err = cartSDK.AddProduct(ctx, cart.ID, "prod_456",
    cart.WithQuantity(3),
    cart.WithVariant("color", "blue"),
)

// Update line item in cart
cart, err = cartSDK.UpdateLineItem(ctx, cart.ID, "item_xyz789",
    cart.WithQuantityForUpdate(5),
)

// Remove line item from cart (returns updated cart)
cart, err = cartSDK.RemoveLineItem(ctx, cart.ID, "item_xyz789")

// Empty the entire cart (returns empty cart)
emptyCart, err := cartSDK.EmptyCart(ctx, cart.ID)

// Delete the entire cart entirely
err = cartSDK.DeleteCart(ctx)

// Checkout the cart and create an order
order, err := cartSDK.Checkout(ctx, cart.ID, cart.CheckoutRequest{
    CustomerEmail:   "user@example.com",
    ShippingAddress: &cart.Address{...},
})

// Cart total is available as int64 (cents/smallest currency unit)
fmt.Printf("Cart total: %d\n", cart.TotalAmount)
```

### Checkout Service

```go
checkoutSDK := checkout.NewCheckoutService(sfClient)

// Create a new checkout session from a cart
session, err := checkoutSDK.Create(ctx, "cart_abc123", checkout.CreateCheckoutRequest{
    CustomerEmail:   "user@example.com",
    ShippingAddress: &checkout.Address{...},
    PaymentMethodID: "pm_stripe_xyz",
})

// Get checkout by ID
session, err := checkoutSDK.Get(ctx, "chk_abc123")

// Update customer info for a checkout
session, err = checkoutSDK.UpdateCustomer(ctx, "chk_abc123", checkout.CustomerInfo{...})

// Process payment for the checkout
session, err = checkoutSDK.ProcessPayment(ctx, "chk_abc123", checkout.PaymentRequest{...})

// Initialize checkout preview (functional options)
preview, err := checkoutSDK.Initialize(ctx,
    checkout.WithGateway("stripe"),
    checkout.WithCustomer("cust_123"),
    checkout.WithCart("cart_456"),
    checkout.WithCash(true),
    checkout.WithTip(500),
)

// Get checkout status (functional options for ID or token)
status, err := checkoutSDK.Status(ctx,
    checkout.WithCheckoutID("chk_abc123"),
)

// Capture checkout as an order (functional options)
captured, err := checkoutSDK.CaptureCheckout(ctx, "token_xyz",
    checkout.WithTransactionDetails(details),
    checkout.WithNotes("Rush delivery"),
)

// Capture QPay checkout
order, err := checkoutSDK.CaptureQPay(ctx, "chk_abc123", true, nil)

// Get delivery service quote
quote, err := checkoutSDK.GetDeliveryServiceQuote(ctx, checkout.ServiceQuoteParams{
    Origin:      "loc_origin",
    Destination: "loc_dest",
    CartID:      "cart_456",
})

// Update payment intent for a checkout
preview, err = checkoutSDK.UpdatePaymentIntent(ctx, "pi_abc123")
```

### Customers Service

```go
customerSDK := customer.NewCustomerService(sfClient)

// Create a new customer
cust, err := customerSDK.Create(ctx, customer.CustomerCreateRequest{
    Identity: "user@example.com",
})

// Get customer by ID
cust, err := customerSDK.Get(ctx, "cust_abc123")

// Login with identity and password
loginResp, err := customerSDK.Login(ctx, customer.LoginRequest{
    Identity: "user@example.com",
    Password: "secure_password_123",
})

// Login via SMS authentication
smsLogin, err := customerSDK.LoginWithSMS(ctx, customer.SMSSignInRequest{
    Identity: "+1234567890",
})

// Verify SMS code to complete login
loginResp, err = customerSDK.VerifySMSCode(ctx, customer.SMSConfirmSignInRequest{
    Identity: "+1234567890",
    Code:     "123456",
})

// Login via social providers
appleLogin, err := customerSDK.LoginWithApple(ctx, identityToken, authCode)
googleLogin, err := customerSDK.LoginWithGoogle(ctx, idToken, clientID)
fbCustomer, err := customerSDK.LoginWithFacebook(ctx, facebookUserID, nil, nil)

// Update customer profile (functional options)
cust, err = customerSDK.Update(ctx, "cust_abc123",
    customer.WithName("John Doe"),
    customer.WithEmail("john@example.com"),
    customer.WithPhone("+1234567890"),
)

// List places for authenticated customer (requires token)
places, err := customerSDK.ListPlaces(ctx, "customer_token_here",
    customer.PlaceWithLimit(20),
)

// List orders for authenticated customer (requires token)
orders, err := customerSDK.ListOrders(ctx, "customer_token_here",
    customer.OrderWithStatus("delivered"),
)

// Request account creation code
err = customerSDK.RequestCreationCode(ctx, customer.RequestCreationCodeRequest{
    Identity: "user@example.com",
})

// Register device for push notifications
resp, err := customerSDK.RegisterDevice(ctx, "customer_token_here", customer.RegisterDeviceRequest{...})

// Phone verification flow
err = customerSDK.RequestPhoneVerification(ctx, "+1234567890")
err = customerSDK.VerifyPhoneNumber(ctx, "123456", "+1234567890")

// Account closure flow
err = customerSDK.InitiateAccountClosure(ctx)
err = customerSDK.ConfirmAccountClosure(ctx, "closure_code_123")

// Stripe integration
ephemeralKey, err := customerSDK.GetStripeEphemeralKey(ctx)
setupIntent, err := customerSDK.CreateStripeSetupIntent(ctx)
```

### Orders Service

```go
orderSDK := order.NewOrderService(sfClient)

// List orders with optional filtering
orders, err := orderSDK.List(ctx,
    order.WithLimit(20),
    order.WithStatus(order.StatusDelivered),
)

// Get single order by ID or order number
ord, err := orderSDK.Get(ctx, "order_abc123")

// Create a new order from a cart
ord, err = orderSDK.Create(ctx, "cart_abc123")

// Mark an order as picked up
err = orderSDK.MarkPickedUp(ctx, "order_abc123")

// Generate a receipt for an order (functional options)
receipt, err := orderSDK.GenerateReceipt(ctx, "order_abc123",
    order.WithEbarimtReceiverType("individual"),
)
```

### Categories Service

```go
categorySDK := category.NewCategoryService(sfClient)

// List categories with optional filtering and pagination
categories, err := categorySDK.List(ctx,
    category.WithLimit(50),
    category.WithSearch("electronics"),
)

// Get category by ID
category, err := categorySDK.Get(ctx, "cat_abc123")
```

### Reviews Service

```go
reviewSDK := review.NewReviewService(sfClient)

// List all reviews with pagination
reviews, err := reviewSDK.List(ctx,
    review.WithLimit(10),
    review.WithOffset(0),
)

// Create a new review (functional options for additional fields)
newReview, err := reviewSDK.Create(ctx, "Great product!", 5,
    review.WithProductID("prod_abc123"),
)

// Get single review by ID
review, err := reviewSDK.Get(ctx, "rev_xyz789")

// Delete a review
err = reviewSDK.Delete(ctx, "rev_xyz789")

// Count reviews for a store
count, err := reviewSDK.CountByStore(ctx, "store_abc123")

// Count reviews by rating (1-5)
count, err = reviewSDK.CountByRating(ctx, 5)
```

### Store Service

```go
storeSDK := store.NewStoreService(sfClient)

// Get about store information
aboutInfo, err := storeSDK.About(ctx)

// List payment gateways with pagination
gateways, err := storeSDK.ListGateways(ctx,
    store.WithLimit(20),
)

// Get specific gateway by ID
gateway, err := storeSDK.GetGateway(ctx, "gw_stripe")
```

### FoodTruck Service

```go
foodtruckSDK := foodtruck.NewService(sfClient)

// List all food trucks with optional filtering
trucks, err := foodtruckSDK.List(ctx,
    foodtruck.WithLimit(20),
    foodtruck.WithSort("name"),
)

// Get single food truck by ID
truck, err := foodtruckSDK.Get(ctx, "ft_abc123")
```

### Location Service

```go
locationSDK := location.NewService(sfClient)

// List locations for a specific store
locations, err := locationSDK.List(ctx, "store_abc123",
    location.WithLimit(50),
)

// Get single location by ID (requires store ID as query param)
loc, err := locationSDK.Get(ctx, "store_abc123", "loc_xyz789")

// List all store locations via dedicated endpoint
storeLocations, err := locationSDK.ListLocations(ctx)
```

## Configuration Options

All configuration uses the functional options pattern:

| Option | Description | Example |
|--------|-------------|---------|
| `WithUserAgent(ua)` | **Required.** Custom User-Agent header. The server rejects requests without one (403). | `WithUserAgent("MyApp/1.0")` |
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
sfClient, err := sf.NewStorefront("invalid_key")
if err != nil {
    if errors.Is(err, sf.ErrInvalidAPIKey) {
        log.Fatal("Please check your API key configuration")
    }
}

// Check for resource not found
products, err := productSDK.List(ctx, opts)
if errors.Is(err, sf.ErrResourceNotFound) {
    log.Println("No products found matching criteria")
}

// Handle cart-specific errors
err = cartSDK.AddItem(ctx, "cart_123", "prod_456", 0)
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
func getListProducts(ctx context.Context) ([]*product.Product, error) {
    products, err := productSDK.List(ctx, opts)
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
        fmt.Fprintln(w, `[{"id": "1", "name": "Test"}]`)
    })

    client := setupTestClient(t, handler)
    
    productSDK := product.NewProductService(client)
    products, err := productSDK.List(context.Background(), product.WithLimit(10))
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
products, err := productSDK.List(product.WithLimit(20))

// ✅ Correct
products, err := productSDK.List(context.Background(), product.WithLimit(20))
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
product, err := productSDK.Get(ctx, "id")
if product.Description != nil {
    fmt.Println(*product.Description)
} else {
    fmt.Println("No description available")
}
```

### 4. Ignoring Context Cancellation

Always check for context cancellation:

```go
products, err := productSDK.List(ctx, opts)
if ctx.Err() != nil {
    return fmt.Errorf("operation cancelled: %w", ctx.Err())
}
```

### 5. Not Using Pointer Types for Optional Request Fields

When sending requests, use pointers for optional fields:

```go
// ❌ Wrong - always sends field even if empty
customerSDK.Update(ctx, "cust_123")

// ✅ Correct - only updates if pointer is non-nil via functional options
customerSDK.Update(ctx, "cust_123", customer.WithName("John"))
```

### 6. Using Client Accessor Methods That Don't Exist

Services are NOT accessed via methods on the client:

```go
// ❌ Wrong - these accessor methods do not exist
products, err := sfClient.Products().List(ctx)
cart, err := sfClient.Cart().Get(ctx, "id")

// ✅ Correct - construct services manually
productSDK := product.NewProductService(sfClient)
products, err := productSDK.List(ctx)
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

- [ ] Client initialized with WithUserAgent option (required by server, returns 403 without it)
- [ ] Client initialized with proper error handling
- [ ] Services constructed via `NewXxxService(client)` — NOT via client accessor methods
- [ ] All service methods called with context.Context as first parameter
- [ ] Errors checked and wrapped appropriately using `%w`
- [ ] Specific errors validated using `errors.Is()` and `errors.As()`
- [ ] Optional fields use pointer types in structs
- [ ] Currency values stored as int64 (not strings)
- [ ] Tests use mocked HTTP clients via setupTestClient()
- [ ] Context timeouts are set for all operations

---

*This skill is your primary reference when working with the Fleetbase Storefront Go SDK. Consult it whenever you need guidance on initialization, service usage, error handling, or testing patterns.*
