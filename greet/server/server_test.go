package main

import (
	"context"
	pb "greet/greet/proto"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
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

func TestUserService_RegisterUser_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	server := Server{db: db}

	mock.ExpectExec("INSERT INTO users").
		WithArgs("user", "password", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	res, err := server.RegisterUser(context.Background(), &pb.RegisterUserRequest{
		Username: "user",
		Password: "password",
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.Token)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserService_UpdatingUserDetails_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	server := Server{db: db}

	mock.ExpectExec("UPDATE users").
		WithArgs("name", 20, "token").
		WillReturnResult(sqlmock.NewResult(0, 1))

	res, err := server.PostDetails(context.Background(), &pb.UserDetailsRequest{
		Name:  "name",
		Age:   20,
		Token: "token",
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.Message)

	assert.NoError(t, mock.ExpectationsWereMet())
}