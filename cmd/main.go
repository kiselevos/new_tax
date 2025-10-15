package main

import (
	"context"
	"log/slog"
	"new_tax/internal/server"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger := slog.Default()

	srv := server.New(":8081", logger)
	go srv.Serve()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	srv.Shutdown(ctx)
}
