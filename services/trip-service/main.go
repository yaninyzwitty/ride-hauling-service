package main

import (
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

var GrpcAddr = ":50051"

func main() {
	lis, err := net.Listen("tcp", GrpcAddr)

	if err != nil {
		slog.Error("failed to listen: ", "error", err)
		return
	}

	grpcServer := grpc.NewServer()

	NewTripGrpcHandler(grpcServer)

	slog.Info("trip service started", "address", GrpcAddr)
	if err := grpcServer.Serve(lis); err != nil {
		slog.Error("failed to serve: ", "error", err)
		return
	}

}
