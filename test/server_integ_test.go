package test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	pb "new_tax/gen/grpc/api"
	taxconnect "new_tax/gen/grpc/api/taxconnect"
	"new_tax/internal/server"
	"new_tax/pkg/logx"

	"connectrpc.com/connect"
)

func Test_Server_Healthz(t *testing.T) {
	logger := logx.New()
	srv := server.New(":0", logger)             // твой ConnectRPC сервер
	ts := httptest.NewServer(srv.HttpHandler()) // создаём тестовый HTTP сервер
	defer ts.Close()

	// Connect client
	client := taxconnect.NewTaxServiceClient(
		http.DefaultClient,
		ts.URL,
		connect.WithGRPC(),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := client.Healthz(ctx, connect.NewRequest(&pb.HealthzRequest{}))
	if err != nil {
		t.Fatalf("Healthz RPC failed: %v", err)
	}

	if res.Msg.GetStatus() != "ok" {
		t.Fatalf("unexpected healthz status: got %q, want %q", res.Msg.GetStatus(), "ok")
	}
}
