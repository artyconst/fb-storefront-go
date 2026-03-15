package config

import (
	"io"
	"strings"
	"time"

	sf "github.com/artyconst/fb-storefront-go/internal/types"
)

// Option is a function that modifies *ClientConfig.
type Option func(*sf.ClientConfig)

// WithAPIHost sets the server URL for API requests. The SDK automatically appends "/storefront/v1".
// Only the server URL should be provided - do not include the API path.
// Examples:
//   - "https://api.fleetbase.io" → https://api.fleetbase.io/storefront/v1
//   - "https://custom.api.com/" → https://custom.api.com/storefront/v1 (trailing slash stripped)
//   - "https://custom.api.com/with/path" → https://custom.api.com/with/path/storefront/v1 (path preserved)
func WithAPIHost(host string) Option {
	return func(c *sf.ClientConfig) {
		if host != "" {
			c.ServerURL = normalizeServerURL(host)
		}
	}
}

// normalizeServerURL strips all trailing slashes from the URL for consistent composition.
func normalizeServerURL(url string) string {
	url = strings.TrimSpace(url)
	for strings.HasSuffix(url, "/") {
		url = strings.TrimSuffix(url, "/")
	}
	return url
}

// WithAPIPath sets the API path that will be appended to the server URL.
// This allows customization of the endpoint path (e.g., "/storefront/v1").
// Example: sf.WithAPIPath("/storefront/v2") to use a different API version.
func WithAPIPath(path string) Option {
	return func(c *sf.ClientConfig) {
		if path != "" && !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
		c.APIPath = strings.TrimRight(path, "/")
	}
}

// WithTimeout sets the request timeout in seconds.
func WithTimeout(seconds time.Duration) Option {
	return func(c *sf.ClientConfig) {
		if seconds > 0 {
			c.Timeout = seconds
		}
	}
}

// WithLogLevel sets the minimum log level.
func WithLogLevel(level sf.LogLevel) Option {
	return func(c *sf.ClientConfig) {
		c.LogLevel = level
	}
}

// WithLoggerOutput sets the output writer for logs.
func WithLoggerOutput(output io.Writer) Option {
	return func(c *sf.ClientConfig) {
		if w, ok := output.(io.Writer); ok {
			c.LoggerOutput = w
		}
	}
}

// WithDebugMode enables debug logging with raw response annotation.
func WithDebugMode() Option {
	return func(c *sf.ClientConfig) {
		c.LogLevel = sf.LevelDebug
		if c.LogConfig == nil {
			c.LogConfig = sf.DefaultLoggingConfig()
		}
		c.LogConfig.AnnotateRawResponses = true
	}
}
