package config

import (
	otel "go.opentelemetry.io/otel"
)

const (
	DefaultTransportType string = "default"
)

var (
	InvalidSpanId otel.SpanID = otel.SpanID{}
)
