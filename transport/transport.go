package transport

import (
	"github.com/stackify/stackify-go-apm/config"
	"github.com/stackify/stackify-go-apm/trace/span"
)

type Transport interface {
	HandleTrace(*span.StackifySpan)
	SendAll()
}

func NewTransport(c *config.Config) Transport {
	if c.TransportType == config.DefaultTransportType {
		return newDefaultTransport(c)
	}

	// fallback to default transport if transport type is unkown
	return newDefaultTransport(c)
}
