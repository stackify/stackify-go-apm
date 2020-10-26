package span

import (
	"fmt"
	"strings"

	apitrace "go.opentelemetry.io/otel/api/trace"
	export "go.opentelemetry.io/otel/sdk/export/trace"
	"go.stackify.com/apm/config"
	"go.stackify.com/apm/utils"
)

var (
	InvalidSpanId       apitrace.SpanID = apitrace.SpanID{}
	instrumentationType                 = map[string]string{
		"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp":                                    "net.http",
		"go.opentelemetry.io/contrib/instrumentation/github.com/bradfitz/gomemcache/memcache/otelmemcache": "gomemcached",
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
	var spanAttributes = map[string]string{}
	for _, attribute := range sd.Attributes {
		spanAttributes[string(attribute.Key)] = fmt.Sprintf("%v", attribute.Value.AsInterface())
	}

	sspan := StackifySpan{
		Id:       utils.SpanIdToString(sd.SpanContext.SpanID[:]),
		ParentId: utils.SpanIdToString(sd.ParentSpanID[:]),
		Call:     sd.Name,
		ReqBegin: utils.TimeToTimestamp(sd.StartTime),
		ReqEnd:   utils.TimeToTimestamp(sd.EndTime),
		Props:    make(map[string]string),
		Stacks:   []*StackifySpan{},
	}

	if sd.ParentSpanID == InvalidSpanId {
		var tracetype string
		if sd.SpanKind.String() == "server" {
			tracetype = "WEBAPP"
		} else {
			tracetype = "TASK"
		}

		if IsHTTPSpan(spanAttributes) && !strings.HasPrefix(sspan.Call, "/") {
			sspan.Call = spanAttributes["http.target"]
		}

		sspan.Props["PROFILER_VERSION"] = "prototype"
		sspan.Props["CATEGORY"] = "Go"
		sspan.Props["TRACE_ID"] = utils.TranceIdToUUID(sd.SpanContext.TraceID[:])
		sspan.Props["TRACE_SOURCE"] = "GO"
		sspan.Props["TRACE_TARGET"] = "RETRACE"
		sspan.Props["TRACE_VERSION"] = "2.0"
		sspan.Props["TRACETYPE"] = tracetype
		sspan.Props["HOST_NAME"] = c.HostName
		sspan.Props["OS_TYPE"] = c.OSType
		sspan.Props["PROCESS_ID"] = c.ProcessID
		sspan.Props["APPLICATION_PATH"] = "/"
		sspan.Props["APPLICATION_FILESYSTEM_PATH"] = c.BaseDIR
		sspan.Props["APPLICATION_NAME"] = c.ApplicationName
		sspan.Props["APPLICATION_ENV"] = c.EnvironmentName
		SetSpanPropsIfAvailable(&sspan, "REPORTING_URL", spanAttributes, "http.target", sspan.Call)
		SetSpanPropsIfAvailable(&sspan, "METHOD", spanAttributes, "http.method", "")
		SetSpanPropsIfAvailable(&sspan, "STATUS", spanAttributes, "http.status_code", "")
		SetSpanPropsIfAvailable(&sspan, "URL", spanAttributes, "http.url", "")
	} else {
		instrumentation, ok := instrumentationType[sd.InstrumentationLibrary.Name]
		if ok {
			spanName := fmt.Sprintf("%s.%s", instrumentation, sd.Name)
			sspan.Call = spanName
		}
		sspan.Props["CATEGORY"] = "Go"

		if IsHTTPSpan(spanAttributes) {
			sspan.Props["CATEGORY"] = "Web External"
			sspan.Props["SUBCATEGORY"] = "Execute"
			sspan.Props["COMPONENT_CATEGORY"] = "Web External"
			sspan.Props["COMPONENT_DETAIL"] = "Execute"
			SetSpanPropsIfAvailable(&sspan, "METHOD", spanAttributes, "http.method", "")
			SetSpanPropsIfAvailable(&sspan, "STATUS", spanAttributes, "http.status_code", "")
			SetSpanPropsIfAvailable(&sspan, "URL", spanAttributes, "http.url", "")
		}

		if IsTemplateSpan(spanAttributes) {
			sspan.Props["CATEGORY"] = "Template"
			sspan.Props["SUBCATEGORY"] = "Render"
			sspan.Props["COMPONENT_CATEGORY"] = "Template"
			sspan.Props["COMPONENT_DETAIL"] = "Template"
		}

		if IsMemcachedSpan(spanAttributes) {
			sspan.Props["CATEGORY"] = "Cache"
			sspan.Props["SUBCATEGORY"] = "Execute"
			sspan.Props["COMPONENT_CATEGORY"] = "Cache"
			sspan.Props["COMPONENT_DETAIL"] = "Execute"
			SetSpanPropsIfAvailable(&sspan, "OPERATION", spanAttributes, "db.operation", "")
			SetSpanPropsIfAvailable(&sspan, "CACHEKEY", spanAttributes, "db.memcached.item", "")
		}
	}

	return sspan
}

func SetSpanPropsIfAvailable(sspan *StackifySpan, sspanKey string, attributes map[string]string, attributeKey string, defaultValue string) {
	value, ok := attributes[attributeKey]
	if ok {
		sspan.Props[sspanKey] = value
	} else if len(defaultValue) > 0 {
		sspan.Props[sspanKey] = defaultValue
	}
}

func IsHTTPSpan(spanAttributes map[string]string) bool {
	return isAttributePresent("http.method", spanAttributes) && isAttributePresent("http.status_code", spanAttributes)
}

func IsTemplateSpan(spanAttributes map[string]string) bool {
	return isAttributePresent("go.template", spanAttributes)
}

func IsMemcachedSpan(spanAttributes map[string]string) bool {
	return isAttributePresent("db.operation", spanAttributes) && isAttributeValueEqualTo("db.system", spanAttributes, "memcached")
}

func isAttributePresent(attrName string, spanAttributes map[string]string) bool {
	_, ok := spanAttributes[attrName]
	if ok {
		return true
	}
	return false
}

func isAttributeValueEqualTo(attrName string, spanAttributes map[string]string, value string) bool {
	val, ok := spanAttributes[attrName]
	if ok {
		return val == value
	}
	return false
}
