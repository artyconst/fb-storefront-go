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
        fmt.Printf("%s - $%s\n", p.Name, p.Price.String())
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

    fmt.Printf("%s - $%s\n", prod.Name, prod.Price.String())
    fmt.Printf("SKU: %s\n", prod.SKU)
    fmt.Printf("Stock: %d\n", prod.StockQuantity)
}
```

#### Search Products

Search for products by query string with optional category filtering:

```go
func searchProducts() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    // Note: Search functionality requires SearchOptions which uses functional options
    products, err := sf.Products().Search(context.Background(), product.SearchQuery{
        Query:      "wireless headphones",
        CategoryID: "cat_electronics", // Optional category filter
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
| `CategoryID` | string | Optional category filter | No |

Note: Pagination parameters are applied via functional options after the search.

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
    ID             string     // Unique product identifier (e.g., "prod_abc123")
    Name           string     // Product name
    Description    *string    // Optional description
    Price          Decimal    // Current price as decimal string
    CompareAtPrice *Decimal   // Original/comparison price for discounts
    SKU            string     // Stock keeping unit identifier
    StockQuantity  int        // Available quantity in stock
    Images         []Image    // Product images array
    Categories     []Category // Associated categories
}

type Image struct {
    ID       string `json:"id"`       // Unique image identifier
    URL      string `json:"url"`      // Image CDN URL
    AltText  string `json:"alt_text,omitempty"` // Optional alt text for accessibility
    Position int    `json:"position"` // Display order position
}

type Category struct {
    ID   string `json:"id"`   // Category identifier
    Name string `json:"name"` // Category name
}
```

**Important Notes:**

- **Price Handling**: Use the `Decimal` type (string alias) for all currency values to maintain precision and avoid floating-point errors. Always use `.String()` when displaying prices.
- **Optional Fields**: Pointer types (`*string`, `*Decimal`) indicate optional fields that may be nil.
- **Images Array**: Products can have multiple images with positions determining display order.

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
