package category

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	sf "github.com/artyconst/fb-storefront-go"
)

func setupTestClient(t *testing.T, handler http.Handler) *sf.StorefrontClient {
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	client, err := sf.NewStorefront("sk_test_key", sf.WithAPIHost(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	return client
}

func TestCategoryService_List(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]*Category{
			{ID: "cat_1", Name: "Electronics"},
			{ID: "cat_2", Name: "Clothing"},
		})
	})

	client := setupTestClient(t, handler)
	service := NewCategoryService(client)

	t.Run("success list categories", func(t *testing.T) {
		categories, err := service.List(context.Background())
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if len(categories) != 2 {
			t.Errorf("Expected 2 categories, got %d", len(categories))
		}
	})

	t.Run("success with options", func(t *testing.T) {
		categories, err := service.List(context.Background(), WithLimit(10), WithOffset(5))
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if len(categories) != 2 {
			t.Errorf("Expected 2 categories, got %d", len(categories))
		}
	})

	t.Run("fails with search option", func(t *testing.T) {
		categories, err := service.List(context.Background(), WithSearch("test"))
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if categories == nil {
			t.Error("Expected non-nil categories")
		}
	})

	t.Run("fails with server error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Server error"})
		})
		client := setupTestClient(t, handler)
		service := NewCategoryService(client)

		categories, err := service.List(context.Background())
		if err == nil {
			t.Fatal("Expected error for server failure")
		}
		if categories != nil {
			t.Error("Expected nil categories on error")
		}
	})
}

func TestCategoryService_Get(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":    "cat_123",
			"name":  "Electronics",
			"count": 42,
		})
	})

	client := setupTestClient(t, handler)
	service := NewCategoryService(client)

	t.Run("success get category", func(t *testing.T) {
		category, err := service.Get(context.Background(), "cat_123")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if category == nil {
			t.Fatal("Expected non-nil category")
		}
		if category.ID != "cat_123" {
			t.Errorf("Expected ID 'cat_123', got '%s'", category.ID)
		}
	})

	t.Run("fails with empty id", func(t *testing.T) {
		category, err := service.Get(context.Background(), "")
		if err == nil {
			t.Fatal("Expected error for empty ID")
		}
		if category != nil {
			t.Error("Expected nil category on error")
		}
	})

	t.Run("fails with server error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Not found"})
		})
		client := setupTestClient(t, handler)
		service := NewCategoryService(client)

		category, err := service.Get(context.Background(), "cat_123")
		if err == nil {
			t.Fatal("Expected error for not found")
		}
		if category != nil {
			t.Error("Expected nil category on error")
		}
	})

	t.Run("fails with invalid id format", func(t *testing.T) {
		category, err := service.Get(context.Background(), "invalid-id-with-dashes")
		if err == nil {
			t.Fatal("Expected error for invalid ID")
		}
		if category != nil {
			t.Error("Expected nil category on error")
		}
	})
}
