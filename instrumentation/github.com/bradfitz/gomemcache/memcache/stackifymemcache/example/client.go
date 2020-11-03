package main

import (
	"context"
	"log"

	"github.com/bradfitz/gomemcache/memcache"

	"go.stackify.com/apm"
	"go.stackify.com/apm/config"
	"go.stackify.com/apm/instrumentation/github.com/bradfitz/gomemcache/memcache/stackifymemcache"
)

func initStackifyTrace() (*apm.StackifyAPM, error) {
	return apm.NewStackifyAPM(
		config.WithApplicationName("Go Application"),
		config.WithEnvironmentName("Test"),
		config.WithDebug(true),
	)
}

func main() {
	stackifyAPM, err := initStackifyTrace()
	if err != nil {
		log.Fatalf("failed to initialize stackifyapm: %v", err)
	}
	defer stackifyAPM.Shutdown()

	tp := stackifyAPM.TraceProvider
	tracer := stackifyAPM.Tracer
	ctx := stackifyAPM.Context

	var host, port = "127.0.0.1", "1117"
	c := stackifymemcache.NewClientWithTracing(
		memcache.New(
			host+":"+port,
		),
		stackifymemcache.WithTracerProvider(tp),
	)

	ctx, span := tracer.Start(ctx, "custom")
	doMemcacheOperations(ctx, c)
	span.End()
}

func doMemcacheOperations(ctx context.Context, c *stackifymemcache.Client) {
	cc := c.WithContext(ctx)

	// Add
	err := cc.Add(&memcache.Item{
		Key:   "foo",
		Value: []byte("bar"),
	})
	if err != nil {
		log.Printf("Add failed: %s", err)
	} else {
		log.Printf("Add successful: %s", "foo")
	}

	item, err := cc.Get("foo")
	if err != nil {
		log.Printf("Get failed: %s", err)
	} else {
		log.Printf("Get successful: %s", item.Value)
	}

	err = cc.Delete("baz")
	if err != nil {
		log.Printf("Delete failed: %s", err)
	}

	err = cc.DeleteAll()
	if err != nil {
		log.Printf("DeleteAll failed: %s", err)
	}

}
