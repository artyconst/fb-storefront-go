package order

import (
	"context"
	"fmt"
	"strings"

	sf "github.com/artyconst/fb-storefront-go"
)

// OrderService handles order-related operations.
type OrderService struct {
	client *sf.StorefrontClient
}

// NewOrderService creates a new Order service instance.
func NewOrderService(client *sf.StorefrontClient) *OrderService {
	return &OrderService{client: client}
}

// Get retrieves an order by ID or order number.
func (s *OrderService) Get(ctx context.Context, identifier string) (*Order, error) {
	path := "/orders/" + identifier
	var order Order
	if err := s.client.GetJSON(ctx, path, &order); err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	return &order, nil
}

// List retrieves orders with optional filtering.
func (s *OrderService) List(ctx context.Context, opts ...ListOption) ([]*Order, error) {
	options := &ListOptions{}
	for _, o := range opts {
		o(options)
	}

	path := "/orders"
	queryParts := []string{}

	if options.Page > 0 {
		queryParts = append(queryParts, "page="+fmt.Sprint(options.Page))
	}
	if options.Limit > 0 {
		queryParts = append(queryParts, "limit="+fmt.Sprint(options.Limit))
	}
	if string(options.Status) != "" {
		queryParts = append(queryParts, "status="+string(options.Status))
	}

	if len(queryParts) > 0 {
		path += "?" + strings.Join(queryParts, "&")
	}

	var result struct {
		Data []*Order `json:"data"`
	}
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}
	return result.Data, nil
}

// Create creates a new order from a completed checkout.
func (s *OrderService) Create(ctx context.Context, cartID string) (*Order, error) {
	path := "/orders"
	body := map[string]string{"cart_id": cartID}
	var order Order
	if err := s.client.PostJSON(ctx, path, body, &order); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}
	return &order, nil
}

// MarkPickedUp marks an order as picked up - PUT /orders/picked-up.
func (s *OrderService) MarkPickedUp(ctx context.Context, orderID string) error {
	if orderID == "" {
		return fmt.Errorf("order ID is required")
	}

	path := "/orders/picked-up"
	body := map[string]any{"id": orderID}

	if err := s.client.PutJSON(ctx, path, body, &struct{}{}); err != nil {
		return fmt.Errorf("failed to mark order as picked up: %w", err)
	}
	return nil
}

// GenerateReceipt generates a receipt for an order - POST /orders/receipt.
func (s *OrderService) GenerateReceipt(ctx context.Context, orderID string, opts ...ReceiptOption) (*Receipt, error) {
	if orderID == "" {
		return nil, fmt.Errorf("order ID is required")
	}

	options := &ReceiptOptions{}
	for _, o := range opts {
		o(options)
	}

	path := "/orders/receipt"
	body := map[string]any{"id": orderID}
	if options.EbarimtReceiverType != "" {
		body["ebarimt_receiver_type"] = options.EbarimtReceiverType
	}
	if options.EbarimtReceiver != nil {
		body["ebarimt_receiver"] = *options.EbarimtReceiver
	}

	var receipt Receipt
	if err := s.client.PostJSON(ctx, path, body, &receipt); err != nil {
		return nil, fmt.Errorf("failed to generate receipt: %w", err)
	}
	return &receipt, nil
}

