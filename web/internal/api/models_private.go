package api

import (
	"fmt"
	"time"

	pb "github.com/kiselevos/new_tax/gen/grpc/api"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type EmploymentType string

const (
	EmploymentTypeTD           EmploymentType = "TD"
	EmploymentTypeGPH          EmploymentType = "GPH"
	EmploymentTypeSelfEmployed EmploymentType = "SELF_EMPLOYED"
)

type NpdIncomeSource string

const (
	NpdIncomeSourceIndividual   NpdIncomeSource = "INDIVIDUAL"
	NpdIncomeSourceLegalEntity  NpdIncomeSource = "LEGAL_ENTITY"
)

// PrivateCalcRequest описывает JSON-запрос к /api/v1/private-calc
type PrivateCalcRequest struct {
	GrossSalary           uint64          `json:"gross_salary"`
	EmploymentType        EmploymentType  `json:"employment_type,omitempty"`
	TerritorialMultiplier *uint64         `json:"territorial_multiplier,omitempty"`
	NorthernCoefficient   *uint64         `json:"northern_coefficient,omitempty"`
	StartDate             *string         `json:"start_date,omitempty"`
	HasTaxPrivilege       *bool           `json:"has_tax_privilege,omitempty"`
	IsNotResident         *bool           `json:"is_not_resident,omitempty"`
	MonthlyBonuses        []uint64        `json:"monthly_bonuses,omitempty"`

	// Налоговые вычеты (ст. 218–220 НК РФ)
	ChildrenCount          *uint32  `json:"children_count,omitempty"`
	DisabledChildrenCount  *uint32  `json:"disabled_children_count,omitempty"`
	HousingExpense         *uint64  `json:"housing_expense,omitempty"`
	MortgageExpense        *uint64  `json:"mortgage_expense,omitempty"`
	SocialExpense          *uint64  `json:"social_expense,omitempty"`
	ChildEduExpense        *uint64  `json:"child_edu_expense,omitempty"`

	// Параметры самозанятого (НПД)
	NpdIncomeSource          NpdIncomeSource `json:"npd_income_source,omitempty"`
	HasRegistrationDeduction *bool           `json:"has_registration_deduction,omitempty"`
}

// DeductionResult описывает итог расчёта налоговых вычетов
type DeductionResult struct {
	ChildrenMonthlyDeduction uint64 `json:"children_monthly_deduction"`
	ChildrenMonths           uint32 `json:"children_months"`
	ChildrenReturn           uint64 `json:"children_return"`
	PropertyReturnThisYear   uint64 `json:"property_return_this_year"`
	PropertyReturnTotal      uint64 `json:"property_return_total"`
	SocialReturn             uint64 `json:"social_return"`
	TotalReturn              uint64 `json:"total_return"`
}

// MonthlyPrivateTax описывает помесячный расчёт в ответе
type MonthlyPrivateTax struct {
	Month              time.Time `json:"month"`
	MonthlyGrossIncome uint64    `json:"monthly_gross_income"`
	MonthlyNetIncome   uint64    `json:"monthly_net_income"`
	MonthlyTaxAmount   uint64    `json:"monthly_tax_amount"`
	TaxRate            uint64    `json:"tax_rate"`

	AnnualGrossIncome uint64 `json:"annual_gross_income"`
	AnnualNetIncome   uint64 `json:"annual_net_income"`
	AnnualTaxAmount   uint64 `json:"annual_tax_amount"`

	MonthlyBaseGrossIncome  uint64 `json:"monthly_base_gross_income"`
	MonthlyNorthGrossIncome uint64 `json:"monthly_north_gross_income"`
	MonthlyBaseTaxAmount    uint64 `json:"monthly_base_tax_amount"`
	MonthlyNorthTaxAmount   uint64 `json:"monthly_north_tax_amount"`

	AnnualBaseGrossIncome  uint64 `json:"annual_base_gross_income"`
	AnnualNorthGrossIncome uint64 `json:"annual_north_gross_income"`
	AnnualBaseTaxAmount    uint64 `json:"annual_base_tax_amount"`
	AnnualNorthTaxAmount   uint64 `json:"annual_north_tax_amount"`

	MonthlyPFR uint64 `json:"monthlyPFR"`
	MonthlyFOMS uint64 `json:"monthlyFOMS"`
	MonthlyFSS uint64 `json:"monthlyFSS"`
	AnnualPFR  uint64 `json:"annualPFR"`
	AnnualFOMS uint64 `json:"annualFOMS"`
	AnnualFSS  uint64 `json:"annualFSS"`

	MonthlyBonus    uint64 `json:"monthly_bonus,omitempty"`
	NpdDeductionUsed uint64 `json:"npd_deduction_used,omitempty"`
}

// PrivateCalcResponse описывает JSON-ответ на /api/v1/private-calc
type PrivateCalcResponse struct {
	GrossSalary           uint64  `json:"gross_salary"`
	TerritorialMultiplier *uint64 `json:"territorial_multiplier,omitempty"`
	NorthernCoefficient   *uint64 `json:"northern_coefficient,omitempty"`

	AnnualGrossIncome uint64 `json:"annual_gross_income"`
	AnnualNetIncome   uint64 `json:"annual_net_income"`
	AnnualTaxAmount   uint64 `json:"annual_tax_amount"`

	AnnualPFR  uint64 `json:"annualPFR"`
	AnnualFOMS uint64 `json:"annualFOMS"`
	AnnualFSS  uint64 `json:"annualFSS"`

	NpdLimitExceeded bool             `json:"npd_limit_exceeded,omitempty"`
	DeductionResult  *DeductionResult `json:"deduction_result,omitempty"`

	MonthlyDetails []MonthlyPrivateTax `json:"monthly_details"`
}

// Validate проверяет поля запроса перед отправкой в бекенд.
func (r *PrivateCalcRequest) Validate() error {
	if r.GrossSalary == 0 {
		return fmt.Errorf("salary must be > 0")
	}
	if r.GrossSalary > 1_000_000_000 {
		return fmt.Errorf("salary must be < 1_000_000_000")
	}

	if r.EmploymentType != "" {
		switch r.EmploymentType {
		case EmploymentTypeTD, EmploymentTypeGPH, EmploymentTypeSelfEmployed:
		default:
			return fmt.Errorf("employment_type must be TD, GPH or SELF_EMPLOYED")
		}
	}

	if r.TerritorialMultiplier != nil {
		v := *r.TerritorialMultiplier
		if v < 100 || v > 200 {
			return fmt.Errorf("territorial_multiplier must be between 100 and 200")
		}
	}

	if r.NorthernCoefficient != nil {
		v := *r.NorthernCoefficient
		if v < 100 || v > 200 {
			return fmt.Errorf("northern_coefficient must be between 100 and 200")
		}
	}

	if r.StartDate != nil {
		if _, err := time.Parse("2006-01-02", *r.StartDate); err != nil {
			return fmt.Errorf("invalid start_date, expected YYYY-MM-DD")
		}
	}

	if len(r.MonthlyBonuses) > 0 && len(r.MonthlyBonuses) != 12 {
		return fmt.Errorf("monthly_bonuses must contain exactly 12 elements")
	}

	if r.NpdIncomeSource != "" {
		switch r.NpdIncomeSource {
		case NpdIncomeSourceIndividual, NpdIncomeSourceLegalEntity:
		default:
			return fmt.Errorf("npd_income_source must be INDIVIDUAL or LEGAL_ENTITY")
		}
	}

	return nil
}

// ToPrivateProto конвертирует DTO в proto-запрос.
func (r *PrivateCalcRequest) ToPrivateProto() *pb.CalculatePrivateRequest {
	req := &pb.CalculatePrivateRequest{
		GrossSalary: r.GrossSalary,
	}

	if r.EmploymentType != "" {
		v := pb.EmploymentType(pb.EmploymentType_value[string(r.EmploymentType)])
		req.EmploymentType = &v
	}

	if r.TerritorialMultiplier != nil {
		v := *r.TerritorialMultiplier
		req.TerritorialMultiplier = &v
	}

	if r.NorthernCoefficient != nil {
		v := *r.NorthernCoefficient
		req.NorthernCoefficient = &v
	}

	if r.StartDate != nil {
		t, _ := time.Parse("2006-01-02", *r.StartDate)
		req.StartDate = timestamppb.New(t)
	}

	if r.HasTaxPrivilege != nil {
		req.HasTaxPrivilege = r.HasTaxPrivilege
	}

	if r.IsNotResident != nil {
		req.IsNotResident = r.IsNotResident
	}

	if len(r.MonthlyBonuses) == 12 {
		req.MonthlyBonuses = r.MonthlyBonuses
	}

	if r.ChildrenCount != nil {
		v := *r.ChildrenCount
		req.ChildrenCount = &v
	}

	if r.DisabledChildrenCount != nil {
		v := *r.DisabledChildrenCount
		req.DisabledChildrenCount = &v
	}

	if r.HousingExpense != nil {
		v := *r.HousingExpense
		req.HousingExpense = &v
	}

	if r.MortgageExpense != nil {
		v := *r.MortgageExpense
		req.MortgageExpense = &v
	}

	if r.SocialExpense != nil {
		v := *r.SocialExpense
		req.SocialExpense = &v
	}

	if r.ChildEduExpense != nil {
		v := *r.ChildEduExpense
		req.ChildEduExpense = &v
	}

	if r.NpdIncomeSource != "" {
		v := pb.NpdIncomeSource(pb.NpdIncomeSource_value[string(r.NpdIncomeSource)])
		req.NpdIncomeSource = &v
	}

	if r.HasRegistrationDeduction != nil {
		req.HasRegistrationDeduction = r.HasRegistrationDeduction
	}

	return req
}

// NewPrivateResponseToJSON конвертирует proto-ответ в JSON DTO.
func NewPrivateResponseToJSON(resp *pb.CalculatePrivateResponse) *PrivateCalcResponse {
	out := &PrivateCalcResponse{
		GrossSalary:           resp.GrossSalary,
		TerritorialMultiplier: resp.TerritorialMultiplier,
		NorthernCoefficient:   resp.NorthernCoefficient,
		AnnualTaxAmount:       resp.AnnualTaxAmount,
		AnnualGrossIncome:     resp.AnnualGrossIncome,
		AnnualNetIncome:       resp.AnnualNetIncome,
		AnnualPFR:             resp.AnnualPFR,
		AnnualFOMS:            resp.AnnualFOMS,
		AnnualFSS:             resp.AnnualFSS,
		NpdLimitExceeded:      resp.NpdLimitExceeded,
	}

	if resp.DeductionResult != nil {
		d := resp.DeductionResult
		out.DeductionResult = &DeductionResult{
			ChildrenMonthlyDeduction: d.ChildrenMonthlyDeduction,
			ChildrenMonths:           d.ChildrenMonths,
			ChildrenReturn:           d.ChildrenReturn,
			PropertyReturnThisYear:   d.PropertyReturnThisYear,
			PropertyReturnTotal:      d.PropertyReturnTotal,
			SocialReturn:             d.SocialReturn,
			TotalReturn:              d.TotalReturn,
		}
	}

	out.MonthlyDetails = make([]MonthlyPrivateTax, 0, len(resp.MonthlyDetails))
	for _, m := range resp.MonthlyDetails {
		out.MonthlyDetails = append(out.MonthlyDetails, MonthlyPrivateTax{
			Month:                   m.GetMonth().AsTime(),
			MonthlyGrossIncome:      m.MonthlyGrossIncome,
			MonthlyNetIncome:        m.MonthlyNetIncome,
			MonthlyTaxAmount:        m.MonthlyTaxAmount,
			TaxRate:                 m.TaxRate,
			AnnualGrossIncome:       m.AnnualGrossIncome,
			AnnualNetIncome:         m.AnnualNetIncome,
			AnnualTaxAmount:         m.AnnualTaxAmount,
			MonthlyBaseGrossIncome:  m.MonthlyBaseGrossIncome,
			MonthlyNorthGrossIncome: m.MonthlyNorthGrossIncome,
			MonthlyBaseTaxAmount:    m.MonthlyBaseTaxAmount,
			MonthlyNorthTaxAmount:   m.MonthlyNorthTaxAmount,
			AnnualBaseGrossIncome:   m.AnnualBaseGrossIncome,
			AnnualNorthGrossIncome:  m.AnnualNorthGrossIncome,
			AnnualBaseTaxAmount:     m.AnnualBaseTaxAmount,
			AnnualNorthTaxAmount:    m.AnnualNorthTaxAmount,
			MonthlyPFR:              m.MonthlyPFR,
			MonthlyFOMS:             m.MonthlyFOMS,
			MonthlyFSS:              m.MonthlyFSS,
			AnnualPFR:               m.AnnualPFR,
			AnnualFOMS:              m.AnnualFOMS,
			AnnualFSS:               m.AnnualFSS,
			MonthlyBonus:            m.MonthlyBonus,
			NpdDeductionUsed:        m.NpdDeductionUsed,
		})
	}

	return out
}
