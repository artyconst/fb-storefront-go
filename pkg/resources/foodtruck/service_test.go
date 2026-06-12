package foodtruck

import (
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

func TestService_List(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]*FoodTruck{
			{ID: "ft_1", Online: true, Status: "active"},
			{ID: "ft_2", Online: false, Status: "offline"},
		})
	})

	client := setupTestClient(t, handler)
	service := NewService(client)

	t.Run("success list food trucks", func(t *testing.T) {
		trucks, err := service.List(t.Context())
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if len(trucks) != 2 {
			t.Errorf("Expected 2 food trucks, got %d", len(trucks))
		}
	})

	t.Run("success list with options", func(t *testing.T) {
		trucks, err := service.List(t.Context(), WithLimit(10), WithOffset(5), WithSort("name"))
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if trucks == nil {
			t.Error("Expected non-nil food trucks")
		}
	})

	t.Run("fails with server error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Server error"})
		})
		client := setupTestClient(t, handler)
		service := NewService(client)

		trucks, err := service.List(t.Context())
		if err == nil {
			t.Fatal("Expected error for server failure")
		}
		if trucks != nil {
			t.Error("Expected nil food trucks on error")
		}
	})
}

func TestService_Get(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(FoodTruck{
			ID:     "ft_123",
			Status: "active",
			Vehicle: &Vehicle{
				ID:   "veh_1",
				Name: "Taco Truck",
			},
		})
	})

	client := setupTestClient(t, handler)
	service := NewService(client)

	t.Run("success get food truck", func(t *testing.T) {
		truck, err := service.Get(t.Context(), "ft_123")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if truck == nil {
			t.Fatal("Expected non-nil food truck")
		}
		if truck.ID != "ft_123" {
			t.Errorf("Expected ID 'ft_123', got '%s'", truck.ID)
		}
	})

	t.Run("fails with empty id", func(t *testing.T) {
		truck, err := service.Get(t.Context(), "")
		if err == nil {
			t.Fatal("Expected error for empty ID")
		}
		if truck != nil {
			t.Error("Expected nil food truck on error")
		}
	})

	t.Run("fails with server error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Not found"})
		})
		client := setupTestClient(t, handler)
		service := NewService(client)

		truck, err := service.Get(t.Context(), "ft_123")
		if err == nil {
			t.Fatal("Expected error for not found")
		}
		if truck != nil {
			t.Error("Expected nil food truck on error")
		}
	})
}

func TestService_ListWithOptions(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]*FoodTruck{
			{ID: "ft_1", Online: true},
			{ID: "ft_2", Online: false},
		})
	})

	client := setupTestClient(t, handler)
	service := NewService(client)

	t.Run("success with limit", func(t *testing.T) {
		trucks, err := service.List(t.Context(), WithLimit(10))
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if len(trucks) != 2 {
			t.Errorf("Expected 2 food trucks, got %d", len(trucks))
		}
	})

	t.Run("success with offset", func(t *testing.T) {
		trucks, err := service.List(t.Context(), WithOffset(5))
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if len(trucks) != 2 {
			t.Errorf("Expected 2 food trucks, got %d", len(trucks))
		}
	})

	t.Run("success with sort", func(t *testing.T) {
		trucks, err := service.List(t.Context(), WithSort("name"))
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if len(trucks) != 2 {
			t.Errorf("Expected 2 food trucks, got %d", len(trucks))
		}
	})

	t.Run("fails with server error and options", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Server error"})
		})
		client := setupTestClient(t, handler)
		service := NewService(client)

		trucks, err := service.List(t.Context(), WithLimit(10))
		if err == nil {
			t.Fatal("Expected error for server failure")
		}
		if trucks != nil {
			t.Error("Expected nil food trucks on error")
		}
	})
}
