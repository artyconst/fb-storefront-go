package checkout

// CheckoutStatus represents the status of a checkout.
type CheckoutStatus string

const (
	CheckoutStatusPending    CheckoutStatus = "pending"
	CheckoutStatusProcessing CheckoutStatus = "processing"
	CheckoutStatusCompleted  CheckoutStatus = "completed"
	CheckoutStatusFailed     CheckoutStatus = "failed"
)

// CustomerInfo contains customer details for checkout.
type CustomerInfo struct {
	ID    string `json:"id,omitempty"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
}

// PaymentMethod represents a payment method for checkout.
type PaymentMethod struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	LastFour string `json:"last_four,omitempty"`
	Brand    string `json:"brand,omitempty"`
	ExpMonth int    `json:"exp_month,omitempty"`
	ExpYear  int    `json:"exp_year,omitempty"`
}

// Checkout represents a checkout session.
type Checkout struct {
	ID              string         `json:"id"`
	CartID          string         `json:"cart_id"`
	Status          CheckoutStatus `json:"status"`
	Customer        *CustomerInfo  `json:"customer,omitempty"`
	ShippingAddress *Address       `json:"shipping_address,omitempty"`
	BillingAddress  *Address       `json:"billing_address,omitempty"`
	PaymentMethod   *PaymentMethod `json:"payment_method,omitempty"`
	Amount          int64          `json:"amount"`
	TaxAmount       *int64         `json:"tax_amount,omitempty"`
	Currency        string         `json:"currency"`
	CreatedAt       string         `json:"created_at"`
	UpdatedAt       string         `json:"updated_at"`
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

// CreateCheckoutRequest contains parameters for creating a checkout.
type CreateCheckoutRequest struct {
	CustomerEmail   string   `json:"customer_email,omitempty"`
	ShippingAddress *Address `json:"shipping_address,omitempty"`
	BillingAddress  *Address `json:"billing_address,omitempty"`
	PaymentMethodID string   `json:"payment_method_id,omitempty"`
}

// PaymentRequest contains payment processing data.
type PaymentRequest struct {
	MethodID string            `json:"method_id"`
	CVV      string            `json:"cvv,omitempty"`
	SaveCard bool              `json:"save_card,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// DeliveryServiceQuote represents a delivery service quote.
type DeliveryServiceQuote struct {
	ID          string         `json:"id"`
	Origin      string         `json:"origin"`
	Destination string         `json:"destination"`
	CartID      string         `json:"cart_id,omitempty"`
	Price       int64          `json:"price"`
	Currency    string         `json:"currency"`
	Status      CheckoutStatus `json:"status"`
	CreatedAt   string         `json:"created_at"`
}

// ServiceQuoteParams for delivery service quote query.
type ServiceQuoteParams struct {
	Origin       string `json:"origin,omitempty"`
	Destination  string `json:"destination,omitempty"`
	CartID       string `json:"cart_id,omitempty"`
	ServiceQuote string `json:"service_quote,omitempty"`
}
