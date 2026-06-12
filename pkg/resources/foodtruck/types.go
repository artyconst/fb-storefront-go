package foodtruck

import "time"

// FoodTruck represents a food truck entity from the Storefront API.
type FoodTruck struct {
	ID          string       `json:"id"`
	Vehicle     *Vehicle     `json:"vehicle,omitempty"`
	ServiceArea *ServiceArea `json:"service_area,omitempty"`
	Zone        *Zone        `json:"zone,omitempty"`
	Catalogs    []*Catalog   `json:"catalogs,omitempty"`
	Location    *Location    `json:"location,omitempty"`
	Online      bool         `json:"online"`
	Status      string       `json:"status"`
	CreatedAt   time.Time    `json:"created_at,omitzero"`
	UpdatedAt   time.Time    `json:"updated_at,omitzero"`
}

// Vehicle represents the vehicle details of a food truck.
type Vehicle struct {
	ID    string  `json:"id"`
	Name  string  `json:"name,omitempty"`
	Type  string  `json:"type,omitempty"`
	Plate *string `json:"plate,omitempty"`
}

// ServiceArea represents the service area for a food truck.
type ServiceArea struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

// Zone represents the zone assignment for a food truck.
type Zone struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

// Catalog represents a catalog associated with a food truck.
type Catalog struct {
	ID        string    `json:"id"`
	Name      string    `json:"name,omitempty"`
	CreatedAt time.Time `json:"created_at,omitzero"`
}

// Location represents the current location of a food truck.
type Location struct {
	ID        string  `json:"id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// ListOptions contains parameters for listing food trucks.
type ListOptions struct {
	Limit  int64  `json:"limit,omitempty"`
	Offset int64  `json:"offset,omitempty"`
	Sort   string `json:"sort,omitempty"`
}
