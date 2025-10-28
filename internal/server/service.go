package server

import (
	"context"
	"time"

	pb "github.com/kiselevos/new_tax/gen/grpc/api"
	"github.com/kiselevos/new_tax/internal/calculate"
	"github.com/kiselevos/new_tax/pkg/logx"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverStruct struct {
	pb.UnimplementedTaxServiceServer
}

func NewGRPCServer() *serverStruct {
	return &serverStruct{}
}

func (s *serverStruct) Healthz(ctx context.Context, req *pb.HealthzRequest) (*pb.HealthzResponse, error) {
	logx.From(ctx).Debug("healthz ok")
	return &pb.HealthzResponse{Status: "ok"}, nil
}

func (s *serverStruct) CalculatePrivate(ctx context.Context, req *pb.CalculatePrivateRequest) (*pb.CalculatePrivateResponse, error) {

	l := logx.From(ctx).With("calc_type", "private")
	start := time.Now()
	l.Info("calc_start")

	input := calculate.FromPrivateRequest(req)

	if err := calculate.ValidateCalculateInput(input); err != nil {
		l.Warn("calc_invalid_arguments", "reason", err.Error())
		return nil, status.Errorf(codes.InvalidArgument, "validate: %v", err)
	}

	months := calculate.CalculateMonthlyTax(input)
	if len(months) == 0 {
		l.Error("calc_no_months_produced")
		return nil, status.Error(codes.Internal, "no data produced")
	}

	resp := &pb.CalculatePrivateResponse{
		MonthlyDetails:        calculate.ToGRPCPrivateResponse(months),
		AnnualTaxAmount:       months[len(months)-1].AnnualTaxAmount,
		AnnualGrossIncome:     months[len(months)-1].AnnualGrossIncome,
		AnnualNetIncome:       months[len(months)-1].AnnualNetIncome,
		GrossSalary:           input.GrossSalary,
		TerritorialMultiplier: &input.TerritorialMultiplier,
		NorthernCoefficient:   &input.NorthernCoefficient,
	}

	l.Info("calc_done",
		"months_count", len(months),
		"duration", time.Since(start).String(),
	)

	return resp, nil
}

func (s *serverStruct) CalculatePublic(ctx context.Context, req *pb.CalculatePublicRequest) (*pb.CalculatePublicResponse, error) {

	l := logx.From(ctx).With("calc_type", "public")
	start := time.Now()
	l.Info("calc_start")

	input := calculate.FromPublicRequest(req)

	if err := calculate.ValidateCalculateInput(input); err != nil {
		l.Warn("calc_invalid_arguments", "reason", err.Error())
		return nil, status.Errorf(codes.InvalidArgument, "validate: %v", err)
	}

	months := calculate.CalculateMonthlyTax(input)
	if len(months) == 0 {
		l.Error("calc_no_months_produced")
		return nil, status.Error(codes.Internal, "no data produced")
	}

	resp := &pb.CalculatePublicResponse{
		MonthlyDetails:        calculate.ToGRPCPublicResponse(months),
		AnnualTaxAmount:       months[len(months)-1].AnnualTaxAmount,
		AnnualGrossIncome:     months[len(months)-1].AnnualGrossIncome,
		AnnualNetIncome:       months[len(months)-1].AnnualNetIncome,
		GrossSalary:           input.GrossSalary,
		TerritorialMultiplier: &input.TerritorialMultiplier,
		NorthernCoefficient:   &input.NorthernCoefficient,
	}

	l.Info("calc_done",
		"months_count", len(months),
		"duration", time.Since(start).String(),
	)

	return resp, nil
}
