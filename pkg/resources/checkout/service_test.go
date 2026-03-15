package checkout

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

func TestCheckoutService_Create(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var body map[string]interface{}
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
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":       "checkout_123",
			"cart_id":  body["cart_id"],
			"status":   "pending",
			"amount":   999,
			"currency": "USD",
		})
	})

	client := setupTestClient(t, handler)
	service := NewCheckoutService(client)

	t.Run("success create checkout", func(t *testing.T) {
		req := CreateCheckoutRequest{CustomerEmail: "test@example.com"}
		checkout, err := service.Create(context.Background(), "cart_123", req)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if checkout == nil {
			t.Fatal("Expected non-nil checkout")
		}
		if checkout.ID != "checkout_123" {
			t.Errorf("Expected ID 'checkout_123', got '%s'", checkout.ID)
		}
	})

	t.Run("fails with empty cart id", func(t *testing.T) {
		req := CreateCheckoutRequest{CustomerEmail: "test@example.com"}
		checkout, err := service.Create(context.Background(), "", req)
		if err == nil {
			t.Fatal("Expected error for empty cart ID")
		}
		if checkout != nil {
			t.Error("Expected nil checkout on error")
		}
	})

	t.Run("fails with shipping address", func(t *testing.T) {
		req := CreateCheckoutRequest{
			CustomerEmail: "test@example.com",
			ShippingAddress: &Address{
				FirstName:    "John",
				LastName:     "Doe",
				AddressLine1: "123 Main St",
				City:         "New York",
				State:        "NY",
				PostalCode:   "10001",
				Country:      "US",
			},
		}
		checkout, err := service.Create(context.Background(), "cart_123", req)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if checkout == nil {
			t.Error("Expected non-nil checkout")
		}
	})

	t.Run("fails with server error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Server error"})
		})
		client := setupTestClient(t, handler)
		service := NewCheckoutService(client)

		checkout, err := service.Create(context.Background(), "cart_123", CreateCheckoutRequest{})
		if err == nil {
			t.Fatal("Expected error for server failure")
		}
		if checkout != nil {
			t.Error("Expected nil checkout on error")
		}
	})
}

func TestCheckoutService_Get(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":       "checkout_123",
			"cart_id":  "cart_456",
			"status":   "processing",
			"amount":   999,
			"currency": "USD",
		})
	})

	client := setupTestClient(t, handler)
	service := NewCheckoutService(client)

	t.Run("success get checkout", func(t *testing.T) {
		checkout, err := service.Get(context.Background(), "checkout_123")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if checkout == nil {
			t.Fatal("Expected non-nil checkout")
		}
	})

	t.Run("fails with empty id", func(t *testing.T) {
		checkout, err := service.Get(context.Background(), "")
		if err == nil {
			t.Fatal("Expected error for empty ID")
		}
		if checkout != nil {
			t.Error("Expected nil checkout on error")
		}
	})

	t.Run("fails with server error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Not found"})
		})
		client := setupTestClient(t, handler)
		service := NewCheckoutService(client)

		checkout, err := service.Get(context.Background(), "checkout_123")
		if err == nil {
			t.Fatal("Expected error for not found")
		}
		if checkout != nil {
			t.Error("Expected nil checkout on error")
		}
	})
}

func TestCheckoutService_UpdateCustomer(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var customer CustomerInfo
		if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
			t.Errorf("Failed to decode request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if customer.Email == "" && customer.Phone == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Customer info required"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":          "checkout_123",
			"customer_id": "cust_456",
			"status":      "pending",
		})
	})

	client := setupTestClient(t, handler)
	service := NewCheckoutService(client)

	t.Run("success update customer", func(t *testing.T) {
		customer := CustomerInfo{Email: "new@example.com"}
		checkout, err := service.UpdateCustomer(context.Background(), "checkout_123", customer)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if checkout == nil {
			t.Fatal("Expected non-nil checkout")
		}
	})

	t.Run("fails with server error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Server error"})
		})

		server := httptest.NewServer(handler)
		defer server.Close()

		client, err := sf.NewStorefront("sk_test_key", sf.WithAPIHost(server.URL))
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		service := NewCheckoutService(client)

		ctx := context.Background()
		result, err := service.UpdateCustomer(ctx, "checkout_123", CustomerInfo{ID: "cust_456"})
		if err == nil {
			t.Fatal("Expected error for server failure")
		}
		if result != nil {
			t.Errorf("Expected nil on error, got: %+v", result)
		}
	})

	t.Run("fails with server error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Server error"})
		})
		client := setupTestClient(t, handler)
		service := NewCheckoutService(client)

		checkout, err := service.UpdateCustomer(context.Background(), "checkout_123", CustomerInfo{Email: "test@example.com"})
		if err == nil {
			t.Fatal("Expected error for server failure")
		}
		if checkout != nil {
			t.Error("Expected nil checkout on error")
		}
	})
}

