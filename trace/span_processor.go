package trace

import (
	"context"
	"runtime"
	"sync"
	"time"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
	export "go.opentelemetry.io/otel/sdk/export/trace"
	"go.stackify.com/apm/trace/span"
)

var (
	invalidTraceID trace.ID
	validSpan      = map[string]bool{
		// HTML
		"GET":    true,
		"POST":   true,
		"PUT":    true,
		"DELETE": true,
		"PATCH":  true,

		// Template
		"gin.renderer.html":     true,
		"beego.render.template": true,

		// Memcache
		"add":        true,
		"cas":        true,
		"decr":       true,
		"delete":     true,
		"delete_all": true,
		"flush_all":  true,
		"get":        true,
		"incr":       true,
		"ping":       true,
		"replace":    true,
		"set":        true,
		"touch":      true,
	}
)

const (
	DefaultTimeout = 500 * time.Millisecond
)

type StackifySpanProcessor struct {
	e                    *StackifySpanExporter
	traces               map[trace.ID][]*export.SpanData
	traces_started_count map[trace.ID]int
	traces_ended_count   map[trace.ID]int
	trace_ids_to_export  []trace.ID
	queue                chan trace.ID
	queueMutex           sync.Mutex
	timer                *time.Timer
	stopWait             sync.WaitGroup
	stopOnce             sync.Once
}

// NewStackifySpanProcessor function creates StackifySpanProcessor and runs queue worker over goroutine.
func NewStackifySpanProcessor(exporter *StackifySpanExporter) *StackifySpanProcessor {
	ssp := &StackifySpanProcessor{
		e:                    exporter,
		traces:               make(map[trace.ID][]*export.SpanData),
		traces_started_count: make(map[trace.ID]int),
		traces_ended_count:   make(map[trace.ID]int),
		trace_ids_to_export:  []trace.ID{},
		timer:                time.NewTimer(DefaultTimeout),
		queue:                make(chan trace.ID, 100),
	}

	ssp.stopWait.Add(1)
	go func() {
		ssp.processQueue()
		ssp.drainQueue()
		defer ssp.stopWait.Done()
	}()
	return ssp
}

// OnStart method counts how many started spans from a specific trace.
func (ssp *StackifySpanProcessor) OnStart(sd *export.SpanData) {
	if !ssp.isSpanValid(sd) {
		return
	}
	ssp.queueMutex.Lock()
	defer ssp.queueMutex.Unlock()

	ssp.traces_started_count[sd.SpanContext.TraceID] += 1
}

// OnEnd method appends spans into store and enqueue spans if trace is finished.
func (ssp *StackifySpanProcessor) OnEnd(sd *export.SpanData) {
	if !ssp.isSpanValid(sd) {
		return
	}
	ssp.queueMutex.Lock()
	defer ssp.queueMutex.Unlock()

	trace_id := sd.SpanContext.TraceID
	ssp.traces[sd.SpanContext.TraceID] = append(ssp.traces[sd.SpanContext.TraceID], sd)
	ssp.traces_ended_count[trace_id] += 1

	if ssp.isTraceExportable(trace_id) {
		ssp.enqueue(trace_id)
	}
}

// Shutdown method cleans up and drain the queue.
func (ssp *StackifySpanProcessor) Shutdown() {
	ssp.stopOnce.Do(func() {
		ssp.enqueue(invalidTraceID)
		ssp.stopWait.Wait()
	})
}

// ForceFlush export remaining traces from queue.
func (ssp *StackifySpanProcessor) ForceFlush() {
	ssp.exportSpans()
}

// processQueue method loops over the queue and process every entry.
func (ssp *StackifySpanProcessor) processQueue() {
	defer ssp.timer.Stop()

	for {
		select {
		case <-ssp.timer.C:
			ssp.exportSpans()
		case trace_id := <-ssp.queue:
			if trace_id == invalidTraceID {
				ssp.exportSpans()
				return
			}
			ssp.queueMutex.Lock()
			ssp.trace_ids_to_export = append(ssp.trace_ids_to_export, trace_id)
			ssp.queueMutex.Unlock()
		}
	}
}

// drainQueue method drains the queue making sure we are processing all spans.
func (ssp *StackifySpanProcessor) drainQueue() {
	for {
		select {
		case trace_id := <-ssp.queue:
			if trace_id == invalidTraceID {
				ssp.exportSpans()
				return
			}
			ssp.queueMutex.Lock()
			ssp.trace_ids_to_export = append(ssp.trace_ids_to_export, trace_id)
			ssp.queueMutex.Unlock()
		default:
			close(ssp.queue)
		}
	}
}

// exportSpans method exports all finished traces.
func (ssp *StackifySpanProcessor) exportSpans() {
	ssp.timer.Reset(DefaultTimeout)

	ssp.queueMutex.Lock()
	defer ssp.queueMutex.Unlock()

	for len(ssp.trace_ids_to_export) > 0 {
		var trace_id trace.ID
		trace_id, ssp.trace_ids_to_export = ssp.trace_ids_to_export[0], ssp.trace_ids_to_export[1:]
		trace := ssp.traces[trace_id]

		if ssp.isTraceExportable(trace_id) && len(trace) > 0 {
			if err := ssp.e.ExportSpans(context.Background(), trace); err != nil {
				global.Handle(err)
			}
			delete(ssp.traces, trace_id)
			delete(ssp.traces_started_count, trace_id)
			delete(ssp.traces_ended_count, trace_id)
		}
	}
}

// enqueue method enqueue finished traces.
func (ssp *StackifySpanProcessor) enqueue(trace_id trace.ID) {
	defer func() {
		x := recover()
		switch err := x.(type) {
		case nil:
			return
		case runtime.Error:
			if err.Error() == "send on closed channel" {
				return
			}
		}
		panic(x)
	}()

	ssp.queue <- trace_id
}

// isTraceExportable method validates in trace is finished or not.
func (ssp *StackifySpanProcessor) isTraceExportable(trace_id trace.ID) bool {
	return ssp.traces_started_count[trace_id]-ssp.traces_ended_count[trace_id] <= 0
}

// isSpanValid method checks if span is a valid stackify span.
func (ssp *StackifySpanProcessor) isSpanValid(sd *export.SpanData) bool {
	_, ok := validSpan[sd.Name]
	if !ok {
		ok = sd.InstrumentationLibrary.Name == span.Otelgocql
	}
	return ok || sd.ParentSpanID == span.InvalidSpanId
}
