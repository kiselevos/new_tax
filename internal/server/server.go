package server

import (
	"context"
	"log/slog"
	"net"
	"new_tax/pkg/logx"

	pb "new_tax/gen/grpc/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	Grpc *grpc.Server
	Lis  net.Listener
}

func New(addr string, logger *slog.Logger) (*Server, error) {

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			UnaryRecovery(logger),
			UnaryLogger(logger),
		),
	)

	pb.RegisterTaxServiceServer(s, NewGRPCServer())
	reflection.Register(s)

	return &Server{Grpc: s, Lis: lis}, nil
}

func (s *Server) Serve() error {
	return s.Grpc.Serve(s.Lis)
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
