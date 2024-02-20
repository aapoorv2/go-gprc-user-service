package main

import (
	"context"
	"fmt"
	pb "greet/greet/proto"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type User struct {
	Username string `json:"username"` 
	Name     string `json:"name"`
	Password string `json:"password"`
	Age      int64  `json:"age"`
	Token    string `json:"token"`
}

func main() {
	connection, err := grpc.Dial("localhost:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer connection.Close()
	c := pb.NewGreeterClient(connection)
	r := gin.Default()

	r.POST("/user", func(ctx *gin.Context) {
		var user User
		if err := ctx.ShouldBindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		res, err := c.RegisterUser(ctx, &pb.RegisterUserRequest{
			Username: user.Username,
			Password: user.Password,
		})

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"token": res.Token})
	})

	r.PUT("/user", func(ctx *gin.Context) {
		var user User
		if err := ctx.ShouldBindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		res, err := c.PostDetails(ctx, &pb.UserDetailsRequest{
			Name:  user.Name,
			Age:   user.Age,
			Token: user.Token,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"success": res.Message})
	})

	r.GET("/user", func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		fmt.Println(token)
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(token, bearerPrefix) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}
		token = strings.TrimPrefix(token, bearerPrefix)
		
		res, err := c.GetDetails(ctx, &pb.FetchUserDetailsRequest{Token: token})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"username": res.Name, "age": res.Age})
	})

	r.Run(":8080")
}

func SayHello(c pb.GreeterClient) {
	var name string
	fmt.Scan(&name)
	res, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("failed to connect")
	}
	fmt.Println(res.Greeting)
}

