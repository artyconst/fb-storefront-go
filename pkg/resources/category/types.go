package category

import "time"

// Category represents a product category.
type Category struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Description  *string                `json:"description,omitempty"`
	IconURL      string                 `json:"icon_url,omitempty"`
	Tags         []string               `json:"tags,omitempty"`
	Translations []interface{}          `json:"translations,omitempty"`
	Meta         map[string]interface{} `json:"meta,omitempty"`
	Order        *int64                 `json:"order,omitempty"`
	Slug         string                 `json:"slug,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    *time.Time             `json:"updated_at,omitempty"`
}

// ListOptions contains parameters for listing categories.
type ListOptions struct {
	Offset   int64   `json:"offset,omitempty"`
	Limit    int64   `json:"limit,omitempty"`
	Search   string  `json:"search,omitempty"`
	ParentID *string `json:"parent_id,omitempty"`
}
