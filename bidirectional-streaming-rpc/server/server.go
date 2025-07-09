package server

import (
	"fmt"
	"io"
	"math/rand/v2"
	"net"

	"example.com/m/pb"
	"google.golang.org/grpc"
)

type StockServiceServer struct {
	pb.UnimplementedStockServiceServer
}

func (*StockServiceServer) StreamStockPrices(stream grpc.BidiStreamingServer[pb.StockRequest, pb.StockResponse]) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		symbol := req.Symbol
		fmt.Printf("Received symbol: %s\n", symbol)

		go func(symbol string) {
			for i := 0; i < 10; i++ {
				price := rand.Float32() * 100

				err := stream.Send(&pb.StockResponse{
					Symbol: symbol,
					Price:  price,
				})
				if err != nil {
					panic(fmt.Sprintf("failed to send stock price: %s", err.Error()))
				}
			}
		}(symbol)
	}
}

func Run() {
	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterStockServiceServer(grpcServer, &StockServiceServer{})
	if err := grpcServer.Serve(listen); err != nil {
		panic(err)
	}
}
