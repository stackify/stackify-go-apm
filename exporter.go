package apm

import (
	"context"

	export "go.opentelemetry.io/otel/sdk/export/trace"
)

type StackifySpanExporter struct {
	c *Config
	t Transport
}

// ExportSpans method convert spans into stackify span format.
func (ssp *StackifySpanExporter) ExportSpans(ctx context.Context, spans []*export.SpanData) error {
	trace, err := ssp.toStackifyTrace(spans)
	if err != nil {
		return nil
	}
	ssp.t.HandleTrace(trace)
	return nil
}

// Shutdown method ensures we send all valid spans from queue.
func (ssp *StackifySpanExporter) Shutdown(context.Context) error {
	ssp.t.SendAll()
	return nil
}

// toStackifyTrace method converts spans to stackify trace format.
func (ssp *StackifySpanExporter) toStackifyTrace(sd []*export.SpanData) (*StackifySpan, error) {
	stackifySpans := []*StackifySpan{}
	stackifySpansMap := make(map[string]*StackifySpan)
	stackifySpan := &StackifySpan{}

	for _, s := range sd {
		stackifySpan := NewSpan(ssp.c, s)
		stackifySpans = append(stackifySpans, &stackifySpan)
		stackifySpansMap[SpanIdToString(s.SpanContext.SpanID[:])] = &stackifySpan
	}

	for _, s := range stackifySpans {
		if s.ParentId != s.Id && s.ParentId != SpanIdToString(InvalidSpanId[:]) {
			stackifySpansMap[s.ParentId].Stacks = append(stackifySpansMap[s.ParentId].Stacks, stackifySpansMap[s.Id])
		} else {
			stackifySpan = s
		}
	}
	return stackifySpan, nil
}

// NewStackifySpanExporter function creates a StackifySpanExporter.
func NewStackifySpanExporter(c *Config, t *Transport) *StackifySpanExporter {
	return &StackifySpanExporter{
		c: c,
		t: *t,
	}
}
