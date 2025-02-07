package http

import (
	"net/http"
	"sync"
	"time"
)

type Client struct {
	*http.Client
}

// Singleton
var (
	httpClient http.Client
	once       sync.Once
)

func NewHTTPClient(timeout time.Duration) *Client {
	once.Do(func() {
		httpClient = http.Client{
			Timeout:   timeout,
			Transport: newConntrackRoundTripper(),
		}
	})

	return &Client{Client: &httpClient}
}

// NewTinkoffHttpClient - отдельный HTTP клиент для API Тинькофф
func NewTinkoffHttpClient(timeout time.Duration) *Client {
	return &Client{
		Client: &http.Client{
			Timeout:   timeout,
			Transport: newTinkoffHttpClientRoundTripper(),
		},
	}
}
