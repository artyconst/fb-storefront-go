### Orders Service

The Orders service allows you to view and manage customer orders. Orders are created when customers complete checkout from their shopping carts.

#### List Orders

Retrieve order history with optional filtering by status:

```go
import (
    "context"
    "fmt"
    "log"

    "github.com/artyconst/fb-storefront-go/pkg/order"
)

func listOrders() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    opts := &order.ListOptions{
        Page:   1,
        Limit:  20,
        Status: order.OrderStatusProcessing, // Filter by status
    }

    orders, err := sf.Orders().List(context.Background(), opts)
    if err != nil {
        log.Fatal(err)
    }

    for _, ord := range orders {
        fmt.Printf("#%s - $%s - %s\n", 
            ord.OrderNumber, 
            ord.TotalAmount.String(), 
            ord.Status)
    }
}
```

**ListOptions Parameters:**

| Parameter | Type | Description | Default |
|-----------|------|-------------|---------|
| `Page` | int | Pagination page number | 1 |
| `Limit` | int | Items per page (default: 20) | 20 |
| `Status` | string | Filter by order status | - |

#### Get Order Details

Retrieve specific order information including line items and shipping details:

```go
func getOrder() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    ord, err := sf.Orders().Get(context.Background(), "ord_abc123")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Order #%s - %s\n", ord.OrderNumber, ord.Status)
    fmt.Printf("Total Amount: $%s\n", ord.TotalAmount.String())
    fmt.Printf("Items: %d\n", len(ord.Items))

    for _, item := range ord.Items {
        fmt.Printf("- %s x%d: $%s\n", 
            item.Name, 
            item.Quantity, 
            item.Price.String())
    }

    if ord.Customer != nil && ord.Customer.Email != nil {
        fmt.Printf("Customer Email: %s\n", *ord.Customer.Email)
    }

    if ord.ShippingAddress != nil {
        addr := ord.ShippingAddress
        fmt.Printf("Shipping Address:\n")
        fmt.Printf("  %s %s\n", addr.FirstName, addr.LastName)
        fmt.Printf("  %s\n", addr.AddressLine1)
        fmt.Printf("  %s, %s %s\n", addr.City, addr.State, addr.PostalCode)
        fmt.Printf("  %s\n", addr.Country)
    }
}
```

#### Order Status Constants

Use predefined constants for filtering by status:

```go
// Filter orders by specific status
opts := &order.ListOptions{
    Status: order.OrderStatusProcessing, // Get processing orders
}

availableStatuses := []string{
    order.OrderStatusPending,      // Pending orders - awaiting payment or confirmation
    order.OrderStatusConfirmed,    // Confirmed orders - payment received, being prepared
    order.OrderStatusProcessing,   // Processing orders - being packed/shipped
    order.OrderStatusShipped,      // Shipped orders - en route to customer
    order.OrderStatusDelivered,    // Delivered orders - successfully delivered
    order.OrderStatusCancelled,    // Cancelled orders - cancelled by customer or system
    order.OrderStatusRefunded,     // Refunded orders - payment refunded
}
```

**Order Status Values:**

| Constant | Value | Description |
|----------|-------|-------------|
| `OrderStatusPending` | "pending" | Awaiting payment confirmation |
| `OrderStatusConfirmed` | "confirmed" | Payment confirmed, processing |
| `OrderStatusProcessing` | "processing" | Being prepared for shipment |
| `OrderStatusShipped` | "shipped" | Shipped with tracking available |
| `OrderStatusDelivered` | "delivered" | Successfully delivered to customer |
| `OrderStatusCancelled` | "cancelled" | Order cancelled before completion |
| `OrderStatusRefunded` | "refunded" | Refund processed for this order |

#### Orders Service Methods

The Orders service provides the following methods:

| Method | Description | Parameters | Returns |
|--------|-------------|------------|---------|
| `List(ctx, opts)` | Get paginated list of orders | ListOptions with filters | Array of Order objects |
| `Get(ctx, orderID)` | Retrieve specific order by ID | Order ID string | Single Order object |

#### Order Structure

The `Order` type contains the following fields:

```go
type Order struct {
    ID             string      // Unique order identifier (e.g., "ord_abc123")
    OrderNumber    string      // Human-readable order number (e.g., "ORD-001234")
    Status         string      // Current order status (pending, confirmed, processing, shipped, delivered, cancelled, refunded)
    Customer       *OrderCustomer  // Customer information for this order
    Items          []*OrderItem   // Array of ordered items
    Subtotal       Decimal     // Subtotal before tax and discounts
    TaxAmount      *Decimal    // Calculated tax amount
    ShippingCost   *Decimal    // Shipping cost if applicable
    Discount       *Decimal    // Applied discount amount
    TotalAmount    Decimal     // Final total amount charged
    Currency       string      // Currency code (e.g., "USD")
    ShippingAddress  *Address   // Complete shipping address
    BillingAddress   *Address    // Billing address (may differ from shipping)
    PaymentMethodID  *string    // Payment method used for this order
    CreatedAt      time.Time   // When order was created
    UpdatedAt      time.Time   // Last update timestamp
}

type OrderCustomer struct {
    ID        *string `json:"id,omitempty"`     // Customer account ID (nil for guest checkout)
    Email     *string `json:"email,omitempty"`  // Customer email address
    Phone     *string `json:"phone,omitempty"`  // Customer phone number
}

type OrderItem struct {
    ID        string  // Unique order item identifier
    ProductID string  // Related product identifier
    Name      string  // Product name (snapshot at time of order)
    Quantity  int     // Ordered quantity
    Price     Decimal // Unit price at time of purchase
    Total     Decimal // Line total (quantity × unit price)
}

type Address struct {
    FirstName     string `json:"first_name"`      // Recipient first name
    LastName      string `json:"last_name"`       // Recipient last name
    Company       string `json:"company,omitempty"`  // Company name (optional)
    AddressLine1  string `json:"address_line_1"`  // Street address line 1
    AddressLine2  string `json:"address_line_2,omitempty"` // Apt, suite, unit (optional)
    City          string `json:"city"`            // City or locality
    State         string `json:"state"`           // State or province
    PostalCode    string `json:"postal_code"`     // ZIP or postal code
    Country       string `json:"country"`         // ISO 3166-1 alpha-2 country code (e.g., "US")
    Phone         string `json:"phone,omitempty"` // Contact phone number (optional)
}
```

