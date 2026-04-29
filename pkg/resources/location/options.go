package location

// ListOptions holds query parameters for location listing.
type ListOptions struct {
	Limit  int64
	Offset int64
}

// ListOption is a functional option for ListOptions.
type ListOption func(*ListOptions)

// WithListLimit sets the limit for list results.
func WithListLimit(limit int64) ListOption {
	return func(o *ListOptions) {
		o.Limit = limit
	}
}

// WithListOffset sets the offset for list results.
func WithListOffset(offset int64) ListOption {
	return func(o *ListOptions) {
		o.Offset = offset
	}
}
