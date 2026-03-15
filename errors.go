package storefront

import (
	"errors"
	"fmt"
)

// Error types for programmatic error handling
var (
	ErrInvalidAPIKey     = errors.New("invalid or missing API key")
	ErrResourceNotFound  = errors.New("resource not found")
	ErrNetworkTimeout    = errors.New("network timeout occurred")
	ErrInvalidRequest    = errors.New("invalid request parameters")
	ErrPaymentFailed     = errors.New("payment processing failed")
	ErrCartEmpty         = errors.New("cart has no items")
	ErrInsufficientStock = errors.New("insufficient product stock")
)

// APIError represents an error response from the Storefront API
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"status,omitempty"`
}

func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("API error (%s): %s", e.Code, e.Message)
	}
	return "API error: " + e.Message
}

// Is checks if the given error is an APIError
func (e *APIError) Is(target error) bool {
	_, ok := target.(*APIError)
	return ok
}
