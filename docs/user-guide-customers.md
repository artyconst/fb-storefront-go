### Customers Service

The Customers service manages customer accounts and profiles, allowing you to create new customer accounts and retrieve or update existing customer information.

#### Create Customer Account

Register a new customer account:

```go
import (
    "context"
    "fmt"
    "log"

    "github.com/artyconst/fb-storefront-go/pkg/customer"
)

func createCustomer() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    customerReq := &customer.CreateCustomerRequest{
        Email:     "john.doe@example.com",
        Phone:     stringPtr("+14155551234"),
        FirstName: stringPtr("John"),
        LastName:  stringPtr("Doe"),
    }

    customer, err := sf.Customer().Create(context.Background(), *customerReq)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Customer created: %s\n", customer.Email)
    fmt.Printf("Customer ID: %s\n", customer.ID)
}

// Helper function to create string pointers for optional fields
func stringPtr(s string) *string {
    return &s
}
```

**CreateCustomerRequest Parameters:**

| Parameter | Type | Description | Required |
|-----------|------|-------------|----------|
| `Email` | string | Customer email address (unique identifier) | Yes |
| `Phone` | string | Customer phone number | No |
| `FirstName` | string | Customer first name | No |
| `LastName` | string | Customer last name | No |

#### Get Customer Details

Retrieve a customer by ID:

```go
func getCustomer() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    customer, err := sf.Customer().Get(context.Background(), "cust_abc123")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("%s %s - %s\n", 
        *customer.FirstName, 
        *customer.LastName, 
        customer.Email)
    
    if customer.Phone != nil {
        fmt.Printf("Phone: %s\n", *customer.Phone)
    }
}
```

#### Update Customer Profile

Modify existing customer information:

```go
func updateCustomer() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    customerReq := &customer.UpdateCustomerRequest{
        Phone:     stringPtr("+14159876543"),
        FirstName: stringPtr("Jonathan"),
    }

    customer, err := sf.Customer().Update(context.Background(), "cust_abc123", *customerReq)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Customer updated successfully")
    fmt.Printf("New name: %s %s\n", *customer.FirstName, *customer.LastName)
}

// Helper function to create string pointers for optional fields
func stringPtr(s string) *string {
    return &s
}
```

**UpdateCustomerRequest Parameters:**

| Parameter | Type | Description | Required |
|-----------|------|-------------|----------|
| `Phone` | string | Customer phone number (update or keep existing) | No |
| `FirstName` | string | Customer first name | No |
| `LastName` | string | Customer last name | No |

**Note:** At least one field must be provided in the update request.

#### Customers Service Methods

The Customers service provides the following methods:

| Method | Description | Parameters | Returns |
|--------|-------------|------------|---------|
| `Create(ctx, req)` | Create new customer account | Customer creation request | New customer |
| `Get(ctx, customerID)` | Retrieve customer by ID | Customer ID string | Customer details |
| `Update(ctx, customerID, req)` | Update customer information | Customer ID and update request | Updated customer |

#### Customer Structure

The `Customer` type contains the following fields:

```go
type Customer struct {
    ID          string  // Unique customer identifier (e.g., "cust_abc123")
    Email       string  // Customer email address (unique, used for authentication)
    Phone       *string // Optional phone number
    FirstName   *string // Optional first name
    LastName    *string // Optional last name
    CreatedAt   time.Time // When customer account was created
    UpdatedAt   time.Time // Last update timestamp
}
```

**Field Details:**

- **ID**: Used to reference the customer in API calls and associate with carts/orders
- **Email**: Unique identifier for the customer; used as username for authentication
- **Phone**: Optional contact number, useful for SMS notifications
- **FirstName/LastName**: Customer's display name components
- **CreatedAt/UpdatedAt**: Timestamps for audit trails

#### Complete Customer Workflow Example

This example demonstrates creating a customer, updating their profile, and using them with carts:

