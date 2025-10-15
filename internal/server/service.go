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
	logx.From(ctx).Info("healthz ok")
	return connect.NewResponse(&pb.HealthzResponse{Status: "ok"}), nil
}

// CalculatePrivate
func (s *taxServiceServer) CalculatePrivate(
	ctx context.Context,
	req *connect.Request[pb.CalculatePrivateRequest],
) (*connect.Response[pb.CalculatePrivateResponse], error) {

	log := logx.From(ctx)
	log.Info("📨 CalculatePrivate called", "req", req)

	input := calculate.FromPrivateRequest(req.Msg)

	if err := calculate.ValidateCalculateInput(input); err != nil {
		log.Info("invalid arguments", "err", err)
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	months := calculate.CalculateMonthlyTax(input)
	if len(months) == 0 {
		return nil, connect.NewError(connect.CodeInternal, errors.New("no months calculated"))
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

	return connect.NewResponse(resp), nil
}

// CalculatePublic
func (s *taxServiceServer) CalculatePublic(
	ctx context.Context,
	req *connect.Request[pb.CalculatePublicRequest],
) (*connect.Response[pb.CalculatePublicResponse], error) {

	log := logx.From(ctx)
	log.Info("📨 CalculatePublic called", "req", req)

	input := calculate.FromPublicRequest(req.Msg)

	if err := calculate.ValidateCalculateInput(input); err != nil {
		log.Info("invalid arguments", "err", err)
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	months := calculate.CalculateMonthlyTax(input)
	if len(months) == 0 {
		return nil, connect.NewError(connect.CodeInternal, errors.New("no months calculated"))
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

	return connect.NewResponse(resp), nil
}
