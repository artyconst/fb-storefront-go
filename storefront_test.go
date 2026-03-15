package storefront

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// normalizeServerURL strips all trailing slashes from the URL for consistent composition.
func normalizeServerURL(url string) string {
	url = strings.TrimSpace(url)
	for strings.HasSuffix(url, "/") {
		url = strings.TrimSuffix(url, "/")
	}
	return url
}

// AboutStoreResponse for testing - same structure as in store package
type AboutStoreResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func TestClientAuthorizationHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test_api_key" {
			t.Errorf("Expected 'Bearer test_api_key', got '%s'", auth)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":"test","name":"Test Store","created_at":"2024-01-01T00:00:00Z","updated_at":"2024-01-01T00:00:00Z"}`))
	}))
	defer server.Close()

	client, err := newStorefrontClient(ClientConfig{
		ServerURL: server.URL,
		APIKey:    "test_api_key",
	})
	if err != nil {
		t.Fatal(err)
	}

	var resp AboutStoreResponse
	err = client.GetJSON(context.Background(), "/about", &resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.ID != "test" {
		t.Errorf("Expected ID 'test', got '%s'", resp.ID)
	}
}

func TestClientMissingAPIKey(t *testing.T) {
	_, err := newStorefrontClient(ClientConfig{
		ServerURL: "https://api.fleetbase.io",
		APIKey:    "",
	})
	if err == nil {
		t.Error("Expected error for missing API key, got nil")
	}
}

func TestAboutStoreResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":"abc123","name":"Test Store","description":"A test store","created_at":"2024-01-01T12:00:00Z","updated_at":"2024-01-02T12:00:00Z"}`))
	}))
	defer server.Close()

	client, err := newStorefrontClient(ClientConfig{
		ServerURL: server.URL,
		APIKey:    "test_key",
	})
	if err != nil {
		t.Fatal(err)
	}

	var resp AboutStoreResponse
	err = client.GetJSON(context.Background(), "/about", &resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.ID != "abc123" {
		t.Errorf("Expected ID 'abc123', got '%s'", resp.ID)
	}
	if resp.Name != "Test Store" {
		t.Errorf("Expected name 'Test Store', got '%s'", resp.Name)
	}
	if resp.Description != "A test store" {
		t.Errorf("Expected description, got '%s'", resp.Description)
	}

	expectedTime, _ := time.Parse(time.RFC3339, "2024-01-01T12:00:00Z")
	if !resp.CreatedAt.Equal(expectedTime) {
		t.Errorf("Expected created_at %v, got %v", expectedTime, resp.CreatedAt)
	}
}

func TestWithAPIHostNormalization(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		expected string
	}{
		{"no trailing slash", "https://api.fleetbase.io", "https://api.fleetbase.io/storefront/v1"},
		{"single trailing slash", "https://api.fleetbase.io/", "https://api.fleetbase.io/storefront/v1"},
		{"multiple trailing slashes", "https://api.fleetbase.io///", "https://api.fleetbase.io/storefront/v1"},
		{"with path component", "https://custom.api.com/storefront", "https://custom.api.com/storefront/storefront/v1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := newStorefrontClient(ClientConfig{
				ServerURL: normalizeServerURL(tt.host),
				APIPath:   "/storefront/v1",
				APIKey:    "test_key",
			})
			if err != nil {
				t.Fatal(err)
			}

			if client.baseURL != tt.expected {
				t.Errorf("Expected baseURL %s, got %s", tt.expected, client.baseURL)
			}
		})
	}
}
