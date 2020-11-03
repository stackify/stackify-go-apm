package stackifyecho

import (
	"github.com/labstack/echo/v4"

	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

func Middleware() echo.MiddlewareFunc {
	return otelecho.Middleware("StackifyMiddleware")
}
