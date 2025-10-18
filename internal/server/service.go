package server

import (
	"context"
	"errors"
	"net/http"

	pb "new_tax/gen/grpc/api"
	taxconnect "new_tax/gen/grpc/api/taxconnect"
	"new_tax/internal/calculate"
	"new_tax/pkg/logx"

	"connectrpc.com/connect"
)

// taxServiceServer реализует интерфейс taxconnect.TaxServiceHandler
type taxServiceServer struct{}

// NewTaxServiceHandler — адаптер для регистрации в mux.Handle(...)
func NewTaxServiceHandler() (string, http.Handler) {
	svc := &taxServiceServer{}
	return taxconnect.NewTaxServiceHandler(svc)
}

// Healthz
func (s *taxServiceServer) Healthz(
	ctx context.Context,
	req *connect.Request[pb.HealthzRequest],
) (*connect.Response[pb.HealthzResponse], error) {
	logx.From(ctx).Debug("Health check passed")
	return connect.NewResponse(&pb.HealthzResponse{Status: "ok"}), nil
}

// CalculatePrivate
func (s *taxServiceServer) CalculatePrivate(
	ctx context.Context,
	req *connect.Request[pb.CalculatePrivateRequest],
) (*connect.Response[pb.CalculatePrivateResponse], error) {

	log := logx.From(ctx)
	input := calculate.FromPrivateRequest(req.Msg)
	log.Debug("Starting private tax calculation", "gross_salary", input.GrossSalary)

	if err := calculate.ValidateCalculateInput(input); err != nil {
		log.Warn("Validation failed", "err", err)
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	months := calculate.CalculateMonthlyTax(input)
	if len(months) == 0 {
		err := errors.New("no months calculated")
		log.Error("Calculation failed", "err", err)
		return nil, connect.NewError(connect.CodeInternal, err)
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

	log.Info("Private tax calculated",
		"annual_tax", resp.AnnualTaxAmount,
		"months", len(months),
	)

	return connect.NewResponse(resp), nil
}

// CalculatePublic
func (s *taxServiceServer) CalculatePublic(
	ctx context.Context,
	req *connect.Request[pb.CalculatePublicRequest],
) (*connect.Response[pb.CalculatePublicResponse], error) {

	log := logx.From(ctx)
	input := calculate.FromPublicRequest(req.Msg)
	log.Debug("Starting public tax calculation", "gross_salary", input.GrossSalary)

	if err := calculate.ValidateCalculateInput(input); err != nil {
		log.Warn("Validation failed", "err", err)
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	months := calculate.CalculateMonthlyTax(input)
	if len(months) == 0 {
		err := errors.New("no months calculated")
		log.Error("Calculation failed", "err", err)
		return nil, connect.NewError(connect.CodeInternal, err)
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

	log.Info("Public tax calculated",
		"annual_tax", resp.AnnualTaxAmount,
		"months", len(months),
	)

	return connect.NewResponse(resp), nil
}
