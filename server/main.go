package main

import (
	"context"
	"io"
	"log"
	"net"
	"time"

	pb "example/proto"

	"google.golang.org/grpc"
)

const (
	port = ":8080"
)

type helloServer struct {
	pb.GreetServiceServer
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to start server %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGreetServiceServer(grpcServer, &helloServer{})

	log.Printf("server started at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to start: %v", err)
	}
}

func (s *helloServer) SayHello(ctx context.Context, req *pb.NoParam) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{
		Message: "Hello",
	}, nil
}

func (s *helloServer) SayHelloServerStreaming(req *pb.NamesList, stream pb.GreetService_SayHelloServerStreamingServer) error {
	log.Printf("Got request with names : %v", req.Names)

	for _, name := range req.Names {
		res := &pb.HelloResponse{
			Message: "Hello " + name,
		}
		if err := stream.Send(res); err != nil {
			return err
		}
		// 2 second delay to simulate a long running process
		time.Sleep(2 * time.Second)
	}
	return nil
}

// client_stream
func (s *helloServer) SayHelloClientStreaming(stream pb.GreetService_SayHelloClientStreamingServer) error {
	var messages []string
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.MessagesList{Messages: messages})
		}
		if err != nil {
			return err
		}
		log.Printf("Got request with name : %v", req.Name)
		messages = append(messages, "Hello "+req.Name)
	}
}

func (s *helloServer) SayHelloBidirectonalStreaming(stream pb.GreetService_SayHelloBidirectonalStreamingServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		log.Printf("Got request with name : %v", req.Name)
		res := &pb.HelloResponse{
			Message: "Hello " + req.Name,
		}
		if err := stream.Send(res); err != nil {
			return err
		}
	}
}
