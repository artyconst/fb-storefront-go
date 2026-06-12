### Cart Service

The Cart service manages shopping cart operations, allowing customers to add items, update quantities, and proceed to checkout. Carts can be associated with customer accounts or used for guest checkout.

#### Get Existing Cart

Retrieve the current cart by ID:

```go
import (
    "context"
    "fmt"
    "log"

    "github.com/artyconst/fb-storefront-go/pkg/cart"
)

func getCart() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    // Get existing cart by ID
    cart, err := sf.Cart().Get(context.Background(), "cart_abc123")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Cart ID: %s\n", cart.ID)
    fmt.Printf("Status: %s\n", cart.Status)
    fmt.Printf("Total items: %d\n", len(cart.Items))
    fmt.Printf("Subtotal: $%d\n", cart.Subtotal)
    fmt.Printf("Total Amount: $%d\n", cart.TotalAmount)

    for _, item := range cart.Items {
        fmt.Printf("- %s x%d: $%d\n", item.Name, item.Quantity, item.Total)
    }
}
```

#### Cart Creation (Implicit via Get)

Carts are created implicitly when you call `Get()` with a cart identifier. If no cart exists for that ID, the server creates one automatically. There is no explicit `Create` method — this matches the Fleetbase Storefront API design where carts are identified by any unique string (e.g., device ID, session token).

```go
func getOrCreateCart() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    // Get cart — creates one implicitly if it doesn't exist for this ID
    cart, err := sf.Cart().Get(context.Background(), "cart_abc123")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Cart ID: %s\n", cart.ID)
    fmt.Printf("Status: %s\n", cart.Status)
}
```

#### Add Item to Cart

Add products to the shopping cart:

```go
func addToCart() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    cart, err := sf.Cart().AddItem(context.Background(), "cart_abc123", "prod_abc123", 2, nil, nil, "", "")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Cart total items: %d\n", len(cart.Items))
    fmt.Printf("Total amount: $%d\n", cart.TotalAmount)
}
```

**AddItem Parameters:**

| Parameter | Type | Description | Required |
|-----------|------|-------------|----------|
| `cartID` | string | Cart identifier | Yes |
| `productID` | string | Product identifier to add | Yes |
| `quantity` | int | Number of units to add | Yes |
| `addons` | []interface{} | Product addons (optional) | No |
| `variants` | map[string]any | Product variants (optional) | No |
| `scheduledAt` | string | Scheduled delivery time (optional) | No |
| `storeLocation` | string | Store location identifier (optional) | No |

#### Update Item

Modify quantities of existing cart items:

```go
func updateItem() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    // Increase quantity
    cart, err := sf.Cart().UpdateItem(context.Background(), "cart_abc123", "item_xyz789", 5, nil, nil)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Updated item quantity to: %d\n", 5)
    fmt.Printf("New total: $%d\n", cart.TotalAmount)
}
```

**UpdateItem Parameters:**

| Parameter | Type | Description | Required |
|-----------|------|-------------|----------|
| `cartID` | string | Cart identifier | Yes |
| `lineItemID` | string | Cart item ID to update | Yes |
| `quantity` | int | New quantity (must be >= 1) | Yes |
| `addons` | []interface{} | Product addons (optional) | No |
| `variants` | map[string]any | Product variants (optional) | No |

#### Remove Item from Cart

Remove a specific item from the cart:

```go
func removeItem() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    cart, err := sf.Cart().RemoveItem(context.Background(), "cart_abc123", "item_xyz789")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Remaining items: %d\n", len(cart.Items))
    fmt.Printf("Updated total: $%d\n", cart.TotalAmount)
}
```

**RemoveItem Parameters:**

| Parameter | Type | Description | Required |
|-----------|------|-------------|----------|
| `cartID` | string | Cart identifier | Yes |
| `itemID` | string | Cart item ID to remove | Yes |

#### Clear Cart

Remove all items from cart but keep the cart itself:

```go
func clearCart() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    err = sf.Cart().Clear(context.Background(), "cart_abc123")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Cart cleared successfully")
}
```

#### Checkout Cart

Complete the purchase and create an order from cart:

```go
func checkout() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    checkoutReq := cart.CheckoutRequest{
        CustomerEmail: "customer@example.com",
        ShippingAddress: &cart.Address{
            FirstName:    "John",
            LastName:     "Doe",
            AddressLine1: "123 Main St",
            City:         "San Francisco",
            State:        "CA",
            PostalCode:   "94105",
            Country:      "US",
            Phone:        "+14155551234",
        },
        PaymentMethodID: "pm_card_visa_ending_4242",
    }

    order, err := sf.Cart().Checkout(context.Background(), "cart_abc123", checkoutReq)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Order created: %s\n", order.OrderNumber)
    fmt.Printf("Total paid: $%d\n", order.TotalAmount)
}
```

