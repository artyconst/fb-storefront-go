### Checkout Service

The Checkout service manages checkout sessions, allowing you to initialize checkout processes, update customer information, and process payments before order creation.

#### Create Checkout Session

Initialize a checkout session from an active cart:

```go
import (
    "context"
    "fmt"
    "log"

    "github.com/artyconst/fb-storefront-go/pkg/checkout"
)

func createCheckout() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    checkoutReq := &checkout.CreateCheckoutRequest{
        CustomerEmail: "customer@example.com",
        ShippingAddress: &checkout.Address{
            FirstName:     "John",
            LastName:      "Doe",
            AddressLine1:  "123 Main St",
            City:          "San Francisco",
            State:         "CA",
            PostalCode:    "94105",
            Country:       "US",
        },
        PaymentMethodID: "pm_card_visa_ending_4242",
    }

    checkout, err := sf.Checkout().Create(context.Background(), "cart_abc123", *checkoutReq)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Checkout ID: %s\n", checkout.ID)
    fmt.Printf("Status: %s\n", checkout.Status)
    fmt.Printf("Amount: $%s\n", checkout.Amount.String())
}
```

**CreateCheckoutRequest Parameters:**

| Parameter | Type | Description | Required |
|-----------|------|-------------|----------|
| `CustomerEmail` | string | Customer email address for the order | Yes |
| `ShippingAddress` | *Address | Complete shipping address details | Yes |
| `PaymentMethodID` | string | Payment method token from payment processor | Yes (for immediate processing) |

**Create Parameters:**

| Parameter | Type | Description | Required |
|-----------|------|-------------|----------|
| `cartID` | string | Active cart ID to convert to checkout | Yes |
| `request` | *CreateCheckoutRequest | Checkout initialization data | Yes |

#### Get Checkout Details

Retrieve a checkout session by ID:

```go
func getCheckout() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    checkout, err := sf.Checkout().Get(context.Background(), "checkout_xyz789")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Status: %s\n", checkout.Status)
    fmt.Printf("Amount: $%s\n", checkout.Amount.String())
    
    if checkout.Customer != nil && checkout.Customer.Email != nil {
        fmt.Printf("Customer Email: %s\n", *checkout.Customer.Email)
    }

    if checkout.ShippingAddress != nil {
        addr := checkout.ShippingAddress
        fmt.Printf("Ship to: %s %s\n", addr.FirstName, addr.LastName)
        fmt.Printf("Address: %s, %s %s\n", addr.AddressLine1, addr.City, addr.PostalCode)
    }
}
```

#### Update Checkout Customer

Update customer information during checkout before payment:

```go
func updateCheckoutCustomer() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    customerInfo := &checkout.CustomerInfo{
        Email: "newemail@example.com",
        Phone: "+14155559876",
    }

    checkout, err := sf.Checkout().UpdateCustomer(context.Background(), "checkout_xyz789", *customerInfo)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Updated customer email: %s\n", *checkout.Customer.Email)
}
```

**CustomerInfo Parameters:**

| Parameter | Type | Description | Required |
|-----------|------|-------------|----------|
| `Email` | string | Customer email address | Yes (at least one field required) |
| `Phone` | string | Customer phone number | No |

#### Process Payment

Process payment for the checkout session:

```go
func processPayment() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    paymentReq := &checkout.PaymentRequest{
        MethodID: "pm_card_visa_ending_4242",
        CVV:      "123",
        SaveCard: true, // Save for future use
    }

    checkout, err := sf.Checkout().ProcessPayment(context.Background(), "checkout_xyz789", *paymentReq)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Payment processed. Status: %s\n", checkout.Status)
    fmt.Printf("Amount charged: $%s\n", checkout.Amount.String())
}
```

**PaymentRequest Parameters:**

| Parameter | Type | Description | Required |
|-----------|------|-------------|----------|
| `MethodID` | string | Payment method token from payment processor | Yes |
| `CVV` | string | CVV/CVC code for card payments | No (depending on processor requirements) |
| `SaveCard` | bool | Whether to save card for future use | No (default: false) |

#### Checkout Structure

The `Checkout` type represents a checkout session:

```go
type Checkout struct {
    ID                string      // Unique checkout identifier (e.g., "checkout_xyz789")
    Status            CheckoutStatus  // Current status of the checkout process
    Amount            Decimal     // Total amount to be charged
    Currency          string      // Currency code (e.g., "USD")
    Customer          *CustomerInfo  // Customer information for this order
    ShippingAddress   *Address    // Complete shipping address details
    BillingAddress    *Address    // Optional billing address (defaults to shipping)
    PaymentMethodID   *string     // Payment method being used
    CreatedAt         time.Time   // When checkout was created
    UpdatedAt         time.Time   // Last update timestamp
}

type CustomerInfo struct {
    Email   *string `json:"email,omitempty"`   // Customer email address
    Phone   *string `json:"phone,omitempty"`   // Customer phone number
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

type CheckoutStatus string

const (
    CheckoutStatusPending   CheckoutStatus = "pending"     // Awaiting payment
    CheckoutStatusProcessing CheckoutStatus = "processing"  // Payment being processed
    CheckoutStatusCompleted CheckoutStatus = "completed"    // Payment successful, order created
    CheckoutStatusFailed    CheckoutStatus = "failed"       // Payment failed or cancelled
)
```

