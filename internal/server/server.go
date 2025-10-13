package server

import (
	"context"
	pb "new_tax/gen/grpc/api"
	"new_tax/pkg/logx"

	"google.golang.org/grpc"
)

type serverStruct struct {
	pb.UnimplementedTaxServiceServer
}

func NewGRPCServer() *serverStruct {
	return &serverStruct{}
}

func (s *serverStruct) Healthz(ctx context.Context, req *pb.HealthzRequest) (*pb.HealthzResponse, error) {
	logx.From(ctx).Info("healthz ok")
	return &pb.HealthzResponse{Status: "ok"}, nil
}

// ShutdownGRPCServer завершает работу gRPC-сервера с использованием GracefulStop.
func ShutdownGRPCServer(ctx context.Context, srv *grpc.Server) {
	log := logx.From(ctx)
	log.Info("Shutting down gRPC server gracefully...")

	done := make(chan struct{})

	go func() {
		srv.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		log.Info("gRPC server graceful shutdown complete")
	case <-ctx.Done():
		log.Warn("graceful shutdown timed out, forcing stop")
		srv.Stop()
	}
}
