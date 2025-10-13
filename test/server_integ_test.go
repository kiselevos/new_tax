package test

import (
	"context"
	"testing"
	"time"

	pb "new_tax/gen/grpc/api"
	"new_tax/internal/server"
	"new_tax/pkg/logx"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Test_Server_Healthz(t *testing.T) {
	logger := logx.New()

	srv, err := server.New(":0", logger)
	if err != nil {
		t.Fatalf("server.New: %v", err)
	}
	t.Logf("listening on %s", srv.Lis.Addr().String())

	go func() {
		_ = srv.Serve()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		srv.Lis.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	cli := pb.NewTaxServiceClient(conn)

	resp, err := cli.Healthz(ctx, &pb.HealthzRequest{})
	if err != nil {
		t.Fatalf("Healthz call failed: %v", err)
	}
	if resp.GetStatus() != "ok" {
		t.Fatalf("unexpected status: %q", resp.GetStatus())
	}

	shCtx, shCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shCancel()
	server.ShutdownGRPCServer(shCtx, srv.Grpc)
}
