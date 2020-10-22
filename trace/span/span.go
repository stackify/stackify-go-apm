package span

import (
	"fmt"

	apitrace "go.opentelemetry.io/otel/api/trace"
	export "go.opentelemetry.io/otel/sdk/export/trace"
	"go.stackify.com/apm/config"
	"go.stackify.com/apm/utils"
)

var (
	InvalidSpanId       apitrace.SpanID = apitrace.SpanID{}
	validSpanAttributes                 = map[string]string{
		"http.method":      "METHOD",
		"http.url":         "URL",
		"http.status_code": "STATUS",
	}
	instrumentationType = map[string]string{
		"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp": "net.http",
	}
	categoryType = map[string]string{
		"GET": "Web External",
	}
	subCategoryType = map[string]string{
		"GET": "Execute",
	}
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

func NewSpan(c *config.Config, sd *export.SpanData) StackifySpan {
	var spanName string = sd.Name
	instrumentation, ok := instrumentationType[sd.InstrumentationLibrary.Name]
	if ok {
		spanName = fmt.Sprintf("%s.%s", instrumentation, sd.Name)
	}
	sspan := StackifySpan{
		Id:       utils.SpanIdToString(sd.SpanContext.SpanID[:]),
		ParentId: utils.SpanIdToString(sd.ParentSpanID[:]),
		Call:     spanName,
		ReqBegin: utils.TimeToTimestamp(sd.StartTime),
		ReqEnd:   utils.TimeToTimestamp(sd.EndTime),
		Props:    make(map[string]string),
		Stacks:   []*StackifySpan{},
	}

	if sd.ParentSpanID == InvalidSpanId {
		sspan.Props["PROFILER_VERSION"] = "v3"
		sspan.Props["CATEGORY"] = "Go"
		sspan.Props["TRACE_ID"] = utils.TranceIdToUUID(sd.SpanContext.TraceID[:])
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
		category, ok := categoryType[sd.Name]
		if !ok {
			category = "Go"
		}
		sspan.Props["CATEGORY"] = category
		sspan.Props["SUBCATEGORY"] = subCategoryType[sd.Name]

		for _, attribute := range sd.Attributes {
			key, ok := validSpanAttributes[string(attribute.Key)]
			if ok {
				sspan.Props[key] = fmt.Sprintf("%v", attribute.Value.AsInterface())
			}
		}
	}
	return sspan
}
