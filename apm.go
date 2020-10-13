package apm // import go.stackify.com/apm

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/api/global"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type StackifyAPM struct {
	config        *Config
	transport     *Transport
	spanExporter  *StackifySpanExporter
	spanProcessor *StackifySpanProcessor
	traceProvider *sdktrace.TracerProvider
	Tracer        Tracer
	Context       context.Context
}

func (sapm *StackifyAPM) Shutdown() {
	time.Sleep(1 * time.Second)
	sapm.spanProcessor.Shutdown()
}

func NewStackifyAPM(opts ...ConfigOptions) (*StackifyAPM, error) {
	fmt.Println("APM Starting...")

	// initialize Config
	c := NewConfig(opts...)

	// initialize Transport
	t := NewTransport(c)

	// initialize stackify span exporter
	sse := NewStackifySpanExporter(c, &t)

	// initialize stackify span processor
	ssp := NewStackifySpanProcessor(sse)

	// initialize OT tracer provider
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(ssp))

	// set tracer provider as global tracer provider
	global.SetTracerProvider(tp)
	tracer := global.Tracer("stackifyapm_tracer")

	// create context
	ctx := context.Background()

	// initialize StackifyAPM
	stackifyAPM := StackifyAPM{
		config:        c,
		transport:     &t,
		spanExporter:  sse,
		spanProcessor: ssp,
		traceProvider: tp,
		Tracer:        tracer,
		Context:       ctx,
	}

	fmt.Println("APM Started...")
	return &stackifyAPM, nil
}
