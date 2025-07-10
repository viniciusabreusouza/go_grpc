package main

import (
	"fmt"
	"net"
	"sync"

	"github.com/viniciusabreusouza/go_grpc/chat"
	"google.golang.org/grpc"
)

type chatServer struct {
	chat.UnimplementedChatServiceServer
	mu       sync.Mutex
	clients  map[chat.ChatService_JoinServer]bool
	messages chan *chat.Message
}

func newServer() *chatServer {
	return &chatServer{
		clients:  make(map[chat.ChatService_JoinServer]bool),
		messages: make(chan *chat.Message),
	}
}

func (c *chatServer) Join(stream chat.ChatService_JoinServer) error {
	c.mu.Lock()
	c.clients[stream] = true
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		delete(c.clients, stream)
		c.mu.Unlock()
	}()

	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				return
			}

			c.messages <- msg
		}
	}()

	for messages := range c.messages {
		c.mu.Lock()
		for client := range c.clients {
			err := client.Send(messages)
			if err != nil {
				fmt.Println("Error sending message to client:", err)
			}
		}
		c.mu.Unlock()
	}

	return nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	chat.RegisterChatServiceServer(grpcServer, newServer())
	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}
