package store

// GatewaysOption is a functional option for gateway listing operations
type GatewaysOption func(*ListGatewaysOpts)

// ListGatewaysOpts contains parameters for gateway listing
type ListGatewaysOpts struct {
	Limit  int64
	Offset int64
}

// WithGatewayLimit sets the maximum number of gateways to return
func WithGatewayLimit(limit int64) GatewaysOption {
	return func(o *ListGatewaysOpts) {
		o.Limit = limit
	}
}

// WithGatewayOffset sets the offset for gateway listing pagination
func WithGatewayOffset(offset int64) GatewaysOption {
	return func(o *ListGatewaysOpts) {
		o.Offset = offset
	}
}
