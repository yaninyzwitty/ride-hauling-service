package trip_client

import (
	"context"
	"fmt"

	pb "github.com/yaninyzwitty/ride-hauling-app/shared/proto/trip"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TripServiceClient struct {
	Client pb.TripServiceClient
	conn   *grpc.ClientConn
}

func NewTripServiceClient(tripServiceAddr string) (*TripServiceClient, error) {
	conn, err := grpc.NewClient(tripServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to trip service: %w", err)
	}

	client := pb.NewTripServiceClient(conn)

	return &TripServiceClient{
		Client: client,
		conn:   conn,
	}, nil
}

func (tsc *TripServiceClient) Close() error {
	if tsc.conn != nil {
		return tsc.conn.Close()
	}
	return nil
}

func (tsc *TripServiceClient) CreateTrip(ctx context.Context, req *pb.CreateTripRequest) (*pb.CreateTripResponse, error) {
	return tsc.Client.CreateTrip(ctx, req)
}
