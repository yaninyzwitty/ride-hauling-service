package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	pb "github.com/yaninyzwitty/ride-hauling-app/shared/proto/trip"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TripGrpcHandler struct {
	pb.UnimplementedTripServiceServer
}

type OSRMResponse struct {
	Routes []struct {
		Distance float64 `json:"distance"`
		Duration float64 `json:"duration"`
		Geometry struct {
			Coordinates [][]float64 `json:"coordinates"`
		} `json:"geometry"`
	} `json:"routes"`
}

func NewTripGrpcHandler(s *grpc.Server) {
	// create a new TripGrpcHandler then register it with the gRPC server
	handler := &TripGrpcHandler{}

	pb.RegisterTripServiceServer(s, handler)

}

func (h *TripGrpcHandler) CreateTrip(ctx context.Context, req *pb.CreateTripRequest) (*pb.CreateTripResponse, error) {
	pickup := req.GetStartLocation()
	destination := req.GetEndLocation()

	url := fmt.Sprintf(
		"https://router.project-osrm.org/route/v1/driving/%f,%f;%f,%f?overview=full&geometries=geojson",
		pickup.Longitude, pickup.Latitude,
		destination.Longitude, destination.Latitude,
	)

	// make http request to the OSRM API
	resp, err := http.Get(url)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to make request to OSRM API: %v", err)
	}

	defer resp.Body.Close()
	// read and parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to read response body: %v", err)
	}

	var routeResp OSRMResponse
	if err := json.Unmarshal(body, &routeResp); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to parse response body: %v", err)
	}

	// parse OSRMResponse to pb.Route
	route := routeResp.Routes[0]

	geometry := route.Geometry.Coordinates

	coordinates := make([]*pb.Coordinate, len(geometry))

	for i, coord := range geometry {
		coordinates[i] = &pb.Coordinate{
			Latitude:  float32(coord[0]),
			Longitude: float32(coord[1]),
		}
	}

	return &pb.CreateTripResponse{
		Route: &pb.Route{
			Geometry: []*pb.Geometry{
				{
					Coordinates: coordinates,
				},
			},
			Distance: float32(route.Distance),
			Duration: float32(route.Duration),
		},
	}, nil

}