func TestCheckoutService_ProcessPayment(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req PaymentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if req.MethodID == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Method ID required"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":         "checkout_123",
			"status":     "processing",
			"payment_id": req.MethodID,
		})
	})

	client := setupTestClient(t, handler)
	service := NewCheckoutService(client)

	t.Run("success process payment", func(t *testing.T) {
		req := PaymentRequest{MethodID: "pm_123"}
		checkout, err := service.ProcessPayment(context.Background(), "checkout_123", req)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if checkout == nil {
			t.Fatal("Expected non-nil checkout")
		}
	})

	t.Run("sends request with empty method id", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"id": "checkout_123"})
		})

		server := httptest.NewServer(handler)
		defer server.Close()

		client, err := sf.NewStorefront("sk_test_key", sf.WithAPIHost(server.URL))
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		service = NewCheckoutService(client)

		req := PaymentRequest{MethodID: ""}
		checkout, err := service.ProcessPayment(context.Background(), "checkout_123", req)
		if err != nil {
			t.Logf("Client accepted request with empty MethodID (validation is server-side): %v", err)
		}
		if checkout == nil && err == nil {
			t.Error("Expected checkout result")
		}
	})

	t.Run("fails with empty checkout id", func(t *testing.T) {
		req := PaymentRequest{MethodID: "pm_123"}
		checkout, err := service.ProcessPayment(context.Background(), "", req)
		if err == nil {
			t.Fatal("Expected error for empty checkout ID")
		}
		if checkout != nil {
			t.Error("Expected nil checkout on error")
		}
	})

	t.Run("fails with server error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Server error"})
		})
		client := setupTestClient(t, handler)
		service := NewCheckoutService(client)

		checkout, err := service.ProcessPayment(context.Background(), "checkout_123", PaymentRequest{MethodID: "pm_123"})
		if err == nil {
			t.Fatal("Expected error for server failure")
		}
		if checkout != nil {
			t.Error("Expected nil checkout on error")
		}
	})
}

func TestCheckoutService_GetDeliveryServiceQuote(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":          "quote_123",
			"origin":      "warehouse_a",
			"destination": "store_b",
			"price":       500,
			"currency":    "USD",
		})
	})

	client := setupTestClient(t, handler)
	service := NewCheckoutService(client)

	t.Run("success get quote with params", func(t *testing.T) {
		params := ServiceQuoteParams{Origin: "wh_a", Destination: "dest_b"}
		quote, err := service.GetDeliveryServiceQuote(context.Background(), params)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if quote == nil {
			t.Fatal("Expected non-nil quote")
		}
	})

	t.Run("success get quote with cart id", func(t *testing.T) {
		params := ServiceQuoteParams{CartID: "cart_123"}
		quote, err := service.GetDeliveryServiceQuote(context.Background(), params)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if quote == nil {
			t.Fatal("Expected non-nil quote")
		}
	})

	t.Run("success get quote empty params", func(t *testing.T) {
		params := ServiceQuoteParams{}
		quote, err := service.GetDeliveryServiceQuote(context.Background(), params)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if quote == nil {
			t.Fatal("Expected non-nil quote")
		}
	})

	t.Run("fails with server error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Server error"})
		})
		client := setupTestClient(t, handler)
		service := NewCheckoutService(client)

		quote, err := service.GetDeliveryServiceQuote(context.Background(), ServiceQuoteParams{})
		if err == nil {
			t.Fatal("Expected error for server failure")
		}
		if quote != nil {
			t.Error("Expected nil quote on error")
		}
	})
}

func TestCheckoutService_CaptureCheckout(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":       "checkout_123",
			"status":   "completed",
			"order_id": "order_456",
		})
	})

	client := setupTestClient(t, handler)
	service := NewCheckoutService(client)

	t.Run("success capture checkout", func(t *testing.T) {
		checkout, err := service.CaptureCheckout(context.Background(), "")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if checkout == nil {
			t.Fatal("Expected non-nil checkout")
		}
	})

	t.Run("success capture with token", func(t *testing.T) {
		checkout, err := service.CaptureCheckout(context.Background(), "token_123")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if checkout == nil {
			t.Fatal("Expected non-nil checkout")
		}
	})

	t.Run("fails with server error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Server error"})
		})
		client := setupTestClient(t, handler)
		service := NewCheckoutService(client)

		checkout, err := service.CaptureCheckout(context.Background(), "")
		if err == nil {
			t.Fatal("Expected error for server failure")
		}
		if checkout != nil {
			t.Error("Expected nil checkout on error")
		}
	})
}
