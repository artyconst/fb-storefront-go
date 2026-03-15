package product

// Product represents a product in the store.
type Product struct {
	ID             string     `json:"id"`
	Name           string     `json:"name"`
	Description    *string    `json:"description,omitempty"`
	Price          int64      `json:"price"`
	CompareAtPrice *int64     `json:"compare_at_price,omitempty"`
	SKU            string     `json:"sku"`
	StockQuantity  int        `json:"stock_quantity"`
	Images         []string   `json:"images"`
	Categories     []Category `json:"categories"`
	CreatedAt      string     `json:"created_at"`
	UpdatedAt      string     `json:"updated_at"`
}

// Image represents a product image.
type Image struct {
	ID       string `json:"id"`
	URL      string `json:"url"`
	AltText  string `json:"alt_text,omitempty"`
	Position int    `json:"position"`
}

// Category represents a category reference.
type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ListOptions contains parameters for listing products.
type ListOptions struct {
	Limit    *uint64 `json:"limit,omitempty"`
	Offset   int64   `json:"offset,omitempty"`
	Category string  `json:"category_id,omitempty"`
	SortBy   string  `json:"sort_by,omitempty"`
	Order    string  `json:"order,omitempty"`
}
