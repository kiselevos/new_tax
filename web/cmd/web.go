package main

import (
	"context"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/kiselevos/new_tax/pkg/logx"
	"github.com/kiselevos/new_tax/web"
	"github.com/kiselevos/new_tax/web/handlers"
	"github.com/kiselevos/new_tax/web/internal/client"
	"github.com/kiselevos/new_tax/web/internal/middleware"
	"github.com/kiselevos/new_tax/web/internal/server"
)

func main() {

	logger := logx.New()

	if err := godotenv.Load(".env"); err != nil {
		logger.Error("file .env not read", "err", err)
		os.Exit(1)
	}

	addr := os.Getenv("WEB_PORT")

	tmpls, err := template.New("").Funcs(web.Funcs).ParseGlob("templates/*.tmpl")
	if err != nil {
		logger.Error("templates_parse_failed", "err", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()

	clientGRPC, conn, err := client.NewTaxClient()
	if err != nil {
		logger.Error("grpc_dial_failed", "err", err)
		os.Exit(1)
	}
	defer conn.Close()

	s := handlers.NewServer(tmpls, clientGRPC)

	s.Routes(mux)

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	httpSrv := server.New(addr, middleware.Logger(mux))

	go func() {
		logger.Info("listening", "addr", addr)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("http_serve_failed", "err", err)
			os.Exit(1)
		}
	}()

	// graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	logger.Info("shutdown_request")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		logger.Warn("http_shutdown_timeout", "err", err)
		_ = httpSrv.Close()
	}
	logger.Info("shutdown_complete")
}
