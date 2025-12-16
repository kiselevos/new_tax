package main

import (
	"context"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/kiselevos/new_tax/pkg/logx"
	"github.com/kiselevos/new_tax/web"
	"github.com/kiselevos/new_tax/web/handlers"
	"github.com/kiselevos/new_tax/web/internal/api"
	"github.com/kiselevos/new_tax/web/internal/client"
	"github.com/kiselevos/new_tax/web/internal/config"
	"github.com/kiselevos/new_tax/web/internal/geoip"
	"github.com/kiselevos/new_tax/web/internal/middleware"
	"github.com/kiselevos/new_tax/web/internal/server"
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		slog.Error("env_file_not_loaded", "err", err)
		os.Exit(1)
	}

	cfg, err := config.Load()
	if err != nil {
		slog.Error("config_load_failed", "err", err)
		os.Exit(1)
	}

	logger := logx.New(cfg.LogMode, cfg.LogLevel)
	slog.SetDefault(logger)

	tmpl, err := template.New("").Funcs(web.Funcs).ParseGlob("templates/*.tmpl")
	if err != nil {
		logger.Error("templates_parse_failed", "err", err)
		os.Exit(1)
	}

	htmlMux := http.NewServeMux()
	apiMux := http.NewServeMux()

	clientGRPC, conn, err := client.NewTaxClient(cfg.Backend)
	if err != nil {
		logger.Error("grpc_dial_failed", "err", err)
		os.Exit(1)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			logger.Error("failed to close connection", "err", err)
		}
	}()

	htmlServer := handlers.NewServer(tmpl, clientGRPC)

	htmlServer.Routes(htmlMux)
	api.RegisterApiRoutes(apiMux, clientGRPC, cfg.APIVersion, tmpl)

	htmlMux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	apiHandler := middleware.Chain(
		apiMux,
		middleware.CORSMiddleware,
	)

	rootMux := http.NewServeMux()
	rootMux.Handle("/api/", apiHandler)
	rootMux.Handle("/", htmlMux)

	// Подключаем GeoDataIP
	geoDB, err := geoip.LoadFromCSV(cfg.GeoIPPath)
	if err != nil {
		logger.Warn("geoip_disabled",
			"path", cfg.GeoIPPath,
			"err", err,
		)
		geoDB = geoip.NewEmpty()
	}

	// подключаем метрики
	rootHandler := middleware.Chain(
		rootMux,
		middleware.RegionMiddleware(geoDB),
		middleware.MetricsMiddleware,
		middleware.Logger,
	)

	httpSrv := server.New(cfg.WebPort, rootHandler)

	go func() {
		logger.Info("listening", "addr", cfg.WebPort)
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
