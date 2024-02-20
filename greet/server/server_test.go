package main

import (
	"context"
	"fmt"
	pb "greet/greet/proto"
	"log"
	"net"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
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

func TestUserService_RegisterUser_Failure(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	server := Server{db: db}

	// Expecting an error during user registration
	mock.ExpectExec("INSERT INTO users").
		WithArgs("existing_user", "password", sqlmock.AnyArg()).
		WillReturnError(fmt.Errorf("username already exists"))

	res, err := server.RegisterUser(context.Background(), &pb.RegisterUserRequest{
		Username: "existing_user",
		Password: "password",
	})

	assert.Error(t, err)
	assert.Nil(t, res)

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

func TestUserService_UpdatingUserDetails_Failure(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	server := Server{db: db}

	mock.ExpectExec("UPDATE users").
		WithArgs("name", 20, "invalid_token").
		WillReturnError(fmt.Errorf("Invalid Token"))

	res, err := server.PostDetails(context.Background(), &pb.UserDetailsRequest{
		Name:  "name",
		Age:   20,
		Token: "invalid_token",
	})

	assert.Error(t, err)
	assert.Nil(t, res)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserService_FetchingUserDetails_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	server := Server{db: db}

	mock.ExpectQuery("SELECT name, age FROM users").
		WithArgs("valid_token").
		WillReturnRows(sqlmock.NewRows([]string{"name", "age"}).AddRow("Test Name", 20))

	res, err := server.GetDetails(context.Background(), &pb.FetchUserDetailsRequest{
		Token: "valid_token",
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "Test Name", res.Name)
	assert.Equal(t, int64(20), res.Age)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserService_FetchingUserDetails_Failure(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	server := Server{db: db}

	mock.ExpectQuery("SELECT name, age FROM users").
		WithArgs("invalid_token").
		WillReturnError(fmt.Errorf("Invalid Token"))

	res, err := server.GetDetails(context.Background(), &pb.FetchUserDetailsRequest{
		Token: "invalid_token",
	})

	assert.Error(t, err)
	assert.Nil(t, res)

	assert.NoError(t, mock.ExpectationsWereMet())
}
