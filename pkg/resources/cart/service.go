package cart

import (
	"context"
	"fmt"

	sf "github.com/artyconst/fb-storefront-go"
)

// CartService handles cart-related operations.
type CartService struct {
	client *sf.StorefrontClient
}

// NewCartService creates a new Cart service instance.
func NewCartService(client *sf.StorefrontClient) *CartService {
	return &CartService{client: client}
}

// Get retrieves a cart by ID.
func (s *CartService) Get(ctx context.Context, id string) (*Cart, error) {
	path := "/carts/" + id
	var cart Cart
	if err := s.client.GetJSON(ctx, path, &cart); err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}
	return &cart, nil
}

// Create creates a new cart.
func (s *CartService) Create(ctx context.Context, customerID string) (*Cart, error) {
	path := "/carts"
	body := map[string]string{"customer_id": customerID}
	var cart Cart
	if err := s.client.PostJSON(ctx, path, body, &cart); err != nil {
		return nil, fmt.Errorf("failed to create cart: %w", err)
	}
	return &cart, nil
}

// AddItem adds a product to the cart per API spec /carts/{id}/{product_id}.
func (s *CartService) AddItem(ctx context.Context, cartID, productID string, quantity int, addons []interface{}, variants []map[string]any, scheduledAt, storeLocation string) (*Cart, error) {
	path := "/carts/" + cartID + "/" + productID
	body := CartItemRequest{
		ProductID:     productID,
		Quantity:      quantity,
		Addons:        addons,
		Variants:      variants,
		ScheduledAt:   scheduledAt,
		StoreLocation: storeLocation,
	}
	var cart Cart
	if err := s.client.PostJSON(ctx, path, body, &cart); err != nil {
		return nil, fmt.Errorf("failed to add item to cart: %w", err)
	}
	return &cart, nil
}

// UpdateItem updates an item in the cart per API spec /carts/{id}/{line_item_id}.
func (s *CartService) UpdateItem(ctx context.Context, cartID, lineItemID string, quantity int, addons []interface{}, variants []map[string]any) (*Cart, error) {
	path := "/carts/" + cartID + "/" + lineItemID
	body := map[string]interface{}{
		"quantity": quantity,
		"addons":   addons,
		"variants": variants,
	}
	var cart Cart
	if err := s.client.PutJSON(ctx, path, body, &cart); err != nil {
		return nil, fmt.Errorf("failed to update item: %w", err)
	}
	return &cart, nil
}

// RemoveItem removes an item from the cart per API spec /carts/{id}/{line_item_id}.
func (s *CartService) RemoveItem(ctx context.Context, cartID, lineItemID string) (*Cart, error) {
	path := "/carts/" + cartID + "/" + lineItemID
	var cart Cart
	if err := s.client.DeleteJSON(ctx, path, &cart); err != nil {
		return nil, fmt.Errorf("failed to remove item from cart: %w", err)
	}
	return &cart, nil
}

// Clear removes all items from the cart per API spec /carts/{id}/empty.
func (s *CartService) Clear(ctx context.Context, cartID string) error {
	path := "/carts/" + cartID + "/empty"
	var resp map[string]any
	if err := s.client.PutJSON(ctx, path, struct{}{}, &resp); err != nil {
		return fmt.Errorf("failed to clear cart: %w", err)
	}
	return nil
}

// Checkout processes the cart checkout and returns an order.
func (s *CartService) Checkout(ctx context.Context, cartID string, req CheckoutRequest) (*Order, error) {
	path := "/carts/" + cartID + "/checkout"
	var order Order
	if err := s.client.PostJSON(ctx, path, req, &order); err != nil {
		return nil, fmt.Errorf("failed to checkout: %w", err)
	}
	return &order, nil
}
