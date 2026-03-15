package cart

// CartStatus represents the status of a cart.
type CartStatus string

const (
	CartStatusActive    CartStatus = "active"
	CartStatusCompleted CartStatus = "completed"
	CartStatusAbandoned CartStatus = "abandoned"
)

// CartItem represents an item in the cart.
type CartItem struct {
	ID        string `json:"id"`
	ProductID string `json:"product_id"`
	Name      string `json:"name"`
	Quantity  int    `json:"quantity"`
	Price     int64  `json:"price"`
	Total     int64  `json:"total"`
}

// Cart represents a shopping cart.
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

// CartItemRequest represents a request to add an item to cart.
type CartItemRequest struct {
	ProductID     string           `json:"product_id"`
	Quantity      int              `json:"quantity"`
	Addons        []interface{}    `json:"addons,omitempty"`
	Variants      []map[string]any `json:"variants,omitempty"`
	ScheduledAt   string           `json:"scheduled_at,omitempty"`
	StoreLocation string           `json:"store_location,omitempty"`
}

// Address represents a shipping or billing address.
type Address struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Company      string `json:"company,omitempty"`
	AddressLine1 string `json:"address_line_1"`
	AddressLine2 string `json:"address_line_2,omitempty"`
	City         string `json:"city"`
	State        string `json:"state,omitempty"`
	PostalCode   string `json:"postal_code"`
	Country      string `json:"country"`
	Phone        string `json:"phone,omitempty"`
}

// CheckoutRequest contains data needed for checkout.
type CheckoutRequest struct {
	CustomerEmail   string                 `json:"customer_email,omitempty"`
	ShippingAddress *Address               `json:"shipping_address,omitempty"`
	BillingAddress  *Address               `json:"billing_address,omitempty"`
	PaymentMethodID string                 `json:"payment_method_id,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// Order represents a minimal order structure for cart checkout.
type Order struct {
	ID          string `json:"id"`
	OrderNumber string `json:"order_number"`
	Status      string `json:"status"`
	TotalAmount int64  `json:"total_amount"`
	Currency    string `json:"currency"`
}
