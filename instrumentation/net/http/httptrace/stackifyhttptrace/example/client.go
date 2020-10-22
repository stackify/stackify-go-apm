package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptrace"

	"go.stackify.com/apm"
	"go.stackify.com/apm/config"
	"go.stackify.com/apm/instrumentation/net/http/httptrace/stackifyhttptrace"
	"go.stackify.com/apm/instrumentation/net/http/stackifyhttp"
)

func initStackifyTrace() (*apm.StackifyAPM, error) {
	return apm.NewStackifyAPM(
		config.WithApplicationName("Jayr GOLANG 11:22"),
		config.WithEnvironmentName("Test"),
		config.WithDebug(true),
	)
}

func main() {
	fmt.Println("Starting simple application.")

	url := flag.String("server", "https://golang.org/", "server url")
	flag.Parse()

	client := http.Client{Transport: stackifyhttp.NewTransport(http.DefaultTransport)}

	var body []byte

	stackifyAPM, err := initStackifyTrace()
	if err != nil {
		log.Fatalf("failed to initialize stackifyapm: %v", err)
	}
	defer stackifyAPM.Shutdown()
	tr := stackifyAPM.Tracer
	ctx := stackifyAPM.Context

	err = func(ctx context.Context) error {
		ctx, span := tr.Start(ctx, "custom")
		defer span.End()

		ctx = httptrace.WithClientTrace(ctx, stackifyhttptrace.NewClientTrace(ctx))
		req, _ := http.NewRequestWithContext(ctx, "GET", *url, nil)

		fmt.Println("Sending request...")
		res, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		body, err = ioutil.ReadAll(res.Body)
		_ = res.Body.Close()

		return err
	}(ctx)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Response Received: %s\n\n\n", body)
	fmt.Println("Application done.")
}
