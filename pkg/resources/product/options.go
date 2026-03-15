package product

// ListOption is a functional option for List operations
type ListOption func(*ListOptions)

// WithOffset sets the offset for pagination
func WithOffset(offset int64) ListOption {
	return func(o *ListOptions) {
		o.Offset = offset
	}
}

// WithCategory filters products by category ID
func WithCategory(categoryID string) ListOption {
	return func(o *ListOptions) {
		o.Category = categoryID
	}
}

// WithSortBy sets the field to sort results by
func WithSortBy(sortBy string) ListOption {
	return func(o *ListOptions) {
		o.SortBy = sortBy
	}
}

// WithOrder sets the sort order (asc/desc)
func WithOrder(order string) ListOption {
	return func(o *ListOptions) {
		o.Order = order
	}
}

// WithLimit sets the maximum number of products to return
func WithLimit(limit uint64) ListOption {
	if limit == 0 {
		limit = 1 // enforce minimum value
	}
	return func(o *ListOptions) {
		o.Limit = &limit
	}
}
