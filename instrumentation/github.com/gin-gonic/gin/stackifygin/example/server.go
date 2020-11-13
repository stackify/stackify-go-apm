package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	apm "github.com/stackify/stackify-go-apm"
	"github.com/stackify/stackify-go-apm/config"
	"github.com/stackify/stackify-go-apm/instrumentation/github.com/gin-gonic/gin/stackifygin"
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

	r := gin.New()
	r.Use(stackifygin.Middleware())

	tmplName := "index"
	tmplStr := "Hello, {{ .name }}!\n"
	tmpl := template.Must(template.New(tmplName).Parse(tmplStr))
	r.SetHTMLTemplate(tmpl)

	r.GET("/index", func(c *gin.Context) {
		stackifygin.HTML(c, http.StatusOK, tmplName, gin.H{
			"name": "World",
		})
	})
	_ = r.Run(":8000")
}
