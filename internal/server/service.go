package server

import (
	"context"
	"log/slog"
	"time"

	pb "github.com/kiselevos/new_tax/gen/grpc/api"
	"github.com/kiselevos/new_tax/internal/calculate"
	"github.com/kiselevos/new_tax/pkg/logx"
	"github.com/redis/go-redis/v9"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverStruct struct {
	pb.UnimplementedTaxServiceServer

	redis    *redis.Client
	logger   *slog.Logger
	cacheTTL time.Duration
}

func NewGRPCServer(rdb *redis.Client, logger *slog.Logger, ttl time.Duration) *serverStruct {
	return &serverStruct{
		redis:    rdb,
		logger:   logger,
		cacheTTL: ttl,
	}
}

func (s *serverStruct) Healthz(ctx context.Context, req *pb.HealthzRequest) (*pb.HealthzResponse, error) {
	logx.From(ctx).Debug("healthz ok")
	return &pb.HealthzResponse{Status: "ok"}, nil
}

func (s *serverStruct) CalculatePrivate(ctx context.Context, req *pb.CalculatePrivateRequest) (*pb.CalculatePrivateResponse, error) {
	log := logx.From(ctx).With("calc_type", "private")

	input := calculate.FromPrivateRequest(req)

	// --- validation ---
	if err := calculate.ValidateCalculateInput(input); err != nil {
		log.Warn("calc_invalid_arguments",
			"reason", err.Error(),
		)
		return nil, status.Errorf(codes.InvalidArgument, "validate: %v", err)
	}

	//--- redis ---
	key := buildPrivateKey(input)

	cached := &pb.CalculatePrivateResponse{}

	ok, err := s.cacheGet(ctx, key, cached)
	if err != nil {
		log.Warn("cache_get_failed", "err", err)
	}
	if ok {
		log.Info("cache_hit")
		return cached, nil
	}

	// --- calculation ---
	months := calculate.CalculateMonthlyTax(input)
	if len(months) == 0 {
		log.Error("calc_no_months_produced")
		return nil, status.Error(codes.Internal, "no data produced")
	}

	last := months[len(months)-1]

	resp := &pb.CalculatePrivateResponse{
		MonthlyDetails:        calculate.ToGRPCPrivateResponse(months),
		AnnualTaxAmount:       last.AnnualTaxAmount,
		AnnualGrossIncome:     last.AnnualGrossIncome,
		AnnualNetIncome:       last.AnnualNetIncome,
		AnnualPFR:             last.AnnualPFR,
		AnnualFSS:             last.AnnualFSS,
		AnnualFOMS:            last.AnnualFOMS,
		GrossSalary:           input.GrossSalary,
		TerritorialMultiplier: &input.TerritorialMultiplier,
		NorthernCoefficient:   &input.NorthernCoefficient,
	}

	err = s.cacheSet(ctx, key, resp, s.cacheTTL)
	if err != nil {
		log.Warn("cache_set_failed", "err", err)
	}

	return resp, nil
}

func (s *serverStruct) CalculatePublic(ctx context.Context, req *pb.CalculatePublicRequest) (*pb.CalculatePublicResponse, error) {
	log := logx.From(ctx).With("calc_type", "public")

	input := calculate.FromPublicRequest(req)

	// --- validation ---
	if err := calculate.ValidateCalculateInput(input); err != nil {
		log.Warn("calc_invalid_arguments",
			"reason", err.Error(),
		)
		return nil, status.Errorf(codes.InvalidArgument, "validate: %v", err)
	}

	// --- redis ---
	key := buildPublicKey(input)

	cached := &pb.CalculatePublicResponse{}

	ok, err := s.cacheGet(ctx, key, cached)
	if err != nil {
		log.Warn("cache_get_failed", "err", err)
	}
	if ok {
		log.Info("cache_hit")
		return cached, nil
	}

	// --- business calculation ---
	months := calculate.CalculateMonthlyTax(input)
	if len(months) == 0 {
		log.Error("calc_no_months_produced")
		return nil, status.Error(codes.Internal, "no data produced")
	}

	last := months[len(months)-1]

	resp := &pb.CalculatePublicResponse{
		MonthlyDetails:        calculate.ToGRPCPublicResponse(months),
		AnnualTaxAmount:       last.AnnualTaxAmount,
		AnnualGrossIncome:     last.AnnualGrossIncome,
		AnnualNetIncome:       last.AnnualNetIncome,
		GrossSalary:           input.GrossSalary,
		TerritorialMultiplier: &input.TerritorialMultiplier,
		NorthernCoefficient:   &input.NorthernCoefficient,
	}

	err = s.cacheSet(ctx, key, resp, s.cacheTTL)
	if err != nil {
		log.Warn("cache_set_failed", "err", err)
	}

	return resp, nil
}
