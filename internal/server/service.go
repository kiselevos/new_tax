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

	input := calculate.FromPrivateRequest(req)

	// --- validation ---
	if err := calculate.ValidateCalculateInput(input); err != nil {
		l.Warn("calc_invalid_arguments",
			"reason", err.Error(),
			"gross_salary", input.GrossSalary,
			"territorial_multiplier", input.TerritorialMultiplier,
			"northern_coefficient", input.NorthernCoefficient,
			"is_not_resident", input.IsNotResident,
			"has_tax_privilege", input.HasTaxPrivilege,
		)
		return nil, status.Errorf(codes.InvalidArgument, "validate: %v", err)
	}

	// --- calculation ---
	months := calculate.CalculateMonthlyTax(input)
	if len(months) == 0 {
		l.Error("calc_no_months_produced")
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

	// --- business success log ---
	l.Info("calc_done",
		"months", len(months),
		"duration_ms", time.Since(start).Milliseconds(),
		"gross_salary", input.GrossSalary,
		"is_not_resident", input.IsNotResident,
		"has_privilege", input.HasTaxPrivilege,
	)

	return resp, nil
}

func (s *serverStruct) CalculatePublic(ctx context.Context, req *pb.CalculatePublicRequest) (*pb.CalculatePublicResponse, error) {
	l := logx.From(ctx).With("calc_type", "public")
	start := time.Now()

	input := calculate.FromPublicRequest(req)

	// --- validation ---
	if err := calculate.ValidateCalculateInput(input); err != nil {
		l.Warn("calc_invalid_arguments",
			"reason", err.Error(),
			"gross_salary", input.GrossSalary,
			"territorial_multiplier", input.TerritorialMultiplier,
			"northern_coefficient", input.NorthernCoefficient,
		)
		return nil, status.Errorf(codes.InvalidArgument, "validate: %v", err)
	}

	// --- business calculation ---
	months := calculate.CalculateMonthlyTax(input)
	if len(months) == 0 {
		l.Error("calc_no_months_produced")
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

	// --- business success log (INFO) ---
	l.Info("calc_done",
		"months", len(months),
		"gross_salary", input.GrossSalary,
	)

	// --- internal performance metric (DEBUG) ---
	l.Debug("calc_duration_ms",
		"duration_ms", time.Since(start).Milliseconds(),
	)

	return resp, nil
}
