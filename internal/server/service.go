package server

import (
	"context"
	"log/slog"

	pb "github.com/kiselevos/new_tax/gen/grpc/api"
	"github.com/kiselevos/new_tax/internal/calculate"
	"github.com/kiselevos/new_tax/pkg/logx"
	"github.com/redis/go-redis/v9"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverStruct struct {
	pb.UnimplementedTaxServiceServer

	redis  *redis.Client
	logger *slog.Logger
}

func NewGRPCServer(rdb *redis.Client, logger *slog.Logger) *serverStruct {
	return &serverStruct{
		redis:  rdb,
		logger: logger,
	}
}

func (s *serverStruct) Healthz(ctx context.Context, req *pb.HealthzRequest) (*pb.HealthzResponse, error) {
	logx.From(ctx).Debug("healthz ok")
	return &pb.HealthzResponse{Status: "ok"}, nil
}

func (s *serverStruct) CalculatePrivate(ctx context.Context, req *pb.CalculatePrivateRequest) (*pb.CalculatePrivateResponse, error) {
	log := logx.From(ctx).With("calc_type", "private")

	if s.redis == nil {
		log.Debug("redis_disabled")
	} else {
		log.Debug("redis_enabled")
	}

	input := calculate.FromPrivateRequest(req)

	// --- validation ---
	if err := calculate.ValidateCalculateInput(input); err != nil {
		log.Warn("calc_invalid_arguments",
			"reason", err.Error(),
		)
		return nil, status.Errorf(codes.InvalidArgument, "validate: %v", err)
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

	return resp, nil
}

func (s *serverStruct) CalculatePublic(ctx context.Context, req *pb.CalculatePublicRequest) (*pb.CalculatePublicResponse, error) {
	log := logx.From(ctx).With("calc_type", "public")

	if s.redis == nil {
		log.Info("redis_disabled")
	} else {
		log.Info("redis_enabled")
	}

	input := calculate.FromPublicRequest(req)

	// --- validation ---
	if err := calculate.ValidateCalculateInput(input); err != nil {
		log.Warn("calc_invalid_arguments",
			"reason", err.Error(),
		)
		return nil, status.Errorf(codes.InvalidArgument, "validate: %v", err)
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

	return resp, nil
}
