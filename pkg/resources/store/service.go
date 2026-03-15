package store

import (
	"context"
	"fmt"
	"net/url"

	sf "github.com/artyconst/fb-storefront-go"
)

// StoreService provides store and gateway operations
type StoreService struct {
	client *sf.StorefrontClient
}

// NewStoreService creates a new StoreService instance
func NewStoreService(client *sf.StorefrontClient) *StoreService {
	return &StoreService{client: client}
}

// About retrieves about store information - GET /about
func (s *StoreService) About(ctx context.Context) (*AboutStoreResponse, error) {
	var resp AboutStoreResponse
	if err := s.client.GetJSON(ctx, "/about", &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListGateways lists payment gateways - GET /gateways
func (s *StoreService) ListGateways(ctx context.Context, opts ...GatewaysOption) (*PaymentGatewaysResponse, error) {
	options := &ListGatewaysOpts{}
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

	endpoint := "/gateways"
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	var gateways []PaymentGateway
	if err := s.client.GetJSON(ctx, endpoint, &gateways); err != nil {
		return nil, fmt.Errorf("failed to list gateways: %w", err)
	}

	resp := PaymentGatewaysResponse{
		Data:   gateways,
		Limit:  options.Limit,
		Offset: options.Offset,
	}
	return &resp, nil
}

// GetGateway retrieves specific gateway - GET /gateways/{id}
func (s *StoreService) GetGateway(ctx context.Context, gatewayID string) (*PaymentGateway, error) {
	if gatewayID == "" {
		return nil, fmt.Errorf("gateway ID is required")
	}

	var resp PaymentGateway
	endpoint := "/gateways/" + gatewayID
	if err := s.client.GetJSON(ctx, endpoint, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
