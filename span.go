package apm

import (
	export "go.opentelemetry.io/otel/sdk/export/trace"
)

type StackifySpan struct {
	Id       string            `json:"id"`
	ParentId string            `json:"-"`
	Call     string            `json:"call"`
	ReqBegin string            `json:"reqBegin"`
	ReqEnd   string            `json:"reqEnd"`
	Props    map[string]string `json:"props"`
	Stacks   []*StackifySpan   `json:"stacks"`
	// Exceptions
}

func NewSpan(c *Config, sd *export.SpanData) StackifySpan {
	sspan := StackifySpan{
		Id:       SpanIdToString(sd.SpanContext.SpanID[:]),
		ParentId: SpanIdToString(sd.ParentSpanID[:]),
		Call:     sd.Name,
		ReqBegin: TimeToTimestamp(sd.StartTime),
		ReqEnd:   TimeToTimestamp(sd.EndTime),
		Props:    make(map[string]string),
		Stacks:   []*StackifySpan{},
	}

	if sd.ParentSpanID == InvalidSpanId {
		sspan.Props["PROFILER_VERSION"] = "v3"
		sspan.Props["CATEGORY"] = "Go"
		sspan.Props["TRACE_ID"] = TranceIdToUUID(sd.SpanContext.TraceID[:])
		sspan.Props["TRACE_SOURCE"] = "GO"
		sspan.Props["TRACE_TARGET"] = "RETRACE"
		sspan.Props["TRACE_VERSION"] = "2.0"
		sspan.Props["TRACETYPE"] = "TASK"
		sspan.Props["HOST_NAME"] = c.HostName
		sspan.Props["OS_TYPE"] = c.OSType
		sspan.Props["PROCESS_ID"] = c.ProcessID
		sspan.Props["APPLICATION_PATH"] = "/"
		sspan.Props["APPLICATION_FILESYSTEM_PATH"] = c.BaseDIR
		sspan.Props["APPLICATION_NAME"] = c.ApplicationName
		sspan.Props["APPLICATION_ENV"] = c.EnvironmentName
		sspan.Props["REPORTING_URL"] = sspan.Call
	} else {
		sspan.Props["CATEGORY"] = "Go"
	}
	return sspan
}
