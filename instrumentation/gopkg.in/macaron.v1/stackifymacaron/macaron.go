package stackifymacaron

import (
	"gopkg.in/macaron.v1"

	"go.opentelemetry.io/contrib/instrumentation/gopkg.in/macaron.v1/otelmacaron"
)

func Middleware() macaron.Handler {
	return otelmacaron.Middleware("StackifyMiddleware")
}
