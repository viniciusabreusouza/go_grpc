package client

import (
	"context"
	"fmt"
	"time"

	"example.com/m/pb"
	"google.golang.org/grpc"
)

func Run() {
	dial, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	defer dial.Close()

	client := pb.NewTemperatureServiceClient(dial)

	stream, err := client.RecordTemperature(context.Background())

	if err != nil {
		panic(err)
	}

	temperatures := []float32{10, 25, 15, 30, 33}

	for _, temp := range temperatures {
		fmt.Printf("Sending temperature: %f\n", temp)
		err := stream.Send(&pb.TemperatureRequest{
			Temperature: temp,
		})

		if err != nil {
			panic(err)
		}

		time.Sleep(1 * time.Second)
	}

	recv, err := stream.CloseAndRecv()

	if err != nil {
		panic(err)
	}

	fmt.Printf("Average temperature: %f\n", recv.AverageTemperature)
}
