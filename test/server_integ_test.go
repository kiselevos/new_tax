package test

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	pb "github.com/kiselevos/new_tax/gen/grpc/api"
	"github.com/kiselevos/new_tax/internal/config"
	"github.com/kiselevos/new_tax/internal/server"
	"github.com/kiselevos/new_tax/pkg/logx"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
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

const testAPIKey = "test-api-key-1234"

func newTestServer(t *testing.T) pb.TaxServiceClient {
	t.Helper()

	logger := logx.NewTest()
	cfg := config.Config{
		ApiKey:   testAPIKey,
		BackPort: "127.0.0.1:0",
		RateLimitCfg: &config.RateLimitConfig{
			PublicRPS:    1000,
			PublicBurst:  1000,
			PrivateRPS:   1000,
			PrivateBurst: 1000,
			TTL:          time.Minute,
			CleanupEvery: 100,
		},
	}

	srv, err := server.New(&cfg, logger)
	require.NoError(t, err, "server.New")

	go func() { _ = srv.Serve() }()
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		server.ShutdownGRPCServer(ctx, srv)
	})

	port := srv.Lis.Addr().(*net.TCPAddr).Port
	conn, err := grpc.NewClient(
		fmt.Sprintf("dns:///127.0.0.1:%d", port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	require.NoError(t, waitForReady(ctx, conn))

	return pb.NewTaxServiceClient(conn)
}

func Test_Server_Healthz(t *testing.T) {
	cli := newTestServer(t)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := cli.Healthz(ctx, &pb.HealthzRequest{})
	require.NoError(t, err)
	assert.Equal(t, "ok", res.GetStatus())
}

func Test_CalculatePublic(t *testing.T) {
	cli := newTestServer(t)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := cli.CalculatePublic(ctx, &pb.CalculatePublicRequest{
		GrossSalary: 100_000_00, // 100 000 ₽
	})
	require.NoError(t, err)

	assert.Len(t, res.GetMonthlyDetails(), 12, "должно быть 12 месяцев")
	assert.Greater(t, res.GetAnnualTaxAmount(), uint64(0), "налог должен быть ненулевым")
	assert.Equal(t, uint64(100_000_00), res.GetGrossSalary())
	// 100 000 * 13% * 12 = 156 000 ₽
	assert.Equal(t, uint64(156_000_00), res.GetAnnualTaxAmount())
}

func Test_CalculatePrivate(t *testing.T) {
	cli := newTestServer(t)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	ctx = metadata.AppendToOutgoingContext(ctx, "x-api-key", testAPIKey)

	res, err := cli.CalculatePrivate(ctx, &pb.CalculatePrivateRequest{
		GrossSalary: 100_000_00, // 100 000 ₽
	})
	require.NoError(t, err)

	assert.Len(t, res.GetMonthlyDetails(), 12)
	assert.Equal(t, uint64(156_000_00), res.GetAnnualTaxAmount())
	assert.Greater(t, res.GetAnnualPFR(), uint64(0), "ПФР должен быть ненулевым")
	assert.Greater(t, res.GetAnnualFOMS(), uint64(0), "ФОМС должен быть ненулевым")
}
