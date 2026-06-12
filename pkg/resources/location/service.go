package location

import (
	"context"
	"fmt"
	"net/url"

	sf "github.com/artyconst/fb-storefront-go"
)

// Service provides location operations.
type Service struct {
	client *sf.StorefrontClient
}

// NewService creates a new Service instance.
func NewService(client *sf.StorefrontClient) *Service {
	return &Service{client: client}
}

// List retrieves all locations for a store - GET /locations?store={storeID}.
func (s *Service) List(ctx context.Context, storeID string, opts ...ListOption) ([]*Location, error) {
	if storeID == "" {
		return nil, fmt.Errorf("store ID is required")
	}

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
	params.Set("store", storeID)

	endpoint := "/locations"
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	var locations []*Location
	if err := s.client.GetJSON(ctx, endpoint, &locations); err != nil {
		return nil, fmt.Errorf("failed to list locations: %w", err)
	}

	return locations, nil
}

// Get retrieves a single location by ID - GET /locations/{id}?store={storeID}.
func (s *Service) Get(ctx context.Context, storeID, locationID string, opts ...ListOption) (*Location, error) {
	if storeID == "" {
		return nil, fmt.Errorf("store ID is required")
	}
	if locationID == "" {
		return nil, fmt.Errorf("location ID is required")
	}

	options := &ListOptions{}
	for _, o := range opts {
		o(options)
	}

	params := url.Values{}
	params.Set("store", storeID)

	endpoint := "/locations/" + locationID
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	var loc Location
	if err := s.client.GetJSON(ctx, endpoint, &loc); err != nil {
		return nil, fmt.Errorf("failed to get location: %w", err)
	}

	return &loc, nil
}

// ListLocations retrieves all store locations via the dedicated endpoint - GET /store-locations.
func (s *Service) ListLocations(ctx context.Context) ([]*StoreLocation, error) {
	var locations []*StoreLocation
	if err := s.client.GetJSON(ctx, "/store-locations", &locations); err != nil {
		return nil, fmt.Errorf("failed to list store locations: %w", err)
	}

	return locations, nil
}
