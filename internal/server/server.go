package server

import (
	"context"
	"log/slog"
	"net/http"

	taxconnect "new_tax/gen/grpc/api/taxconnect"

	"connectrpc.com/grpcreflect"
)

// Server — структура HTTP сервера
type Server struct {
	httpServer *http.Server
}

// New создаёт новый HTTP сервер для ConnectRPC
func New(addr string, logger *slog.Logger) *Server {
	mux := http.NewServeMux()

	// ✅ Регистрируем сервис Connect
	service := &taxServiceServer{}
	path, handler := taxconnect.NewTaxServiceHandler(service)
	mux.Handle(path, handler)

	// ✅ Включаем рефлексию (для дебага и buf curl)
	reflector := grpcreflect.NewStaticReflector("tax.TaxService")
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	// ✅ Оборачиваем в middleware Connect (CORS и logging можно добавить здесь)
	handlerWithCORS := withCORS(mux)

	s := &http.Server{
		Addr:    addr,
		Handler: handlerWithCORS,
	}
	return &Server{httpServer: s}
}

// Serve запускает сервер
func (s *Server) Serve() error {
	slog.Info("🌐 Connect server listening on", "addr", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

// Shutdown завершает сервер
func (s *Server) Shutdown(ctx context.Context) error {
	slog.Info("🧩 Connect server shutting down...")
	return s.httpServer.Shutdown(ctx)
}

// withCORS — разрешает фронтенду обращаться к backend
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Connect-Protocol-Version, Connect-Timeout-Ms")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
