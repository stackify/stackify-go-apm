package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"go.stackify.com/apm"
	"go.stackify.com/apm/config"
	"go.stackify.com/apm/instrumentation/net/http/stackifyhttp"
)

func initStackifyTrace() (*apm.StackifyAPM, error) {
	return apm.NewStackifyAPM(
		config.WithApplicationName("Go Application"),
		config.WithEnvironmentName("Test"),
		config.WithDebug(true),
	)
}

func indexHandler(w http.ResponseWriter, req *http.Request) {
	_, _ = io.WriteString(w, "Hello, world!\n")
}

func main() {
	stackifyAPM, err := initStackifyTrace()
	if err != nil {
		log.Fatalf("failed to initialize stackifyapm: %v", err)
	}
	defer stackifyAPM.Shutdown()

	stackifyHandler := stackifyhttp.NewHandler(http.HandlerFunc(indexHandler), "Index")
	http.Handle("/index", stackifyHandler)

	fmt.Println("Starting server.")
	err = http.ListenAndServe(":8000", nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Server closed.")
}
