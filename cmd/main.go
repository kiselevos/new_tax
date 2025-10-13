package main

import (
	"context"
	"errors"
	"new_tax/internal/server"
	"new_tax/pkg/helpers"
	"new_tax/pkg/logx"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
)

func main() {

	logger := logx.New()
	addr := helpers.AddrChecker(os.Getenv("PORT"))

	srv, err := server.New(addr, logger)
	if err != nil {
		logger.Error("init", "err", err)
		os.Exit(1)
	}

	logger.Info("listening", "addr", addr)

	grpcErrCh := make(chan error, 1)
	go func() {
		grpcErrCh <- srv.Serve()
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Ловим сигналы для gracefull shoutdown
	select {
	case sig := <-sigCh:
		logger.Info("signal received, shutting down", "signal", sig.String())
	case err := <-grpcErrCh:
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			logger.Error("gRPC serve failed", "error", err)
		} else {
			logger.Info("gRPC server stopped")
		}
	}

	// Делаем graceful с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	server.ShutdownGRPCServer(ctx, srv.Grpc)
}
