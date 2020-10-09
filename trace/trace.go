package trace

import (
	otel "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
)

type TraceID otel.TraceID
type Tracer otel.Tracer
type Span otel.Span
type Key label.Key
