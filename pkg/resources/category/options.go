package category

// ListOption is a functional option for List operations
type ListOption func(*ListOptions)

// WithOffset sets the offset for pagination
func WithOffset(offset int64) ListOption {
	return func(o *ListOptions) {
		o.Offset = offset
	}
}

// WithLimit sets the maximum number of categories to return
func WithLimit(limit int64) ListOption {
	return func(o *ListOptions) {
		o.Limit = limit
	}
}

// WithSearch searches categories by name or description
func WithSearch(search string) ListOption {
	return func(o *ListOptions) {
		o.Search = search
	}
}

// WithParentID filters categories by parent ID (for nested categories)
func WithParentID(parentID string) ListOption {
	return func(o *ListOptions) {
		o.ParentID = &parentID
	}
}
