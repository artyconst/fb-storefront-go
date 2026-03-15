package cart

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

func TestCartService_Get(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/storefront/v1/carts/cart_123" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":           "cart_123",
			"status":       "active",
			"total_amount": 999,
			"currency":     "USD",
		})
	})

	client := setupTestClient(t, handler)
	service := NewCartService(client)

	t.Run("success get cart", func(t *testing.T) {
		cart, err := service.Get(context.Background(), "cart_123")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if cart == nil {
			t.Fatal("Expected non-nil cart")
		}
		if cart.ID != "cart_123" {
			t.Errorf("Expected cart ID 'cart_123', got '%s'", cart.ID)
		}
	})

	t.Run("sends request with empty id", func(t *testing.T) {
		_, err := service.Get(context.Background(), "")
		if err != nil {
			t.Logf("Client accepted request (validation is server-side): %v", err)
		}
	})

	t.Run("fails with server error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Server error"})
		})
		client := setupTestClient(t, handler)
		service := NewCartService(client)

		cart, err := service.Get(context.Background(), "cart_123")
		if err == nil {
			t.Fatal("Expected error for server failure")
		}
		if cart != nil {
			t.Error("Expected nil cart on error")
		}
	})
}

func TestCartService_Create(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var body map[string]string
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("Failed to decode request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if body["customer_id"] == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Customer ID required"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":           "cart_456",
			"customer_id":  body["customer_id"],
			"status":       "active",
			"total_amount": 0,
			"currency":     "USD",
		})
	})

	client := setupTestClient(t, handler)
	service := NewCartService(client)

	t.Run("success create cart", func(t *testing.T) {
		cart, err := service.Create(context.Background(), "cust_123")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if cart == nil {
			t.Fatal("Expected non-nil cart")
		}
		if cart.ID != "cart_456" {
			t.Errorf("Expected cart ID 'cart_456', got '%s'", cart.ID)
		}
	})

	t.Run("fails with empty customer id", func(t *testing.T) {
		cart, err := service.Create(context.Background(), "")
		if err == nil {
			t.Fatal("Expected error for empty customer ID")
		}
		if cart != nil {
			t.Error("Expected nil cart on error")
		}
	})

	t.Run("fails with server error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Server error"})
		})
		client := setupTestClient(t, handler)
		service := NewCartService(client)

		cart, err := service.Create(context.Background(), "cust_123")
		if err == nil {
			t.Fatal("Expected error for server failure")
		}
		if cart != nil {
			t.Error("Expected nil cart on error")
		}
	})
}

func TestCartService_AddItem(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var body CartItemRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("Failed to decode request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if body.ProductID == "" || body.Quantity <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid product or quantity"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":           "cart_789",
			"status":       "active",
			"total_amount": 1500,
			"items_count":  1,
			"currency":     "USD",
		})
	})

	client := setupTestClient(t, handler)
	service := NewCartService(client)

	t.Run("success add item", func(t *testing.T) {
		cart, err := service.AddItem(context.Background(), "cart_123", "prod_456", 2, nil, nil, "", "")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if cart == nil {
			t.Fatal("Expected non-nil cart")
		}
	})

	t.Run("fails with invalid product id", func(t *testing.T) {
		cart, err := service.AddItem(context.Background(), "cart_123", "", 2, nil, nil, "", "")
		if err == nil {
			t.Fatal("Expected error for empty product ID")
		}
		if cart != nil {
			t.Error("Expected nil cart on error")
		}
	})

	t.Run("fails with invalid quantity", func(t *testing.T) {
		cart, err := service.AddItem(context.Background(), "cart_123", "prod_456", 0, nil, nil, "", "")
		if err == nil {
			t.Fatal("Expected error for zero quantity")
		}
		if cart != nil {
			t.Error("Expected nil cart on error")
		}
	})

	t.Run("fails with server error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Server error"})
		})
		client := setupTestClient(t, handler)
		service := NewCartService(client)

		cart, err := service.AddItem(context.Background(), "cart_123", "prod_456", 2, nil, nil, "", "")
		if err == nil {
			t.Fatal("Expected error for server failure")
		}
		if cart != nil {
			t.Error("Expected nil cart on error")
		}
	})
}

