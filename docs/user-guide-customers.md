### Customers Service

The Customers service manages customer accounts and profiles, allowing you to create new customer accounts and retrieve existing customer information.

#### Create Customer Account

Register a new customer account:

```go
import (
    "context"
    "fmt"
    "log"

    customerSDK "github.com/artyconst/fb-storefront-go/pkg/resources/customer"
)

func createCustomer() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    customerService := customerSDK.NewCustomerService(sf)

    customerReq := &customerSDK.CustomerCreateRequest{
        Name:   stringPtr("John Doe"),
        Email:  stringPtr("john.doe@example.com"),
        Phone:  stringPtr("+14155551234"),
        Type:   stringPtr("customer"),
    }

    customer, err := customerService.Create(context.Background(), *customerReq)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Customer created: %s\n", derefString(customer.Email))
    fmt.Printf("Customer ID: %s\n", customer.ID)
}

// Helper function to create string pointers for optional fields
func stringPtr(s string) *string {
    return &s
}

func derefString(s *string) string {
    if s != nil {
        return *s
    }
    return ""
}
```

**CustomerCreateRequest Parameters:**

| Parameter | Type | Description | Required |
|-----------|------|-------------|----------|
| `Name` | `*string` | Customer full name (display name) | No |
| `Type` | `*string` | Customer type (e.g., "customer", "merchant") | No |
| `Identity` | `string` | Unique identity identifier (email or phone) | No |
| `Email` | `*string` | Customer email address (unique identifier) | No |
| `Phone` | `*string` | Customer phone number | No |
| `Code` | `*string` | Customer code | No |
| `Title` | `*string` | Customer title | No |
| `Meta` | `map[string]interface{}` | Additional metadata | No |

#### Get Customer Details

Retrieve a customer by ID:

```go
func getCustomer() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    customerService := customerSDK.NewCustomerService(sf)

    customer, err := customerService.Get(context.Background(), "cust_abc123")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("%s - %s\n", 
        derefString(customer.Name), 
        derefString(customer.Email))
    
    if customer.Phone != nil {
        fmt.Printf("Phone: %s\n", *customer.Phone)
    }
    
    fmt.Printf("Type: %s\n", derefString(customer.Type))
    fmt.Printf("Created at: %s\n", customer.CreatedAt)
    fmt.Printf("Updated at: %s\n", customer.UpdatedAt)
}
```

#### Customers Service Methods

The Customers service provides the following methods:

| Method | Description | Parameters | Returns |
|--------|-------------|------------|---------|
| `Create(ctx, req)` | Create new customer account | Customer creation request | New customer |
| `Get(ctx, customerID)` | Retrieve customer by ID | Customer ID string | Customer details |
| `Login(ctx, req)` | Authenticate customer | Login request (identity + password) | Login response with token |
| `LoginWithSMS(ctx, req)` | Initiate SMS authentication | SMS sign-in request (identity) | Login response |
| `VerifySMSCode(ctx, req)` | Confirm SMS code | SMS confirm request (identity + code) | Login response |
| `ListPlaces(ctx, token, opts)` | List saved locations | Auth token + optional pagination | List of places |
| `ListOrders(ctx, token, opts)` | List customer orders | Auth token + optional pagination | List of orders |
| `RequestCreationCode(ctx, req)` | Request account creation code | Request creation code request | Error (nil on success) |
| `RegisterDevice(ctx, token, req)` | Register push notification device | Auth token + device request | Registration response |

#### Customer Structure

The `Customer` type contains the following fields:

```go
type Customer struct {
    ID        string  // Unique customer identifier (e.g., "cust_abc123")
    Name      *string // Optional full name
    Email     *string // Optional email address (unique, used for authentication)
    Phone     *string // Optional phone number
    Type      *string // Optional customer type (e.g., "customer", "merchant")
    CreatedAt string  // When customer account was created (ISO 8601 format string)
    UpdatedAt string  // Last update timestamp (ISO 8601 format string)
}
```

**Field Details:**

- **ID**: Used to reference the customer in API calls and associate with carts/orders
- **Name**: Customer's full display name (combined first and last name)
- **Email**: Unique identifier for the customer; used as username for authentication
- **Phone**: Optional contact number, useful for SMS notifications
- **Type**: Customer classification (e.g., "customer", "merchant")
- **CreatedAt/UpdatedAt**: Timestamps as ISO 8601 format strings for audit trails

#### Customer Authentication

Authenticate an existing customer with identity and password:

```go
func authenticateCustomer() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    customerService := customerSDK.NewCustomerService(sf)

    loginResp, err := customerService.Login(context.Background(), customerSDK.LoginRequest{
        Identity: "john.doe@example.com",
        Password: "securepassword",
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Login successful! Token: %s...\n", loginResp.Token[:min(len(loginResp.Token), 20)])
    fmt.Printf("Customer ID: %s\n", loginResp.Customer.ID)
    fmt.Printf("Expires at: %s\n", loginResp.ExpiresAt)
}
```

**LoginRequest Parameters:**

| Parameter | Type | Description | Required |
|-----------|------|-------------|----------|
| `Identity` | `string` | Customer email or phone number | Yes |
| `Password` | `string` | Customer password | Yes |

