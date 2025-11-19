package scraper

import (
	"net/http"
	"time"
)

// HTTPClient interface for making HTTP requests
type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

// defaultHTTPClient provides a configured HTTP client with timeouts
type defaultHTTPClient struct {
	client *http.Client
}

// NewHTTPClient creates a new HTTP client with proper timeouts
func NewHTTPClient() HTTPClient {
	return &defaultHTTPClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Get performs an HTTP GET request
func (c *defaultHTTPClient) Get(url string) (*http.Response, error) {
	return c.client.Get(url)
}

// GRPC Client interface for gRPC communication
type GRPCClient interface {
	Connect(address string) error
	Close() error
}
