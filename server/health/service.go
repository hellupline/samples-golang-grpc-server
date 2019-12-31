package health

import (
	"context"

	"google.golang.org/grpc/health/grpc_health_v1"
)

// Service implements the gRPC health checking protocol.
type Service struct{}

// New creates a new Service.
func New() *Service {
	return &Service{}
}

// Check implements the Check RPC of the gRPC health checking protocol.
func (o *Service) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

// Watch implements the Watch RPC of the gRPC health checking protocol.
func (o *Service) Watch(req *grpc_health_v1.HealthCheckRequest, server grpc_health_v1.Health_WatchServer) error {
	return server.Send(&grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	})
}
