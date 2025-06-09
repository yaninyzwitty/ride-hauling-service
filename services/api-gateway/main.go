package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	trippb "github.com/yaninyzwitty/ride-hauling-app/services/api-gateway/grpc_clients/trip_client"
	"github.com/yaninyzwitty/ride-hauling-app/shared/pkg"
)

func main() {
	config := &pkg.Config{}

	if err := config.LoadConfig("config.yaml"); err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// create a trip service url preferably (use string.builder for performance)
	tripServiceUrl := fmt.Sprintf(":%d", config.TripService.Port)

	// create a trip client
	tripClient, err := trippb.NewTripServiceClient(tripServiceUrl)
	if err != nil {
		slog.Error("failed to create trip service client", "error", err)
		return
	}
	defer tripClient.Close()

	// Create a new server
	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", config.APIGateway.Port),
	}

	// Setup routes
	http.HandleFunc("/trip", EnableCors(func(w http.ResponseWriter, r *http.Request) {
		handleCreateTrip(w, r, tripClient)
	}))

	// Graceful shutdown setup
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		slog.Info("API Gateway running", "port", config.APIGateway.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	<-quit
	slog.Info("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("Server exited gracefully")
}

// try this one
// {
// 	"start_location": {
// 	  "latitude": 34.0522,
// 	  "longitude": -118.2437
// 	},
// 	"end_location": {
// 	  "latitude": 37.7749,
// 	  "longitude": -122.4194
// 	}
//   }
