package product

import (
	"context"
	"fmt"
	"net/url"

	sf "github.com/artyconst/fb-storefront-go"
)

// ProductService handles product-related operations.
type ProductService struct {
	client *sf.StorefrontClient
}

// NewProductService creates a new Product service instance.
func NewProductService(client *sf.StorefrontClient) *ProductService {
	return &ProductService{client: client}
}

// List retrieves a list of products with optional filtering.
func (s *ProductService) List(ctx context.Context, opts ...ListOption) ([]*Product, error) {
	options := &ListOptions{}
	for _, o := range opts {
		o(options)
	}

	params := url.Values{}
	if options.Offset > 0 {
		params.Set("offset", fmt.Sprintf("%d", options.Offset))
	}
	if options.Category != "" {
		params.Set("category_id", options.Category)
	}
	if options.SortBy != "" {
		params.Set("sort_by", options.SortBy)
	}
	if options.Order != "" {
		params.Set("order", options.Order)
	}

	endpoint := "/products"
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	var products []*Product
	if err := s.client.GetJSON(ctx, endpoint, &products); err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	return products, nil
}

// Get retrieves a single product by ID.
func (s *ProductService) Get(ctx context.Context, id string) (*Product, error) {
	if id == "" {
		return nil, fmt.Errorf("product ID is required")
	}

	var product Product
	endpoint := "/products/" + id
	if err := s.client.GetJSON(ctx, endpoint, &product); err != nil {
		return nil, err
	}
	return &product, nil
}

// FindByCategory retrieves products in a specific category.
func (s *ProductService) FindByCategory(ctx context.Context, categoryID string) ([]*Product, error) {
	if categoryID == "" {
		return nil, fmt.Errorf("category ID cannot be empty")
	}

	return s.List(ctx, WithCategory(categoryID))
}
