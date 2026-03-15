package customer

import (
	"context"
	"fmt"
	"net/url"

	sf "github.com/artyconst/fb-storefront-go"
)

// CustomerService handles customer-related operations including authentication.
type CustomerService struct {
	client *sf.StorefrontClient
}

// NewCustomerService creates a new Customer service instance.
func NewCustomerService(client *sf.StorefrontClient) *CustomerService {
	return &CustomerService{client: client}
}

// Get retrieves a customer by ID.
func (s *CustomerService) Get(ctx context.Context, id string) (*Customer, error) {
	if id == "" {
		return nil, fmt.Errorf("customer ID is required")
	}

	var customer Customer
	endpoint := "/customers/" + id
	if err := s.client.GetJSON(ctx, endpoint, &customer); err != nil {
		return nil, err
	}
	return &customer, nil
}

// Create creates a new customer.
func (s *CustomerService) Create(ctx context.Context, req CustomerCreateRequest) (*Customer, error) {
	var payload struct {
		Name     *string                `json:"name,omitempty"`
		Type     *string                `json:"type,omitempty"`
		Identity string                 `json:"identity,omitempty"`
		Code     *string                `json:"code,omitempty"`
		Title    *string                `json:"title,omitempty"`
		Email    *string                `json:"email,omitempty"`
		Phone    *string                `json:"phone,omitempty"`
		Meta     map[string]interface{} `json:"meta,omitempty"`
	}

	if req.Name != nil {
		payload.Name = req.Name
	}
	if req.Type != nil {
		payload.Type = req.Type
	}
	payload.Identity = req.Identity
	if req.Code != nil {
		payload.Code = req.Code
	}
	if req.Title != nil {
		payload.Title = req.Title
	}
	if req.Email != nil {
		payload.Email = req.Email
	}
	if req.Phone != nil {
		payload.Phone = req.Phone
	}
	if req.Meta != nil {
		payload.Meta = req.Meta
	}

	var customer Customer
	if err := s.client.PostJSON(ctx, "/customers", payload, &customer); err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}
	return &customer, nil
}

// Login authenticates a customer with identity and password.
func (s *CustomerService) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	if req.Identity == "" || req.Password == "" {
		return nil, fmt.Errorf("identity and password are required")
	}

	var payload struct {
		Identity string `json:"identity"`
		Password string `json:"password"`
	}
	payload.Identity = req.Identity
	payload.Password = req.Password

	var response LoginResponse
	if err := s.client.PostJSON(ctx, "/customers/login", payload, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// LoginWithSMS initiates SMS-based authentication.
func (s *CustomerService) LoginWithSMS(ctx context.Context, req SMSSignInRequest) (*LoginResponse, error) {
	if req.Identity == "" {
		return nil, fmt.Errorf("identity is required")
	}

	var response LoginResponse
	if err := s.client.PostJSON(ctx, "/customers/login-with-sms", req, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// VerifySMSCode confirms the SMS code and completes authentication.
func (s *CustomerService) VerifySMSCode(ctx context.Context, req SMSConfirmSignInRequest) (*LoginResponse, error) {
	if req.Identity == "" || req.Code == "" {
		return nil, fmt.Errorf("identity and code are required")
	}

	var response LoginResponse
	if err := s.client.PostJSON(ctx, "/customers/verify-code", req, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// ListPlaces retrieves places for the authenticated customer.
func (s *CustomerService) ListPlaces(ctx context.Context, token string, opts ...PlaceOption) ([]*Place, error) {
	if token == "" {
		return nil, ErrCustomerTokenRequired
	}

	var placeOpts ListPlacesOptions
	for _, opt := range opts {
		opt(&placeOpts)
	}

	var response struct {
		Items []*Place `json:"items"`
	}

	endpoint := "/customers/places"
	queryParams := make(url.Values)
	if placeOpts.Page > 0 {
		queryParams.Set("page", fmt.Sprintf("%d", placeOpts.Page))
	}
	if placeOpts.Limit > 0 {
		queryParams.Set("limit", fmt.Sprintf("%d", placeOpts.Limit))
	}
	if placeOpts.Sort != "" {
		queryParams.Set("sort", placeOpts.Sort)
	}

	if len(queryParams) > 0 {
		endpoint += "?" + queryParams.Encode()
	}

	var requestOpts []sf.RequestOption
	requestOpts = append(requestOpts, sf.WithCustomerToken(token))
	if err := s.client.GetJSON(ctx, endpoint, &response, requestOpts...); err != nil {
		return nil, fmt.Errorf("failed to list places: %w", err)
	}

	return response.Items, nil
}

// ListOrders retrieves orders for the authenticated customer.
func (s *CustomerService) ListOrders(ctx context.Context, token string, opts ...OrderOption) ([]*Order, error) {
	if token == "" {
		return nil, ErrCustomerTokenRequired
	}

	var orderOpts ListOrdersOptions
	for _, opt := range opts {
		opt(&orderOpts)
	}

	var response struct {
		Items []*Order `json:"items"`
	}

	endpoint := "/customers/orders"
	queryParams := make(url.Values)
	if orderOpts.Limit > 0 {
		queryParams.Set("limit", fmt.Sprintf("%d", orderOpts.Limit))
	}
	if orderOpts.Offset > 0 {
		queryParams.Set("offset", fmt.Sprintf("%d", orderOpts.Offset))
	}
	if orderOpts.Status != nil {
		queryParams.Set("status", *orderOpts.Status)
	}
	if orderOpts.Sort != "" {
		queryParams.Set("sort", orderOpts.Sort)
	}

	if len(queryParams) > 0 {
		endpoint += "?" + queryParams.Encode()
	}

	var requestOpts []sf.RequestOption
	requestOpts = append(requestOpts, sf.WithCustomerToken(token))
	if err := s.client.GetJSON(ctx, endpoint, &response, requestOpts...); err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}

	return response.Items, nil
}

// RequestCreationCode initiates the process for a new customer to receive an account creation code.
func (s *CustomerService) RequestCreationCode(ctx context.Context, req RequestCreationCodeRequest) error {
	if req.Identity == "" {
		return fmt.Errorf("identity is required")
	}

	var response struct {
		Message string `json:"message"`
	}

	err := s.client.PostJSON(ctx, "/customers/request-creation-code", req, &response)
	if err != nil {
		return fmt.Errorf("failed to request creation code: %w", err)
	}

	return nil
}

// RegisterDevice registers a device for the authenticated customer.
func (s *CustomerService) RegisterDevice(ctx context.Context, token string, req RegisterDeviceRequest) (*RegisterDeviceResponse, error) {
	if token == "" {
		return nil, ErrCustomerTokenRequired
	}

	var response RegisterDeviceResponse

	var requestOpts []sf.RequestOption
	requestOpts = append(requestOpts, sf.WithCustomerToken(token))

	if err := s.client.PostJSON(ctx, "/customers/register-device", req, &response, requestOpts...); err != nil {
		return nil, fmt.Errorf("failed to register device: %w", err)
	}

	return &response, nil
}
