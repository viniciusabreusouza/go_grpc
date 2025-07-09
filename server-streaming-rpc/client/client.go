package client

import (
	"context"
	"fmt"
	"io"

	"example.com/m/pb"
	"google.golang.org/grpc"
)

func Run() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	client := pb.NewStatusServiceClient(conn)

	stream, err := client.StreamStatus(context.Background(), &pb.StreamRequest{
		TaskId: "123",
	})

	if err != nil {
		panic(err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}

		fmt.Printf("received status: %s, progress: %d%% \n", res.GetMessage(), res.GetProgress())
	}
}
