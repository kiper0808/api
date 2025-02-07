package tracing

import (
	"context"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// TracerProvider is an interface that wraps the trace.TracerProvider and adds the Shutdown method
// that is not part of the interface but is part of the implementation.
type TracerProvider interface {
	trace.TracerProvider
	Shutdown(context.Context) error
	RegisterSpanProcessor(sp tracesdk.SpanProcessor)
}

// UpdateHTTPClient updates the http client with the necessary otel transport.
func UpdateHTTPClient(client *http.Client, tracerProvider trace.TracerProvider) {
	client.Transport = otelhttp.NewTransport(
		client.Transport,
		otelhttp.WithTracerProvider(tracerProvider),
		otelhttp.WithSpanNameFormatter(func(_ string, r *http.Request) string {
			return r.Method + " " + r.URL.Path
		}),
	)
}
