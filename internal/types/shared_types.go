package types

import (
	"io"
	"time"
)

// LogLevel represents the severity level of log messages.
type LogLevel int

const (
	// LevelError logs error messages only.
	LevelError LogLevel = iota
	// LevelWarn logs warnings and errors.
	LevelWarn
	// LevelInfo logs informational messages, warnings, and errors.
	LevelInfo
	// LevelDebug logs all messages including debug details.
	LevelDebug
)

func (l LogLevel) String() string {
	switch l {
	case LevelError:
		return "ERROR"
	case LevelWarn:
		return "WARN"
	case LevelInfo:
		return "INFO"
	case LevelDebug:
		return "DEBUG"
	default:
		return "UNKNOWN"
	}
}

// LoggingConfig holds configuration for logging behavior.
type LoggingConfig struct {
	LogLevel             LogLevel // Minimum log level to output
	AnnotateRawResponses bool     // When true, adds visual separators for raw response logging (console-friendly)
	EnableRequestBody    bool     // Whether to log request bodies in debug mode
}

// DefaultLoggingConfig returns a default configuration.
func DefaultLoggingConfig() *LoggingConfig {
	return &LoggingConfig{
		LogLevel:             LevelInfo,
		AnnotateRawResponses: false,
		EnableRequestBody:    true,
	}
}

// ClientConfig configuration for client initialization
type ClientConfig struct {
	ServerURL    string         // Server URL without API path, e.g., "https://api.fleetbase.io" or custom endpoint
	APIPath      string         // API path appended to server URL (default: "/storefront/v1")
	APIKey       string         // Required: Bearer token
	Timeout      time.Duration  // Default 30s
	LogLevel     LogLevel       // Minimum log level (default: LevelInfo)
	LoggerOutput io.Writer      // Output writer (default: os.Stdout)
	LogConfig    *LoggingConfig // Custom logging configuration
}
