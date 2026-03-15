package order

// ListOption is a functional option for List operations
type ListOption func(*ListOptions)

// WithPage sets the page number for pagination (page-based, not offset-based)
func WithPage(page int) ListOption {
	return func(o *ListOptions) {
		o.Page = page
	}
}

// WithLimit sets the maximum number of orders per page
func WithLimit(limit int) ListOption {
	return func(o *ListOptions) {
		o.Limit = limit
	}
}

// WithStatus filters orders by status
func WithStatus(status OrderStatus) ListOption {
	return func(o *ListOptions) {
		o.Status = status
	}
}
