package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"go.stackify.com/apm"
	"go.stackify.com/apm/config"
	"go.stackify.com/apm/instrumentation/github.com/gorilla/mux/stackifymux"
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

	r := mux.NewRouter()
	r.Use(stackifymux.Middleware())
	r.HandleFunc("/index", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reply := fmt.Sprintln("Hello World!")
		_, _ = w.Write(([]byte)(reply))
	}))
	http.Handle("/", r)
	_ = http.ListenAndServe(":8000", nil)
}
