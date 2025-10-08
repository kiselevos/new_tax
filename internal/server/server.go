package server

import (
	"context"
	pb "new_tax/gen/grpc/api"
)

type serverStruct struct {
	pb.UnimplementedTaxServiceServer
}

func NewGRPCServer() *serverStruct {
	return &serverStruct{}
}

func (s *serverStruct) Healthz(ctx context.Context, req *pb.HealthzRequest) (*pb.HealthzResponse, error) {
	return &pb.HealthzResponse{Status: "ok"}, nil
}
