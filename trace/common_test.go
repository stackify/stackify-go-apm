package trace_test

import (
	"time"

	"github.com/stackify/stackify-go-apm/config"
	"github.com/stackify/stackify-go-apm/trace/span"
	apitrace "go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/codes"
	export "go.opentelemetry.io/otel/sdk/export/trace"
)

var (
	invalidSpanId = apitrace.SpanID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	parentSpanID  = apitrace.SpanID{0xFF, 0xFE, 0xFD, 0xFC, 0xFB, 0xFA, 0xF9, 0xF8}
	childSpanID   = apitrace.SpanID{0xEF, 0xEE, 0xED, 0xEC, 0xEB, 0xEA, 0xE9, 0xE8}
)

func createConfig() *config.Config {
	return config.NewConfig(
		config.WithApplicationName("TestApp"),
		config.WithEnvironmentName("test"),
	)
}

func createOtelSpanData(c *config.Config, name string, kind apitrace.SpanKind, spanId apitrace.SpanID, parentId apitrace.SpanID) *export.SpanData {
	c.ProcessID = "ProcessID"
	c.HostName = "HostName"
	c.OSType = "OSType"
	c.BaseDIR = "BaseDIR"
	startTime := time.Unix(1585674086, 1234)
	endTime := startTime.Add(10 * time.Second)
	return &export.SpanData{
		SpanContext: apitrace.SpanContext{
			TraceID: apitrace.ID{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
			SpanID:  spanId,
		},
		SpanKind:     kind,
		ParentSpanID: parentId,
		Name:         name,
		StartTime:    startTime,
		EndTime:      endTime,
		StatusCode:   codes.Error,
	}
}

type testTransport struct {
	HandleTraceCallCount int
	SendAllCallCount     int
}

func (dt *testTransport) HandleTrace(stackifySpan *span.StackifySpan) {
	dt.HandleTraceCallCount += 1
}
func (dt *testTransport) SendAll() {
	dt.SendAllCallCount += 1
}
