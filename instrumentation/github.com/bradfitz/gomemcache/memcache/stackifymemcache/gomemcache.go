package stackifymemcache

import (
	"github.com/bradfitz/gomemcache/memcache"

	"go.opentelemetry.io/contrib/instrumentation/github.com/bradfitz/gomemcache/memcache/otelmemcache"
)

type Client struct {
	*otelmemcache.Client
}

func NewClientWithTracing(client *memcache.Client, opts ...otelmemcache.Option) *Client {
	otelclient := otelmemcache.NewClientWithTracing(client, opts...)
	return &Client{
		otelclient,
	}
}
