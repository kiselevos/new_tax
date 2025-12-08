package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/kiselevos/new_tax/internal/config"
	"github.com/kiselevos/new_tax/internal/server"
	"github.com/kiselevos/new_tax/pkg/helpers"
	"github.com/kiselevos/new_tax/pkg/logx"

	"google.golang.org/grpc"
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("env_file_not_loaded", "err", err)
	}

	conf, err := config.Load()
	if err != nil {
		log.Fatal("can't load config:", err)
	}

	logger := logx.New(conf.LogMode, conf.LogLevel)
	slog.SetDefault(logger)

	addr := helpers.AddrChecker(conf.BackPort)

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
