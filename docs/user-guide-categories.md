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
        categorySDK.WithLimit(50),
    )
    if err != nil {
        log.Fatal(err)
    }

    for _, cat := range categories {
        fmt.Printf("%s (%d products)\n", cat.Name, cat.ProductCount)
    }
}
```

**ListOptions Parameters:**

Use functional options like `WithLimit()` to set parameters. The category package uses pagination via limit and offset rather than page-based pagination.

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
    fmt.Printf("Description: %s\n", *cat.Description)
    fmt.Printf("Products: %d\n", cat.ProductCount)
    
    if cat.ParentID != nil {
        fmt.Printf("Parent Category ID: %s\n", *cat.ParentID)
    }
}
```

#### Category Structure

The `Category` type contains the following fields:

```go
type Category struct {
    ID           string   // Unique category identifier (e.g., "cat_electronics")
    Name         string   // Category display name
    Slug         string   // URL-friendly slug for routing
    Description  *string  // Optional detailed description
    ParentID     *string  // Parent category ID for hierarchical structures
    ImageURL     *string  // Category banner/product image CDN URL
    ProductCount int      // Number of products in this category
}
```

**Field Details:**

- **ID**: Used to reference the category in API calls and product filtering
- **Slug**: SEO-friendly identifier used in URLs (e.g., `/categories/electronics`)
- **ParentID**: Enables nested category hierarchies (e.g., Electronics → Audio → Headphones)
- **ProductCount**: Cached count of products in this category for efficient display

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
    fmt.Printf("Products in this category: %d\n", cat.ProductCount)

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

    // Check if it has a parent
    if headphonesCat.ParentID != nil {
        parentCat, err := sf.Categories().Get(context.Background(), *headphonesCat.ParentID)
        if err == nil {
            fmt.Printf("Parent Category: %s\n", parentCat.Name)
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
        categorySDK.WithLimit(50),
    )
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Total categories: %d\n", len(categories))

    // Get details for each category and its products
    for _, cat := range categories {
        fmt.Printf("\n%s (%d products)\n", cat.Name, cat.ProductCount)
        
        if cat.Description != nil {
            fmt.Printf("  Description: %s\n", *cat.Description)
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
