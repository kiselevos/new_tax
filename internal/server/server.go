package server

import (
	"context"
	"log/slog"
	"net"

	"github.com/kiselevos/new_tax/internal/config"
	"github.com/kiselevos/new_tax/internal/middleware"
	"github.com/kiselevos/new_tax/pkg/logx"

	pb "github.com/kiselevos/new_tax/gen/grpc/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	Grpc *grpc.Server
	Lis  net.Listener
}

func New(cfg *config.Config, logger *slog.Logger) (*Server, error) {

	lis, err := net.Listen("tcp", cfg.BackPort)
	if err != nil {
		return nil, err
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.UnaryRecovery(logger),
			middleware.UnaryLogger(logger),
			middleware.Auth(cfg.ApiKey),
			middleware.RateLimitInterceptor(cfg.RateLimitCfg),
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
