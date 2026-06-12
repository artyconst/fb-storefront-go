package customer

import (
	"errors"
	"time"
)

// Customer represents a customer in the system.
type Customer struct {
	ID        string    `json:"id"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	Phone     *string   `json:"phone,omitempty"`
	AvatarURL *string   `json:"avatar_url,omitempty"`
	CreatedAt time.Time `json:"created_at,omitzero"`
	UpdatedAt time.Time `json:"updated_at,omitzero"`
}

// ErrCustomerTokenRequired is returned when a customer token is required but not provided.
var ErrCustomerTokenRequired = errors.New("customer token is required")

// LoginRequest holds credentials for the Login operation.
type LoginRequest struct {
	Identity string `json:"identity"`
	Password string `json:"password"`
}

// LoginResponse represents the response from a login operation.
type LoginResponse struct {
	Customer *Customer `json:"customer,omitempty"`
	Token    string    `json:"token,omitempty"`
}

// SMSSignInRequest holds parameters for SMS-based authentication initiation.
type SMSSignInRequest struct {
	Identity string `json:"identity"`
}

// SMSConfirmSignInRequest holds parameters for confirming an SMS code.
type SMSConfirmSignInRequest struct {
	Identity string `json:"identity"`
	Code     string `json:"code"`
}

// LoginOptions holds optional parameters for the Login operation.
type LoginOptions struct{}

// VerifyCodeOptions holds optional parameters for the VerifyCode operation.
type VerifyCodeOptions struct{}

// OAuthOptions holds optional parameters for OAuth login operations.
type OAuthOptions struct{}

// --- Non-P0 types (existing, kept for compatibility) ---

// CustomerCreateRequest for creating customers.
type CustomerCreateRequest struct {
	Name     *string                `json:"name,omitempty"`
	Type     *string                `json:"type,omitempty"`
	Identity string                 `json:"identity,omitempty"`
	Code     *string                `json:"code,omitempty"`
	Title    *string                `json:"title,omitempty"`
	Email    *string                `json:"email,omitempty"`
	Phone    *string                `json:"phone,omitempty"`
	Meta     map[string]any         `json:"meta,omitempty"`
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

// StripeSetupIntent represents a Stripe setup intent response.
type StripeSetupIntent struct {
	ClientSecret string `json:"client_secret"`
	IntentID     string `json:"intent_id,omitempty"`
}

// UpdateCustomerOptions holds optional parameters for the Update operation.
type UpdateCustomerOptions struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
	Phone *string `json:"phone,omitempty"`
}

// UpdateOption is a functional option for the Update method.
type UpdateOption func(*UpdateCustomerOptions)

// WithName sets the customer's name for update.
func WithName(name string) UpdateOption {
	return func(o *UpdateCustomerOptions) { o.Name = &name }
}

// WithEmail sets the customer's email for update.
func WithEmail(email string) UpdateOption {
	return func(o *UpdateCustomerOptions) { o.Email = &email }
}

// WithPhone sets the customer's phone for update.
func WithPhone(phone string) UpdateOption {
	return func(o *UpdateCustomerOptions) { o.Phone = &phone }
}