```go
func completeCustomerWorkflow() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    // Step 1: Create new customer account
    createReq := &customer.CreateCustomerRequest{
        Email:     "john.doe@example.com",
        Phone:     stringPtr("+14155551234"),
        FirstName: stringPtr("John"),
        LastName:  stringPtr("Doe"),
    }

    customer, err := sf.Customer().Create(context.Background(), *createReq)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Customer created: %s\n", customer.Email)
    fmt.Printf("Customer ID: %s\n", customer.ID)

    // Step 2: Update customer information
    updateReq := &customer.UpdateCustomerRequest{
        FirstName: stringPtr("Jonathan"),
        Phone:     stringPtr("+14159876543"),
    }

    customer, err = sf.Customer().Update(context.Background(), customer.ID, *updateReq)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Updated name: %s\n", *customer.FirstName)
    fmt.Printf("Updated phone: %s\n", *customer.Phone)

    // Step 3: Create cart for this customer
    cart, err := sf.Cart().Create(context.Background(), customer.ID)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Created cart: %s for customer\n", cart.ID)

    // Step 4: Add items to cart (example with product)
    itemReq := &cart.CartItemRequest{
        ProductID: "prod_wireless_headphones",
        Quantity:  1,
    }

    cart, err = sf.Cart().AddItem(context.Background(), cart.ID, *itemReq)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Added item to cart. Total items: %d\n", len(cart.Items))

    // Step 5: Retrieve customer again to verify data persistence
    refreshedCustomer, err := sf.Customer().Get(context.Background(), customer.ID)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Retrieved customer: %s %s\n", 
        *refreshedCustomer.FirstName, 
        *refreshedCustomer.LastName)
}

// Helper function to create string pointers for optional fields
func stringPtr(s string) *string {
    return &s
}
```

#### Error Handling

```go
import (
    "errors"
    "log"

    "github.com/artyconst/fb-storefront-go"
)

func handleCustomerErrors() {
    // Duplicate email error
    createReq := &customer.CreateCustomerRequest{
        Email:     "existing@example.com",  // Already registered
        FirstName: stringPtr("New"),
        LastName:  stringPtr("User"),
    }

    _, err := sf.Customer().Create(context.Background(), *createReq)
    if err != nil {
        log.Printf("Customer creation error: %v", err)
    }

    // Customer not found
    _, err = sf.Customer().Get(context.Background(), "invalid-customer-id")
    if err != nil {
        if errors.Is(err, storefront.ErrResourceNotFound) {
            log.Println("Customer not found")
        } else {
            log.Printf("API Error: %v", err)
        }
    }

    // Update with no fields provided
    _, err = sf.Customer().Update(context.Background(), "cust_abc123", 
        customer.UpdateCustomerRequest{})
    if err != nil {
        log.Printf("Update error (no fields): %v", err)
    }
}
```

#### Customer Email Uniqueness

Email addresses must be unique across all customers. Attempting to create a duplicate email will result in an error:

```go
func handleDuplicateEmail() {
    // First customer creation succeeds
    _, err := sf.Customer().Create(context.Background(), customer.CreateCustomerRequest{
        Email: "unique@example.com",
    })
    if err != nil {
        log.Fatal(err)
    }

    // Second attempt with same email fails
    _, err = sf.Customer().Create(context.Background(), customer.CreateCustomerRequest{
        Email: "unique@example.com",  // Duplicate!
    })
    
    if err != nil {
        // Handle duplicate email error
        log.Printf("Email already registered: %v", err)
    }
}
```

#### Using Customers with Carts and Orders

Customers are typically created first, then associated with carts which can later be converted to orders:

```go
func customerShoppingFlow() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    // 1. Create or retrieve existing customer
    customerReq := &customer.CreateCustomerRequest{
        Email:     "shopper@example.com",
        FirstName: stringPtr("Jane"),
        LastName:  stringPtr("Smith"),
    }

    customer, err := sf.Customer().Create(context.Background(), *customerReq)
    if err != nil {
        // Check for duplicate email
        log.Printf("Customer error: %v", err)
        return
    }

    fmt.Printf("Welcome, %s!\n", *customer.FirstName)

    // 2. Create cart associated with customer
    cart, err := sf.Cart().Create(context.Background(), customer.ID)
    if err != nil {
        log.Fatal(err)
    }

    // 3. Add products to cart
    itemReq := &cart.CartItemRequest{
        ProductID: "prod_example_product",
        Quantity:  2,
    }

    cart, err = sf.Cart().AddItem(context.Background(), cart.ID, *itemReq)
    if err != nil {
        log.Fatal(err)
    }

    // 4. Checkout converts cart to order with customer attached
    checkoutReq := &cart.CheckoutRequest{
        CustomerEmail: customer.Email,
        ShippingAddress: &cart.Address{
            FirstName:    "Jane",
            LastName:     "Smith",
            AddressLine1: "456 Oak Ave",
            City:         "Los Angeles",
            State:        "CA",
            PostalCode:   "90001",
            Country:      "US",
        },
        PaymentMethodID: "pm_card_mastercard_ending_5555",
    }

    order, err := sf.Cart().Checkout(context.Background(), cart.ID, *checkoutReq)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Order #%s placed successfully!\n", order.OrderNumber)
}
```
