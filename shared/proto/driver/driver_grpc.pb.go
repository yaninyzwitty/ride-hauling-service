// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.0--rc3
// source: driver.proto

package driver

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	DriverService_FindNearbyDrivers_FullMethodName = "/driver.DriverService/FindNearbyDrivers"
)

// DriverServiceClient is the client API for DriverService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DriverServiceClient interface {
	FindNearbyDrivers(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[FindNearbyDriversRequest, StreamDriversResponse], error)
}

type driverServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDriverServiceClient(cc grpc.ClientConnInterface) DriverServiceClient {
	return &driverServiceClient{cc}
}

func (c *driverServiceClient) FindNearbyDrivers(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[FindNearbyDriversRequest, StreamDriversResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &DriverService_ServiceDesc.Streams[0], DriverService_FindNearbyDrivers_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[FindNearbyDriversRequest, StreamDriversResponse]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type DriverService_FindNearbyDriversClient = grpc.BidiStreamingClient[FindNearbyDriversRequest, StreamDriversResponse]

// DriverServiceServer is the server API for DriverService service.
// All implementations must embed UnimplementedDriverServiceServer
// for forward compatibility.
type DriverServiceServer interface {
	FindNearbyDrivers(grpc.BidiStreamingServer[FindNearbyDriversRequest, StreamDriversResponse]) error
	mustEmbedUnimplementedDriverServiceServer()
}

// UnimplementedDriverServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedDriverServiceServer struct{}

func (UnimplementedDriverServiceServer) FindNearbyDrivers(grpc.BidiStreamingServer[FindNearbyDriversRequest, StreamDriversResponse]) error {
	return status.Errorf(codes.Unimplemented, "method FindNearbyDrivers not implemented")
}
func (UnimplementedDriverServiceServer) mustEmbedUnimplementedDriverServiceServer() {}
func (UnimplementedDriverServiceServer) testEmbeddedByValue()                       {}

// UnsafeDriverServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DriverServiceServer will
// result in compilation errors.
type UnsafeDriverServiceServer interface {
	mustEmbedUnimplementedDriverServiceServer()
}

func RegisterDriverServiceServer(s grpc.ServiceRegistrar, srv DriverServiceServer) {
	// If the following call pancis, it indicates UnimplementedDriverServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&DriverService_ServiceDesc, srv)
}

func _DriverService_FindNearbyDrivers_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DriverServiceServer).FindNearbyDrivers(&grpc.GenericServerStream[FindNearbyDriversRequest, StreamDriversResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type DriverService_FindNearbyDriversServer = grpc.BidiStreamingServer[FindNearbyDriversRequest, StreamDriversResponse]

// DriverService_ServiceDesc is the grpc.ServiceDesc for DriverService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DriverService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "driver.DriverService",
	HandlerType: (*DriverServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "FindNearbyDrivers",
			Handler:       _DriverService_FindNearbyDrivers_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "driver.proto",
}