#### Complete Order Management Example

This example demonstrates listing orders with filters and retrieving detailed information:

```go
func manageOrders() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    // Step 1: List all pending orders
    fmt.Println("=== Pending Orders ===")
    pendingOpts := &order.ListOptions{
        Status: order.OrderStatusPending,
        Limit:  20,
        Page:   1,
    }

    pendingOrders, err := sf.Orders().List(context.Background(), pendingOpts)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Found %d pending orders\n", len(pendingOrders))
    for _, ord := range pendingOrders {
        fmt.Printf("#%s - $%s - #%s\n", 
            ord.OrderNumber, 
            ord.TotalAmount.String(),
            ord.ID)
    }

    // Step 2: List all shipped orders
    fmt.Println("\n=== Shipped Orders ===")
    shippedOpts := &order.ListOptions{
        Status: order.OrderStatusShipped,
        Limit:  20,
        Page:   1,
    }

    shippedOrders, err := sf.Orders().List(context.Background(), shippedOpts)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Found %d shipped orders\n", len(shippedOrders))

    // Step 3: Get detailed information for first pending order
    if len(pendingOrders) > 0 {
        firstOrder := pendingOrders[0]
        
        ordered, err := sf.Orders().Get(context.Background(), firstOrder.ID)
        if err != nil {
            log.Fatal(err)
        }

        fmt.Printf("\n=== Order Details: #%s ===\n", ordered.OrderNumber)
        fmt.Printf("Status: %s\n", ordered.Status)
        fmt.Printf("Total: $%s\n", ordered.TotalAmount.String())
        
        // Display order items
        fmt.Println("Items:")
        for _, item := range ordered.Items {
            fmt.Printf("  - %s x%d @ $%s = $%s\n", 
                item.Name, 
                item.Quantity, 
                item.Price.String(), 
                item.Total.String())
        }

        // Display shipping address
        if ordered.ShippingAddress != nil {
            addr := ordered.ShippingAddress
            fmt.Printf("\nShipping to:\n")
            fmt.Printf("  %s %s\n", addr.FirstName, addr.LastName)
            fmt.Printf("  %s %s\n", addr.AddressLine1, addr.City)
            fmt.Printf("  %s, %s %s\n", addr.State, addr.PostalCode, addr.Country)
        }

        // Display customer information if available
        if ordered.Customer != nil && ordered.Customer.Email != nil {
            fmt.Printf("\nCustomer: %s\n", *ordered.Customer.Email)
            if ordered.Customer.Phone != nil {
                fmt.Printf("Phone: %s\n", *ordered.Customer.Phone)
            }
        }
    }
}
```

#### Order Number Format

Order numbers are human-readable identifiers that differ from the internal order ID:

```go
// Internal ID (used in API calls): "ord_abc123xyz"
orderID := ord.ID  // e.g., "ord_abc123xyz"

// Human-readable number (displayed to customers): "ORD-001234"
orderNumber := ord.OrderNumber  // e.g., "ORD-001234"
```

#### Error Handling

```go
import (
    "errors"
    "log"

    "github.com/artyconst/fb-storefront-go"
)

func handleOrderErrors() {
    // Order not found
    _, err := sf.Orders().Get(context.Background(), "invalid-order-id")
    if err != nil {
        if errors.Is(err, storefront.ErrResourceNotFound) {
            log.Println("Order not found")
        } else {
            log.Printf("API Error: %v", err)
        }
    }

    // Invalid status filter (use constants instead of arbitrary strings)
    opts := &order.ListOptions{
        Status: "invalid_status",  // This may cause unexpected behavior
    }
    
    orders, err := sf.Orders().List(context.Background(), opts)
    if err != nil {
        log.Printf("Filter error: %v", err)
    }

    // Use proper constants
    validOpts := &order.ListOptions{
        Status: order.OrderStatusDelivered,  // Always use constants
    }
    
    orders, err = sf.Orders().List(context.Background(), validOpts)
}
```

#### Order Workflow Context

Orders are created automatically when customers complete checkout from carts. The Orders service is primarily used for viewing and managing existing orders:

```go
func orderCreationFlow() {
    // Step 1: Customer adds items to cart (using Cart service)
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    cart, err := sf.Cart().Get(context.Background(), "cart_abc123")
    if err != nil {
        log.Fatal(err)
    }

    // Step 2: Customer proceeds to checkout (using Cart.Checkout method)
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
        },
        PaymentMethodID: "pm_card_visa_ending_4242",
    }

    // This creates the order automatically
    order, err := sf.Cart().Checkout(context.Background(), cart.ID, *checkoutReq)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Order created: #%s\n", order.OrderNumber)

    // Step 3: Retrieve order using Orders service
    retrievedOrder, err := sf.Orders().Get(context.Background(), order.ID)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Retrieved order # %s - Status: %s\n", 
        retrievedOrder.OrderNumber, 
        retrievedOrder.Status)
}
```