**CheckoutRequest Parameters:**

| Parameter | Type | Description | Required |
|-----------|------|-------------|----------|
| `CustomerEmail` | string | Customer email address | No |
| `ShippingAddress` | *Address | Shipping address details | No |
| `BillingAddress` | *Address | Billing address details (defaults to shipping) | No |
| `PaymentMethodID` | string | Payment method token from payment processor | No |
| `Metadata` | map[string]interface{} | Additional metadata (optional) | No |

**Address Structure:**

```go
type Address struct {
    FirstName    string `json:"first_name"`     // Recipient first name
    LastName     string `json:"last_name"`      // Recipient last name
    Company      string `json:"company,omitempty"`  // Company name (optional)
    AddressLine1 string `json:"address_line_1"` // Street address line 1
    AddressLine2 string `json:"address_line_2,omitempty"` // Apt, suite, etc. (optional)
    City         string `json:"city"`           // City or locality
    State        string `json:"state"`          // State or province
    PostalCode   string `json:"postal_code"`    // ZIP or postal code
    Country      string `json:"country"`        // ISO 3166-1 alpha-2 country code (e.g., "US")
    Phone        string `json:"phone,omitempty"` // Contact phone number (optional)
}
```

#### Cart Structure

The `Cart` type represents the shopping cart:

```go
type Cart struct {
    ID          string      `json:"id"`
    Status      CartStatus  `json:"status"`
    CustomerID  *string     `json:"customer_id,omitempty"`
    Items       []*CartItem `json:"items"`
    Subtotal    int64       `json:"subtotal,omitempty"`
    TaxAmount   *int64      `json:"tax_amount,omitempty"`
    Discount    *int64      `json:"discount,omitempty"`
    TotalAmount int64       `json:"total_amount"`
    Currency    string      `json:"currency"`
    CreatedAt   string      `json:"created_at"`
    UpdatedAt   string      `json:"updated_at"`
}

type CartItem struct {
    ID        string `json:"id"`
    ProductID string `json:"product_id"`
    Name      string `json:"name"`
    Quantity  int    `json:"quantity"`
    Price     int64  `json:"price"`
    Total     int64  `json:"total"`
}

type CartStatus string

const (
    CartStatusActive       CartStatus = "active"        // Cart has items, not completed
    CartStatusCompleted    CartStatus = "completed"     // Cart converted to order
    CartStatusAbandoned    CartStatus = "abandoned"     // Customer abandoned cart
)
```

#### Cart Mutations (Functional Options)

The cart service provides a set of mutation methods that use the functional options pattern. These are parallel implementations to the existing methods, offering more flexibility through optional parameters passed as functions rather than positional arguments.

##### AddProduct

Add a product using functional options instead of positional parameters:

```go
func addProductWithOptions() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    // Add a product using functional options
    cart, err := sf.Cart().AddProduct(context.Background(), "cart_abc123", "prod_01J...",
        cart.WithQuantity(2),
        cart.WithVariant("color", "red"),
    )
    if err != nil {
        log.Fatalf("Failed to add product: %v", err)
    }

    fmt.Printf("Cart total items: %d\n", len(cart.Items))
}
```

**Available AddProduct Options:**

| Option | Type | Description |
|--------|------|-------------|
| `WithQuantity(int64)` | int64 | Number of units to add |
| `WithVariant(string, any)` | key-value pair | Product variant (can be called multiple times) |
| `WithScheduledAt(string)` | string | Scheduled delivery time |
| `WithStoreLocation(string)` | string | Store location identifier |

##### UpdateLineItem

Update a line item using functional options instead of a quantity parameter:

```go
func updateLineItemWithOptions() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    // Update a line item using functional options
    cart, err := sf.Cart().UpdateLineItem(context.Background(), "cart_abc123", "line_01J...",
        cart.WithQuantityForUpdate(5),
    )
    if err != nil {
        log.Fatalf("Failed to update line item: %v", err)
    }

    fmt.Printf("Updated total: $%d\n", cart.TotalAmount)
}
```

**Available UpdateLineItem Options:**

| Option | Type | Description |
|--------|------|-------------|
| `WithQuantityForUpdate(int64)` | int64 | New quantity (must be >= 1) |
| `WithVariantForUpdate(string, any)` | key-value pair | Product variant update (can be called multiple times) |

##### RemoveLineItem

Remove a line item — returns the updated cart:

```go
func removeLineItem() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    // Remove a line item — returns the updated cart
    cart, err := sf.Cart().RemoveLineItem(context.Background(), "cart_abc123", "line_01J...")
    if err != nil {
        log.Fatalf("Failed to remove line item: %v", err)
    }

    fmt.Printf("Cart now has %d items\n", len(cart.Items))
}
```

