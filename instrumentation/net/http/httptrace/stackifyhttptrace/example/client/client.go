package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptrace"

	apm "github.com/stackify/stackify-go-apm"
	"github.com/stackify/stackify-go-apm/config"
	"github.com/stackify/stackify-go-apm/instrumentation/net/http/httptrace/stackifyhttptrace"
	"github.com/stackify/stackify-go-apm/instrumentation/net/http/stackifyhttp"
)

func initStackifyTrace() (*apm.StackifyAPM, error) {
	return apm.NewStackifyAPM(
		config.WithApplicationName("Go Application"),
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
