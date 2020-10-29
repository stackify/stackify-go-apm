package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"go.stackify.com/apm"
	"go.stackify.com/apm/config"
	"go.stackify.com/apm/instrumentation/google.golang.org/grpc/stackifygrpc"
	"go.stackify.com/apm/instrumentation/google.golang.org/grpc/stackifygrpc/example/proto"
	"google.golang.org/grpc"
)

const (
	port = ":8000"
)

type server struct {
	proto.MessageServiceServer
}

func (s *server) SendMessage(ctx context.Context, in *proto.MessageRequest) (*proto.MessageResponse, error) {
	log.Printf("Received: %v\n", in.GetMessage())
	time.Sleep(50 * time.Millisecond)

	return &proto.MessageResponse{Reply: in.Message + "World"}, nil
}

func (s *server) SendMessageClientStream(stream proto.MessageService_SendMessageClientStreamServer) error {
	i := 0

	for {
		in, err := stream.Recv()

		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("Non EOF error: %v\n", err)
			return err
		}

		log.Printf("Received: %v\n", in.GetMessage())
		i++
	}

	time.Sleep(50 * time.Millisecond)

	return stream.SendAndClose(&proto.MessageResponse{Reply: fmt.Sprintf("World (%v times)", i)})
}

func (s *server) SendMessageServerStream(in *proto.MessageRequest, out proto.MessageService_SendMessageServerStreamServer) error {
	log.Printf("Received: %v\n", in.GetMessage())

	for i := 0; i < 5; i++ {
		err := out.Send(&proto.MessageResponse{Reply: in.Message + "World"})
		if err != nil {
			return err
		}

		time.Sleep(time.Duration(i*50) * time.Millisecond)
	}

	return nil
}

func (s *server) SendMessageBidiStream(stream proto.MessageService_SendMessageBidiStreamServer) error {
	for {
		in, err := stream.Recv()

		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("Non EOF error: %v\n", err)
			return err
		}

		time.Sleep(50 * time.Millisecond)

		log.Printf("Received: %v\n", in.GetMessage())
		err = stream.Send(&proto.MessageResponse{Reply: in.Message + "World"})

		if err != nil {
			return err
		}
	}

	return nil
}

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

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(stackifygrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(stackifygrpc.StreamServerInterceptor()),
	)

	proto.RegisterMessageServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
