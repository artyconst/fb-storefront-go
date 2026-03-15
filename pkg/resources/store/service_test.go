package store

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	sf "github.com/artyconst/fb-storefront-go"
)

func setupTestClient(t *testing.T, endpoint, method string, statusCode int, responseBody interface{}) (*sf.StorefrontClient, func()) {
	t.Helper()

	responseBytes, err := json.Marshal(responseBody)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == endpoint && r.Method == method {
			w.WriteHeader(statusCode)
			w.Header().Set("Content-Type", "application/json")
			w.Write(responseBytes)
		} else {
			http.NotFound(w, r)
		}
	}))

	client, err := sf.NewStorefront("sk_test_key", sf.WithAPIHost(server.URL))
	if err != nil {
		server.Close()
		t.Fatalf("Failed to create test client: %v", err)
	}

	return client, func() {
		server.Close()
	}
}

func TestStoreService_About_Success(t *testing.T) {
	resp := &AboutStoreResponse{
		ID:          "store_123",
		Name:        "Test Store",
		Description: "A test store",
		Currency:    "USD",
		Country:     "US",
		Rating:      4.5,
		Online:      true,
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/storefront/v1/about" && r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		} else {
			http.NotFound(w, r)
		}
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client, err := sf.NewStorefront("sk_test_key", sf.WithAPIHost(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	service := NewStoreService(client)

	ctx := context.Background()
	result, err := service.About(ctx)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result.ID != "store_123" {
		t.Errorf("Expected ID store_123, got %s", result.ID)
	}

	if result.Name != "Test Store" {
		t.Errorf("Expected name Test Store, got %s", result.Name)
	}

	if result.Rating != 4.5 {
		t.Errorf("Expected rating 4.5, got %f", result.Rating)
	}
}

func TestStoreService_About_Error(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/storefront/v1/about" && r.Method == http.MethodGet {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "server error"})
		} else {
			http.NotFound(w, r)
		}
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client, err := sf.NewStorefront("sk_test_key", sf.WithAPIHost(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	service := NewStoreService(client)

	ctx := context.Background()
	result, err := service.About(ctx)
	if err == nil {
		t.Fatal("Expected error for server failure, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil on error, got: %+v", result)
	}
}
