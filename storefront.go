package storefront

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	sf "github.com/artyconst/fb-storefront-go/internal/types"
	sfconfig "github.com/artyconst/fb-storefront-go/pkg/config"
)

// ConfigOption is a function that modifies *ClientConfig.
type ConfigOption = sfconfig.Option

// QueryParams common pagination parameters (alias for pkg/config.QueryParams).
type QueryParams = sfconfig.QueryParams

// SearchQuery for store search operations (alias for pkg/config.SearchQuery).
type SearchQuery = sfconfig.SearchQuery

// ClientConfig configuration for client initialization.
type ClientConfig = sf.ClientConfig

// LogLevel represents the severity level of log messages.
type LogLevel = sf.LogLevel

// LoggingConfig holds configuration for logging behavior.
type LoggingConfig = sf.LoggingConfig

// DefaultLoggingConfig returns a default configuration.
func DefaultLoggingConfig() *sf.LoggingConfig {
	return sf.DefaultLoggingConfig()
}

// Config functions re-exported from pkg/config for backward compatibility.
var (
	WithAPIHost      = sfconfig.WithAPIHost
	WithAPIPath      = sfconfig.WithAPIPath
	WithTimeout      = sfconfig.WithTimeout
	WithLogLevel     = sfconfig.WithLogLevel
	WithLoggerOutput = sfconfig.WithLoggerOutput
	WithDebugMode    = sfconfig.WithDebugMode
)

// StorefrontClient is the main entry point for API interactions
type StorefrontClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	logger     *StdLogger
	rawLogger  *RawResponseLogger
}

// RequestOptions holds optional parameters for HTTP requests.
type RequestOptions struct {
	CustomerToken string // Customer authentication token for authenticated customer endpoints
}

// RequestOption is a functional option for modifying request options.
type RequestOption func(*RequestOptions)

// WithCustomerToken returns a RequestOption that sets the Customer-Token header.
// Use this for endpoints requiring customer authentication (e.g., ListPlaces, ListOrders).
func WithCustomerToken(token string) RequestOption {
	return func(opts *RequestOptions) {
		opts.CustomerToken = token
	}
}

// NewStorefront is a convenience function that creates a StorefrontClient
// using the functional options pattern for configuration.
//
// Usage:
//
//	sf, err := storefront.NewStorefront("sk_test_key",
//		storefront.WithAPIHost("https://custom.api.com"),
//		storefront.WithTimeout(60),
//	)
func NewStorefront(apiKey string, opts ...ConfigOption) (*StorefrontClient, error) {
	config := ClientConfig{
		APIKey: apiKey,
	}

	for _, opt := range opts {
		opt(&config)
	}

	return newStorefrontClient(config)
}

// newStorefrontClient creates a configured client instance.
func newStorefrontClient(config ClientConfig) (*StorefrontClient, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	if config.ServerURL == "" {
		config.ServerURL = "https://api.fleetbase.io"
	}

	if config.APIPath == "" {
		config.APIPath = "/storefront/v1"
	}

	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	// Compose base URL from server URL and API path
	baseURL := config.ServerURL + config.APIPath
	// Set up logging
	output := config.LoggerOutput
	if output == nil {
		output = os.Stdout
	}

	logger := NewStdLogger(output, config.LogLevel)

	var rawLogger *RawResponseLogger
	if config.LogConfig != nil {
		rawLogger = NewRawResponseLogger(logger, config.LogConfig)
	} else if config.LogLevel >= sf.LevelDebug {
		// Create default logging config for debug mode
		rawLogger = NewRawResponseLogger(logger, DefaultLoggingConfig())
	}

	return &StorefrontClient{
		baseURL:    baseURL,
		apiKey:     config.APIKey,
		httpClient: &http.Client{Timeout: config.Timeout},
		logger:     logger,
		rawLogger:  rawLogger,
	}, nil
}

