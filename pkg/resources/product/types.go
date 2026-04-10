package product

// Product represents a product in the store.
type Product struct {
	ID              string          `json:"id"`
	Name            string          `json:"name"`
	Description     *string         `json:"description,omitempty"`
	SKU             *string         `json:"sku,omitempty"`
	PrimaryImageURL *string         `json:"primary_image_url,omitempty"`
	Price           int64           `json:"price"`
	SalePrice       *int64          `json:"sale_price,omitempty"`
	Currency        string          `json:"currency"`
	IsOnSale        bool            `json:"is_on_sale"`
	IsRecommended   bool            `json:"is_recommended"`
	IsService       bool            `json:"is_service"`
	IsBookable      bool            `json:"is_bookable"`
	IsAvailable     bool            `json:"is_available"`
	Tags            []string        `json:"tags"`
	Status          string          `json:"status"`
	Meta            interface{}     `json:"meta"`
	Slug            string          `json:"slug"`
	Translations    []Translation   `json:"translations"`
	AddonCategories []AddonCategory `json:"addon_categories"`
	Variants        []Variant       `json:"variants"`
	Images          []string        `json:"images"`
	Videos          []string        `json:"videos"`
	Hours           []Hour          `json:"hours"`
	YouTubeURLs     []string        `json:"youtube_urls"`
	CreatedAt       string          `json:"created_at"`
	UpdatedAt       string          `json:"updated_at"`
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

// Translation represents a product translation for different languages.
type Translation struct {
	Locale      string         `json:"locale,omitempty"`
	Name        string         `json:"name,omitempty"`
	Description map[string]any `json:"description,omitempty"`
}

// AddonCategory represents a category of add-ons for a product.
type AddonCategory struct {
	ID           string  `json:"id,omitempty"`
	Name         string  `json:"name,omitempty"`
	Description  *string `json:"description,omitempty"`
	MinSelection int     `json:"min_selection,omitempty"`
	MaxSelection int     `json:"max_selection,omitempty"`
	Addons       []Addon `json:"addons,omitempty"`
}

// Addon represents an add-on option within an addon category.
type Addon struct {
	ID            string  `json:"id,omitempty"`
	Name          string  `json:"name,omitempty"`
	Description   *string `json:"description,omitempty"`
	Price         int64   `json:"price,omitempty"`
	PriceModifier *string `json:"price_modifier,omitempty"`
}

// Variant represents a product variant (e.g., size, color).
type Variant struct {
	ID            string   `json:"id,omitempty"`
	Name          string   `json:"name,omitempty"`
	SKU           *string  `json:"sku,omitempty"`
	Price         int64    `json:"price,omitempty"`
	SalePrice     *int64   `json:"sale_price,omitempty"`
	StockQuantity int      `json:"stock_quantity,omitempty"`
	ImageURL      *string  `json:"image_url,omitempty"`
	Options       []Option `json:"options,omitempty"`
}

// Option represents a variant option (e.g., "Size: Large").
type Option struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// Hour represents service hours for bookable/service products.
type Hour struct {
	Day       int    `json:"day,omitempty"`
	StartTime string `json:"start_time,omitempty"`
	EndTime   string `json:"end_time,omitempty"`
	IsOpen    bool   `json:"is_open,omitempty"`
}

// ListOptions contains parameters for listing products.
type ListOptions struct {
	Limit    *uint64 `json:"limit,omitempty"`
	Offset   int64   `json:"offset,omitempty"`
	Category string  `json:"category_id,omitempty"`
	SortBy   string  `json:"sort_by,omitempty"`
	Order    string  `json:"order,omitempty"`
}
