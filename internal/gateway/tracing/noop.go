package tracing

import (
	"context"

	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

type noopProvider struct {
	trace.TracerProvider
}

func NewNoopProvider() TracerProvider {
	return &noopProvider{
		TracerProvider: noop.NewTracerProvider(),
	}
}

func (p noopProvider) Shutdown(context.Context) error {
	return nil
}

func (p noopProvider) RegisterSpanProcessor(sp tracesdk.SpanProcessor) {}
