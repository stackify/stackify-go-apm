package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.stackify.com/apm"
	"go.stackify.com/apm/config"
	"go.stackify.com/apm/instrumentation/github.com/labstack/echo/stackifyecho"
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

	r := echo.New()
	r.Use(stackifyecho.Middleware())
	r.GET("/index", func(c echo.Context) error {
		reply := fmt.Sprintln("Hello World!")
		return c.JSON(http.StatusOK, reply)
	})
	_ = r.Start(":8000")
}
