package trace

import (
	"context"

	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/label"
)

type Tracer trace.Tracer
type Span trace.Span
type Key label.Key

func SpanFromContext(ctx context.Context) Span {
	return trace.SpanFromContext(ctx)
}
