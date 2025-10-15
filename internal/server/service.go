package server

import (
	"context"
	pb "new_tax/gen/grpc/api"
	"new_tax/internal/calculate"
	"new_tax/pkg/logx"

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
	logx.From(ctx).Info("healthz ok")
	return &pb.HealthzResponse{Status: "ok"}, nil
}

func (s *serverStruct) CalculatePrivate(ctx context.Context, req *pb.CalculatePrivateRequest) (*pb.CalculatePrivateResponse, error) {

	log := logx.From(ctx)

	logx.From(ctx).Info("📨 CalculatePrivate called", "req", req)

	input := calculate.FromPrivateRequest(req)

	if err := calculate.ValidateCalculateInput(input); err != nil {
		log.Info("invalid arguments", "err", err)
		return nil, status.Errorf(codes.InvalidArgument, "validate: %v", err)
	}

	months := calculate.CalculateMonthlyTax(input)

	resp := &pb.CalculatePrivateResponse{
		MonthlyDetails:        calculate.ToGRPCPrivateResponse(months),
		AnnualTaxAmount:       months[len(months)-1].AnnualTaxAmount,
		AnnualGrossIncome:     months[len(months)-1].AnnualGrossIncome,
		AnnualNetIncome:       months[len(months)-1].AnnualNetIncome,
		GrossSalary:           input.GrossSalary,
		TerritorialMultiplier: &input.TerritorialMultiplier,
		NorthernCoefficient:   &input.NorthernCoefficient,
	}
	return resp, nil
}

func (s *serverStruct) CalculatePublic(ctx context.Context, req *pb.CalculatePublicRequest) (*pb.CalculatePublicResponse, error) {

	log := logx.From(ctx)

	input := calculate.FromPublicRequest(req)

	if err := calculate.ValidateCalculateInput(input); err != nil {
		log.Info("invalid arguments", "err", err)
		return nil, status.Errorf(codes.InvalidArgument, "validate: %v", err)
	}

	months := calculate.CalculateMonthlyTax(input)

	resp := &pb.CalculatePublicResponse{
		MonthlyDetails:        calculate.ToGRPCPublicResponse(months),
		AnnualTaxAmount:       months[len(months)-1].AnnualTaxAmount,
		AnnualGrossIncome:     months[len(months)-1].AnnualGrossIncome,
		AnnualNetIncome:       months[len(months)-1].AnnualNetIncome,
		GrossSalary:           input.GrossSalary,
		TerritorialMultiplier: &input.TerritorialMultiplier,
		NorthernCoefficient:   &input.NorthernCoefficient,
	}
	return resp, nil
}
