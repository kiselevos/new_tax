package server

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"new_tax/pkg/logx"

	pb "new_tax/gen/grpc/api"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
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

// Serve запускает классический gRPC сервер
func (s *Server) Serve() error {
	slog.Info("gRPC listening on", "addr", s.Lis.Addr().String())
	return s.Grpc.Serve(s.Lis)
}

// ServeGRPCWeb поднимает HTTP сервер, совместимый с gRPC-Web
func (s *Server) ServeGRPCWeb(addr string) error {
	wrapped := grpcweb.WrapServer(
		s.Grpc,
		grpcweb.WithOriginFunc(func(origin string) bool { return true }), // Разрешаем все CORS
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if wrapped.IsGrpcWebRequest(r) ||
			wrapped.IsAcceptableGrpcCorsRequest(r) ||
			wrapped.IsGrpcWebSocketRequest(r) {
			wrapped.ServeHTTP(w, r)
			return
		}
		http.NotFound(w, r)
	})

	slog.Info("gRPC-Web listening on", "addr", addr)
	return http.ListenAndServe(addr, mux)
}

// ShutdownGRPCServer завершает работу gRPC сервера с GracefulStop
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
