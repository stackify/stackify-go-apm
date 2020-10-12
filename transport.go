package apm

type Transport interface {
	HandleTrace(*StackifySpan)
	SendAll()
}

func NewTransport(c *Config) Transport {
	if c.TransportType == DefaultTransportType {
		return newDefaultTransport(c)
	}

	// fallback to default transport if transport type is unkown
	return newDefaultTransport(c)
}
