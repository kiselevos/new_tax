package main

import (
	"context"
	"log/slog"
	"net/http"
	"new_tax/internal/server"
	"new_tax/pkg/helpers"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger := slog.Default()

	srv := server.New(helpers.GetGRPCWebPort(), logger)

	go func() {
		logger.Info("Server starting",
			"addr", helpers.GetGRPCWebPort(),
			"log_level", os.Getenv("LOG_LEVEL"),
			"mode", os.Getenv("LOG_MODE"),
		)
		if err := srv.Serve(); err != nil && err != http.ErrServerClosed {
			logger.Error("HTTP server failed", "error", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	logger.Info("Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("HTTP server shutdown failed", "error", err)
	} else {
		logger.Info("Server stopped gracefully")
	}
}
