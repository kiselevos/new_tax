package main

import (
	"context"
	"errors"
	"flag"
	"net"
	"new_tax/internal/server"
	"new_tax/pkg/logx"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "new_tax/gen/grpc/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	addr := flag.String("addr", ":50051", "listen address")
	flag.Parse()

	logger := logx.New()

	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		logger.Error("listen failed", "error", err)
		os.Exit(1)
	}

	srvGRPC := grpc.NewServer(
		grpc.ChainUnaryInterceptor(server.UnaryLogger(logger)),
	)
	srv := server.NewGRPCServer()
	pb.RegisterTaxServiceServer(srvGRPC, srv)

	reflection.Register(srvGRPC)

	logger.Info("TaxService listening", "info", *addr)

	grpcErrCh := make(chan error, 1)
	go func() {
		grpcErrCh <- srvGRPC.Serve(lis)
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Ловим сигналы для gracefull shoutdown
	select {
	case sig := <-sigCh:
		logger.Info("signal received, shutting down", "signal", sig.String())

		// Делаем graceful с таймаутом
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		server.ShutdownGRPCServer(ctx, srvGRPC)
	case err := <-grpcErrCh:
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			logger.Error("gRPC serve failed", "error", err)
		} else {
			logger.Info("gRPC server stopped")
		}
	}
}