#### Complete Checkout Workflow Example

This example demonstrates a complete checkout flow with error handling:

```go
func completeCheckoutFlow() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    // Step 1: Get existing cart and verify it has items
    cart, err := sf.Cart().Get(context.Background(), "cart_abc123")
    if err != nil {
        log.Fatal(err)
    }

    if len(cart.Items) == 0 {
        log.Fatal("Cannot checkout empty cart")
    }

    fmt.Printf("Cart has %d items, total: $%s\n", 
        len(cart.Items), cart.TotalAmount.String())

    // Step 2: Create checkout session with customer and shipping details
    checkoutReq := &checkout.CreateCheckoutRequest{
        CustomerEmail: "john.doe@example.com",
        ShippingAddress: &checkout.Address{
            FirstName:     "John",
            LastName:      "Doe",
            AddressLine1:  "123 Main St",
            AddressLine2:  "Apt 4B",
            City:          "San Francisco",
            State:         "CA",
            PostalCode:    "94105",
            Country:       "US",
            Phone:         "+14155551234",
        },
    }

    checkout, err := sf.Checkout().Create(context.Background(), cart.ID, *checkoutReq)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Checkout created: %s\n", checkout.ID)
    fmt.Printf("Status: %s\n", checkout.Status)

    // Step 3: Update customer phone number during checkout
    customerInfo := &checkout.CustomerInfo{
        Phone: "+14155559876",
    }

    checkout, err = sf.Checkout().UpdateCustomer(context.Background(), checkout.ID, *customerInfo)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Updated customer phone to: %s\n", *checkout.Customer.Phone)

    // Step 4: Process payment
    paymentReq := &checkout.PaymentRequest{
        MethodID: "pm_card_visa_ending_4242",
        CVV:      "123",
        SaveCard: true,
    }

    checkout, err = sf.Checkout().ProcessPayment(context.Background(), checkout.ID, *paymentReq)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Payment processed successfully!\n")
    fmt.Printf("Status: %s\n", checkout.Status)
    fmt.Printf("Amount charged: $%s\n", checkout.Amount.String())
}
```

#### Payment Flow Patterns

**Pattern 1: Create with immediate payment**

```go
// Customer provides payment details upfront during checkout creation
checkoutReq := &checkout.CreateCheckoutRequest{
    CustomerEmail:     "customer@example.com",
    ShippingAddress:   &address,
    PaymentMethodID:   "pm_card_visa_ending_4242", // Include here for immediate processing
}

checkout, err := sf.Checkout().Create(ctx, cartID, *checkoutReq)
```

**Pattern 2: Create then payment separately (3DS flow)**

```go
// Step 1: Create checkout without payment
checkoutReq := &checkout.CreateCheckoutRequest{
    CustomerEmail:     "customer@example.com",
    ShippingAddress:   &address,
}

checkout, err := sf.Checkout().Create(ctx, cartID, *checkoutReq)

// Step 2: Redirect customer to 3DS authentication (handled by payment processor)
// After successful 3DS, capture the payment token and process it

// Step 3: Process payment after 3DS completes
paymentReq := &checkout.PaymentRequest{
    MethodID: "pm_token_from_3ds_flow",
}

checkout, err = sf.Checkout().ProcessPayment(ctx, checkout.ID, *paymentReq)
```

#### Error Handling

```go
import (
    "errors"
    "log"

    "github.com/artyconst/fb-storefront-go"
)

func handleCheckoutErrors() {
    // Checkout not found
    _, err := sf.Checkout().Get(context.Background(), "invalid-checkout-id")
    if err != nil {
        if errors.Is(err, storefront.ErrResourceNotFound) {
            log.Println("Checkout session not found or expired")
        } else {
            log.Printf("API Error: %v", err)
        }
    }

    // Payment processing failed
    _, err = sf.Checkout().ProcessPayment(context.Background(), "checkout_xyz", 
        checkout.PaymentRequest{
            MethodID: "pm_expired_card",
        })
    if err != nil {
        log.Printf("Payment failed: %v", err)
    }

    // Cart already converted to order
    _, err = sf.Checkout().Create(context.Background(), "cart_completed", 
        checkout.CreateCheckoutRequest{})
    if err != nil {
        log.Printf("Cannot create checkout: %v", err)
    }
}
```
