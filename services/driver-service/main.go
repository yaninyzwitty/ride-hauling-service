package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	"google.golang.org/grpc"
)

var GrpcAddr = "50051"

func main() {
	// start http server in a go routine
	// go func() {
	// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

	// 		writeJSON(w, http.StatusOK, map[string]string{"message": "Hello there💖"})
	// 	})
	// 	slog.Info("Starting driver service HTTP server", "port", HttpAddr)

	// 	// avoid fatal log
	// 	if err := http.ListenAndServe(HttpAddr, nil); err != nil {
	// 		slog.Error("HTTP server failed", "error", err)
	// 		os.Exit(1)
	// 	}
	// }()

	// start grpc server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", GrpcAddr))
	if err != nil {
		slog.Error("Failed to listen", "error", err)
		os.Exit(1)
	}
	grpcServer := grpc.NewServer()
	service := NewService()
	NewDriverGrpcHandler(grpcServer, service)

	slog.Info("starting grpc server on port",
		"grpcAddr", lis.Addr().String(),
	)

	if err := grpcServer.Serve(lis); err != nil {
		slog.Error("ffailed to serve incoming connections",
			"error", err,
		)
	}

}
