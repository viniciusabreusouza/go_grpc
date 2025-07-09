package client

import (
	"context"
	"fmt"
	"io"
	"time"

	"example.com/m/pb"
	"google.golang.org/grpc"
)

func RunClient() {
	dial, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	defer dial.Close()

	client := pb.NewStockServiceClient(dial)

	stream, err := client.StreamStockPrices(context.Background())

	if err != nil {
		panic(err)
	}

	done := make(chan bool)
	go func() {
		for {
			response, err := stream.Recv()

			if err != nil {
				break
			}

			if err == io.EOF {
				break
			}

			fmt.Printf("Received symbol: %s, price: %f\n", response.Symbol, response.Price)
		}

		close(done)
	}()

	symbols := []string{"AAPL", "MSFT", "GOOG", "AMZN", "FB"}

	for _, symbol := range symbols {
		err := stream.Send(&pb.StockRequest{
			Symbol: symbol,
		})

		if err != nil {
			panic(err)
		}

		time.Sleep(2 * time.Second)
	}

	<-done
}
