package location

import (
	"time"
)

// StoreLocation represents a store location with its place details and working hours.
type StoreLocation struct {
	ID        string      `json:"id"`
	Place     *Place      `json:"place"`
	Hours     []StoreHour `json:"hours"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// Place represents a location with address and coordinates.
type Place struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// StoreHour represents a store's working hours for a specific day.
type StoreHour struct {
	Day   string  `json:"day"`
	Start *string `json:"start"`
	End   *string `json:"end"`
}

// Location represents a standalone location resource.
type Location struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Type      string                 `json:"type"`
	Slug      string                 `json:"slug,omitempty"`
	Place     *Place                 `json:"place,omitempty"`
	Hours     []StoreHour            `json:"hours,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// ListLocationsResponse is the paginated list response for locations.
type ListLocationsResponse struct {
	Data   []Location `json:"data"`
	Total  int64      `json:"total"`
	Limit  int64      `json:"limit"`
	Offset int64      `json:"offset"`
}
