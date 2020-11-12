package span_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/label"
	export "go.opentelemetry.io/otel/sdk/export/trace"
	"go.opentelemetry.io/otel/sdk/instrumentation"

	"github.com/stackify/stackify-go-apm/config"
	"github.com/stackify/stackify-go-apm/trace/span"
)

var (
	ParentSpan = trace.SpanID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	ChildSpan  = trace.SpanID{0xEF, 0xEE, 0xED, 0xEC, 0xEB, 0xEA, 0xE9, 0xE8}
	InsLibTest = instrumentation.Library{
		Name:    "go.opentelemetry.io/contrib/instrumentation/test",
		Version: "v0.0.1",
	}
	InsLibHttp = instrumentation.Library{
		Name:    "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp",
		Version: "v0.0.1",
	}
	InsLibMemcached = instrumentation.Library{
		Name:    "go.opentelemetry.io/contrib/instrumentation/github.com/bradfitz/gomemcache/memcache/otelmemcache",
		Version: "v0.0.1",
	}
)

func createConfig() *config.Config {
	return config.NewConfig(
		config.WithApplicationName("TestApp"),
		config.WithEnvironmentName("test"),
	)
}

func createOtelSpanData(c *config.Config, name string, kind trace.SpanKind, parentId trace.SpanID, attributes []label.KeyValue, inslib instrumentation.Library) *export.SpanData {
	c.ProcessID = "ProcessID"
	c.HostName = "HostName"
	c.OSType = "OSType"
	c.BaseDIR = "BaseDIR"
	startTime := time.Unix(1585674086, 1234)
	endTime := startTime.Add(10 * time.Second)
	return &export.SpanData{
		SpanContext: trace.SpanContext{
			TraceID: trace.ID{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
			SpanID:  trace.SpanID{0xFF, 0xFE, 0xFD, 0xFC, 0xFB, 0xFA, 0xF9, 0xF8},
		},
		SpanKind:               kind,
		ParentSpanID:           parentId,
		Name:                   name,
		StartTime:              startTime,
		EndTime:                endTime,
		StatusCode:             codes.Error,
		Attributes:             attributes,
		InstrumentationLibrary: inslib,
	}
}

func TestCustomParentSpan(t *testing.T) {
	c := createConfig()
	sd := createOtelSpanData(c, "custom", trace.SpanKindClient, ParentSpan, []label.KeyValue{}, InsLibTest)

	stackifySpan := span.NewSpan(c, sd)

	assert.Equal(t, stackifySpan.Call, "custom")
	assert.NotEmpty(t, stackifySpan.ReqBegin)
	assert.NotEmpty(t, stackifySpan.ReqEnd)
	assert.Equal(t, stackifySpan.Props["APPLICATION_NAME"], "TestApp")
	assert.Equal(t, stackifySpan.Props["APPLICATION_ENV"], "test")
	assert.Equal(t, stackifySpan.Props["APPLICATION_FILESYSTEM_PATH"], "BaseDIR")
	assert.Equal(t, stackifySpan.Props["CATEGORY"], "Go")
	assert.Equal(t, stackifySpan.Props["HOST_NAME"], "HostName")
	assert.Equal(t, stackifySpan.Props["OS_TYPE"], "OSType")
	assert.Equal(t, stackifySpan.Props["PROCESS_ID"], "ProcessID")
	assert.Equal(t, stackifySpan.Props["REPORTING_URL"], "custom")
	assert.Equal(t, stackifySpan.Props["TRACETYPE"], "TASK")
	assert.Equal(t, stackifySpan.Props["TRACE_ID"], "00010203-0405-0607-0809-0a0b0c0d0e0f")
	assert.Equal(t, stackifySpan.Props["TRACE_SOURCE"], "GO")
	assert.Equal(t, stackifySpan.Props["TRACE_TARGET"], "RETRACE")
	assert.Equal(t, stackifySpan.Props["TRACE_VERSION"], "2.0")
	assert.Equal(t, stackifySpan.Props["PROFILER_VERSION"], "prototype")
}

func TestHTTPParentSpan(t *testing.T) {
	c := createConfig()
	attributes := []label.KeyValue{
		label.String("http.target", "target"),
		label.String("http.method", "GET"),
		label.String("http.status_code", "200"),
		label.String("http.url", "testurl"),
	}
	sd := createOtelSpanData(c, "index", trace.SpanKindServer, ParentSpan, attributes, InsLibTest)

	stackifySpan := span.NewSpan(c, sd)

	assert.Equal(t, stackifySpan.Call, "target")
	assert.NotEmpty(t, stackifySpan.ReqBegin)
	assert.NotEmpty(t, stackifySpan.ReqEnd)
	assert.Equal(t, stackifySpan.Props["APPLICATION_NAME"], "TestApp")
	assert.Equal(t, stackifySpan.Props["APPLICATION_ENV"], "test")
	assert.Equal(t, stackifySpan.Props["APPLICATION_FILESYSTEM_PATH"], "BaseDIR")
	assert.Equal(t, stackifySpan.Props["CATEGORY"], "Go")
	assert.Equal(t, stackifySpan.Props["HOST_NAME"], "HostName")
	assert.Equal(t, stackifySpan.Props["OS_TYPE"], "OSType")
	assert.Equal(t, stackifySpan.Props["PROCESS_ID"], "ProcessID")
	assert.Equal(t, stackifySpan.Props["REPORTING_URL"], "target")
	assert.Equal(t, stackifySpan.Props["TRACETYPE"], "WEBAPP")
	assert.Equal(t, stackifySpan.Props["TRACE_ID"], "00010203-0405-0607-0809-0a0b0c0d0e0f")
	assert.Equal(t, stackifySpan.Props["TRACE_SOURCE"], "GO")
	assert.Equal(t, stackifySpan.Props["TRACE_TARGET"], "RETRACE")
	assert.Equal(t, stackifySpan.Props["TRACE_VERSION"], "2.0")
	assert.Equal(t, stackifySpan.Props["PROFILER_VERSION"], "prototype")
	assert.Equal(t, stackifySpan.Props["METHOD"], "GET")
	assert.Equal(t, stackifySpan.Props["STATUS"], "200")
	assert.Equal(t, stackifySpan.Props["URL"], "testurl")
}

func TestHTTPChildSpan(t *testing.T) {
	c := createConfig()
	attributes := []label.KeyValue{
		label.String("http.target", "target"),
		label.String("http.method", "GET"),
		label.String("http.status_code", "200"),
		label.String("http.url", "testurl"),
	}
	sd := createOtelSpanData(c, "GET", trace.SpanKindClient, ChildSpan, attributes, InsLibHttp)

	stackifySpan := span.NewSpan(c, sd)

	assert.Equal(t, stackifySpan.Call, "net.http.GET")
	assert.NotEmpty(t, stackifySpan.ReqBegin)
	assert.NotEmpty(t, stackifySpan.ReqEnd)
	assert.Equal(t, stackifySpan.Props["CATEGORY"], "Web External")
	assert.Equal(t, stackifySpan.Props["SUBCATEGORY"], "Execute")
	assert.Equal(t, stackifySpan.Props["COMPONENT_CATEGORY"], "Web External")
	assert.Equal(t, stackifySpan.Props["COMPONENT_DETAIL"], "Execute")
	assert.Equal(t, stackifySpan.Props["METHOD"], "GET")
	assert.Equal(t, stackifySpan.Props["STATUS"], "200")
	assert.Equal(t, stackifySpan.Props["URL"], "testurl")
}

func TestTemplateChildSpan(t *testing.T) {
	c := createConfig()
	attributes := []label.KeyValue{
		label.String("go.template", "value"),
	}
	sd := createOtelSpanData(c, "template.render", trace.SpanKindClient, ChildSpan, attributes, InsLibTest)

	stackifySpan := span.NewSpan(c, sd)

	assert.Equal(t, stackifySpan.Call, "template.render")
	assert.NotEmpty(t, stackifySpan.ReqBegin)
	assert.NotEmpty(t, stackifySpan.ReqEnd)
	assert.Equal(t, stackifySpan.Props["CATEGORY"], "Template")
	assert.Equal(t, stackifySpan.Props["SUBCATEGORY"], "Render")
	assert.Equal(t, stackifySpan.Props["COMPONENT_CATEGORY"], "Template")
	assert.Equal(t, stackifySpan.Props["COMPONENT_DETAIL"], "Template")
}

func TestMemcachedChildSpan(t *testing.T) {
	c := createConfig()
	attributes := []label.KeyValue{
		label.String("db.system", "memcached"),
		label.String("db.operation", "get"),
		label.String("db.memcached.item", "foo"),
	}
	sd := createOtelSpanData(c, "get", trace.SpanKindClient, ChildSpan, attributes, InsLibMemcached)

	stackifySpan := span.NewSpan(c, sd)

	assert.Equal(t, stackifySpan.Call, "gomemcached.get")
	assert.NotEmpty(t, stackifySpan.ReqBegin)
	assert.NotEmpty(t, stackifySpan.ReqEnd)
	assert.Equal(t, stackifySpan.Props["CATEGORY"], "Cache")
	assert.Equal(t, stackifySpan.Props["SUBCATEGORY"], "Execute")
	assert.Equal(t, stackifySpan.Props["COMPONENT_CATEGORY"], "Cache")
	assert.Equal(t, stackifySpan.Props["COMPONENT_DETAIL"], "Execute")
	assert.Equal(t, stackifySpan.Props["OPERATION"], "get")
	assert.Equal(t, stackifySpan.Props["CACHEKEY"], "foo")
}

func TestCasandraConnectSpan(t *testing.T) {
	c := createConfig()
	attributes := []label.KeyValue{
		label.String("db.system", "cassandra"),
		label.String("db.operation", "db.cassandra.connect"),
	}
	sd := createOtelSpanData(c, "test", trace.SpanKindClient, ChildSpan, attributes, InsLibTest)

	stackifySpan := span.NewSpan(c, sd)

	assert.Equal(t, stackifySpan.Call, "db.cassandra.connect")
	assert.NotEmpty(t, stackifySpan.ReqBegin)
	assert.NotEmpty(t, stackifySpan.ReqEnd)
	assert.Equal(t, stackifySpan.Props["CATEGORY"], "Cassandra")
	assert.Equal(t, stackifySpan.Props["SUBCATEGORY"], "Connect")
}

func TestCasandraQuerySpan(t *testing.T) {
	c := createConfig()
	attributes := []label.KeyValue{
		label.String("db.system", "cassandra"),
		label.String("db.statement", "query"),
		label.String("db.cassandra.rows.returned", "1"),
	}
	sd := createOtelSpanData(c, "test", trace.SpanKindClient, ChildSpan, attributes, InsLibTest)

	stackifySpan := span.NewSpan(c, sd)

	assert.Equal(t, stackifySpan.Call, "db.cassandra.query")
	assert.NotEmpty(t, stackifySpan.ReqBegin)
	assert.NotEmpty(t, stackifySpan.ReqEnd)
	assert.Equal(t, stackifySpan.Props["CATEGORY"], "Cassandra")
	assert.Equal(t, stackifySpan.Props["SUBCATEGORY"], "Execute")
	assert.Equal(t, stackifySpan.Props["SQL"], "query")
	assert.Equal(t, stackifySpan.Props["ROW_COUNT"], "1")
}

func TestMongoDBSpan(t *testing.T) {
	c := createConfig()
	attributes := []label.KeyValue{
		label.String("db.system", "mongodb"),
		label.String("db.instance", "collection"),
		label.String("db.operation", "query"),
		label.String("db.statement", "{\"insert\": \"collection\"}"),
		label.String("db.instance", "db"),
	}
	sd := createOtelSpanData(c, "test", trace.SpanKindClient, ChildSpan, attributes, InsLibTest)

	stackifySpan := span.NewSpan(c, sd)

	assert.Equal(t, stackifySpan.Call, "db.mongodb.query")
	assert.NotEmpty(t, stackifySpan.ReqBegin)
	assert.NotEmpty(t, stackifySpan.ReqEnd)
	assert.Equal(t, stackifySpan.Props["CATEGORY"], "MongoDB")
	assert.Equal(t, stackifySpan.Props["SUBCATEGORY"], "Execute")
	assert.Equal(t, stackifySpan.Props["COMPONENT_CATEGORY"], "DB Query")
	assert.Equal(t, stackifySpan.Props["COMPONENT_DETAIL"], "Execute SQL Query")
	assert.Equal(t, stackifySpan.Props["PROVIDER"], "mongodb")
	assert.Equal(t, stackifySpan.Props["MONGODB_COLLECTION"], "db.collection")
	assert.Equal(t, stackifySpan.Props["OPERATION"], "query")
}

func TestGRPCSpan(t *testing.T) {
	c := createConfig()
	attributes := []label.KeyValue{
		label.String("rpc.system", "grpc"),
		label.String("rpc.service", "service"),
		label.String("rpc.method", "method"),
	}
	sd := createOtelSpanData(c, "grpc", trace.SpanKindClient, ChildSpan, attributes, InsLibTest)

	stackifySpan := span.NewSpan(c, sd)

	assert.Equal(t, stackifySpan.Call, "grpc")
	assert.NotEmpty(t, stackifySpan.ReqBegin)
	assert.NotEmpty(t, stackifySpan.ReqEnd)
	assert.Equal(t, stackifySpan.Props["CATEGORY"], "RPC")
	assert.Equal(t, stackifySpan.Props["SUBCATEGORY"], "Execute")
	assert.Equal(t, stackifySpan.Props["PROVIDER"], "grpc")
	assert.Equal(t, stackifySpan.Props["SERVICE"], "service")
	assert.Equal(t, stackifySpan.Props["METHOD"], "method")
}
