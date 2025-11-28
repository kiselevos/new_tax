package testutils

import (
	"context"

	"google.golang.org/grpc"

	pb "github.com/kiselevos/new_tax/gen/grpc/api"
)

type FakeTaxClient struct {
	PublicResp  *pb.CalculatePublicResponse
	PublicErr   error
	PrivateResp *pb.CalculatePrivateResponse
	PrivateErr  error
	HealthzResp *pb.HealthzResponse
	HealthzErr  error
}

func (f *FakeTaxClient) CalculatePublic(ctx context.Context, req *pb.CalculatePublicRequest, opts ...grpc.CallOption) (*pb.CalculatePublicResponse, error) {
	if f.PublicResp != nil || f.PublicErr != nil {
		return f.PublicResp, f.PublicErr
	}
	return &pb.CalculatePublicResponse{}, nil
}

func (f *FakeTaxClient) CalculatePrivate(ctx context.Context, req *pb.CalculatePrivateRequest, opts ...grpc.CallOption) (*pb.CalculatePrivateResponse, error) {
	if f.PrivateResp != nil || f.PrivateErr != nil {
		return f.PrivateResp, f.PrivateErr
	}
	return &pb.CalculatePrivateResponse{}, nil
}
func (f *FakeTaxClient) Healthz(ctx context.Context, req *pb.HealthzRequest, opts ...grpc.CallOption) (*pb.HealthzResponse, error) {
	if f.HealthzResp != nil || f.HealthzErr != nil {
		return f.HealthzResp, f.HealthzErr
	}
	return &pb.HealthzResponse{}, nil
}
