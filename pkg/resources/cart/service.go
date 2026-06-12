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

// derefString returns the value of a string pointer, or empty string if nil.
func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
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

// AddItem adds a product to the cart per API spec /cart/{id}/{product_id}.
func (s *CartService) AddItem(ctx context.Context, cartID string, productID string, quantity int, opts ...AddItemOption) (*Cart, error) {
	options := &AddItemOptions{}
	for _, o := range opts {
		o(options)
	}

	path := "/carts/" + cartID + "/" + productID
	body := CartItemRequest{
		Quantity:      quantity,
		Addons:        options.Addons,
		Variants:      options.Variants,
		ScheduledAt:   derefString(options.ScheduledAt),
		StoreLocation: derefString(options.StoreLocation),
	}

	var cart Cart
	if err := s.client.PostJSON(ctx, path, body, &cart); err != nil {
		return nil, fmt.Errorf("failed to add item to cart: %w", err)
	}
	return &cart, nil
}

// UpdateItem updates an item in the cart per API spec /cart/{id}/{line_item_id}.
func (s *CartService) UpdateItem(ctx context.Context, cartID string, lineItemID string, quantity int, opts ...UpdateItemOption) (*Cart, error) {
	options := &UpdateItemOptions{}
	for _, o := range opts {
		o(options)
	}

	path := "/carts/" + cartID + "/" + lineItemID
	body := UpdateItemRequest{
		Quantity: quantity,
		Addons:   options.Addons,
		Variants: options.Variants,
	}

	var cart Cart
	if err := s.client.PutJSON(ctx, path, body, &cart); err != nil {
		return nil, fmt.Errorf("failed to update item: %w", err)
	}
	return &cart, nil
}

// RemoveItem removes an item from the cart per API spec /cart/{id}/{line_item_id}.
func (s *CartService) RemoveItem(ctx context.Context, cartID, lineItemID string) (*Cart, error) {
	path := "/carts/" + cartID + "/" + lineItemID
	var cart Cart
	if err := s.client.DeleteJSON(ctx, path, &cart); err != nil {
		return nil, fmt.Errorf("failed to remove item from cart: %w", err)
	}
	return &cart, nil
}

// Clear removes all items from the cart per API spec /cart/{id}/empty.
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

// AddProduct adds a product to the cart per API spec /carts/{cart_id}/{product_id}.
func (s *CartService) AddProduct(ctx context.Context, cartID string, productID string, opts ...AddItemOption) (*Cart, error) {
	if cartID == "" || productID == "" {
		return nil, fmt.Errorf("cart ID and product ID are required: %w", sf.ErrInvalidRequest)
	}

	options := &AddItemOptions{}
	for _, o := range opts {
		o(options)
	}

	path := "/carts/" + cartID + "/" + productID
	body := CartItemRequest{
		Quantity:      int(options.Quantity),
		Addons:        options.Addons,
		Variants:      options.Variants,
		ScheduledAt:   derefString(options.ScheduledAt),
		StoreLocation: derefString(options.StoreLocation),
	}

	var cart Cart
	if err := s.client.PostJSON(ctx, path, body, &cart); err != nil {
		return nil, fmt.Errorf("failed to add product to cart: %w", err)
	}
	return &cart, nil
}

// UpdateLineItem updates a line item in the cart per API spec /carts/{cart_id}/{line_item_id}.
func (s *CartService) UpdateLineItem(ctx context.Context, cartID string, lineItemID string, opts ...UpdateItemOption) (*Cart, error) {
	if cartID == "" || lineItemID == "" {
		return nil, fmt.Errorf("cart ID and line item ID are required: %w", sf.ErrInvalidRequest)
	}

	options := &UpdateItemOptions{}
	for _, o := range opts {
		o(options)
	}

	path := "/carts/" + cartID + "/" + lineItemID
	body := UpdateItemRequest{
		Quantity: int(options.Quantity),
		Addons:   options.Addons,
		Variants: options.Variants,
	}

	var cart Cart
	if err := s.client.PutJSON(ctx, path, body, &cart); err != nil {
		return nil, fmt.Errorf("failed to update line item: %w", err)
	}
	return &cart, nil
}

// RemoveLineItem removes a line item from the cart per API spec /carts/{cart_id}/{line_item_id}.
func (s *CartService) RemoveLineItem(ctx context.Context, cartID string, lineItemID string) (*Cart, error) {
	if cartID == "" || lineItemID == "" {
		return nil, fmt.Errorf("cart ID and line item ID are required: %w", sf.ErrInvalidRequest)
	}

	path := "/carts/" + cartID + "/" + lineItemID
	var cart Cart
	if err := s.client.DeleteJSON(ctx, path, &cart); err != nil {
		return nil, fmt.Errorf("failed to remove line item from cart: %w", err)
	}
	return &cart, nil
}

// EmptyCart removes all items from the cart per API spec /carts/{cart_id}/empty.
func (s *CartService) EmptyCart(ctx context.Context, cartID string) (*Cart, error) {
	if cartID == "" {
		return nil, fmt.Errorf("cart ID is required: %w", sf.ErrInvalidRequest)
	}

	path := "/carts/" + cartID + "/empty"
	var cart Cart
	if err := s.client.PutJSON(ctx, path, struct{}{}, &cart); err != nil {
		return nil, fmt.Errorf("failed to empty cart: %w", err)
	}
	return &cart, nil
}

// DeleteCart deletes the entire cart - DELETE /carts
func (s *CartService) DeleteCart(ctx context.Context) error {
	if err := s.client.DeleteJSON(ctx, "/carts", nil); err != nil {
		return fmt.Errorf("failed to delete cart: %w", err)
	}
	return nil
}
