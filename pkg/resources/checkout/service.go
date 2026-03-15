package checkout

import (
	"context"
	"fmt"
	"net/url"

	sf "github.com/artyconst/fb-storefront-go"
)

// CheckoutService handles checkout-related operations.
type CheckoutService struct {
	client *sf.StorefrontClient
}

// NewCheckoutService creates a new Checkout service instance.
func NewCheckoutService(client *sf.StorefrontClient) *CheckoutService {
	return &CheckoutService{client: client}
}

// Create creates a new checkout session from a cart.
func (s *CheckoutService) Create(ctx context.Context, cartID string, req CreateCheckoutRequest) (*Checkout, error) {
	path := "/checkouts"
	body := map[string]interface{}{"cart_id": cartID}
	if req.CustomerEmail != "" {
		body["customer_email"] = req.CustomerEmail
	}
	if req.ShippingAddress != nil {
		body["shipping_address"] = req.ShippingAddress
	}
	if req.BillingAddress != nil {
		body["billing_address"] = req.BillingAddress
	}
	if req.PaymentMethodID != "" {
		body["payment_method_id"] = req.PaymentMethodID
	}

	var checkout Checkout
	if err := s.client.PostJSON(ctx, path, body, &checkout); err != nil {
		return nil, fmt.Errorf("failed to create checkout: %w", err)
	}
	return &checkout, nil
}

// Get retrieves a checkout by ID.
func (s *CheckoutService) Get(ctx context.Context, id string) (*Checkout, error) {
	path := "/checkouts/" + id
	var checkout Checkout
	if err := s.client.GetJSON(ctx, path, &checkout); err != nil {
		return nil, fmt.Errorf("failed to get checkout: %w", err)
	}
	return &checkout, nil
}

// UpdateCustomer updates the customer information for a checkout.
func (s *CheckoutService) UpdateCustomer(ctx context.Context, checkoutID string, customer CustomerInfo) (*Checkout, error) {
	path := "/checkouts/" + checkoutID + "/customer"
	var checkout Checkout
	if err := s.client.PutJSON(ctx, path, customer, &checkout); err != nil {
		return nil, fmt.Errorf("failed to update customer: %w", err)
	}
	return &checkout, nil
}

// ProcessPayment processes a payment for the checkout.
func (s *CheckoutService) ProcessPayment(ctx context.Context, checkoutID string, req PaymentRequest) (*Checkout, error) {
	path := "/checkouts/" + checkoutID + "/payment"
	var checkout Checkout
	if err := s.client.PostJSON(ctx, path, req, &checkout); err != nil {
		return nil, fmt.Errorf("failed to process payment: %w", err)
	}
	return &checkout, nil
}

// GetDeliveryServiceQuote retrieves a delivery service quote based on origin/destination/cart.
func (s *CheckoutService) GetDeliveryServiceQuote(ctx context.Context, params ServiceQuoteParams) (*DeliveryServiceQuote, error) {
	path := "/service-quotes/from-cart"
	urlPath := path
	queryParts := []string{}

	if params.Origin != "" {
		queryParts = append(queryParts, "origin="+params.Origin)
	}
	if params.Destination != "" {
		queryParts = append(queryParts, "destination="+params.Destination)
	}
	if params.CartID != "" {
		queryParts = append(queryParts, "cart="+params.CartID)
	}
	if params.ServiceQuote != "" {
		queryParts = append(queryParts, "service_quote="+params.ServiceQuote)
	}

	if len(queryParts) > 0 {
		urlPath += "?" + joinStrings(queryParts, "&")
	}

	var quote DeliveryServiceQuote
	if err := s.client.GetJSON(ctx, urlPath, &quote); err != nil {
		return nil, fmt.Errorf("failed to get delivery service quote: %w", err)
	}
	return &quote, nil
}

// CaptureCheckout captures a checkout as an order.
func (s *CheckoutService) CaptureCheckout(ctx context.Context, token string) (*Checkout, error) {
	path := "/checkouts/capture"
	urlPath := path
	if token != "" {
		urlPath += "?token=" + url.QueryEscape(token)
	}

	var checkout Checkout
	if err := s.client.PostJSON(ctx, urlPath, nil, &checkout); err != nil {
		return nil, fmt.Errorf("failed to capture checkout: %w", err)
	}
	return &checkout, nil
}

// joinStrings joins string slice with separator.
func joinStrings(parts []string, sep string) string {
	if len(parts) == 0 {
		return ""
	}
	if len(parts) == 1 {
		return parts[0]
	}
	result := parts[0]
	for i := 1; i < len(parts); i++ {
		result += sep + parts[i]
	}
	return result
}
