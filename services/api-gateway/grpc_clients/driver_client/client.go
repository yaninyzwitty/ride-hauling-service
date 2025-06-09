package driver_client

import (
	"context"
	"fmt"

	pb "github.com/yaninyzwitty/ride-hauling-app/shared/proto/driver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DriverServiceClient struct {
	Client              pb.DriverServiceClient
	conn                *grpc.ClientConn
	grpcDriverClientUrl string
}

func NewDriverServiceClient(grpcDriverClientUrl string) (*DriverServiceClient, error) {

	conn, err := grpc.NewClient(grpcDriverClientUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to driver service: %w", err)
	}

	client := pb.NewDriverServiceClient(conn)

	return &DriverServiceClient{
		Client: client,
		conn:   conn,
	}, nil
}

func (dsc *DriverServiceClient) Close() error {
	if dsc.conn != nil {
		return dsc.conn.Close()
	}
	return nil
}

func (dsc *DriverServiceClient) FindNearbyDrivers(ctx context.Context) (pb.DriverService_FindNearbyDriversClient, error) {
	return dsc.Client.FindNearbyDrivers(ctx)

}
