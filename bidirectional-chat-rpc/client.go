package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/viniciusabreusouza/go_grpc/chat"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal("Could not connect:", err)
	}

	defer connection.Close()

	chatClient := chat.NewChatServiceClient(connection)

	stream, err := chatClient.Join(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("Enter your name: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	user := scanner.Text()

	go func() {
		msg, err := stream.Recv()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("[%s] %s: %s\n", time.Unix(msg.Timestamp, 0).Format("15:04:05"), msg.User, msg.Text)
	}()

	for scanner.Scan() {
		msg := &chat.Message{
			User:      user,
			Text:      scanner.Text(),
			Timestamp: time.Now().Unix(),
		}

		if err := stream.Send(msg); err != nil {
			log.Fatalf("Failed to send message: %v", err)
		}
	}
}
