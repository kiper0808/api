package http

import (
	"net/http"
	"sync"
	"time"
)

type GotenbergHTTPClient struct {
	*http.Client
}

// Singleton
var (
	gotenbergHTTPClient     http.Client
	onceGotenbergHTTPClient sync.Once
)

func NewGotenbergHTTPClient(timeout time.Duration) *Client {
	onceGotenbergHTTPClient.Do(func() {
		gotenbergHTTPClient = http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				ForceAttemptHTTP2:     true,
				MaxIdleConns:          1024,
				MaxIdleConnsPerHost:   1024,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			}}
	})

	return &Client{Client: &gotenbergHTTPClient}
}
