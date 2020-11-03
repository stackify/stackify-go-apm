package utils_test

import (
	"testing"
	"time"

	"go.opentelemetry.io/otel/api/trace"
	apitrace "go.opentelemetry.io/otel/api/trace"

	"go.stackify.com/apm/utils"
)

func TestTimeToTimestamp(t *testing.T) {
	layout := "2006-01-02T15:04:05.000Z"
	str := "2014-11-12T11:45:26.371Z"
	myTime, _ := time.Parse(layout, str)

	timestamp := utils.TimeToTimestamp(myTime)
	if timestamp != "1415792726371.0000" {
		t.Errorf("Error converting time to timestamp string\n")
	}
}

func TestTraceIdToUUID(t *testing.T) {
	traceId := trace.ID{}

	uuid := utils.TraceIdToUUID(traceId[:])
	if uuid != "00000000-0000-0000-0000-000000000000" {
		t.Errorf("Error converting traceId to UUID\n")
	}
}

func TestSpanIdToString(t *testing.T) {
	spanId := apitrace.SpanID{}

	strSpanId := utils.SpanIdToString(spanId[:])
	if strSpanId != "0" {
		t.Errorf("Error converting spasnId to string\n")
	}
}
