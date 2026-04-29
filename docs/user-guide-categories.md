### Categories Service

The Categories service allows you to organize and retrieve product categories for your store. Categories help structure your product catalog and enable customers to browse products by classification.

#### List All Categories

Get all available categories with optional pagination:

```go
import (
    "context"
    "fmt"
    "log"

    "github.com/artyconst/fb-storefront-go/pkg/category"
)

func listCategories() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    categories, err := sf.Categories().List(context.Background(),
        category.WithLimit(50),
    )
    if err != nil {
        log.Fatal(err)
    }

    for _, cat := range categories {
        fmt.Printf("%s\n", cat.Name)
    }
}
```

**ListOptions Parameters:**

The category package uses functional options for all list parameters. Pagination uses limit and offset rather than page-based pagination.

| Option | Type | Description |
|--------|------|-------------|
| `WithLimit(limit int64)` | `int64` | Maximum number of categories to return |
| `WithOffset(offset int64)` | `int64` | Offset for pagination (skip N categories) |
| `WithSearch(search string)` | `string` | Search categories by name or description |
| `WithParentID(parentID string)` | `string` | Filter categories by parent ID (for nested categories) |

**Example with all options:**

```go
categories, err := sf.Categories().List(context.Background(),
    category.WithLimit(50),
    category.WithOffset(0),
    category.WithSearch("electronics"),
    category.WithParentID("cat_electronics"),
)
```

#### Get Category Details

Retrieve a specific category with metadata including product count and hierarchy information:

```go
func getCategory() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    cat, err := sf.Categories().Get(context.Background(), "cat_electronics")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("%s - %s\n", cat.Name, cat.Slug)
    if cat.Description != nil {
        fmt.Printf("Description: %s\n", *cat.Description)
    }
    if cat.IconURL != "" {
        fmt.Printf("Icon URL: %s\n", cat.IconURL)
    }
}
```

#### Category Structure

The `Category` type contains the following fields:

```go
type Category struct {
    ID           string                 // Unique category identifier (e.g., "cat_electronics")
    Name         string                 // Category display name
    Description  *string                // Optional detailed description (nil-safe: check before dereferencing)
    IconURL      string                 // Category icon image CDN URL
    Tags         []string               // Category tags
    Translations []interface{}          // Translations
    Meta         map[string]interface{} // Metadata
    Order        *int64                 // Optional display order
    Slug         string                 // URL-friendly slug for routing
    CreatedAt    time.Time              // Creation timestamp
    UpdatedAt    *time.Time             // Last update timestamp (nil-safe: check before use)
}
```

**Field Details:**

- **ID**: Used to reference the category in API calls and product filtering
- **Slug**: SEO-friendly identifier used in URLs (e.g., `/categories/electronics`)
- **IconURL**: Direct string value (not a pointer). Use `if cat.IconURL != ""` to check for presence.
- **Description**: Pointer to string. Always nil-check before dereferencing: `if cat.Description != nil { fmt.Println(*cat.Description) }`
- **Order**: Optional pointer to int64. Nil-check before use.
- **UpdatedAt**: Optional pointer to time.Time. Nil-check before use.

#### Using Categories with Products

Categories are commonly used to filter product listings:

```go
func browseCategoryWithProducts() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    // First get the category details
    cat, err := sf.Categories().Get(context.Background(), "cat_electronics")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Category: %s\n", cat.Name)

    // Then get products filtered by this category using functional options
    products, err := sf.Products().List(context.Background(),
        product.WithCategory(cat.ID),
        product.WithLimit(20),
        product.WithSortBy("name"),
        product.WithOrder("asc"),
    )
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Retrieved %d products\n", len(products))
}
```

#### Hierarchical Categories Example

Categories support parent-child relationships for multi-level hierarchies:

```go
func navigateCategoryHierarchy() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    // Get subcategory
    headphonesCat, err := sf.Categories().Get(context.Background(), "cat_headphones")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Category: %s\n", headphonesCat.Name)

    // Use WithParentID to find parent categories
    // (The Category struct does not include a ParentID field)
    parentCategories, err := sf.Categories().List(context.Background(),
        category.WithParentID("cat_electronics"),
    )
    if err == nil {
        for _, pc := range parentCategories {
            fmt.Printf("Parent: %s\n", pc.Name)
        }
    }

    // Get products in this subcategory using functional options
    products, err := sf.Products().List(context.Background(),
        product.WithCategory(headphonesCat.ID),
        product.WithLimit(20),
    )
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Found %d headphone products\n", len(products))
}
```

#### Error Handling

```go
import (
    "errors"
    "log"

    "github.com/artyconst/fb-storefront-go"
)

func handleCategoryErrors() {
    _, err := sf.Categories().Get(context.Background(), "invalid-category-id")
    if err != nil {
        if errors.Is(err, storefront.ErrResourceNotFound) {
            log.Println("Category not found")
        } else {
            log.Printf("API Error: %v", err)
        }
    }
}
```

#### Complete Example with All Operations

```go
func completeCategoryWorkflow() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    // List all top-level categories using functional options
    categories, err := sf.Categories().List(context.Background(),
        category.WithLimit(50),
    )
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Total categories: %d\n", len(categories))

    // Get details for each category and its products
    for _, cat := range categories {
        fmt.Printf("\n%s\n", cat.Name)
        
        if cat.Description != nil {
            fmt.Printf("  Description: %s\n", *cat.Description)
        }
        if cat.IconURL != "" {
            fmt.Printf("  Icon URL: %s\n", cat.IconURL)
        }

        // Get products in this category using functional options
        products, err := sf.Products().List(context.Background(),
            product.WithCategory(cat.ID),
            product.WithLimit(5),
        )
        if err != nil {
            log.Printf("Error fetching products for %s: %v", cat.Name, err)
            continue
        }

        fmt.Printf("  Sample products:\n")
        for _, p := range products {
            fmt.Printf("    - %s ($%s)\n", p.Name, p.Price.String())
        }
    }
}
```
