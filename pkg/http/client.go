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
			Timeout: timeout,
		}
	})

	return &Client{Client: &httpClient}
}
