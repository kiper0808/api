package http

import (
	"net"
	"net/http"
	"time"

	"github.com/mwitkow/go-conntrack"
)

func newDefaultHTTPTransport() *http.Transport {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = 100
	transport.MaxIdleConnsPerHost = 100

	return transport
}

func newConntrackRoundTripper() http.RoundTripper {
	transport := newDefaultHTTPTransport()
	transport.DialContext = conntrack.NewDialContextFunc(
		conntrack.DialWithName("dostavkeeClient"),
		conntrack.DialWithDialer(&net.Dialer{
			Timeout:   7 * time.Second,
			KeepAlive: 30 * time.Second,
		}),
	)

	return transport
}

// newTinkoffHttpClientRoundTripper - транспорт для API Тинькофф
func newTinkoffHttpClientRoundTripper() http.RoundTripper {
	transport := newDefaultHTTPTransport()

	transport.DialContext = conntrack.NewDialContextFunc(
		conntrack.DialWithName("bankAPIClient"),
		conntrack.DialWithDialer(&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}),
	)

	return transport
}
