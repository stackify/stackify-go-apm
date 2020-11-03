package stackifygocql

import (
	"context"

	"github.com/gocql/gocql"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gocql/gocql/otelgocql"
)

func NewSessionWithTracing(ctx context.Context, cluster *gocql.ClusterConfig, options ...otelgocql.TracedSessionOption) (*gocql.Session, error) {
	return otelgocql.NewSessionWithTracing(ctx, cluster, options...)
}
