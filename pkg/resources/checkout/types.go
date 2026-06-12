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

// CheckoutPreview represents the preview response from checkout initialization.
type CheckoutPreview struct {
	Token         string  `json:"token"`
	PaymentIntent *string `json:"paymentIntent,omitempty"`
	ClientSecret  *string `json:"clientSecret,omitempty"`
	EphemeralKey  *string `json:"ephemeralKey,omitempty"`
	CustomerID    *string `json:"customerId,omitempty"`
}

// CheckoutStatusResponse represents the status response from a checkout.
type CheckoutStatusResponse struct {
	Status string `json:"status"`
	ID     string `json:"id,omitempty"`
	Token  string `json:"token,omitempty"`
}

// InitOptions contains optional parameters for initializing a checkout preview.
type InitOptions struct {
	Gateway      string `json:"gateway,omitempty"`
	Customer     string `json:"customer,omitempty"`
	Cart         string `json:"cart,omitempty"`
	ServiceQuote string `json:"service_quote,omitempty"`
	Cash         bool   `json:"cash,omitempty"`
	Pickup       bool   `json:"pickup,omitempty"`
	Tip          int64  `json:"tip,omitempty"`
	DeliveryTip  int64  `json:"delivery_tip,omitempty"`
}

// InitOption is a functional option for initializing a checkout preview.
type InitOption func(*InitOptions)

// WithGateway sets the payment gateway for checkout initialization.
func WithGateway(gateway string) InitOption {
	return func(opts *InitOptions) {
		opts.Gateway = gateway
	}
}

// WithCustomer sets the customer ID for checkout initialization.
func WithCustomer(customerID string) InitOption {
	return func(opts *InitOptions) {
		opts.Customer = customerID
	}
}

// WithCart sets the cart ID for checkout initialization.
func WithCart(cartID string) InitOption {
	return func(opts *InitOptions) {
		opts.Cart = cartID
	}
}

// WithServiceQuote sets the service quote for checkout initialization.
func WithServiceQuote(quote string) InitOption {
	return func(opts *InitOptions) {
		opts.ServiceQuote = quote
	}
}

// WithCash enables cash payment for checkout initialization.
func WithCash() InitOption {
	return func(opts *InitOptions) {
		opts.Cash = true
	}
}

// WithPickup enables pickup mode for checkout initialization.
func WithPickup() InitOption {
	return func(opts *InitOptions) {
		opts.Pickup = true
	}
}

// WithTip sets the tip amount for checkout initialization.
func WithTip(tip int64) InitOption {
	return func(opts *InitOptions) {
		opts.Tip = tip
	}
}

// WithDeliveryTip sets the delivery tip amount for checkout initialization.
func WithDeliveryTip(deliveryTip int64) InitOption {
	return func(opts *InitOptions) {
		opts.DeliveryTip = deliveryTip
	}
}

// StatusOptions contains optional parameters for checking checkout status.
type StatusOptions struct {
	CheckoutID *string `json:"checkout_id,omitempty"`
	Token      *string `json:"token,omitempty"`
}

// StatusOption is a functional option for checking checkout status.
type StatusOption func(*StatusOptions)

// WithCheckoutID sets the checkout ID for status check.
func WithCheckoutID(checkoutID string) StatusOption {
	return func(opts *StatusOptions) {
		opts.CheckoutID = &checkoutID
	}
}

// WithToken sets the token for status check.
func WithToken(token string) StatusOption {
	return func(opts *StatusOptions) {
		opts.Token = &token
	}
}

// UpdateIntentOptions contains optional parameters for updating a payment intent.
type UpdateIntentOptions struct{}

// UpdateIntentOption is a functional option for updating a payment intent.
type UpdateIntentOption func(*UpdateIntentOptions)

// CaptureOptions contains parameters for capturing a checkout as an order.
type CaptureOptions struct {
	TransactionDetails map[string]any `json:"transactionDetails,omitzero"`
	Notes              *string        `json:"notes,omitempty"`
}

// Order represents an order returned from capture operations.
type Order struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// CaptureOption is a functional option for capturing a checkout.
type CaptureOption func(*CaptureOptions)

// WithTransactionDetails sets the transaction details for capture.
func WithTransactionDetails(details map[string]any) CaptureOption {
	return func(opts *CaptureOptions) {
		opts.TransactionDetails = details
	}
}

// WithNotes sets the notes for capture.
func WithNotes(notes string) CaptureOption {
	return func(opts *CaptureOptions) {
		opts.Notes = &notes
	}
}
