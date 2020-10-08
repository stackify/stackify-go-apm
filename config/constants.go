package config

import (
	apitrace "go.opentelemetry.io/otel/api/trace"
)

const (
	DefaultTransportType string = "default"
)

var (
	InvalidSpanId apitrace.SpanID = apitrace.SpanID{}
)
