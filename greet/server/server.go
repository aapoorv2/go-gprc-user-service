package main

import (
	"context"
	"database/sql"
	pb "greet/greet/proto"
	"log"
	"net"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


type Server struct {
	pb.GreeterServer
	db *sql.DB
}

func main() {
	db, err := sql.Open("postgres", "postgresql://username:password@localhost:5432/hello?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer db.Close()

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to listen on port 9000: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterGreeterServer(grpcServer, &Server{db : db})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port 9000: %v", err)
	}
	log.Println("Server started running on port 9000")
}

func (s *Server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Greeting:"Hello" + req.Name}, nil
}

func (s *Server) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	var existingUsername string
	err := s.db.QueryRow("SELECT username FROM users WHERE username = $1", req.Username).Scan(&existingUsername)
	if err == nil {
		return nil, status.Error(codes.AlreadyExists, "User with the given username already exists")
	}
	token := uuid.New().String()

	_, err = s.db.Exec("INSERT INTO users (username, password, token) VALUES ($1, $2, $3)", req.Username, req.Password, token)
	if err != nil {
		log.Printf("Failed to insert user into PostgreSQL: %v", err)
		return nil, status.Error(codes.Internal, "Failed to register user")
	}
	log.Printf("User registered successfully: %s", req.Username)
	return &pb.RegisterUserResponse{Token: token}, nil
}

func (s *Server) PostDetails(c context.Context, req *pb.UserDetailsRequest) (*pb.UserDetailsResponse, error) {
	var existingToken string
	token := req.Token
	err := s.db.QueryRow("SELECT token FROM users WHERE token = $1", req.Token).Scan(&existingToken)
	if err == sql.ErrNoRows {
		return nil, status.Error(codes.Unauthenticated, "Invalid Token")
	}
	age := req.Age
	name := req.Name
	_, err = s.db.Exec("UPDATE users SET name = $1, age = $2 WHERE token = $3", name, age, token)

	if err != nil {
		log.Printf("Failed to update user details in PostgreSQL: %v", err)
		return nil, status.Error(codes.Internal, "Failed to update user details")
	}

	log.Printf("User details updated successfully")
	return &pb.UserDetailsResponse{Message: "Successfully updated the user information"}, nil
}
func (s *Server) GetDetails(c context.Context, req *pb.FetchUserDetailsRequest) (*pb.FetchUserDetailsResponse, error) {
	var existingToken string
	err := s.db.QueryRow("SELECT token FROM users WHERE token = $1", req.Token).Scan(&existingToken)
	if err == sql.ErrNoRows {
		return nil, status.Error(codes.Unauthenticated, "Invalid Token")
	}
	var name string
	var age int64
	err = s.db.QueryRow("SELECT name, age FROM users WHERE token = $1", req.Token).Scan(&name, &age)
	if err != nil {
		log.Printf("Failed to fetch user details from PostgreSQL: %v", err)
		return nil, status.Error(codes.NotFound, "User not found")
	}

	log.Printf("User details fetched successfully")
	return &pb.FetchUserDetailsResponse{Name: name, Age: age}, nil
}
