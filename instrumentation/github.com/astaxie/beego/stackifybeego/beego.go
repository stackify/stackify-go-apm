package stackifybeego

import (
	"github.com/astaxie/beego"
	"go.opentelemetry.io/contrib/instrumentation/github.com/astaxie/beego/otelbeego"
)

func NewStackifyBeegoMiddleWare() beego.MiddleWare {
	return otelbeego.NewOTelBeegoMiddleWare("stackifyMiddleware")
}

func Render(c *beego.Controller) error {
	return otelbeego.Render(c)
}

func RenderString(c *beego.Controller) (string, error) {
	return otelbeego.RenderString(c)
}

func RenderBytes(c *beego.Controller) ([]byte, error) {
	return otelbeego.RenderBytes(c)
}
