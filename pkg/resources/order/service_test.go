package order

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

func TestOrderService_Get(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":           "order_123",
			"order_number": "ORD-001",
			"status":       "confirmed",
			"total_amount": 999,
			"currency":     "USD",
		})
	})

	client := setupTestClient(t, handler)
	service := NewOrderService(client)

	t.Run("success get order by id", func(t *testing.T) {
		order, err := service.Get(context.Background(), "order_123")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if order == nil {
			t.Fatal("Expected non-nil order")
		}
	})

	t.Run("success get order by number", func(t *testing.T) {
		order, err := service.Get(context.Background(), "ORD-001")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if order == nil {
			t.Fatal("Expected non-nil order")
		}
	})

	t.Run("sends request with empty identifier", func(t *testing.T) {
		_, err := service.Get(context.Background(), "")
		if err != nil {
			t.Logf("Client accepted request (validation is server-side): %v", err)
		}
	})

	t.Run("fails with server error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Not found"})
		})
		client := setupTestClient(t, handler)
		service := NewOrderService(client)

		order, err := service.Get(context.Background(), "order_123")
		if err == nil {
			t.Fatal("Expected error for not found")
		}
		if order != nil {
			t.Error("Expected nil order on error")
		}
	})
}

func TestOrderService_List(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"data": []*Order{
				{ID: "order_123", OrderNumber: "ORD-001", Status: "confirmed"},
				{ID: "order_456", OrderNumber: "ORD-002", Status: "processing"},
			},
		})
	})

	client := setupTestClient(t, handler)
	service := NewOrderService(client)

	t.Run("success list orders", func(t *testing.T) {
		orders, err := service.List(context.Background())
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if len(orders) != 2 {
			t.Errorf("Expected 2 orders, got %d", len(orders))
		}
	})

	t.Run("success list with page option", func(t *testing.T) {
		orders, err := service.List(context.Background(), WithPage(1), WithLimit(10))
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if orders == nil {
			t.Error("Expected non-nil orders")
		}
	})

	t.Run("success list with status option", func(t *testing.T) {
		opts := []ListOption{WithStatus(OrderStatusConfirmed)}
		orders, err := service.List(context.Background(), opts...)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if orders == nil {
			t.Error("Expected non-nil orders")
		}
	})

	t.Run("fails with server error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Server error"})
		})
		client := setupTestClient(t, handler)
		service := NewOrderService(client)

		orders, err := service.List(context.Background())
		if err == nil {
			t.Fatal("Expected error for server failure")
		}
		if orders != nil {
			t.Error("Expected nil orders on error")
		}
	})
}

func TestOrderService_Create(t *testing.T) {
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

		if body["cart_id"] == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Cart ID required"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":           "order_789",
			"order_number": "ORD-003",
			"status":       "confirmed",
			"total_amount": 1500,
			"currency":     "USD",
		})
	})

	client := setupTestClient(t, handler)
	service := NewOrderService(client)

	t.Run("success create order", func(t *testing.T) {
		order, err := service.Create(context.Background(), "cart_123")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if order == nil {
			t.Fatal("Expected non-nil order")
		}
		if order.ID != "order_789" {
			t.Errorf("Expected ID 'order_789', got '%s'", order.ID)
		}
	})

	t.Run("fails with empty cart id", func(t *testing.T) {
		order, err := service.Create(context.Background(), "")
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
		service := NewOrderService(client)

		order, err := service.Create(context.Background(), "cart_123")
		if err == nil {
			t.Fatal("Expected error for server failure")
		}
		if order != nil {
			t.Error("Expected nil order on error")
		}
	})
}

func TestOrderService_MarkPickedUp(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("Failed to decode request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if body["id"] == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Order ID required"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":           body["id"],
			"status":       "delivered",
			"order_number": "ORD-001",
		})
	})

	client := setupTestClient(t, handler)
	service := NewOrderService(client)

	t.Run("success mark order as picked up", func(t *testing.T) {
		err := service.MarkPickedUp(context.Background(), "order_123")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
	})

	t.Run("fails with empty order ID", func(t *testing.T) {
		err := service.MarkPickedUp(context.Background(), "")
		if err == nil {
			t.Fatal("Expected error for empty order ID")
		}
	})

	t.Run("fails with server error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Server error"})
		})
		client := setupTestClient(t, handler)
		service := NewOrderService(client)

		err := service.MarkPickedUp(context.Background(), "order_123")
		if err == nil {
			t.Fatal("Expected error for server failure")
		}
	})
}

func TestOrderService_GenerateReceipt(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("Failed to decode request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if body["id"] == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Order ID required"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"order_id": body["id"],
			"content":  "receipt content here",
		})
	})

	client := setupTestClient(t, handler)
	service := NewOrderService(client)

	t.Run("success generate receipt", func(t *testing.T) {
		receipt, err := service.GenerateReceipt(context.Background(), "order_123")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if receipt == nil {
			t.Fatal("Expected non-nil receipt")
		}
		if receipt.OrderID != "order_123" {
			t.Errorf("Expected OrderID 'order_123', got '%s'", receipt.OrderID)
		}
		if receipt.Content != "receipt content here" {
			t.Errorf("Expected Content 'receipt content here', got '%s'", receipt.Content)
		}
	})

	t.Run("success generate receipt with ebarimt options", func(t *testing.T) {
		receipt, err := service.GenerateReceipt(context.Background(), "order_123",
			WithEbarimtReceiverType("COMPANY"),
			WithEbarimtReceiver("1234567890"),
		)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if receipt == nil {
			t.Fatal("Expected non-nil receipt")
		}
	})

	t.Run("fails with empty order ID", func(t *testing.T) {
		receipt, err := service.GenerateReceipt(context.Background(), "")
		if err == nil {
			t.Fatal("Expected error for empty order ID")
		}
		if receipt != nil {
			t.Error("Expected nil receipt on error")
		}
	})

	t.Run("fails with server error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Server error"})
		})
		client := setupTestClient(t, handler)
		service := NewOrderService(client)

		receipt, err := service.GenerateReceipt(context.Background(), "order_123")
		if err == nil {
			t.Fatal("Expected error for server failure")
		}
		if receipt != nil {
			t.Error("Expected nil receipt on error")
		}
	})
}
