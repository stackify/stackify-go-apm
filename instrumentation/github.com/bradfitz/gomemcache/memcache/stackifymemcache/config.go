package stackifymemcache

import (
	"go.opentelemetry.io/contrib/instrumentation/github.com/bradfitz/gomemcache/memcache/otelmemcache"
	oteltrace "go.opentelemetry.io/otel/api/trace"
)

func WithTracerProvider(provider oteltrace.TracerProvider) otelmemcache.Option {
	return otelmemcache.WithTracerProvider(provider)
}
