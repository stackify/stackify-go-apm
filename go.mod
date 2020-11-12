module github.com/stackify/stackify-go-apm

go 1.15

require (
	github.com/astaxie/beego v1.12.2
	github.com/bradfitz/gomemcache v0.0.0-20190913173617-a41fca850d0b
	github.com/emicklei/go-restful/v3 v3.3.3
	github.com/gin-gonic/gin v1.6.3
	github.com/gocql/gocql v0.0.0-20201024154641-5913df4d474e
	github.com/golang/protobuf v1.4.3
	github.com/google/uuid v1.1.2
	github.com/gorilla/mux v1.8.0
	github.com/labstack/echo/v4 v4.1.17
	github.com/stretchr/testify v1.6.1
	go.mongodb.org/mongo-driver v1.4.2
	go.opentelemetry.io/contrib/instrumentation/github.com/astaxie/beego/otelbeego v0.13.0
	go.opentelemetry.io/contrib/instrumentation/github.com/bradfitz/gomemcache/memcache/otelmemcache v0.13.0
	go.opentelemetry.io/contrib/instrumentation/github.com/emicklei/go-restful/otelrestful v0.13.0
	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.13.0
	go.opentelemetry.io/contrib/instrumentation/github.com/gocql/gocql/otelgocql v0.13.0
	go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux v0.13.0
	go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho v0.13.0
	go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo v0.13.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.13.0
	go.opentelemetry.io/contrib/instrumentation/gopkg.in/macaron.v1/otelmacaron v0.13.0
	go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace v0.13.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.13.0
	go.opentelemetry.io/otel v0.13.0
	go.opentelemetry.io/otel/sdk v0.13.0
	google.golang.org/grpc v1.33.1
	google.golang.org/protobuf v1.25.0
	gopkg.in/macaron.v1 v1.3.9
)