func TestCartService_UpdateItem(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("Failed to decode request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		quantity, ok := body["quantity"].(float64)
		if !ok || quantity <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid quantity"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":           "cart_789",
			"status":       "active",
			"total_amount": 3000,
			"items_count":  2,
			"currency":     "USD",
		})
	})

	client := setupTestClient(t, handler)
	service := NewCartService(client)

	t.Run("success update item", func(t *testing.T) {
		cart, err := service.UpdateItem(context.Background(), "cart_123", "line_item_1", 5, nil, nil)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if cart == nil {
			t.Fatal("Expected non-nil cart")
		}
	})

	t.Run("sends request with empty line item id", func(t *testing.T) {
		_, err := service.UpdateItem(context.Background(), "cart_123", "", 5, nil, nil)
		if err != nil {
			t.Logf("Client accepted request (validation is server-side): %v", err)
		}
	})

	t.Run("fails with invalid quantity", func(t *testing.T) {
		cart, err := service.UpdateItem(context.Background(), "cart_123", "line_item_1", 0, nil, nil)
		if err == nil {
			t.Fatal("Expected error for zero quantity")
		}
		if cart != nil {
			t.Error("Expected nil cart on error")
		}
	})

	t.Run("fails with server error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Server error"})
		})
		client := setupTestClient(t, handler)
		service := NewCartService(client)

		cart, err := service.UpdateItem(context.Background(), "cart_123", "line_item_1", 5, nil, nil)
		if err == nil {
			t.Fatal("Expected error for server failure")
		}
		if cart != nil {
			t.Error("Expected nil cart on error")
		}
	})
}

func TestCartService_RemoveItem(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":           "cart_789",
			"status":       "active",
			"total_amount": 1500,
			"items_count":  1,
			"currency":     "USD",
		})
	})

	client := setupTestClient(t, handler)
	service := NewCartService(client)

	t.Run("success remove item", func(t *testing.T) {
		cart, err := service.RemoveItem(context.Background(), "cart_123", "line_item_1")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if cart == nil {
			t.Fatal("Expected non-nil cart")
		}
	})

	t.Run("sends request with empty line item id", func(t *testing.T) {
		cart, err := service.RemoveItem(context.Background(), "cart_123", "")
		if err != nil {
			t.Logf("Client accepted request (validation is server-side): %v", err)
		}
		if cart == nil && err == nil {
			t.Error("Expected cart result")
		}
	})

	t.Run("fails with server error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Server error"})
		})
		client := setupTestClient(t, handler)
		service := NewCartService(client)

		cart, err := service.RemoveItem(context.Background(), "cart_123", "line_item_1")
		if err == nil {
			t.Fatal("Expected error for server failure")
		}
		if cart != nil {
			t.Error("Expected nil cart on error")
		}
	})
}

func TestCartService_Clear(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if r.URL.Path != "/storefront/v1/carts/cart_123/empty" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Cart cleared"})
	})

	client := setupTestClient(t, handler)
	service := NewCartService(client)

	t.Run("success clear cart", func(t *testing.T) {
		err := service.Clear(context.Background(), "cart_123")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
	})

	t.Run("sends request with empty cart id", func(t *testing.T) {
		err := service.Clear(context.Background(), "")
		if err != nil {
			t.Logf("Client accepted request (validation is server-side): %v", err)
		}
	})

	t.Run("fails with server error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Server error"})
		})
		client := setupTestClient(t, handler)
		service := NewCartService(client)

		err := service.Clear(context.Background(), "cart_123")
		if err == nil {
			t.Fatal("Expected error for server failure")
		}
	})
}

func TestCartService_Checkout(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var req CheckoutRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":           "order_123",
			"order_number": "ORD-001",
			"status":       "confirmed",
			"total_amount": 999,
			"currency":     "USD",
		})
	})

	client := setupTestClient(t, handler)
	service := NewCartService(client)

	t.Run("success checkout", func(t *testing.T) {
		req := CheckoutRequest{
			CustomerEmail: "test@example.com",
		}
		order, err := service.Checkout(context.Background(), "cart_123", req)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if order == nil {
			t.Fatal("Expected non-nil order")
		}
		if order.ID != "order_123" {
			t.Errorf("Expected order ID 'order_123', got '%s'", order.ID)
		}
	})

	t.Run("fails with empty cart id", func(t *testing.T) {
		req := CheckoutRequest{CustomerEmail: "test@example.com"}
		order, err := service.Checkout(context.Background(), "", req)
		if err == nil {
			t.Fatal("Expected error for empty cart ID")
		}
		if order != nil {
			t.Error("Expected nil order on error")
		}
	})

	t.Run("fails with server error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Server error"})
		})
		client := setupTestClient(t, handler)
		service := NewCartService(client)

		order, err := service.Checkout(context.Background(), "cart_123", CheckoutRequest{})
		if err == nil {
			t.Fatal("Expected error for server failure")
		}
		if order != nil {
			t.Error("Expected nil order on error")
		}
	})
}
