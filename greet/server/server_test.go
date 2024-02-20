package main

import (
	"context"
	pb "greet/greet/proto"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestSayHello(t *testing.T) {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		t.Fatalf("Failed to listen on port 9000: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterGreeterServer(grpcServer, &Server{})

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC server over port 9000: %v", err)
		}
	}()
	defer grpcServer.Stop()

	conn, err := grpc.Dial("localhost:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial gRPC server: %v", err)
	}
	defer conn.Close()

	client := pb.NewGreeterClient(conn)

	response, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "Test"})
	if err != nil {
		t.Fatalf("SayHello failed: %v", err)
	}

	expected := "HelloTest"
	if response.Greeting != expected {
		t.Errorf("Expected greeting %s, got %s", expected, response.Greeting)
	}
}