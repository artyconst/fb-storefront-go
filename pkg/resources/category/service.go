package category

import (
	"context"
	"fmt"
	"net/url"

	sf "github.com/artyconst/fb-storefront-go"
)

// CategoryService handles category-related operations.
type CategoryService struct {
	client *sf.StorefrontClient
}

// NewCategoryService creates a new Category service instance.
func NewCategoryService(client *sf.StorefrontClient) *CategoryService {
	return &CategoryService{client: client}
}

// List retrieves a list of categories with optional pagination and filtering.
func (s *CategoryService) List(ctx context.Context, opts ...ListOption) ([]*Category, error) {
	options := &ListOptions{}
	for _, o := range opts {
		o(options)
	}

	params := url.Values{}
	if options.Offset > 0 {
		params.Set("offset", fmt.Sprintf("%d", options.Offset))
	}
	if options.Limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", options.Limit))
	}
	if options.Search != "" {
		params.Set("search", options.Search)
	}
	if options.ParentID != nil && *options.ParentID != "" {
		params.Set("parent_id", *options.ParentID)
	}

	endpoint := "/categories"
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	var categories []*Category
	if err := s.client.GetJSON(ctx, endpoint, &categories); err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}

	return categories, nil
}

// Get retrieves a category by ID.
func (s *CategoryService) Get(ctx context.Context, id string) (*Category, error) {
	if id == "" {
		return nil, fmt.Errorf("category ID is required")
	}

	var category Category
	endpoint := "/categories/" + id
	if err := s.client.GetJSON(ctx, endpoint, &category); err != nil {
		return nil, err
	}
	return &category, nil
}
