package stackifyhttp

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func NewTransport(base http.RoundTripper, opts ...otelhttp.Option) *otelhttp.Transport {
	return otelhttp.NewTransport(base, opts...)
}
