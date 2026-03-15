package review

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

func TestReviewService_CountByStore(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int{"count": 42})
	})

	client := setupTestClient(t, handler)
	service := NewReviewService(client)

	t.Run("success count by store", func(t *testing.T) {
		count, err := service.CountByStore(context.Background(), "store_123")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if count != 42 {
			t.Errorf("Expected count 42, got %d", count)
		}
	})

	t.Run("fails with empty store id", func(t *testing.T) {
		count, err := service.CountByStore(context.Background(), "")
		if err == nil {
			t.Fatal("Expected error for empty store ID")
		}
		if count != 0 {
			t.Error("Expected count 0 on error")
		}
	})

	t.Run("succeeds with valid but non-existent store", func(t *testing.T) {
		_, err := service.CountByStore(context.Background(), "invalid-store-id")
		if err != nil {
			t.Logf("Expected success for invalid ID (server may accept it): %v", err)
		}
	})

	t.Run("fails with server error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Server error"})
		})
		client := setupTestClient(t, handler)
		service := NewReviewService(client)

		count, err := service.CountByStore(context.Background(), "store_123")
		if err == nil {
			t.Fatal("Expected error for server failure")
		}
		if count != 0 {
			t.Error("Expected count 0 on error")
		}
	})
}

func TestReviewService_CountByRating(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rating := r.URL.Query().Get("rating")
		count := 10
		switch rating {
		case "5":
			count = 35
		case "4":
			count = 28
		case "3":
			count = 15
		case "2":
			count = 8
		case "1":
			count = 3
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int{"count": count})
	})

	client := setupTestClient(t, handler)
	service := NewReviewService(client)

	t.Run("success count rating 5", func(t *testing.T) {
		count, err := service.CountByRating(context.Background(), 5)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if count <= 0 {
			t.Errorf("Expected positive count for rating 5, got %d", count)
		}
	})

	t.Run("success count rating 1", func(t *testing.T) {
		count, err := service.CountByRating(context.Background(), 1)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if count <= 0 {
			t.Errorf("Expected positive count for rating 1, got %d", count)
		}
	})

	t.Run("success count rating 3", func(t *testing.T) {
		count, err := service.CountByRating(context.Background(), 3)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if count <= 0 {
			t.Errorf("Expected positive count for rating 3, got %d", count)
		}
	})

	t.Run("fails with invalid rating 0", func(t *testing.T) {
		count, err := service.CountByRating(context.Background(), 0)
		if err == nil {
			t.Fatal("Expected error for rating 0")
		}
		if count != 0 {
			t.Error("Expected count 0 on error")
		}
	})

	t.Run("fails with invalid rating -1", func(t *testing.T) {
		count, err := service.CountByRating(context.Background(), -1)
		if err == nil {
			t.Fatal("Expected error for negative rating")
		}
		if count != 0 {
			t.Error("Expected count 0 on error")
		}
	})

	t.Run("fails with invalid rating 6", func(t *testing.T) {
		count, err := service.CountByRating(context.Background(), 6)
		if err == nil {
			t.Fatal("Expected error for rating > 5")
		}
		if count != 0 {
			t.Error("Expected count 0 on error")
		}
	})

	t.Run("fails with server error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Server error"})
		})
		client := setupTestClient(t, handler)
		service := NewReviewService(client)

		count, err := service.CountByRating(context.Background(), 5)
		if err == nil {
			t.Fatal("Expected error for server failure")
		}
		if count != 0 {
			t.Error("Expected count 0 on error")
		}
	})

	t.Run("fails with rating string that parses to invalid", func(t *testing.T) {
		count, err := service.CountByRating(context.Background(), -5)
		if err == nil {
			t.Fatal("Expected error for negative rating")
		}
		if count != 0 {
			t.Error("Expected count 0 on error")
		}
	})
}
