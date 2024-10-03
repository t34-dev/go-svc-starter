package jsonplaceholder_impl

import (
	"crypto/tls"
	"time"
)

// ServiceOption defines a function type for configuring ServiceConfig
type ServiceOption func(*ServiceConfig)

// ServiceConfig contains configuration parameters for the service
type ServiceConfig struct {
	Timeout   time.Duration
	TLSConfig *tls.Config
	APIKey    string
	UserAgent string
	RateLimit int
}

// WithTimeout sets the timeout for the HTTP client
func WithTimeout(timeout time.Duration) ServiceOption {
	return func(c *ServiceConfig) {
		c.Timeout = timeout
	}
}

// WithTLSConfig sets a custom TLS configuration
func WithTLSConfig(config *tls.Config) ServiceOption {
	return func(c *ServiceConfig) {
		c.TLSConfig = config
	}
}

// WithAPIKey sets the API key for authentication
func WithAPIKey(apiKey string) ServiceOption {
	return func(c *ServiceConfig) {
		c.APIKey = apiKey
	}
}

// WithUserAgent sets a custom User-Agent
func WithUserAgent(userAgent string) ServiceOption {
	return func(c *ServiceConfig) {
		c.UserAgent = userAgent
	}
}

// WithRateLimit sets the limit on the number of requests per second
func WithRateLimit(requestsPerSecond int) ServiceOption {
	return func(c *ServiceConfig) {
		c.RateLimit = requestsPerSecond
	}
}
