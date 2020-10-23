package stackifygin

import (
	"github.com/gin-gonic/gin"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func Middleware() gin.HandlerFunc {
	return otelgin.Middleware("stackifyMiddleware")
}

func HTML(c *gin.Context, code int, name string, obj interface{}) {
	otelgin.HTML(c, code, name, obj)
}
