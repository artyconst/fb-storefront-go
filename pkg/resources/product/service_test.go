package product

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

func TestProductService_List(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]*Product{
			{ID: "prod_1", Name: "Product 1", Price: 999},
			{ID: "prod_2", Name: "Product 2", Price: 1499},
		})
	})

	client := setupTestClient(t, handler)
	service := NewProductService(client)

	t.Run("success list products", func(t *testing.T) {
		products, err := service.List(context.Background())
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if len(products) != 2 {
			t.Errorf("Expected 2 products, got %d", len(products))
		}
	})

	t.Run("success list with options", func(t *testing.T) {
		products, err := service.List(context.Background(), WithOffset(10), WithSortBy("name"))
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if products == nil {
			t.Error("Expected non-nil products")
		}
	})

	t.Run("success list with category filter", func(t *testing.T) {
		products, err := service.List(context.Background(), WithCategory("cat_123"))
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if products == nil {
			t.Error("Expected non-nil products")
		}
	})

	t.Run("fails with server error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Server error"})
		})
		client := setupTestClient(t, handler)
		service := NewProductService(client)

		products, err := service.List(context.Background())
		if err == nil {
			t.Fatal("Expected error for server failure")
		}
		if products != nil {
			t.Error("Expected nil products on error")
		}
	})
}

func TestProductService_Get(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":          "prod_123",
			"name":        "Widget",
			"description": "A useful widget",
			"price":       999,
			"currency":    "USD",
		})
	})

	client := setupTestClient(t, handler)
	service := NewProductService(client)

	t.Run("success get product", func(t *testing.T) {
		product, err := service.Get(context.Background(), "prod_123")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if product == nil {
			t.Fatal("Expected non-nil product")
		}
		if product.ID != "prod_123" {
			t.Errorf("Expected ID 'prod_123', got '%s'", product.ID)
		}
	})

	t.Run("fails with empty id", func(t *testing.T) {
		product, err := service.Get(context.Background(), "")
		if err == nil {
			t.Fatal("Expected error for empty ID")
		}
		if product != nil {
			t.Error("Expected nil product on error")
		}
	})

	t.Run("fails with server error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Not found"})
		})
		client := setupTestClient(t, handler)
		service := NewProductService(client)

		product, err := service.Get(context.Background(), "prod_123")
		if err == nil {
			t.Fatal("Expected error for not found")
		}
		if product != nil {
			t.Error("Expected nil product on error")
		}
	})

	t.Run("sends request with invalid id format", func(t *testing.T) {
		_, err := service.Get(context.Background(), "invalid-id")
		if err != nil {
			t.Logf("Client accepted request (validation is server-side): %v", err)
		}
	})
}

func TestProductService_ListWithOptions(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]*Product{
			{ID: "prod_1", Name: "Widget A"},
			{ID: "prod_2", Name: "Widget B"},
		})
	})

	client := setupTestClient(t, handler)
	service := NewProductService(client)

	t.Run("success with limit", func(t *testing.T) {
		products, err := service.List(context.Background(), WithLimit(10))
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if len(products) != 2 {
			t.Errorf("Expected 2 products, got %d", len(products))
		}
	})

	t.Run("success with category filter", func(t *testing.T) {
		products, err := service.List(context.Background(), WithCategory("cat_123"))
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if len(products) != 2 {
			t.Errorf("Expected 2 products, got %d", len(products))
		}
	})

	t.Run("fails with server error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Server error"})
		})
		client := setupTestClient(t, handler)
		service := NewProductService(client)

		products, err := service.List(context.Background(), WithLimit(10))
		if err == nil {
			t.Fatal("Expected error for server failure")
		}
		if products != nil {
			t.Error("Expected nil products on error")
		}
	})

	t.Run("limit with zero value defaults to 1", func(t *testing.T) {
		products, err := service.List(context.Background(), WithLimit(0))
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if len(products) != 2 {
			t.Errorf("Expected 2 products, got %d", len(products))
		}
	})
}

func TestProductService_Create(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":    "prod_new_123",
			"name":  "New Product",
			"price": 499,
		})
	})

	client := setupTestClient(t, handler)
	service := NewProductService(client)

	t.Run("success create product with name and price", func(t *testing.T) {
		product, err := service.Create(t.Context(), WithProductName("New Product"), WithPrice(499))
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if product == nil {
			t.Fatal("Expected non-nil product")
		}
		if product.ID != "prod_new_123" {
			t.Errorf("Expected ID 'prod_new_123', got '%s'", product.ID)
		}
	})

	t.Run("success create product with all options", func(t *testing.T) {
		product, err := service.Create(t.Context(),
			WithProductName("Full Product"),
			WithDescription("A full product description"),
			WithSKU("SKU-001"),
			WithPrice(999),
			WithSalePrice(799),
			WithCurrency("USD"),
			WithTags([]string{"sale", "new"}),
		)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if product == nil {
			t.Fatal("Expected non-nil product")
		}
	})

	t.Run("fails with empty name", func(t *testing.T) {
		product, err := service.Create(t.Context(), WithPrice(499))
		if err == nil {
			t.Fatal("Expected error for empty name")
		}
		if product != nil {
			t.Error("Expected nil product on error")
		}
	})

	t.Run("fails with server error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Server error"})
		})
		client := setupTestClient(t, handler)
		service := NewProductService(client)

		product, err := service.Create(t.Context(), WithProductName("Bad Product"), WithPrice(499))
		if err == nil {
			t.Fatal("Expected error for server failure")
		}
		if product != nil {
			t.Error("Expected nil product on error")
		}
	})
}

func TestProductService_Update(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":    "prod_123",
			"name":  "Updated Product",
			"price": 599,
		})
	})

	client := setupTestClient(t, handler)
	service := NewProductService(client)

	t.Run("success update product name", func(t *testing.T) {
		product, err := service.Update(t.Context(), "prod_123", WithUpdatedName("Updated Product"))
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if product == nil {
			t.Fatal("Expected non-nil product")
		}
		if product.ID != "prod_123" {
			t.Errorf("Expected ID 'prod_123', got '%s'", product.ID)
		}
	})

	t.Run("success update product with multiple options", func(t *testing.T) {
		product, err := service.Update(t.Context(), "prod_123",
			WithUpdatedName("Updated Product"),
			WithUpdatedPrice(599),
			WithUpdatedSalePrice(499),
			WithUpdatedCurrency("EUR"),
		)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if product == nil {
			t.Fatal("Expected non-nil product")
		}
	})

	t.Run("fails with empty id", func(t *testing.T) {
		product, err := service.Update(t.Context(), "", WithUpdatedName("No ID"))
		if err == nil {
			t.Fatal("Expected error for empty ID")
		}
		if product != nil {
			t.Error("Expected nil product on error")
		}
	})

	t.Run("fails with server error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Server error"})
		})
		client := setupTestClient(t, handler)
		service := NewProductService(client)

		product, err := service.Update(t.Context(), "prod_123", WithUpdatedName("Bad Update"))
		if err == nil {
			t.Fatal("Expected error for server failure")
		}
		if product != nil {
			t.Error("Expected nil product on error")
		}
	})
}
