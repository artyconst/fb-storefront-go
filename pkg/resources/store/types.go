package store

import (
	"time"
)

// StoreOptions represents store configuration options
type StoreOptions struct {
	PickupEnabled          bool  `json:"pickup_enabled,omitempty"`
	RequiredCheckoutMin    bool  `json:"required_checkout_min,omitempty"`
	RequiredCheckoutMinAmt int64 `json:"required_checkout_min_amount,omitempty"`
}

// AboutStoreResponse from GET /about
type AboutStoreResponse struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Description  string        `json:"description,omitempty"`
	Translations []interface{} `json:"translations,omitempty"`
	Website      *string       `json:"website,omitempty"`
	Facebook     *string       `json:"facebook,omitempty"`
	Instagram    *string       `json:"instagram,omitempty"`
	Twitter      *string       `json:"twitter,omitempty"`
	Email        *string       `json:"email,omitempty"`
	Phone        *string       `json:"phone,omitempty"`
	Tags         []string      `json:"tags,omitempty"`
	Currency     string        `json:"currency,omitempty"`
	Country      string        `json:"country,omitempty"`
	Options      StoreOptions  `json:"options,omitempty"`
	LogoURL      string        `json:"logo_url,omitempty"`
	BackdropURL  string        `json:"backdrop_url,omitempty"`
	Rating       float64       `json:"rating,omitempty"`
	Online       bool          `json:"online,omitempty"`
	Alertable    *bool         `json:"alertable,omitempty"`
	IsNetwork    bool          `json:"is_network,omitempty"`
	IsStore      bool          `json:"is_store,omitempty"`
	Slug         string        `json:"slug,omitempty"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

// PaymentGateway from API response
type PaymentGateway struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Type          string                 `json:"type"`
	IsActive      bool                   `json:"is_active"`
	Configuration map[string]interface{} `json:"configuration,omitempty"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

// PaymentGatewaysResponse list response
type PaymentGatewaysResponse struct {
	Data   []PaymentGateway `json:"data"`
	Total  int64            `json:"total"`
	Limit  int64            `json:"limit"`
	Offset int64            `json:"offset"`
}
