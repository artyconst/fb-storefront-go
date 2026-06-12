package location

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	sf "github.com/artyconst/fb-storefront-go"
)

func setupTestClient(t *testing.T, handler http.Handler) *sf.StorefrontClient {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	client, err := sf.NewStorefront("sk_test_key", sf.WithAPIHost(server.URL))
	if err != nil {
		t.Fatalf("Failed to create test client: %v", err)
	}
	return client
}

func TestService_List_Success(t *testing.T) {
	locations := []*Location{
		{
			ID:   "loc_001",
			Name: "Test Location",
			Type: "store",
			Place: &Place{
				ID:        "place_001",
				Name:      "Test Location",
				Address:   "123 Main St",
				Latitude:  40.7128,
				Longitude: -74.0060,
			},
			Hours: []StoreHour{
				{Day: "monday", Start: strPtr("09:00"), End: strPtr("17:00")},
			},
		},
		{
			ID:   "loc_002",
			Name: "Another Location",
			Type: "warehouse",
		},
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/storefront/v1/locations" && r.Method == http.MethodGet {
			if r.URL.Query().Get("store") != "store_123" {
				t.Errorf("Expected store query param 'store_123', got '%s'", r.URL.Query().Get("store"))
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(locations)
		} else {
			http.NotFound(w, r)
		}
	})

	client := setupTestClient(t, handler)
	service := NewService(client)

	ctx := context.Background()
	result, err := service.List(ctx, "store_123")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 locations, got %d", len(result))
	}

	if result[0].ID != "loc_001" {
		t.Errorf("Expected first location ID loc_001, got %s", result[0].ID)
	}

	if result[0].Name != "Test Location" {
		t.Errorf("Expected first location name 'Test Location', got %s", result[0].Name)
	}

	if result[0].Place == nil {
		t.Error("Expected first location Place to be non-nil")
	} else if result[0].Place.Address != "123 Main St" {
		t.Errorf("Expected Place address '123 Main St', got %s", result[0].Place.Address)
	}

	if result[1].ID != "loc_002" {
		t.Errorf("Expected second location ID loc_002, got %s", result[1].ID)
	}
}

func TestService_List_WithPagination(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/storefront/v1/locations" && r.Method == http.MethodGet {
			if r.URL.Query().Get("store") != "store_123" {
				t.Errorf("Expected store query param 'store_123', got '%s'", r.URL.Query().Get("store"))
			}
			if r.URL.Query().Get("limit") != "10" {
				t.Errorf("Expected limit query param '10', got '%s'", r.URL.Query().Get("limit"))
			}
			if r.URL.Query().Get("offset") != "20" {
				t.Errorf("Expected offset query param '20', got '%s'", r.URL.Query().Get("offset"))
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]*Location{})
		} else {
			http.NotFound(w, r)
		}
	})

	client := setupTestClient(t, handler)
	service := NewService(client)

	ctx := context.Background()
	_, err := service.List(ctx, "store_123", WithListLimit(10), WithListOffset(20))
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestService_List_EmptyStoreID(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	client := setupTestClient(t, handler)
	service := NewService(client)

	ctx := context.Background()
	result, err := service.List(ctx, "")
	if err == nil {
		t.Fatal("Expected error for empty store ID, got nil")
	}

	if !strings.Contains(err.Error(), "store ID is required") {
		t.Errorf("Expected error containing 'store ID is required', got: %v", err)
	}

	if result != nil {
		t.Errorf("Expected nil on error, got: %+v", result)
	}
}

func TestService_Get_Success(t *testing.T) {
	location := &Location{
		ID:   "loc_001",
		Name: "Test Location",
		Type: "store",
		Slug: "test-location",
		Place: &Place{
			ID:        "place_001",
			Name:      "Test Location",
			Address:   "123 Main St",
			Latitude:  40.7128,
			Longitude: -74.0060,
		},
		Hours: []StoreHour{
			{Day: "monday", Start: strPtr("09:00"), End: strPtr("17:00")},
		},
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/storefront/v1/locations/loc_001" && r.Method == http.MethodGet {
			if r.URL.Query().Get("store") != "store_123" {
				t.Errorf("Expected store query param 'store_123', got '%s'", r.URL.Query().Get("store"))
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(location)
		} else {
			http.NotFound(w, r)
		}
	})

	client := setupTestClient(t, handler)
	service := NewService(client)

	ctx := context.Background()
	result, err := service.Get(ctx, "store_123", "loc_001")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result.ID != "loc_001" {
		t.Errorf("Expected ID loc_001, got %s", result.ID)
	}

	if result.Name != "Test Location" {
		t.Errorf("Expected name Test Location, got %s", result.Name)
	}

	if result.Slug != "test-location" {
		t.Errorf("Expected slug test-location, got %s", result.Slug)
	}

	if result.Place == nil {
		t.Fatal("Expected Place to be non-nil")
	}

	if result.Place.Address != "123 Main St" {
		t.Errorf("Expected Place address '123 Main St', got %s", result.Place.Address)
	}

	if result.Place.Latitude != 40.7128 {
		t.Errorf("Expected Place latitude 40.7128, got %f", result.Place.Latitude)
	}

	if result.Place.Longitude != -74.0060 {
		t.Errorf("Expected Place longitude -74.0060, got %f", result.Place.Longitude)
	}

	if len(result.Hours) != 1 {
		t.Errorf("Expected 1 hour entry, got %d", len(result.Hours))
	}
}

