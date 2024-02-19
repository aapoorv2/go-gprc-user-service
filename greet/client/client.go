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
	CreateUser(c)
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

func CreateUser(c pb.GreeterClient) {
	res, err := c.CreateUser(context.Background(), &pb.Createrequest{Userid : 2, Name: "abhi2", Address: "mumbai"})
	if err != nil {
		log.Fatalf("failed to connect")
	}
	fmt.Println(res.Userid, res.Name, res.Address)

}