#### Customer Authentication with SMS

Initiate SMS-based authentication:

```go
func authenticateWithSMS() {
    sf, err := storefront.NewStorefront(YOUR_API_KEY)
    if err != nil {
        log.Fatal(err)
    }

    customerService := customerSDK.NewCustomerService(sf)

    // Step 1: Initiate SMS sign-in
    loginResp, err := customerService.LoginWithSMS(context.Background(), customerSDK.SMSSignInRequest{
        Identity: "+14155551234",
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("SMS code sent. Please enter the code you received.")

    // Step 2: Verify SMS code
    loginResp, err = customerService.VerifySMSCode(context.Background(), customerSDK.SMSConfirmSignInRequest{
        Identity: "+14155551234",
        Code:     "123456",
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("SMS authentication successful! Token: %s...\n", loginResp.Token[:min(len(loginResp.Token), 20)])
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

    customerService := customerSDK.NewCustomerService(sf)

    // 1. Create new customer
    customerReq := &customerSDK.CustomerCreateRequest{
        Name:   stringPtr("Jane Smith"),
        Email:  stringPtr("shopper@example.com"),
        Type:   stringPtr("customer"),
    }

    customer, err := customerService.Create(context.Background(), *customerReq)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Welcome, %s!\n", derefString(customer.Name))

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
        CustomerEmail: derefString(customer.Email),
        ShippingAddress: &cart.Address{
            Name:       "Jane Smith",
            AddressLine1: "456 Oak Ave",
            City:       "Los Angeles",
            State:      "CA",
            PostalCode: "90001",
            Country:    "US",
        },
        PaymentMethodID: "pm_card_mastercard_ending_5555",
    }

    order, err := sf.Cart().Checkout(context.Background(), cart.ID, *checkoutReq)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Order #%s placed successfully!\n", order.OrderNumber)
}

func stringPtr(s string) *string {
    return &s
}

func derefString(s *string) string {
    if s != nil {
        return *s
    }
    return ""
}
```

#### Error Handling

```go
import (
    "errors"
    "log"

    "github.com/artyconst/fb-storefront-go"
    customerSDK "github.com/artyconst/fb-storefront-go/pkg/resources/customer"
)

func handleCustomerErrors() {
    sf, _ := storefront.NewStorefront(YOUR_API_KEY)
    customerService := customerSDK.NewCustomerService(sf)

    // Duplicate email error
    createReq := &customerSDK.CustomerCreateRequest{
        Email: stringPtr("existing@example.com"),  // Already registered
    }

    _, err := customerService.Create(context.Background(), *createReq)
    if err != nil {
        log.Printf("Customer creation error: %v", err)
    }

    // Customer not found
    _, err = customerService.Get(context.Background(), "invalid-customer-id")
    if err != nil {
        if errors.Is(err, storefront.ErrResourceNotFound) {
            log.Println("Customer not found")
        } else {
            log.Printf("API Error: %v", err)
        }
    }
}
```

#### Customer Email Uniqueness

Email addresses must be unique across all customers. Attempting to create a duplicate email will result in an error:

```go
func handleDuplicateEmail() {
    sf, _ := storefront.NewStorefront(YOUR_API_KEY)
    customerService := customerSDK.NewCustomerService(sf)

    // First customer creation succeeds
    _, err := customerService.Create(context.Background(), customerSDK.CustomerCreateRequest{
        Email: stringPtr("unique@example.com"),
    })
    if err != nil {
        log.Fatal(err)
    }

    // Second attempt with same email fails
    _, err = customerService.Create(context.Background(), customerSDK.CustomerCreateRequest{
        Email: stringPtr("unique@example.com"),  // Duplicate!
    })
    
    if err != nil {
        // Handle duplicate email error
        log.Printf("Email already registered: %v", err)
    }
}
```

#### Retrieving Customer's Places (Saved Locations)

After authenticating a customer, you can retrieve their saved locations:

```go
func listCustomerPlaces(token string) {
    sf, _ := storefront.NewStorefront(YOUR_API_KEY)
    customerService := customerSDK.NewCustomerService(sf)

    places, err := customerService.ListPlaces(context.Background(), token,
        customerSDK.WithPage(1),
        customerSDK.WithLimit(10),
    )
    if err != nil {
        log.Fatal(err)
    }

    for _, place := range places {
        fmt.Printf("Place: %s at %s\n", place.Name, place.Address)
    }
}
```

#### Retrieving Customer's Orders

After authenticating a customer, you can retrieve their order history:

```go
func listCustomerOrders(token string) {
    sf, _ := storefront.NewStorefront(YOUR_API_KEY)
    customerService := customerSDK.NewCustomerService(sf)

    orders, err := customerService.ListOrders(context.Background(), token,
        customerSDK.WithOrderLimit(20),
        customerSDK.WithOrderSort("-created_at"),
    )
    if err != nil {
        log.Fatal(err)
    }

    for _, order := range orders {
        fmt.Printf("Order #%s - Status: %s - Total: %d %s\n",
            order.ID, order.Status, order.Total, order.Currency)
    }
}
```
