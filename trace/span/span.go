package span

import (
	"encoding/json"
	"fmt"
	"strings"

	apitrace "go.opentelemetry.io/otel/api/trace"
	export "go.opentelemetry.io/otel/sdk/export/trace"
	"go.stackify.com/apm/config"
	"go.stackify.com/apm/utils"
)

const (
	Otelhttp     = "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	Otelmemcache = "go.opentelemetry.io/contrib/instrumentation/github.com/bradfitz/gomemcache/memcache/otelmemcache"
	Otelgocql    = "go.opentelemetry.io/contrib/instrumentation/github.com/gocql/gocql/otelgocql"
	Otelgrpc     = "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
)

var (
	InvalidSpanId       apitrace.SpanID = apitrace.SpanID{}
	instrumentationType                 = map[string]string{
		Otelhttp:     "net.http",
		Otelmemcache: "gomemcached",
		Otelgocql:    "gocql",
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
		SetSpanPropsIfAvailable(&sspan, "METHOD", spanAttributes, "rpc.method", "")
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

		if IsCasandraSpan(spanAttributes) {
			subcategory := "Execute"
			if isAttributeValueEqualTo("db.operation", spanAttributes, "db.cassandra.connect") {
				subcategory = "Connect"
			}

			operation := ""
			if isAttributePresent("db.operation", spanAttributes) {
				operation, _ = spanAttributes["db.operation"]
			} else if isAttributePresent("db.statement", spanAttributes) {
				operation = "db.cassandra.query"
			}

			spanCall := fmt.Sprintf("%s.%s", instrumentation, operation)
			sspan.Call = spanCall

			sspan.Props["CATEGORY"] = "Cassandra"
			sspan.Props["SUBCATEGORY"] = subcategory
			SetSpanPropsIfAvailable(&sspan, "SQL", spanAttributes, "db.statement", "")
			SetSpanPropsIfAvailable(&sspan, "ROW_COUNT", spanAttributes, "db.cassandra.rows.returned", "")
		}

		if IsMongoDBSpan(spanAttributes) {
			sspan.Call = "db.mongodb.query"
			sspan.Props["CATEGORY"] = "MongoDB"
			sspan.Props["SUBCATEGORY"] = "Execute"
			sspan.Props["COMPONENT_CATEGORY"] = "DB Query"
			sspan.Props["COMPONENT_DETAIL"] = "Execute SQL Query"
			SetSpanPropsIfAvailable(&sspan, "PROVIDER", spanAttributes, "db.system", "")
			SetSpanPropsIfAvailable(&sspan, "MONGODB_COLLECTION", spanAttributes, "db.instance", "")
			SetSpanPropsIfAvailable(&sspan, "OPERATION", spanAttributes, "db.operation", "")

			if isAttributePresent("db.statement", spanAttributes) {
				database := spanAttributes["db.instance"]
				statement := []byte(spanAttributes["db.statement"])
				var raw map[string]interface{}
				json.Unmarshal(statement, &raw)
				collection := raw["insert"]
				sspan.Props["MONGODB_COLLECTION"] = fmt.Sprintf("%s.%s", database, collection)
			}
		}

		if IsGRPCSpan(spanAttributes) {
			sspan.Props["CATEGORY"] = "RPC"
			sspan.Props["SUBCATEGORY"] = "Execute"
			SetSpanPropsIfAvailable(&sspan, "PROVIDER", spanAttributes, "rpc.system", "")
			SetSpanPropsIfAvailable(&sspan, "SERVICE", spanAttributes, "rpc.service", "")
			SetSpanPropsIfAvailable(&sspan, "METHOD", spanAttributes, "rpc.method", "")
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

func IsCasandraSpan(spanAttributes map[string]string) bool {
	return (isAttributePresent("db.operation", spanAttributes) || isAttributePresent("db.statement", spanAttributes)) && isAttributeValueEqualTo("db.system", spanAttributes, "cassandra")
}

func IsMongoDBSpan(spanAttributes map[string]string) bool {
	return (isAttributePresent("db.operation", spanAttributes) || isAttributePresent("db.statement", spanAttributes)) && isAttributeValueEqualTo("db.system", spanAttributes, "mongodb")
}

func IsGRPCSpan(spanAttributes map[string]string) bool {
	return isAttributeValueEqualTo("rpc.system", spanAttributes, "grpc")
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
