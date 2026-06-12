package foodtruck

import (
	"context"
	"fmt"
	"net/url"

	sf "github.com/artyconst/fb-storefront-go"
)

// Service handles food truck-related operations.
type Service struct {
	client *sf.StorefrontClient
}

// NewService creates a new FoodTruck service instance.
func NewService(client *sf.StorefrontClient) *Service {
	return &Service{client: client}
}

// List retrieves all food trucks with optional filtering - GET /food-trucks.
func (s *Service) List(ctx context.Context, opts ...ListOption) ([]*FoodTruck, error) {
	options := &ListOptions{}
	for _, o := range opts {
		o(options)
	}

	params := url.Values{}
	if options.Limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", options.Limit))
	}
	if options.Offset > 0 {
		params.Set("offset", fmt.Sprintf("%d", options.Offset))
	}
	if options.Sort != "" {
		params.Set("sort", options.Sort)
	}

	endpoint := "/food-trucks"
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	var foodTrucks []*FoodTruck
	if err := s.client.GetJSON(ctx, endpoint, &foodTrucks); err != nil {
		return nil, fmt.Errorf("failed to list food trucks: %w", err)
	}

	return foodTrucks, nil
}

// Get retrieves a single food truck by ID - GET /food-trucks/{id}.
func (s *Service) Get(ctx context.Context, id string) (*FoodTruck, error) {
	if id == "" {
		return nil, fmt.Errorf("food truck ID is required")
	}

	var foodTruck FoodTruck
	endpoint := "/food-trucks/" + id
	if err := s.client.GetJSON(ctx, endpoint, &foodTruck); err != nil {
		return nil, err
	}
	return &foodTruck, nil
}
