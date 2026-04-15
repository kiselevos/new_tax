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

// newTestServer поднимает gRPC-сервер с переданным RedisConfig и возвращает клиент и cleanup.
func newTestServer(t *testing.T, redisCfg *config.RedisConfig) pb.TaxServiceClient {
	t.Helper()

	logger := logx.NewTest()
	cfg := config.Config{
		ApiKey:   testAPIKey,
		BackPort: "127.0.0.1:0",
		RedisCfg: redisCfg,
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
	t.Helper()

	logger := logx.NewTest()
	cfg := config.Config{
		ApiKey:   "1",
		BackPort: "127.0.0.1:0",
		RedisCfg: &config.RedisConfig{Enabled: false},

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
	if err != nil {
		t.Fatalf("server.New: %v", err)
	}
	go func() {
		_ = srv.Serve()
	}()
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		server.ShutdownGRPCServer(ctx, srv)
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

// -------------------------------------------------------------------
// Graceful degradation: Redis недоступен
// -------------------------------------------------------------------

// Test_Server_StartsWhenRedisUnavailable проверяет, что сервер поднимается
// корректно, если Redis включён в конфиге, но физически недоступен.
func Test_Server_StartsWhenRedisUnavailable(t *testing.T) {
	cli := newTestServer(t, &config.RedisConfig{
		Enabled: true,
		Addr:    "127.0.0.1:1", // заведомо нерабочий адрес
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := cli.Healthz(ctx, &pb.HealthzRequest{})
	require.NoError(t, err, "сервер должен отвечать даже без Redis")
	assert.Equal(t, "ok", res.GetStatus())
}

// Test_CalculatePublic_WorksWithoutRedis проверяет, что CalculatePublic
// возвращает корректный расчёт при redis = nil (Enabled: false).
func Test_CalculatePublic_WorksWithoutRedis(t *testing.T) {
	cli := newTestServer(t, &config.RedisConfig{Enabled: false})

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

// Test_CalculatePrivate_WorksWithoutRedis проверяет, что CalculatePrivate
// возвращает полный ответ (включая взносы работодателя) при redis = nil.
func Test_CalculatePrivate_WorksWithoutRedis(t *testing.T) {
	cli := newTestServer(t, &config.RedisConfig{Enabled: false})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	ctx = metadata.AppendToOutgoingContext(ctx, "x-api-key", testAPIKey)

	res, err := cli.CalculatePrivate(ctx, &pb.CalculatePrivateRequest{
		GrossSalary: 100_000_00, // 100 000 ₽
	})
	require.NoError(t, err)

	assert.Len(t, res.GetMonthlyDetails(), 12)
	assert.Equal(t, uint64(156_000_00), res.GetAnnualTaxAmount())
	// Взносы работодателя тоже должны быть заполнены
	assert.Greater(t, res.GetAnnualPFR(), uint64(0), "ПФР должен быть ненулевым")
	assert.Greater(t, res.GetAnnualFOMS(), uint64(0), "ФОМС должен быть ненулевым")
}
