package client

import (
	"context"
	"crypto/tls"
	"fmt"

	"example.com/m/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

func Run() {
	creds := credentials.NewTLS(&tls.Config{})

	dial, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(creds))
	if err != nil {
		panic(err)
	}

	defer dial.Close()

	userClient := pb.NewUserClient(dial)

	login, err := userClient.Login(context.Background(), &pb.LoginRequest{
		Username: "test",
		Password: "test",
	})

	if err != nil {
		fmt.Printf("login failed: %v\n", err)
		return
	}

	md := metadata.New(map[string]string{
		"Authorization": "Bearer " + login.GetToken(),
	})

	ctx := metadata.NewOutgoingContext(context.Background(), md)

	user, err := userClient.AddUser(ctx, &pb.AddUserRequest{
		Id:   "1",
		Name: "test",
		Age:  10,
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("user created on AddUser: %v\n", user)

	getUserResponse, err := userClient.GetUser(ctx, &pb.GetUserRequest{
		Id: "1",
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("user get on GetUser: %v\n", getUserResponse)
}
