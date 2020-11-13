package trace_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stackify/stackify-go-apm/trace"
	"github.com/stackify/stackify-go-apm/transport"
	apitrace "go.opentelemetry.io/otel/api/trace"
)

func TestSpanProcessorOnStartAndOnEnd(t *testing.T) {
	c := createConfig()
	tt := &testTransport{}
	tr := transport.Transport(tt)
	spe := trace.NewStackifySpanExporter(c, &tr)
	ssp := trace.NewStackifySpanProcessor(spe)
	sd := createOtelSpanData(c, "custom", apitrace.SpanKindClient, parentSpanID, invalidSpanId)

	ssp.OnStart(sd)
	ssp.OnEnd(sd)

	// sleep making sure we processes span convertion
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, tt.HandleTraceCallCount, 1)
}

func TestSpanProcessorShutDown(t *testing.T) {
	c := createConfig()
	tt := &testTransport{}
	tr := transport.Transport(tt)
	spe := trace.NewStackifySpanExporter(c, &tr)
	ssp := trace.NewStackifySpanProcessor(spe)
	sd := createOtelSpanData(c, "custom", apitrace.SpanKindClient, parentSpanID, invalidSpanId)

	ssp.OnStart(sd)
	ssp.OnEnd(sd)
	ssp.Shutdown()

	// sleep making sure we processes span convertion
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, tt.HandleTraceCallCount, 1)
	assert.Equal(t, tt.SendAllCallCount, 1)
}

func TestSpanProcessorForceFlush(t *testing.T) {
	c := createConfig()
	tt := &testTransport{}
	tr := transport.Transport(tt)
	spe := trace.NewStackifySpanExporter(c, &tr)
	ssp := trace.NewStackifySpanProcessor(spe)
	sd := createOtelSpanData(c, "custom", apitrace.SpanKindClient, parentSpanID, invalidSpanId)

	ssp.OnStart(sd)
	ssp.OnEnd(sd)
	ssp.ForceFlush()

	// sleep making sure we processes span convertion
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, tt.HandleTraceCallCount, 1)
}
