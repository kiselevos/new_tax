package main

import (
	"context"
	"errors"
	"net/http"
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

	// Определяем адрес gRPC сервера
	addr := helpers.AddrChecker(os.Getenv("GRPC_PORT"))
	srv, err := server.New(addr, logger)
	if err != nil {
		logger.Error("failed to init server", "err", err)
		os.Exit(1)
	}

	logger.Info("starting servers", "grpc_addr", addr, "grpc_web_addr", ":8081")

	// Каналы для ошибок и сигналов
	grpcErrCh := make(chan error, 1)
	grpcWebErrCh := make(chan error, 1)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Запуск gRPC сервера
	go func() {
		grpcErrCh <- srv.Serve()
	}()

	// Запуск gRPC-Web сервера
	go func() {
		grpcWebErrCh <- srv.ServeGRPCWeb(":8081")
	}()

	// Основной селект: ждём ошибок или сигналов
	select {
	case sig := <-sigCh:
		logger.Info("signal received, shutting down", "signal", sig.String())

	case err := <-grpcErrCh:
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			logger.Error("gRPC serve failed", "error", err)
		} else {
			logger.Info("gRPC server stopped")
		}

	case err := <-grpcWebErrCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("gRPC-Web serve failed", "error", err)
		} else {
			logger.Info("gRPC-Web server stopped")
		}
	}

	// Плавное завершение
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	server.ShutdownGRPCServer(ctx, srv.Grpc)
	logger.Info("shutdown complete ✅")
}