func TestService_Get_EmptyStoreID(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	client := setupTestClient(t, handler)
	service := NewService(client)

	ctx := context.Background()
	result, err := service.Get(ctx, "", "loc_001")
	if err == nil {
		t.Fatal("Expected error for empty store ID, got nil")
	}

	if !strings.Contains(err.Error(), "store ID is required") {
		t.Errorf("Expected error containing 'store ID is required', got: %v", err)
	}

	if result != nil {
		t.Errorf("Expected nil on error, got: %+v", result)
	}
}

func TestService_Get_EmptyLocationID(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	client := setupTestClient(t, handler)
	service := NewService(client)

	ctx := context.Background()
	result, err := service.Get(ctx, "store_123", "")
	if err == nil {
		t.Fatal("Expected error for empty location ID, got nil")
	}

	if !strings.Contains(err.Error(), "location ID is required") {
		t.Errorf("Expected error containing 'location ID is required', got: %v", err)
	}

	if result != nil {
		t.Errorf("Expected nil on error, got: %+v", result)
	}
}

func TestService_Get_NotFound(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/storefront/v1/locations/loc_404" && r.Method == http.MethodGet {
			w.WriteHeader(http.StatusNotFound)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"error": "location not found"})
		} else {
			http.NotFound(w, r)
		}
	})

	client := setupTestClient(t, handler)
	service := NewService(client)

	ctx := context.Background()
	result, err := service.Get(ctx, "store_123", "loc_404")
	if err == nil {
		t.Fatal("Expected error for not found location, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil on error, got: %+v", result)
	}
}

func TestService_List_NetworkError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
	})

	client := setupTestClient(t, handler)
	service := NewService(client)

	ctx := context.Background()
	result, err := service.List(ctx, "store_123")
	if err == nil {
		t.Fatal("Expected error for server failure, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil on error, got: %+v", result)
	}
}

func TestService_ListLocations_Success(t *testing.T) {
	storeLocations := []*StoreLocation{
		{
			ID:   "sl_001",
			Name: "Main Store",
			Address: &Address{
				Line1:   "123 Main St",
				City:    "New York",
				State:   "NY",
				Zip:     "10001",
				Country: "US",
			},
			Hours: &Hours{
				Monday: &DayHours{Open: "09:00", Close: "17:00"},
				Friday: &DayHours{Open: "09:00", Close: "20:00"},
			},
		},
		{
			ID:   "sl_002",
			Name: "Downtown Branch",
			Address: &Address{
				Line1:   "456 Market St",
				City:    "San Francisco",
				State:   "CA",
				Zip:     "94102",
				Country: "US",
			},
		},
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/storefront/v1/store-locations" && r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(storeLocations)
		} else {
			http.NotFound(w, r)
		}
	})

	client := setupTestClient(t, handler)
	service := NewService(client)

	ctx := t.Context()
	result, err := service.ListLocations(ctx)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 store locations, got %d", len(result))
	}

	if result[0].ID != "sl_001" {
		t.Errorf("Expected first location ID sl_001, got %s", result[0].ID)
	}

	if result[0].Name != "Main Store" {
		t.Errorf("Expected first location name 'Main Store', got %s", result[0].Name)
	}

	if result[0].Address == nil {
		t.Error("Expected first location Address to be non-nil")
	} else if result[0].Address.City != "New York" {
		t.Errorf("Expected Address city 'New York', got %s", result[0].Address.City)
	}

	if result[0].Hours == nil {
		t.Error("Expected first location Hours to be non-nil")
	} else if result[0].Hours.Monday == nil {
		t.Error("Expected Monday hours to be non-nil")
	} else if result[0].Hours.Monday.Open != "09:00" {
		t.Errorf("Expected Monday open '09:00', got %s", result[0].Hours.Monday.Open)
	}

	if result[1].ID != "sl_002" {
		t.Errorf("Expected second location ID sl_002, got %s", result[1].ID)
	}

	if result[1].Name != "Downtown Branch" {
		t.Errorf("Expected second location name 'Downtown Branch', got %s", result[1].Name)
	}
}

func TestService_ListLocations_ServerError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/storefront/v1/store-locations" && r.Method == http.MethodGet {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		} else {
			http.NotFound(w, r)
		}
	})

	client := setupTestClient(t, handler)
	service := NewService(client)

	ctx := t.Context()
	result, err := service.ListLocations(ctx)
	if err == nil {
		t.Fatal("Expected error for server failure, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil on error, got: %+v", result)
	}

	if !strings.Contains(err.Error(), "failed to list store locations") {
		t.Errorf("Expected error containing 'failed to list store locations', got: %v", err)
	}
}

func strPtr(s string) *string {
	return &s
}
