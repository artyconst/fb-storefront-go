package customer

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

func TestCustomerService_ListPlaces(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Customer-Token") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Missing customer token"})
			return
		}

		if r.URL.Path != "/storefront/v1/customers/places" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"items": []map[string]interface{}{
				{"id": "1", "name": "Place 1"},
				{"id": "2", "name": "Place 2"},
			},
		})
	})

	client := setupTestClient(t, handler)
	service := NewCustomerService(client)

	t.Run("success with valid token", func(t *testing.T) {
		places, err := service.ListPlaces(context.Background(), "customer_token_123")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if len(places) != 2 {
			t.Errorf("Expected 2 places, got %d", len(places))
		}
	})

	t.Run("fails without token", func(t *testing.T) {
		places, err := service.ListPlaces(context.Background(), "")
		if err == nil {
			t.Fatal("Expected error for missing token")
		}
		if places != nil {
			t.Error("Expected nil places on error")
		}
	})
}

func TestCustomerService_ListOrders(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Customer-Token") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Missing customer token"})
			return
		}

		queryParams := r.URL.Query()
		status := queryParams.Get("status")
		limit := queryParams.Get("limit")
		offset := queryParams.Get("offset")

		response := map[string]interface{}{
			"items": []map[string]interface{}{
				{"id": "1", "status": status, "limit": limit, "offset": offset},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	client := setupTestClient(t, handler)
	service := NewCustomerService(client)

	t.Run("success with valid token and options", func(t *testing.T) {
		status := "completed"
		opts := []OrderOption{WithOrderLimit(10), WithOffset(5), WithStatus(status)}
		orders, err := service.ListOrders(context.Background(), "customer_token_123", opts...)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if len(orders) != 1 {
			t.Errorf("Expected 1 order, got %d", len(orders))
		}
	})

	t.Run("fails without token", func(t *testing.T) {
		orders, err := service.ListOrders(context.Background(), "")
		if err == nil {
			t.Fatal("Expected error for missing token")
		}
		if orders != nil {
			t.Error("Expected nil orders on error")
		}
	})

	t.Run("fails with valid status pointer", func(t *testing.T) {
		status := "pending"
		opts := []OrderOption{WithStatus(status)}
		orders, err := service.ListOrders(context.Background(), "token_123", opts...)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if orders == nil {
			t.Error("Expected non-nil orders")
		}
	})
}

func TestCustomerService_RequestCreationCode(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/storefront/v1/customers/request-creation-code" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var req RequestCreationCodeRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Identity == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Identity is required"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Creation code sent successfully",
		})
	})

	client := setupTestClient(t, handler)
	service := NewCustomerService(client)

	t.Run("success with valid identity", func(t *testing.T) {
		req := RequestCreationCodeRequest{Identity: "user@example.com"}
		err := service.RequestCreationCode(context.Background(), req)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
	})

	t.Run("fails without identity", func(t *testing.T) {
		req := RequestCreationCodeRequest{Identity: ""}
		err := service.RequestCreationCode(context.Background(), req)
		if err == nil {
			t.Fatal("Expected error for missing identity")
		}
	})
}

func TestCustomerService_RegisterDevice(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Customer-Token") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Missing customer token"})
			return
		}

		var req RegisterDeviceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if req.DeviceID == "" || req.Platform == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Device ID and platform are required"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(RegisterDeviceResponse{Message: "Device registered successfully"})
	})

	client := setupTestClient(t, handler)
	service := NewCustomerService(client)

	t.Run("success with valid token and request", func(t *testing.T) {
		req := RegisterDeviceRequest{
			DeviceID:  "device_123",
			Platform:  "ios",
			PushToken: "apns_token_xyz",
		}

		response, err := service.RegisterDevice(context.Background(), "customer_token_123", req)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if response == nil {
			t.Fatal("Expected non-nil response")
		}
	})

	t.Run("fails without token", func(t *testing.T) {
		req := RegisterDeviceRequest{DeviceID: "device_123", Platform: "ios"}
		response, err := service.RegisterDevice(context.Background(), "", req)
		if err == nil {
			t.Fatal("Expected error for missing token")
		}
		if response != nil {
			t.Error("Expected nil response on error")
		}
	})

	t.Run("fails with invalid request", func(t *testing.T) {
		req := RegisterDeviceRequest{DeviceID: "", Platform: ""}
		response, err := service.RegisterDevice(context.Background(), "token_123", req)
		if err == nil {
			t.Fatal("Expected error for missing device ID and platform")
		}
		if response != nil {
			t.Error("Expected nil response on error")
		}
	})
}

func TestCustomerService_Get(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/storefront/v1/customers/cust_123"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":    "cust_123",
			"email": "user@example.com",
		})
	})

	client := setupTestClient(t, handler)
	service := NewCustomerService(client)

	customer, err := service.Get(context.Background(), "cust_123")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if customer == nil {
		t.Fatal("Expected non-nil customer")
	}
}

func TestCustomerService_Get_InvalidID(t *testing.T) {
	client, _ := sf.NewStorefront("sk_test_key")
	service := NewCustomerService(client)

	customer, err := service.Get(context.Background(), "")
	if err == nil {
		t.Fatal("Expected error for empty ID")
	}
	if customer != nil {
		t.Error("Expected nil customer on error")
	}
}
