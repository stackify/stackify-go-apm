package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.stackify.com/apm"
	"go.stackify.com/apm/config"
	"go.stackify.com/apm/instrumentation/go.mongodb.org/mongo-driver/mongo/stackifymongo"
)

func initStackifyTrace() (*apm.StackifyAPM, error) {
	return apm.NewStackifyAPM(
		config.WithApplicationName("Go Application"),
		config.WithEnvironmentName("Test"),
		config.WithDebug(true),
	)
}

func main() {
	stackifyAPM, err := initStackifyTrace()
	if err != nil {
		log.Fatalf("failed to initialize stackifyapm: %v", err)
	}
	defer stackifyAPM.Shutdown()

	opts := options.Client()
	opts.Monitor = stackifymongo.NewMonitor()
	opts.ApplyURI("mongodb://localhost:1113")

	client, err := mongo.Connect(stackifyAPM.Context, opts)
	if err != nil {
		panic(err)
	}

	ctx, span := stackifyAPM.Tracer.Start(stackifyAPM.Context, "custom")
	doMongoOperations(ctx, client)
	span.End()
}

func doMongoOperations(ctx context.Context, client *mongo.Client) {
	db := client.Database("exampleDB")
	inventory := db.Collection("exampleCollection")

	_, err := inventory.InsertOne(ctx, bson.D{
		{Key: "foo", Value: "bar"},
	})
	if err != nil {
		panic(err)
	}
}
