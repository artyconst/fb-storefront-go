package order

// OrderStatus represents the status of an order.
type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusConfirmed  OrderStatus = "confirmed"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
	OrderStatusRefunded   OrderStatus = "refunded"
)

// OrderItem represents an item in an order.
type OrderItem struct {
	ID        string `json:"id"`
	ProductID string `json:"product_id"`
	Name      string `json:"name"`
	SKU       string `json:"sku"`
	Quantity  int    `json:"quantity"`
	Price     int64  `json:"price"`
	Total     int64  `json:"total"`
}

// Order represents an order in the system.
type Order struct {
	ID              string       `json:"id"`
	OrderNumber     string       `json:"order_number"`
	Status          OrderStatus  `json:"status"`
	CustomerID      string       `json:"customer_id"`
	CustomerEmail   string       `json:"customer_email"`
	Items           []*OrderItem `json:"items"`
	Subtotal        int64        `json:"subtotal"`
	TaxAmount       int64        `json:"tax_amount"`
	ShippingCost    int64        `json:"shipping_cost"`
	Discount        *int64       `json:"discount,omitempty"`
	TotalAmount     int64        `json:"total_amount"`
	Currency        string       `json:"currency"`
	BillingAddress  *Address     `json:"billing_address,omitempty"`
	ShippingAddress *Address     `json:"shipping_address,omitempty"`
	PaymentStatus   string       `json:"payment_status"`
	Notes           *string      `json:"notes,omitempty"`
	CreatedAt       string       `json:"created_at"`
	UpdatedAt       string       `json:"updated_at"`
}

// Address represents an address for orders.
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

// ListOptions contains parameters for listing orders.
type ListOptions struct {
	Page   int         `json:"page,omitempty"`
	Limit  int         `json:"limit,omitempty"`
	Status OrderStatus `json:"status,omitempty"`
}