##### EmptyCart

Empty the cart — returns the updated cart with zero items:

```go
func emptyCart() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    // Empty the cart — returns the updated cart with zero items
    cart, err := sf.Cart().EmptyCart(context.Background(), "cart_abc123")
    if err != nil {
        log.Fatalf("Failed to empty cart: %v", err)
    }

    fmt.Printf("Cart emptied: %d items remaining\n", len(cart.Items))
}
```

##### DeleteCart

Permanently delete the entire cart from the server:

```go
func deleteCart() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    // Permanently delete the cart from the server
    err := sf.Cart().DeleteCart(context.Background())
    if err != nil {
        log.Fatalf("Failed to delete cart: %v", err)
    }

    fmt.Println("Cart deleted")
}
```

##### Old vs. New Method Comparison

The following table compares the original methods with their functional-options counterparts, so you can choose which API style fits your use case:

| Old Method | New Method | Key Difference |
|---|---|---|
| `AddItem(ctx, cartID, productID, quantity)` | `AddProduct(ctx, cartID, productID, opts...)` | Quantity via option instead of parameter |
| `UpdateItem(ctx, cartID, lineItemID, quantity)` | `UpdateLineItem(ctx, cartID, lineItemID, opts...)` | Quantity via option instead of parameter |
| `RemoveItem(ctx, cartID, lineItemID) (*Cart, error)` | `RemoveLineItem(ctx, cartID, lineItemID) (*Cart, error)` | Input validation; returns updated cart |
| `Clear(ctx, cartID) error` | `EmptyCart(ctx, cartID) (*Cart, error)` | Returns updated cart instead of just error |

#### Complete Cart Workflow Example

```go
func completeCartWorkflow() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    // Step 1: Get or create cart (creates implicitly if it doesn't exist)
    cart, err := sf.Cart().Get(context.Background(), "cart_abc123")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Cart ID: %s\n", cart.ID)

    // Step 2: Add first product
    cart, err = sf.Cart().AddItem(context.Background(), cart.ID, "prod_wireless_headphones", 1, nil, nil, "", "")
    if err != nil {
        log.Fatal(err)
    }

    // Step 3: Add second product
    cart, err = sf.Cart().AddItem(context.Background(), cart.ID, "prod_phone_case", 2, nil, nil, "", "")
    if err != nil {
        log.Fatal(err)
    }

    // Step 4: Update quantity of first item
    cart, err = sf.Cart().UpdateItem(context.Background(), cart.ID, cart.Items[0].ID, 3, nil, nil)
    if err != nil {
        log.Fatal(err)
    }

    // Step 5: Display cart contents
    fmt.Printf("Cart Summary:\n")
    for _, item := range cart.Items {
        fmt.Printf("- %s x%d @ $%d = $%d\n", 
            item.Name, item.Quantity, item.Price, item.Total)
    }
    fmt.Printf("Subtotal: $%d\n", cart.Subtotal)
    fmt.Printf("Total: $%d\n", cart.TotalAmount)

    // Step 6: Remove one item
    cart, err = sf.Cart().RemoveItem(context.Background(), cart.ID, cart.Items[1].ID)
    if err != nil {
        log.Fatal(err)
    }

    // Step 7: Checkout with shipping and payment details
    checkoutReq := cart.CheckoutRequest{
        CustomerEmail: "john.doe@example.com",
        ShippingAddress: &cart.Address{
            FirstName:    "John",
            LastName:     "Doe",
            AddressLine1: "123 Main St",
            City:         "San Francisco",
            State:        "CA",
            PostalCode:   "94105",
            Country:      "US",
            Phone:        "+14155551234",
        },
        PaymentMethodID: "pm_card_visa_ending_4242",
    }

    order, err := sf.Cart().Checkout(context.Background(), cart.ID, checkoutReq)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Order #%s created successfully!\n", order.OrderNumber)
}
```

#### Error Handling

```go
import (
    "errors"
    "log"

    "github.com/artyconst/fb-storefront-go"
)

func handleCartErrors() {
    // Cart not found
    _, err := sf.Cart().Get(context.Background(), "invalid-cart-id")
    if err != nil {
        if errors.Is(err, storefront.ErrResourceNotFound) {
            log.Println("Cart not found. Create a new cart first.")
        } else {
            log.Printf("API Error: %v", err)
        }
    }

    // Invalid quantity (must be >= 1)
    _, err = sf.Cart().UpdateItem(context.Background(), "cart_abc", "item_xyz", 0, nil, nil)
    if err != nil {
        log.Printf("Invalid quantity error: %v", err)
    }
}
```
