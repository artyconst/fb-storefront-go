package customer

import "fmt"

// Customer represents a customer in the system.
type Customer struct {
	ID        string  `json:"id"`
	Name      *string `json:"name,omitempty"`
	Email     *string `json:"email,omitempty"`
	Phone     *string `json:"phone,omitempty"`
	Type      *string `json:"type,omitempty"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

// LoginRequest for customer authentication.
type LoginRequest struct {
	Identity string `json:"identity"`
	Password string `json:"password"`
}

// SMSSignInRequest for SMS-based authentication initiation.
type SMSSignInRequest struct {
	Identity string `json:"identity"`
}

// SMSConfirmSignInRequest for confirming SMS code.
type SMSConfirmSignInRequest struct {
	Identity string `json:"identity"`
	Code     string `json:"code"`
}

// LoginResponse contains authentication tokens and customer info.
type LoginResponse struct {
	Token     string    `json:"token"`
	Customer  *Customer `json:"customer"`
	ExpiresAt string    `json:"expires_at,omitempty"`
}

// CustomerCreateRequest for creating customers.
type CustomerCreateRequest struct {
	Name     *string                `json:"name,omitempty"`
	Type     *string                `json:"type,omitempty"`
	Identity string                 `json:"identity,omitempty"`
	Code     *string                `json:"code,omitempty"`
	Title    *string                `json:"title,omitempty"`
	Email    *string                `json:"email,omitempty"`
	Phone    *string                `json:"phone,omitempty"`
	Meta     map[string]interface{} `json:"meta,omitempty"`
}

// Place represents a customer's saved location/destination.
type Place struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Order represents a customer order in the system.
type Order struct {
	ID         string `json:"id"`
	Status     string `json:"status"`
	Total      int64  `json:"total"`
	Currency   string `json:"currency"`
	ItemsCount int    `json:"items_count"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// ListPlacesOptions represents pagination and sorting parameters for place listing.
type ListPlacesOptions struct {
	Page  int    `json:"page,omitempty"`
	Limit int    `json:"limit,omitempty"`
	Sort  string `json:"sort,omitempty"`
}

// PlaceOption is a functional option for ListPlaces.
type PlaceOption func(*ListPlacesOptions)

// WithPage sets the page number for pagination.
func WithPage(page int) PlaceOption {
	return func(o *ListPlacesOptions) {
		o.Page = page
	}
}

// WithLimit sets the limit for results per page.
func WithLimit(limit int) PlaceOption {
	return func(o *ListPlacesOptions) {
		o.Limit = limit
	}
}

// WithSort sets the sort field.
func WithSort(sort string) PlaceOption {
	return func(o *ListPlacesOptions) {
		o.Sort = sort
	}
}

// ListOrdersOptions represents pagination and filtering parameters for order listing.
type ListOrdersOptions struct {
	Limit  int     `json:"limit,omitempty"`
	Offset int     `json:"offset,omitempty"`
	Status *string `json:"status,omitempty"`
	Sort   string  `json:"sort,omitempty"`
}

// OrderOption is a functional option for ListOrders.
type OrderOption func(*ListOrdersOptions)

// WithOrderLimit sets the limit for results per page.
func WithOrderLimit(limit int) OrderOption {
	return func(o *ListOrdersOptions) {
		o.Limit = limit
	}
}

// WithOffset sets the offset for pagination.
func WithOffset(offset int) OrderOption {
	return func(o *ListOrdersOptions) {
		o.Offset = offset
	}
}

// WithStatus filters by order status.
func WithStatus(status string) OrderOption {
	return func(o *ListOrdersOptions) {
		o.Status = &status
	}
}

// WithOrderSort sets the sort field.
func WithOrderSort(sort string) OrderOption {
	return func(o *ListOrdersOptions) {
		o.Sort = sort
	}
}

// RequestCreationCodeRequest represents POST payload for verification code request.
type RequestCreationCodeRequest struct {
	Identity string `json:"identity"` // Email address or phone number to send code to
	Mode     string `json:"mode"`     // Verification mode: "email" or "phone" (required)
}

// RegisterDeviceRequest represents device registration payload for push notifications.
type RegisterDeviceRequest struct {
	DeviceID  string `json:"device_id"`  // Unique device identifier (UUID recommended)
	Platform  string `json:"platform"`   // Target platform: "ios" or "android" (required)
	PushToken string `json:"push_token"` // APNs token (iOS) or FCM token (Android) for push notifications
}

// RegisterDeviceResponse represents the response for device registration.
type RegisterDeviceResponse struct {
	Message string `json:"message"`
}

// Customer service specific errors.
var (
	ErrCustomerTokenRequired = fmt.Errorf("customer authentication token is required")
)
