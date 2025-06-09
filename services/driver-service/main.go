package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/yaninyzwitty/ride-hauling-app/shared/pkg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

const (
	shutdownTimeout = 3 * time.Second
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Load configuration
	config := &pkg.Config{}
	if err := config.LoadConfig("config.yaml"); err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	driverServiceAddr := fmt.Sprintf(":%d", config.DriverService.Port)

	lis, err := net.Listen("tcp", driverServiceAddr)
	if err != nil {
		slog.Error("failed to listen", "error", err, "address", driverServiceAddr)
		os.Exit(1)
	}
	// Configure gRPC server options
	serverOpts := []grpc.ServerOption{
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     5 * time.Minute,
			MaxConnectionAge:      10 * time.Minute,
			MaxConnectionAgeGrace: 5 * time.Second,
			Time:                  2 * time.Minute,
			Timeout:               20 * time.Second,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             10 * time.Second,
			PermitWithoutStream: true,
		}),
	}

	// Create gRPC server
	grpcServer := grpc.NewServer(serverOpts...)
	// Register services

	driverService := NewService()

	NewDriverGrpcHandler(grpcServer, driverService)
	reflection.Register(grpcServer) //   Enable reflection for development tools
	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		sig := <-quit
		slog.Info("received shutdown signal", "signal", sig)

		// Create shutdown context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		// Create shutdown channel
		shutdownCh := make(chan struct{})

		go func() {
			grpcServer.GracefulStop()
			close(shutdownCh)
		}()

		select {
		case <-ctx.Done():
			slog.Warn("shutdown timeout reached, forcing stop")
			grpcServer.Stop()
		case <-shutdownCh:
			slog.Info("server stopped gracefully")
		}
	}()

	// Start server
	slog.Info("starting gRPC server",
		"address", driverServiceAddr,
		"pid", os.Getpid(),
	)

	if err := grpcServer.Serve(lis); err != nil {
		slog.Error("failed to serve gRPC server", "error", err)
		os.Exit(1)
	}

	wg.Wait()
	slog.Info("service shutdown complete")

}
