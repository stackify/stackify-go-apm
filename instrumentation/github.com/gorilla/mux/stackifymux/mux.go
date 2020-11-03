package stackifymux

import (
	"github.com/gorilla/mux"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

func Middleware() mux.MiddlewareFunc {
	return otelmux.Middleware("StackifyMiddleware")
}
