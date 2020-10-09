package apm

import (
	"context"
	"fmt"
	"time"

	"go.stackify.com/apm/config"
	"go.stackify.com/apm/trace"
	"go.stackify.com/apm/transport"

	"go.opentelemetry.io/otel/api/baggage"
	"go.opentelemetry.io/otel/api/global"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type StackifyAPM struct {
	config        *config.Config
	transport     *transport.Transport
	spanExporter  *trace.StackifySpanExporter
	spanProcessor *trace.StackifySpanProcessor
	traceProvider *sdktrace.TracerProvider
	Tracer        trace.Tracer
	Context       context.Context
}

func (sapm *StackifyAPM) Shutdown() {
	time.Sleep(1 * time.Second)
	sapm.spanProcessor.Shutdown()
}

func NewStackifyAPM(opts ...config.ConfigOptions) (*StackifyAPM, error) {
	fmt.Println("APM Starting...")

	// initialize Config
	c := config.NewConfig(opts...)

	// initialize Transport
	t := transport.NewTransport(c)

	// initialize stackify span exporter
	sse := trace.NewStackifySpanExporter(c, &t)

	// initialize stackify span processor
	ssp := trace.NewStackifySpanProcessor(sse)

	// initialize OT tracer provider
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(ssp))

	// set tracer provider as global tracer provider
	global.SetTracerProvider(tp)
	tracer := global.Tracer("stackifyapm_tracer")

	// create context
	ctx := context.Background()
	ctx = baggage.NewContext(ctx)

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
