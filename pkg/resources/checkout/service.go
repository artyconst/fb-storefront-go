package checkout

import (
	"context"
	"fmt"
	"net/url"
	"strings"

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
	body := map[string]any{"cart_id": cartID}
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
		urlPath += "?" + strings.Join(queryParts, "&")
	}

	var quote DeliveryServiceQuote
	if err := s.client.GetJSON(ctx, urlPath, &quote); err != nil {
		return nil, fmt.Errorf("failed to get delivery service quote: %w", err)
	}
	return &quote, nil
}

// Initialize creates a checkout preview by calling the /before endpoint.
// It accepts functional options for gateway, customer, cart, service quote, cash, pickup, tip, and delivery tip.
func (s *CheckoutService) Initialize(ctx context.Context, opts ...InitOption) (*CheckoutPreview, error) {
	options := &InitOptions{}
	for _, o := range opts {
		o(options)
	}

	path := "/checkouts/before"
	urlPath := path
	queryParts := []string{}

	if options.Gateway != "" {
		queryParts = append(queryParts, "gateway="+url.QueryEscape(options.Gateway))
	}
	if options.Customer != "" {
		queryParts = append(queryParts, "customer="+url.QueryEscape(options.Customer))
	}
	if options.Cart != "" {
		queryParts = append(queryParts, "cart="+url.QueryEscape(options.Cart))
	}
	if options.ServiceQuote != "" {
		queryParts = append(queryParts, "service_quote="+url.QueryEscape(options.ServiceQuote))
	}
	if options.Cash {
		queryParts = append(queryParts, "cash=true")
	}
	if options.Pickup {
		queryParts = append(queryParts, "pickup=true")
	}
	if options.Tip != 0 {
		queryParts = append(queryParts, fmt.Sprintf("tip=%d", options.Tip))
	}
	if options.DeliveryTip != 0 {
		queryParts = append(queryParts, fmt.Sprintf("delivery_tip=%d", options.DeliveryTip))
	}

	if len(queryParts) > 0 {
		urlPath += "?" + strings.Join(queryParts, "&")
	}

	var preview CheckoutPreview
	if err := s.client.GetJSON(ctx, urlPath, &preview); err != nil {
		return nil, fmt.Errorf("failed to initialize checkout: %w", err)
	}
	return &preview, nil
}

// Status retrieves the status of a checkout. It accepts functional options for checkout ID or token.
func (s *CheckoutService) Status(ctx context.Context, opts ...StatusOption) (*CheckoutStatusResponse, error) {
	options := &StatusOptions{}
	for _, o := range opts {
		o(options)
	}

	path := "/checkouts/status"
	urlPath := path
	queryParts := []string{}

	if options.CheckoutID != nil && *options.CheckoutID != "" {
		queryParts = append(queryParts, "checkout_id="+url.QueryEscape(*options.CheckoutID))
	}
	if options.Token != nil && *options.Token != "" {
		queryParts = append(queryParts, "token="+url.QueryEscape(*options.Token))
	}

	if len(queryParts) > 0 {
		urlPath += "?" + strings.Join(queryParts, "&")
	}

	var status CheckoutStatusResponse
	if err := s.client.GetJSON(ctx, urlPath, &status); err != nil {
		return nil, fmt.Errorf("failed to get checkout status: %w", err)
	}
	return &status, nil
}

// UpdatePaymentIntent updates the payment intent for a checkout.
func (s *CheckoutService) UpdatePaymentIntent(ctx context.Context, paymentIntentID string, opts ...UpdateIntentOption) (*CheckoutPreview, error) {
	if paymentIntentID == "" {
		return nil, fmt.Errorf("failed to update payment intent: %w", sf.ErrInvalidRequest)
	}

	options := &UpdateIntentOptions{}
	for _, o := range opts {
		o(options)
	}

	path := "/checkouts/stripe-update-payment-intent"
	body := map[string]any{
		"payment_intent_id": paymentIntentID,
	}

	var preview CheckoutPreview
	if err := s.client.PutJSON(ctx, path, body, &preview); err != nil {
		return nil, fmt.Errorf("failed to update payment intent: %w", err)
	}
	return &preview, nil
}

// CaptureCheckout captures a checkout as an order. Token is passed as query parameter,
// request body can contain optional transaction details and notes for the order.
func (s *CheckoutService) CaptureCheckout(ctx context.Context, token string, opts ...CaptureOption) (*Checkout, error) {
	options := &CaptureOptions{}
	for _, o := range opts {
		o(options)
	}

	path := "/checkouts/capture"
	urlPath := path
	if token != "" {
		urlPath += "?token=" + url.QueryEscape(token)
	}

	body := map[string]any{}
	if options.TransactionDetails != nil {
		body["transactionDetails"] = options.TransactionDetails
	}
	if options.Notes != nil {
		body["notes"] = *options.Notes
	}

	var checkout Checkout
	if err := s.client.PostJSON(ctx, urlPath, body, &checkout); err != nil {
		return nil, fmt.Errorf("failed to capture checkout: %w", err)
	}
	return &checkout, nil
}

// CaptureQPay captures a checkout via QPay callback - POST /checkouts/capture-qpay
func (s *CheckoutService) CaptureQPay(ctx context.Context, checkoutID string, respond bool, test *string) (*Order, error) {
	if checkoutID == "" {
		return nil, fmt.Errorf("checkout ID is required")
	}

	body := map[string]any{
		"id":      checkoutID,
		"respond": respond,
	}
	if test != nil {
		body["test"] = *test
	}

	var order Order
	if err := s.client.PostJSON(ctx, "/checkouts/capture-qpay", body, &order); err != nil {
		return nil, fmt.Errorf("failed to capture QPay: %w", err)
	}
	return &order, nil
}
