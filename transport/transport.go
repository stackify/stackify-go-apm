package transport

import (
	"go.stackifyapm.com/apm/config"
	"go.stackifyapm.com/apm/trace/span"
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
