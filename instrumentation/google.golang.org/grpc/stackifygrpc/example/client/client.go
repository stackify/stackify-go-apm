package main

import (
	"context"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"go.stackify.com/apm"
	"go.stackify.com/apm/config"
	"go.stackify.com/apm/instrumentation/google.golang.org/grpc/stackifygrpc"
	"go.stackify.com/apm/instrumentation/google.golang.org/grpc/stackifygrpc/example/proto"
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

	var conn *grpc.ClientConn
	conn, err = grpc.Dial(":8000", grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(stackifygrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(stackifygrpc.StreamClientInterceptor()),
	)
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer func() { _ = conn.Close() }()

	client := proto.NewMessageServiceClient(conn)

	ctx, span := stackifyAPM.Tracer.Start(stackifyAPM.Context, "custom")
	SendMessage(client, ctx)
	SendMessageClientStream(client, ctx)
	SendMessageServerStream(client, ctx)
	SendMessageBidiStream(client, ctx)
	span.End()
}

func SendMessage(c proto.MessageServiceClient, ctx context.Context) {
	md := metadata.Pairs(
		"timestamp", time.Now().Format(time.StampNano),
	)

	ctx = metadata.NewOutgoingContext(ctx, md)
	response, err := c.SendMessage(ctx, &proto.MessageRequest{Message: "Hello"})
	if err != nil {
		log.Fatalf("Error when calling SendMessage: %s", err)
	}
	log.Printf("Response from server: %s", response.Reply)
}

func SendMessageClientStream(c proto.MessageServiceClient, ctx context.Context) {
	md := metadata.Pairs(
		"timestamp", time.Now().Format(time.StampNano),
	)

	ctx = metadata.NewOutgoingContext(ctx, md)
	stream, err := c.SendMessageClientStream(ctx)
	if err != nil {
		log.Fatalf("Error when opening SendMessageClientStream: %s", err)
	}

	for i := 0; i < 5; i++ {
		err := stream.Send(&proto.MessageRequest{Message: "Hello"})

		time.Sleep(time.Duration(i*50) * time.Millisecond)

		if err != nil {
			log.Fatalf("Error when sending to SendMessageClientStream: %s", err)
		}
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error when closing SendMessageClientStream: %s", err)
	}

	log.Printf("Response from server: %s", response.Reply)
}

func SendMessageServerStream(c proto.MessageServiceClient, ctx context.Context) {
	md := metadata.Pairs(
		"timestamp", time.Now().Format(time.StampNano),
	)

	ctx = metadata.NewOutgoingContext(ctx, md)
	stream, err := c.SendMessageServerStream(ctx, &proto.MessageRequest{Message: "Hello"})
	if err != nil {
		log.Fatalf("Error when opening SendMessageServerStream: %s", err)
	}

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("Error when receiving from SendMessageServerStream: %s", err)
		}

		log.Printf("Response from server: %s", response.Reply)
		time.Sleep(50 * time.Millisecond)
	}
}

func SendMessageBidiStream(c proto.MessageServiceClient, ctx context.Context) {
	md := metadata.Pairs(
		"timestamp", time.Now().Format(time.StampNano),
	)

	ctx = metadata.NewOutgoingContext(ctx, md)
	stream, err := c.SendMessageBidiStream(ctx)
	if err != nil {
		log.Fatalf("Error when opening SendMessageBidiStream: %s", err)
	}

	serverClosed := make(chan struct{})
	clientClosed := make(chan struct{})

	go func() {
		for i := 0; i < 5; i++ {
			err := stream.Send(&proto.MessageRequest{Message: "Hello"})

			if err != nil {
				log.Fatalf("Error when sending to SendMessageBidiStream: %s", err)
			}

			time.Sleep(50 * time.Millisecond)
		}

		err := stream.CloseSend()
		if err != nil {
			log.Fatalf("Error when closing SendMessageBidiStream: %s", err)
		}

		clientClosed <- struct{}{}
	}()

	go func() {
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatalf("Error when receiving from SendMessageBidiStream: %s", err)
			}

			log.Printf("Response from server: %s", response.Reply)
			time.Sleep(50 * time.Millisecond)
		}

		serverClosed <- struct{}{}
	}()

	// Wait until client and server both closed the connection.
	<-clientClosed
	<-serverClosed
}
