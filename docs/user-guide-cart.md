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
    fmt.Printf("Subtotal: $%s\n", cart.Subtotal.String())
    fmt.Printf("Total Amount: $%s\n", cart.TotalAmount.String())

    for _, item := range cart.Items {
        fmt.Printf("- %s x%d: $%s\n", item.Name, item.Quantity, item.Total.String())
    }
}
```

#### Create New Cart

Create a new empty cart associated with a customer:

```go
func createCart() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    // Create cart for existing customer
    cart, err := sf.Cart().Create(context.Background(), "cust_12345")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Created cart: %s\n", cart.ID)
    fmt.Printf("Customer ID: %s\n", *cart.CustomerID)
    fmt.Printf("Status: %s\n", cart.Status)
}
```

**Create Parameters:**

| Parameter | Type | Description | Required |
|-----------|------|-------------|----------|
| `customerID` | string | Customer ID to associate with cart | Yes |

#### Add Item to Cart

Add products to the shopping cart:

```go
func addToCart() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    itemReq := &cart.CartItemRequest{
        ProductID: "prod_abc123",
        Quantity:  2,
    }

    cart, err := sf.Cart().AddItem(context.Background(), "cart_abc123", *itemReq)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Cart total items: %d\n", len(cart.Items))
    fmt.Printf("Total amount: $%s\n", cart.TotalAmount.String())
}
```

**CartItemRequest Parameters:**

| Parameter | Type | Description | Required |
|-----------|------|-------------|----------|
| `ProductID` | string | Product identifier to add | Yes |
| `Quantity` | int | Number of units to add | Yes |

#### Update Item Quantity

Modify quantities of existing cart items:

```go
func updateQuantity() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    // Increase quantity
    cart, err := sf.Cart().UpdateQuantity(context.Background(), "cart_abc123", "item_xyz789", 5)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Updated item quantity to: %d\n", 5)
    fmt.Printf("New total: $%s\n", cart.TotalAmount.String())
}
```

**UpdateQuantity Parameters:**

| Parameter | Type | Description | Required |
|-----------|------|-------------|----------|
| `cartID` | string | Cart identifier | Yes |
| `itemID` | string | Cart item ID to update | Yes |
| `quantity` | int | New quantity (must be >= 1) | Yes |

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
    fmt.Printf("Updated total: $%s\n", cart.TotalAmount.String())
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

    checkoutReq := &cart.CheckoutRequest{
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

    order, err := sf.Cart().Checkout(context.Background(), "cart_abc123", *checkoutReq)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Order created: %s\n", order.OrderNumber)
    fmt.Printf("Total paid: $%s\n", order.TotalAmount.String())
}
```

**CheckoutRequest Parameters:**

| Parameter | Type | Description | Required |
|-----------|------|-------------|----------|
| `CustomerEmail` | string | Customer email address | Yes |
| `ShippingAddress` | *Address | Shipping address details | Yes |
| `PaymentMethodID` | string | Payment method token from payment processor | Yes |

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
    ID          string      // Unique cart identifier (e.g., "cart_abc123")
    Status      CartStatus  // Current status: active, completed, or abandoned
    CustomerID  *string     // Associated customer ID (nil for guest carts)
    Items       []*CartItem // Array of cart items
    Subtotal    Decimal    // Subtotal before tax and discounts
    TaxAmount   *Decimal    // Optional calculated tax amount
    Discount    *Decimal    // Optional discount amount applied
    TotalAmount Decimal     // Final total (subtotal + tax - discount)
    Currency    string      // Currency code (e.g., "USD", "EUR")
}

type CartItem struct {
    ID        string  // Unique cart item identifier
    ProductID string  // Related product identifier
    Name      string  // Product name (snapshot at time of adding to cart)
    Quantity  int     // Quantity ordered
    Price     Decimal // Unit price per item
    Total     Decimal // Line total (quantity × unit price)
}

type CartStatus string

const (
    CartStatusActive       CartStatus = "active"        // Cart has items, not completed
    CartStatusCompleted    CartStatus = "completed"     // Cart converted to order
    CartStatusAbandoned    CartStatus = "abandoned"     // Customer abandoned cart
)
```

#### Complete Cart Workflow Example

```go
func completeCartWorkflow() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    // Step 1: Create new cart for customer
    cart, err := sf.Cart().Create(context.Background(), "cust_12345")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Created cart: %s\n", cart.ID)

    // Step 2: Add first product
    item1 := &cart.CartItemRequest{
        ProductID: "prod_wireless_headphones",
        Quantity:  1,
    }
    cart, err = sf.Cart().AddItem(context.Background(), cart.ID, *item1)
    if err != nil {
        log.Fatal(err)
    }

    // Step 3: Add second product
    item2 := &cart.CartItemRequest{
        ProductID: "prod_phone_case",
        Quantity:  2,
    }
    cart, err = sf.Cart().AddItem(context.Background(), cart.ID, *item2)
    if err != nil {
        log.Fatal(err)
    }

    // Step 4: Update quantity of first item
    cart, err = sf.Cart().UpdateQuantity(context.Background(), cart.ID, cart.Items[0].ID, 3)
    if err != nil {
        log.Fatal(err)
    }

    // Step 5: Display cart contents
    fmt.Printf("Cart Summary:\n")
    for _, item := range cart.Items {
        fmt.Printf("- %s x%d @ $%s = $%s\n", 
            item.Name, item.Quantity, item.Price.String(), item.Total.String())
    }
    fmt.Printf("Subtotal: $%s\n", cart.Subtotal.String())
    fmt.Printf("Total: $%s\n", cart.TotalAmount.String())

    // Step 6: Remove one item
    cart, err = sf.Cart().RemoveItem(context.Background(), cart.ID, cart.Items[1].ID)
    if err != nil {
        log.Fatal(err)
    }

    // Step 7: Checkout with shipping and payment details
    checkoutReq := &cart.CheckoutRequest{
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

    order, err := sf.Cart().Checkout(context.Background(), cart.ID, *checkoutReq)
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
    _, err = sf.Cart().UpdateQuantity(context.Background(), "cart_abc", "item_xyz", 0)
    if err != nil {
        log.Printf("Invalid quantity error: %v", err)
    }
}
```
