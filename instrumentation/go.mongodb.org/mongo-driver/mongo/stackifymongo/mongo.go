package stackifymongo

import (
	"go.mongodb.org/mongo-driver/event"

	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

func NewMonitor() *event.CommandMonitor {
	return otelmongo.NewMonitor("stackify")
}
