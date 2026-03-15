package utils

import (
	"context"
)

// HTTPClient defines the interface for making HTTP requests.
type HTTPClient interface {
	Get(ctx context.Context, path string) (*HTTPResponse, error)
	Post(ctx context.Context, path string, data any) (*HTTPResponse, error)
	Put(ctx context.Context, path string, data any) (*HTTPResponse, error)
	Delete(ctx context.Context, path string) (*HTTPResponse, error)
}

// HTTPResponse represents an HTTP response from the API.
type HTTPResponse struct {
	StatusCode int
	Body       []byte
}
