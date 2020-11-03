package main

import (
	"log"

	"github.com/astaxie/beego"

	apm "github.com/stackify/stackify-go-apm"
	"github.com/stackify/stackify-go-apm/config"
	"github.com/stackify/stackify-go-apm/instrumentation/github.com/astaxie/beego/stackifybeego"
)

func initStackifyTrace() (*apm.StackifyAPM, error) {
	return apm.NewStackifyAPM(
		config.WithApplicationName("Go Application"),
		config.WithEnvironmentName("Test"),
		config.WithDebug(true),
	)
}

type ExampleController struct {
	beego.Controller
}

func (c *ExampleController) Template() {
	c.TplName = "index.tpl"
	// Render the template file with tracing enabled
	if err := stackifybeego.Render(&c.Controller); err != nil {
		c.Abort("500")
	}
}

func main() {
	stackifyAPM, err := initStackifyTrace()
	if err != nil {
		log.Fatalf("failed to initialize stackifyapm: %v", err)
	}
	defer stackifyAPM.Shutdown()

	// To enable tracing on template rendering, disable autorender
	beego.BConfig.WebConfig.AutoRender = false
	beego.BConfig.WebConfig.ViewsPath = "instrumentation/github.com/astaxie/beego/stackifybeego/example/views"

	beego.Router("/index", &ExampleController{}, "get:Template")

	mware := stackifybeego.NewStackifyBeegoMiddleWare()
	beego.RunWithMiddleWares(":8000", mware)
}
