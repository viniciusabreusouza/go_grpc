package server

import (
	"io"
	"net"

	"example.com/m/pb"
	"google.golang.org/grpc"
)

type TemperatureServer struct {
	pb.UnsafeTemperatureServiceServer
}

func (s *TemperatureServer) RecordTemperature(stream grpc.ClientStreamingServer[pb.TemperatureRequest, pb.TemperatureResponse]) error {
	var sum float32
	var count int32

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&pb.TemperatureResponse{
				AverageTemperature: sum / float32(count),
			})
		}

		if err != nil {
			return err
		}

		sum += req.Temperature
		count++
	}
}

func Run() {
	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterTemperatureServiceServer(grpcServer, &TemperatureServer{})
	if err := grpcServer.Serve(listen); err != nil {
		panic(err)
	}
}
