package api

import (
	"fmt"
	"time"

	pb "github.com/kiselevos/new_tax/gen/grpc/api"
)

// PublicCalcRequest описывает JSON-запрос к /api/v1/calc
type PublicCalcRequest struct {
	GrossSalary           uint64  `json:"gross_salary"`                     // оклад
	TerritorialMultiplier *uint64 `json:"territorial_multiplier,omitempty"` // тер. коэффициент (например 115 = 1.15)
	NorthernCoefficient   *uint64 `json:"northern_coefficient,omitempty"`   // северная надбавка (например 50 = 0.5)
}

// MonthlyPublicTax описыает JSON ответ по месяцым
type MonthlyPublicTax struct {
	Month              time.Time `json:"month"`
	MonthlyGrossIncome uint64    `json:"monthly_gross_income"`
	MonthlyNetIncome   uint64    `json:"monthly_net_income"`
	MonthlyTaxAmount   uint64    `json:"monthly_tax_amount"`

	AnnualGrossIncome uint64 `json:"annual_gross_income"`
	AnnualNetIncome   uint64 `json:"annual_net_income"`
	AnnualTaxAmount   uint64 `json:"annual_tax_amount"`
}

// PublicCalcResponse описывает JSON ответ
type PublicCalcResponse struct {
	MonthlyDetails []MonthlyPublicTax `json:"monthly_details"`

	AnnualTaxAmount   uint64 `json:"annual_tax_amount"`
	AnnualGrossIncome uint64 `json:"annual_gross_income"`
	AnnualNetIncome   uint64 `json:"annual_net_income"`

	GrossSalary           uint64  `json:"gross_salary"`
	TerritorialMultiplier *uint64 `json:"territorial_multiplier,omitempty"`
	NorthernCoefficient   *uint64 `json:"northern_coefficient,omitempty"`
}

// Validate DTO
func (r *PublicCalcRequest) Validate() error {

	if r.GrossSalary <= 0 {
		return fmt.Errorf("salary must be > 0")
	}

	if r.GrossSalary > 1_000_000_000 {
		return fmt.Errorf("salary must be < 1_000_000_000")
	}

	if r.TerritorialMultiplier != nil {
		v := *r.TerritorialMultiplier
		if v < 100 || v > 200 {
			return fmt.Errorf("territorial multiplier must be between 100 and 200")
		}
	}

	if r.NorthernCoefficient != nil {
		v := *r.NorthernCoefficient
		if v < 100 || v > 200 {
			return fmt.Errorf("northern coefficient must be between 100 and 200")
		}
	}

	return nil
}

// ToProto convert DTO -> gRPC
func (r *PublicCalcRequest) ToProto() *pb.CalculatePublicRequest {

	req := &pb.CalculatePublicRequest{
		GrossSalary: r.GrossSalary,
	}

	if r.TerritorialMultiplier != nil {
		v := *r.TerritorialMultiplier
		req.TerritorialMultiplier = &v
	}

	if r.NorthernCoefficient != nil {
		v := *r.NorthernCoefficient
		req.NorthernCoefficient = &v
	}

	return req
}

// NewPublicResponseToJSON сonvert public gRPC -> DTO
func NewPublicResponseToJSON(resp *pb.CalculatePublicResponse) *PublicCalcResponse {

	out := &PublicCalcResponse{
		AnnualTaxAmount:       resp.AnnualTaxAmount,
		AnnualGrossIncome:     resp.AnnualGrossIncome,
		AnnualNetIncome:       resp.AnnualNetIncome,
		GrossSalary:           resp.GrossSalary,
		TerritorialMultiplier: resp.TerritorialMultiplier,
		NorthernCoefficient:   resp.NorthernCoefficient,
	}

	out.MonthlyDetails = make([]MonthlyPublicTax, 0, len(resp.MonthlyDetails))
	for _, m := range resp.MonthlyDetails {
		out.MonthlyDetails = append(out.MonthlyDetails, MonthlyPublicTax{
			Month:              m.GetMonth().AsTime(),
			MonthlyGrossIncome: m.MonthlyGrossIncome,
			MonthlyNetIncome:   m.MonthlyNetIncome,
			MonthlyTaxAmount:   m.MonthlyTaxAmount,
			AnnualGrossIncome:  m.AnnualGrossIncome,
			AnnualNetIncome:    m.AnnualNetIncome,
			AnnualTaxAmount:    m.AnnualTaxAmount,
		})
	}

	return out
}