// doRequest executes HTTP request with Bearer authentication and optional Customer-Token header.
func (c *StorefrontClient) doRequest(ctx context.Context, method, endpoint string, body io.Reader, opts ...RequestOption) (*http.Response, error) {

	// URL composition is validated once at initialization time via WithAPIHost/WithAPIPath
	// No per-request validation needed: baseURL and endpoint are guaranteed correct format
	url := c.baseURL + endpoint

	// Log request details in debug mode before making request
	if c.rawLogger != nil {
		c.rawLogger.LogRequest(method, url, body)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	// Apply Customer-Token header if provided via options
	var reqOpts RequestOptions
	for _, opt := range opts {
		opt(&reqOpts)
	}
	if reqOpts.CustomerToken != "" {
		req.Header.Set("Customer-Token", reqOpts.CustomerToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

// GetJSON performs GET request and unmarshals JSON response.
func (c *StorefrontClient) GetJSON(ctx context.Context, endpoint string, v interface{}, opts ...RequestOption) error {
	resp, err := c.doRequest(ctx, http.MethodGet, endpoint, nil, opts...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Log response in debug mode regardless of status code
	if c.rawLogger != nil && resp.StatusCode > 0 {
		c.rawLogger.LogResponse(resp.StatusCode, body)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiErr struct {
			Message string `json:"message"`
			Code    string `json:"code,omitempty"`
		}
		json.Unmarshal(body, &apiErr)
		errMsg := apiErr.Message
		if errMsg == "" {
			errMsg = string(body)
		}
		c.logger.Debug("API ERROR: %s", errMsg)
		return fmt.Errorf("API error %d: %s (code: %s)", resp.StatusCode, errMsg, apiErr.Code)
	}

	return json.Unmarshal(body, v)
}

// PostJSON performs POST request with JSON body and unmarshals response.
func (c *StorefrontClient) PostJSON(ctx context.Context, endpoint string, data interface{}, v interface{}, opts ...RequestOption) error {
	bodyBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodPost, endpoint, bytes.NewReader(bodyBytes), opts...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Log response in debug mode regardless of status code
	if c.rawLogger != nil && resp.StatusCode > 0 {
		c.rawLogger.LogResponse(resp.StatusCode, responseBody)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiErr struct {
			Message string `json:"message"`
			Code    string `json:"code,omitempty"`
		}
		json.Unmarshal(responseBody, &apiErr)
		errMsg := apiErr.Message
		if errMsg == "" {
			errMsg = string(responseBody)
		}
		c.logger.Debug("API ERROR: %s", errMsg)
		return fmt.Errorf("API error %d: %s (code: %s)", resp.StatusCode, errMsg, apiErr.Code)
	}

	return json.Unmarshal(responseBody, v)
}

// PutJSON performs PUT request with JSON body and unmarshals response.
func (c *StorefrontClient) PutJSON(ctx context.Context, endpoint string, data interface{}, v interface{}, opts ...RequestOption) error {
	bodyBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodPut, endpoint, bytes.NewReader(bodyBytes), opts...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Log response in debug mode regardless of status code
	if c.rawLogger != nil && resp.StatusCode > 0 {
		c.rawLogger.LogResponse(resp.StatusCode, responseBody)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiErr struct {
			Message string `json:"message"`
			Code    string `json:"code,omitempty"`
		}
		json.Unmarshal(responseBody, &apiErr)
		errMsg := apiErr.Message
		if errMsg == "" {
			errMsg = string(responseBody)
		}
		c.logger.Debug("API ERROR: %s", errMsg)
		return fmt.Errorf("API error %d: %s (code: %s)", resp.StatusCode, errMsg, apiErr.Code)
	}

	return json.Unmarshal(responseBody, v)
}

// DeleteJSON performs DELETE request and unmarshals JSON response.
func (c *StorefrontClient) DeleteJSON(ctx context.Context, endpoint string, v interface{}, opts ...RequestOption) error {
	resp, err := c.doRequest(ctx, http.MethodDelete, endpoint, nil, opts...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Log response in debug mode regardless of status code
	if c.rawLogger != nil && resp.StatusCode > 0 {
		c.rawLogger.LogResponse(resp.StatusCode, body)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiErr struct {
			Message string `json:"message"`
			Code    string `json:"code,omitempty"`
		}
		json.Unmarshal(body, &apiErr)
		errMsg := apiErr.Message
		if errMsg == "" {
			errMsg = string(body)
		}
		c.logger.Debug("API ERROR: %s", errMsg)
		return fmt.Errorf("API error %d: %s (code: %s)", resp.StatusCode, errMsg, apiErr.Code)
	}

	return json.Unmarshal(body, v)
}
