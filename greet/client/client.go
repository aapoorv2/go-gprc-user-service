package main

import (
	"context"
	"fmt"
	pb "greet/greet/proto"
	"log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	connection, err := grpc.Dial("localhost:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect %v", err)
	}
	c := pb.NewGreeterClient(connection)
	// SayHello(c)
	// RegisterUser(c)
	// PostDetails(c)
	GetDetails(c)
}

func SayHello(c pb.GreeterClient){
	var name string
	fmt.Scan(&name)
	res, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("failed to connect")
	}
	fmt.Println(res.Greeting)
}
func RegisterUser(c pb.GreeterClient) {
	res, err := c.RegisterUser(context.Background(), &pb.RegisterUserRequest{Username: "abhi", Password: "pass"})
	if err != nil {
		log.Fatalf("failed to connect")
	}
	fmt.Println(res.Token)
}
func PostDetails(c pb.GreeterClient) {
	res, err := c.PostDetails(context.Background(), &pb.UserDetailsRequest{Name: "abhyudaya", Age: 22, Token: "4e9188bf-3d66-497b-8d20-dd22689b1e53"})
	if err != nil {
		log.Fatalf("failed to connect")
	}
	fmt.Println(res.Message)
}

func GetDetails(c pb.GreeterClient) {
	res, err := c.GetDetails(context.Background(), &pb.FetchUserDetailsRequest{Token: "4e9188bf-3d66-497b-8d20-dd22689b1e53"})
	if err != nil {
		log.Fatalf("failed to connect")
	}
	fmt.Println(res)
}
