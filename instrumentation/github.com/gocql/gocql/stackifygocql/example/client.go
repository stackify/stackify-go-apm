package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gocql/gocql"
	"go.stackify.com/apm"
	"go.stackify.com/apm/config"
	"go.stackify.com/apm/instrumentation/github.com/gocql/gocql/stackifygocql"
)

const keyspace = "sample"

var wg sync.WaitGroup

func initStackifyTrace() (*apm.StackifyAPM, error) {
	return apm.NewStackifyAPM(
		config.WithApplicationName("Go Application"),
		config.WithEnvironmentName("Test"),
		config.WithDebug(true),
	)
}

func initDB() {
	cluster := getCluster("system")

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}

	statement := fmt.Sprintf(
		"create keyspace if not exists %s with replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 }",
		keyspace,
	)
	if err := session.Query(statement).Exec(); err != nil {
		log.Fatal(err)
	}

	cluster = getCluster(keyspace)
	session, err = cluster.CreateSession()

	statement = "create table if not exists book(id UUID, title text, author_first_name text, author_last_name text, PRIMARY KEY(id))"
	if err = session.Query(statement).Exec(); err != nil {
		log.Fatal(err)
	}

	if err := session.Query("create index if not exists on book(author_last_name)").Exec(); err != nil {
		log.Fatal(err)
	}
}

func getCluster(ks string) *gocql.ClusterConfig {
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Port = 1111
	cluster.Keyspace = ks
	cluster.Consistency = gocql.LocalQuorum
	cluster.ProtoVersion = 3
	cluster.Timeout = 2 * time.Second
	return cluster
}

func main() {
	stackifyAPM, err := initStackifyTrace()
	if err != nil {
		log.Fatalf("failed to initialize stackifyapm: %v", err)
	}
	defer stackifyAPM.Shutdown()

	initDB()

	tracer := stackifyAPM.Tracer
	ctx := stackifyAPM.Context

	ctx, span := tracer.Start(ctx, "custom")
	defer span.End()

	cluster := getCluster(keyspace)
	session, err := stackifygocql.NewSessionWithTracing(
		ctx,
		cluster,
	)
	if err != nil {
		log.Fatalf("failed to create a session, %v", err)
	}
	defer session.Close()

	// batch
	batch := session.NewBatch(gocql.LoggedBatch)
	for i := 0; i < 5; i++ {
		batch.Query(
			"INSERT INTO book (id, title, author_first_name, author_last_name) VALUES (?, ?, ?, ?)",
			gocql.TimeUUID(),
			fmt.Sprintf("Example Book %d", i),
			"firstName",
			"lastName",
		)
	}
	if err := session.ExecuteBatch(batch.WithContext(ctx)); err != nil {
		log.Printf("failed to batch insert, %v", err)
	}

	res := session.Query(
		"SELECT title, author_first_name, author_last_name from book WHERE author_last_name = ?",
		"lastName",
	).WithContext(ctx).PageSize(100).Iter()

	var (
		title     string
		firstName string
		lastName  string
	)

	for res.Scan(&title, &firstName, &lastName) {
		res.Scan(&title, &firstName, &lastName)
	}
	res.Close()

	if err = session.Query("truncate table book").WithContext(ctx).Exec(); err != nil {
		log.Printf("failed to delete data, %v", err)
	}
}
