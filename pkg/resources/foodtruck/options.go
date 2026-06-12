package foodtruck

// ListOption is a functional option for List operations.
type ListOption func(*ListOptions)

// WithLimit sets the maximum number of food trucks to return.
func WithLimit(limit int64) ListOption {
	return func(o *ListOptions) {
		o.Limit = limit
	}
}

// WithOffset sets the offset for pagination.
func WithOffset(offset int64) ListOption {
	return func(o *ListOptions) {
		o.Offset = offset
	}
}

// WithSort sets the sort field for results.
func WithSort(sort string) ListOption {
	return func(o *ListOptions) {
		o.Sort = sort
	}
}
