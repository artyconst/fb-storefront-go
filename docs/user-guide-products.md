### Products Service

The Products service allows you to browse, search, and retrieve product information from your Fleetbase store.

#### List All Products

Get all products with optional filtering and pagination:

```go
import (
    "context"
    "fmt"
    "log"

    "github.com/artyconst/fb-storefront-go/pkg/product"
)

func listProducts() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    products, err := sf.Products().List(context.Background(),
        product.WithCategory("cat_electronics"),
        product.WithSortBy("name"),
        product.WithOrder("asc"),
    )
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Found %d products\n", len(products))
    for _, p := range products {
        fmt.Printf("%s - $%d\n", p.Name, p.Price)
    }
}
```

**ListOptions Parameters:**

| Parameter | Type | Description | Default |
|-----------|------|-------------|---------|
| `Limit` | uint64 | Maximum items to return (use WithLimit) | unlimited |
| `Offset` | int64 | Pagination offset (use WithOffset) | 0 |
| `Category` | string | Filter by category ID (use WithCategory) | - |
| `SortBy` | string | Sort field (use WithSortBy) | - |
| `Order` | string | Sort order: asc or desc (use WithOrder) | - |

Use functional options like `WithLimit()`, `WithOffset()`, etc. to set parameters.

#### Get Single Product

Retrieve a specific product by ID:

```go
func getProduct() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    prod, err := sf.Products().Get(context.Background(), "prod_abc123")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("%s - $%d\n", prod.Name, prod.Price)
    fmt.Printf("SKU: %s\n", prod.SKU)
}
```

#### Search Products

Search for products by query string:

```go
func searchProducts() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    products, err := sf.Products().Search(context.Background(), product.SearchQuery{
        Query:  "wireless headphones",
        Limit:  20,
        Offset: 0,
        Store:  "store_default",
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Found %d search results\n", len(products))
}
```

**SearchQuery Parameters:**

| Parameter | Type | Description | Required |
|-----------|------|-------------|----------|
| `Query` | string | Search query string | Yes |
| `Limit` | int64 | Maximum items to return | No |
| `Offset` | int64 | Pagination offset | No |
| `Store` | string | Store identifier | No |

#### Find Products by Category

Retrieve all products in a specific category using List with WithCategory:

```go
func getCategoryProducts() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    // Use functional options to filter by category and set limit
    products, err := sf.Products().List(context.Background(),
        product.WithCategory("cat_electronics"),
        product.WithLimit(50),
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Found %d electronics products\n", len(products))
}
```

#### Product Structure

The `Product` type contains the following fields:

```go
type Product struct {
    ID              string                 `json:"id"`
    Name            string                 `json:"name"`
    Description     *string                `json:"description,omitempty"`
    Price           int64                  `json:"price"`
    SalePrice       *int64                 `json:"sale_price,omitempty"`
    Currency        string                 `json:"currency,omitempty"`
    SKU             *string                `json:"sku,omitempty"`
    PrimaryImageURL string                 `json:"primary_image_url,omitempty"`
    Tags            []string               `json:"tags,omitempty"`
    Status          string                 `json:"status,omitempty"`
    Meta            map[string]interface{} `json:"meta,omitempty"`
    Slug            string                 `json:"slug,omitempty"`
    IsOnSale        bool                   `json:"is_on_sale,omitempty"`
    IsService       bool                   `json:"is_service,omitempty"`
    IsBookable      bool                   `json:"is_bookable,omitempty"`
    IsAvailable     bool                   `json:"is_available,omitempty"`
    CreatedAt       string                 `json:"created_at"`
    UpdatedAt       string                 `json:"updated_at"`
}
```

**Important Notes:**

- **Monetary Values**: All currency fields (`Price`, `SalePrice`) use `int64` to represent the price in the smallest currency unit (e.g., cents). Use `%d` format verb when displaying prices.
- **Optional Fields**: Pointer types (`*string`, `*int64`) indicate optional fields that may be nil.
- **Timestamps**: `CreatedAt` and `UpdatedAt` are ISO 8601 date strings, not `time.Time` objects.

#### Error Handling

```go
import (
    "errors"
    "log"

    "github.com/artyconst/fb-storefront-go"
)

func handleProductErrors() {
    _, err := sf.Products().Get(context.Background(), "nonexistent-id")
    if err != nil {
        if errors.Is(err, storefront.ErrResourceNotFound) {
            log.Println("Product not found")
        } else {
            log.Printf("API Error: %v", err)
        }
    }
}
```

## Product Write Operations

### Create (Create a new product using functional options)

```go
// Create a new product with required fields
product, err := service.Create(ctx, 
    productSDK.WithProductName("Premium Widget"),
    productSDK.WithPrice(2999), // $29.99 in cents
)
if err != nil {
    log.Fatalf("Failed to create product: %v", err)
}
fmt.Printf("Created product: %s (ID: %s)\n", product.Name, product.ID)

// Create with optional fields
product, err = service.Create(ctx, 
    productSDK.WithProductName("Premium Widget"),
    productSDK.WithPrice(2999),
    productSDK.WithDescription("A high-quality widget for all occasions"),
    productSDK.WithCategoryID("cat_01J..."),
)
if err != nil {
    log.Fatalf("Failed to create product: %v", err)
}
```

### Update (Update an existing product using functional options)

```go
// Update the price of a product
product, err := service.Update(ctx, "prod_01J...", 
    productSDK.WithUpdatePrice(3499), // $34.99 in cents
)
if err != nil {
    log.Fatalf("Failed to update product: %v", err)
}
fmt.Printf("Updated price to: %d\n", product.Price)

// Update multiple fields at once
product, err = service.Update(ctx, "prod_01J...", 
    productSDK.WithUpdateName("Premium Widget Pro"),
    productSDK.WithUpdatePrice(3499),
    productSDK.WithUpdateDescription("The upgraded version of our best-selling widget"),
)
if err != nil {
    log.Fatalf("Failed to update product: %v", err)
}
```

#### Complete Example with Context Timeout

```go
func listProductsWithTimeout() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    products, err := sf.Products().List(ctx,
        product.WithLimit(20),
        product.WithSortBy("created_at"),
        product.WithOrder("desc"),
    )
    if ctx.Err() == context.DeadlineExceeded {
        log.Fatal("Request timed out")
    }
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Found %d products\n", len(products))
}
```
