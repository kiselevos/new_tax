package server

import (
	"context"
	"log/slog"
	"net"

	"github.com/kiselevos/new_tax/internal/config"
	"github.com/kiselevos/new_tax/internal/middleware"
	"github.com/kiselevos/new_tax/internal/redisconn"
	"github.com/kiselevos/new_tax/pkg/logx"
	"github.com/redis/go-redis/v9"

	pb "github.com/kiselevos/new_tax/gen/grpc/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	Grpc  *grpc.Server
	Lis   net.Listener
	Redis *redis.Client
}

func New(cfg *config.Config, logger *slog.Logger) (*Server, error) {

	lis, err := net.Listen("tcp", cfg.BackPort)
	if err != nil {
		return nil, err
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.UnaryRecovery(),
			middleware.UnaryLogger(),
			middleware.Auth(cfg.ApiKey),
			middleware.RateLimitInterceptor(cfg.RateLimitCfg),
		),
	)

	rdb, err := redisconn.Connect(cfg.RedisCfg, logger)
	if err != nil {
		rdb = nil
	}

	pb.RegisterTaxServiceServer(s, NewGRPCServer(rdb, logger))
	reflection.Register(s)

	return &Server{Grpc: s, Lis: lis, Redis: rdb}, nil
}

func (s *Server) Serve() error {
	return s.Grpc.Serve(s.Lis)
}

// ShutdownGRPCServer завершает работу gRPC-сервера с использованием GracefulStop.
func ShutdownGRPCServer(ctx context.Context, srv *Server) {
	log := logx.From(ctx)
	log.Info("Shutting down gRPC server gracefully...")

	if srv.Redis != nil {
		if err := srv.Redis.Close(); err != nil {
			log.Warn("redis_close_failed", "err", err)
		} else {
			log.Info("redis_closed")
		}
	}

	done := make(chan struct{})

	go func() {
		srv.Grpc.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		log.Info("gRPC server graceful shutdown complete")
	case <-ctx.Done():
		log.Warn("graceful shutdown timed out, forcing stop")
		srv.Grpc.Stop()
	}
}
