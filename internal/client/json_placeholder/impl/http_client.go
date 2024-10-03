package jsonplaceholder_impl

import (
	"net/http"
	"time"
)

const baseURL = "https://jsonplaceholder.typicode.com"

// newHTTPClient creates a new HTTP client with the given configuration
func newHTTPClient(config *ServiceConfig) *http.Client {
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
		ForceAttemptHTTP2:   true,
		TLSClientConfig:     config.TLSConfig,
	}

	return &http.Client{
		Timeout:   config.Timeout,
		Transport: transport,
	}
}
