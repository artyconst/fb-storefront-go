package location

import (
	"time"
)

// StoreLocation represents a store location returned by the dedicated /store-locations endpoint.
type StoreLocation struct {
	ID      string  `json:"id"`
	Name    string  `json:"name,omitempty"`
	Address *Address `json:"address,omitempty"`
	Hours   *Hours  `json:"hours,omitempty"`
}

// Address represents a physical address for a store location.
type Address struct {
	Line1   string `json:"line_1,omitempty"`
	Line2   string `json:"line_2,omitempty"`
	City    string `json:"city,omitempty"`
	State   string `json:"state,omitempty"`
	Zip     string `json:"zip,omitempty"`
	Country string `json:"country,omitempty"`
}

// Hours represents the weekly working hours for a store location.
type Hours struct {
	Monday    *DayHours `json:"monday,omitempty"`
	Tuesday   *DayHours `json:"tuesday,omitempty"`
	Wednesday *DayHours `json:"wednesday,omitempty"`
	Thursday  *DayHours `json:"thursday,omitempty"`
	Friday    *DayHours `json:"friday,omitempty"`
	Saturday  *DayHours `json:"saturday,omitempty"`
	Sunday    *DayHours `json:"sunday,omitempty"`
}

// DayHours represents the open and close times for a single day.
type DayHours struct {
	Open  string `json:"open,omitempty"`
	Close string `json:"close,omitempty"`
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
	Metadata  map[string]any `json:"metadata,omitempty"`
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
