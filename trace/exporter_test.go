package trace_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	apitrace "go.opentelemetry.io/otel/api/trace"
	export "go.opentelemetry.io/otel/sdk/export/trace"

	"github.com/stackify/stackify-go-apm/trace"
	"github.com/stackify/stackify-go-apm/transport"
)

func TestTransportExportSpans(t *testing.T) {
	c := createConfig()
	tt := &testTransport{}
	tr := transport.Transport(tt)
	exporter := trace.NewStackifySpanExporter(c, &tr)
	sd := createOtelSpanData(c, "custom", apitrace.SpanKindClient, parentSpanID, invalidSpanId)

	exporter.ExportSpans(context.Background(), []*export.SpanData{sd})

	assert.Equal(t, tt.HandleTraceCallCount, 1)
	assert.Equal(t, tt.SendAllCallCount, 0)
}

func TestTransportShutdown(t *testing.T) {
	c := createConfig()
	tt := &testTransport{}
	tr := transport.Transport(tt)
	exporter := trace.NewStackifySpanExporter(c, &tr)

	exporter.Shutdown(context.Background())

	assert.Equal(t, tt.HandleTraceCallCount, 0)
	assert.Equal(t, tt.SendAllCallCount, 1)
}

func TestParentChildSpans(t *testing.T) {
	c := createConfig()
	tt := &testTransport{}
	tr := transport.Transport(tt)
	exporter := trace.NewStackifySpanExporter(c, &tr)
	parent := createOtelSpanData(c, "custom", apitrace.SpanKindClient, parentSpanID, invalidSpanId)
	child := createOtelSpanData(c, "custom", apitrace.SpanKindClient, childSpanID, parentSpanID)

	exporter.ExportSpans(context.Background(), []*export.SpanData{parent, child})

	assert.Equal(t, tt.HandleTraceCallCount, 1)
	assert.Equal(t, tt.SendAllCallCount, 0)
}
