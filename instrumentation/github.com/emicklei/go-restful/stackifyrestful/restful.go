package stackifyrestful

import (
	"github.com/emicklei/go-restful/v3"

	"go.opentelemetry.io/contrib/instrumentation/github.com/emicklei/go-restful/otelrestful"
)

func StackifyFilter() restful.FilterFunction {
	return otelrestful.OTelFilter("Stackify")
}
