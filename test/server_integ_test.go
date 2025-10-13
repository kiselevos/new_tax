package test

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	pb "new_tax/gen/grpc/api"
	"new_tax/internal/server"
	"new_tax/pkg/logx"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

func waitForReady(ctx context.Context, conn *grpc.ClientConn) error {
	conn.Connect()
	for {
		st := conn.GetState()
		if st == connectivity.Ready {
			return nil
		}
		if !conn.WaitForStateChange(ctx, st) {
			return ctx.Err()
		}
	}
}

func Test_Server_Healthz(t *testing.T) {
	t.Helper()

	logger := logx.New()

	addr := "127.0.0.1:0"
	srv, err := server.New(addr, logger)
	if err != nil {
		t.Fatalf("server.New: %v", err)
	}
	go func() {
		_ = srv.Serve()
	}()
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		server.ShutdownGRPCServer(ctx, srv.Grpc)
	})

	laddr := srv.Lis.Addr()
	t.Logf("listening on %s", laddr.String())

	var port int
	if ta, ok := laddr.(*net.TCPAddr); ok {
		port = ta.Port
	} else {
		t.Fatalf("unexpected listener addr type: %T", laddr)
	}

	target := fmt.Sprintf("dns:///127.0.0.1:%d", port)

	conn, err := grpc.NewClient(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("grpc.NewClient: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := waitForReady(ctx, conn); err != nil {
		t.Fatalf("wait for ready: %v (state=%v)", err, conn.GetState())
	}

	cli := pb.NewTaxServiceClient(conn)
	res, err := cli.Healthz(ctx, &pb.HealthzRequest{})
	if err != nil {
		t.Fatalf("Healthz RPC failed: %v", err)
	}
	if res.GetStatus() != "ok" {
		t.Fatalf("unexpected healthz: %q", res.GetStatus())
	}
}
