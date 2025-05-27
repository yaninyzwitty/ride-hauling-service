package main

import (
	"fmt"
	"log/slog"
	"time"

	pb "github.com/yaninyzwitty/ride-hauling-app/shared/proto/driver"
	"google.golang.org/grpc"
)

// DriverGrpcHandler implements the DriverService gRPC server
type DriverGrpcHandler struct {
	pb.UnimplementedDriverServiceServer
	service *Service
}

// NewDriverGrpcHandler registers the gRPC server handler for driver services
func NewDriverGrpcHandler(s *grpc.Server, service *Service) {
	handler := &DriverGrpcHandler{
		service: service,
	}
	pb.RegisterDriverServiceServer(s, handler)
}

// FindNearbyDrivers streams nearby driver updates every 3 seconds
func (h *DriverGrpcHandler) FindNearbyDrivers(reqStream pb.DriverService_FindNearbyDriversServer) error {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-reqStream.Context().Done():
			slog.Info("gRPC client disconnected")
			return reqStream.Context().Err()

		case <-ticker.C:
			drivers := h.service.FindNearbyDrivers()

			err := reqStream.Send(&pb.StreamDriversResponse{
				NearbyDrivers: drivers,
			})
			if err != nil {
				slog.Error("error streaming drivers", "error", err)
				return fmt.Errorf("streaming error: %w", err)
			}
		}
	}
}
