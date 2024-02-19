package main

import (
	"log"
	"net"
	pb"greet/greet/proto"
	"google.golang.org/grpc"
	"context"
)
type Server struct {
	pb.GreeterServer
}
type User struct {
	userid int64
	name string
	address string
}
var Users []User
func main() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to listen on port 9000: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterGreeterServer(grpcServer, &Server{})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port 9000: %v", err)
	}
}
func (s *Server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Greeting:"Hello" + req.Name}, nil
}
func (s *Server) CreateUser(c context.Context, req *pb.Createrequest) (*pb.UserResponse, error) {
	usr := User{userid: req.Userid, name: req.Name, address: req.Address}
	Users = append(Users, usr)
	return &pb.UserResponse{Userid: req.Userid, Name : req.Name, Address: req.Address, Message: "success"}, nil
}
