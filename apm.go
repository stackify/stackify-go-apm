package apm

import (
	"context"
	"time"

	"github.com/stackify/stackify-go-apm/config"
	"github.com/stackify/stackify-go-apm/trace"
	"github.com/stackify/stackify-go-apm/transport"

	"go.opentelemetry.io/otel/api/global"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type StackifyAPM struct {
	config        *config.Config
	transport     *transport.Transport
	spanExporter  *trace.StackifySpanExporter
	spanProcessor *trace.StackifySpanProcessor
	TraceProvider *sdktrace.TracerProvider
	Tracer        trace.Tracer
	Context       context.Context
}

func (sapm *StackifyAPM) Shutdown() {
	time.Sleep(1 * time.Second)
	sapm.spanProcessor.Shutdown()
}

func NewStackifyAPM(opts ...config.ConfigOptions) (*StackifyAPM, error) {
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
	tracer := global.Tracer(config.StackifyInstrumentationName)

	// create context
	ctx := context.Background()

	// initialize StackifyAPM
	stackifyAPM := StackifyAPM{
		config:        c,
		transport:     &t,
		spanExporter:  sse,
		spanProcessor: ssp,
		TraceProvider: tp,
		Tracer:        tracer,
		Context:       ctx,
	}

	return &stackifyAPM, nil
}
