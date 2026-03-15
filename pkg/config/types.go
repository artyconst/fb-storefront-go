package config

import sf "github.com/artyconst/fb-storefront-go/internal/types"

// QueryParams common pagination parameters
type QueryParams struct {
	Limit  int64 `json:"limit,omitempty"`
	Offset int64 `json:"offset,omitempty"`
}

// SearchQuery for store search operations
type SearchQuery struct {
	Query  string `json:"query,omitempty"`
	Limit  int64  `json:"limit,omitempty"`
	Offset int64  `json:"offset,omitempty"`
	Store  string `json:"store,omitempty"`
}

// LogLevel represents the severity level of log messages (alias to internal/types).
type LogLevel = sf.LogLevel

const (
	// LevelError logs error messages only.
	LevelError = sf.LevelError
	// LevelWarn logs warnings and errors.
	LevelWarn = sf.LevelWarn
	// LevelInfo logs informational messages, warnings, and errors.
	LevelInfo = sf.LevelInfo
	// LevelDebug logs all messages including debug details.
	LevelDebug = sf.LevelDebug
)

// LoggingConfig holds configuration for logging behavior (alias to internal/types).
type LoggingConfig = sf.LoggingConfig

// DefaultLoggingConfig returns a default logging configuration.
func DefaultLoggingConfig() *LoggingConfig {
	return sf.DefaultLoggingConfig()
}
