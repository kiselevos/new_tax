package api

import (
	"fmt"
	"time"

	pb "github.com/kiselevos/new_tax/gen/grpc/api"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// PublicCalcRequest описывает JSON-запрос к /api/v1/privatecalc
type PrivateCalcRequest struct {
	GrossSalary           uint64  `json:"gross_salary"`                     // оклад
	TerritorialMultiplier *uint64 `json:"territorial_multiplier,omitempty"` // тер. коэффициент (например 115 = 1.15)
	NorthernCoefficient   *uint64 `json:"northern_coefficient,omitempty"`   // северная надбавка (например 50 = 0.5)
	StartDate             *string `json:"start_date,omitempty"`             // Месяц начала рассчета
	HasTaxPrivilege       *bool   `json:"has_tax_privilege,omitempty"`      // Льготы для силовых структур
	IsNotResident         *bool   `json:"is_not_resident,omitempty"`        // Статус налогового нерезидент
}

// MonthlyPrivateTax описыает JSON ответ по месяцым
type MonthlyPrivateTax struct {
	Month              time.Time `json:"month"`
	MonthlyGrossIncome uint64    `json:"monthly_gross_income"`
	MonthlyNetIncome   uint64    `json:"monthly_net_income"`
	MonthlyTaxAmount   uint64    `json:"monthly_tax_amount"`

	AnnualGrossIncome uint64 `json:"annual_gross_income"`
	AnnualNetIncome   uint64 `json:"annual_net_income"`
	AnnualTaxAmount   uint64 `json:"annual_tax_amount"`

	TaxRate uint64 `json:"tax_rate"` // Процент налога на данный момент

	MonthlyNorthGrossIncome uint64 `json:"monthly_north_gross_income"` // Гросс доход северной надбавки
	MonthlyBaseGrossIncome  uint64 `json:"monthly_base_gross_income"`  // Гросс доход оклад + территориальный коэффициент

	MonthlyNorthTaxAmount uint64 `json:"monthly_north_tax_amount"` // Месячный налог с северной надбавки
	MonthlyBaseTaxAmount  uint64 `json:"monthly_base_tax_amount"`  // Месячный налог с оклада + территориального коэффициента

	AnnualNorthGrossIncome uint64 `json:"annual_north_gross_income"` // Доход с начала года северной надбавки
	AnnualBaseGrossIncome  uint64 `json:"annual_base_gross_income"`  // Доход с начала года оклад + территориальный коэффициент

	AnnualNorthTaxAmount uint64 `json:"annual_north_tax_amount"` // Налог с начала года с северной надбавки
	AnnualBaseTaxAmount  uint64 `json:"annual_base_tax_amount"`  // Налог с начала года с оклада + территориального коэффициента

	MonthlyPFR  uint64 `json:"monthlyPFR"`  //Месячный налог с работодателя в Пенсионный фонд России
	MonthlyFOMS uint64 `json:"monthlyFOMS"` //Месячный налог с работодателя в Фонд Обязательного медиционского страхования
	MonthlyFSS  uint64 `json:"monthlyFSS"`  //Месячный налог с работодателя в Фонд социального страхования (Больничные, декреты)

	AnnualPFR  uint64 `json:"annualPFR"`  // Налог с начала выбранного периода работодателя в Пенсионный фонд России
	AnnualFOMS uint64 `json:"annualFOMS"` //Налог с начала выбранного периода работодателя в Фонд Обязательного медиционского страхования
	AnnualFSS  uint64 `json:"annualFSS"`  //Налог с начала выбранного периода работодателя в Фонд социального страхования (Больничные, декреты)
}

// PublicCalcResponse описывает JSON ответ
type PrivateCalcResponse struct {
	MonthlyDetails []MonthlyPrivateTax `json:"monthly_details"`

	AnnualTaxAmount   uint64 `json:"annual_tax_amount"`
	AnnualGrossIncome uint64 `json:"annual_gross_income"`
	AnnualNetIncome   uint64 `json:"annual_net_income"`

	GrossSalary           uint64  `json:"gross_salary"`
	TerritorialMultiplier *uint64 `json:"territorial_multiplier,omitempty"`
	NorthernCoefficient   *uint64 `json:"northern_coefficient,omitempty"`

	AnnualPFR  uint64 `json:"annualPFR"`  // Налог с начала выбранного периода работодателя в Пенсионный фонд России
	AnnualFOMS uint64 `json:"annualFOMS"` //Налог с начала выбранного периода работодателя в Фонд Обязательного медиционского страхования
	AnnualFSS  uint64 `json:"annualFSS"`  //Налог с начала выбранного периода работодателя в Фонд социального страхования (Больничные, декреты)
}

// Validate DTO
func (r *PrivateCalcRequest) Validate() error {

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

	if r.StartDate != nil {
		_, err := time.Parse("2006-01-02", *r.StartDate)
		if err != nil {
			return fmt.Errorf("invalid start_date")
		}
	}

	return nil
}

// ToProto convert DTO -> gRPC
func (r *PrivateCalcRequest) ToPrivateProto() *pb.CalculatePrivateRequest {

	req := &pb.CalculatePrivateRequest{
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

	return req
}

// NewPublicResponseToJSON сonvert public gRPC -> DTO
func NewPrivateResponseToJSON(resp *pb.CalculatePrivateResponse) *PrivateCalcResponse {

	out := &PrivateCalcResponse{
		AnnualTaxAmount:       resp.AnnualTaxAmount,
		AnnualGrossIncome:     resp.AnnualGrossIncome,
		AnnualNetIncome:       resp.AnnualNetIncome,
		GrossSalary:           resp.GrossSalary,
		TerritorialMultiplier: resp.TerritorialMultiplier,
		NorthernCoefficient:   resp.NorthernCoefficient,
		AnnualPFR:             resp.AnnualPFR,
		AnnualFOMS:            resp.AnnualFOMS,
		AnnualFSS:             resp.AnnualFSS,
	}

	out.MonthlyDetails = make([]MonthlyPrivateTax, 0, len(resp.MonthlyDetails))
	for _, m := range resp.MonthlyDetails {
		out.MonthlyDetails = append(out.MonthlyDetails, MonthlyPrivateTax{
			Month:                   m.GetMonth().AsTime(),
			MonthlyGrossIncome:      m.MonthlyGrossIncome,
			MonthlyNetIncome:        m.MonthlyNetIncome,
			MonthlyTaxAmount:        m.MonthlyTaxAmount,
			AnnualGrossIncome:       m.AnnualGrossIncome,
			AnnualNetIncome:         m.AnnualNetIncome,
			AnnualTaxAmount:         m.AnnualTaxAmount,
			TaxRate:                 m.TaxRate,
			MonthlyNorthGrossIncome: m.MonthlyNorthGrossIncome,
			MonthlyBaseGrossIncome:  m.MonthlyBaseGrossIncome,
			MonthlyNorthTaxAmount:   m.MonthlyNorthTaxAmount,
			MonthlyBaseTaxAmount:    m.MonthlyBaseTaxAmount,
			AnnualNorthGrossIncome:  m.AnnualNorthGrossIncome,
			AnnualBaseGrossIncome:   m.AnnualBaseGrossIncome,
			AnnualNorthTaxAmount:    m.AnnualNorthTaxAmount,
			AnnualBaseTaxAmount:     m.AnnualBaseTaxAmount,
			MonthlyPFR:              m.MonthlyPFR,
			MonthlyFOMS:             m.MonthlyFOMS,
			MonthlyFSS:              m.MonthlyFSS,
			AnnualPFR:               m.AnnualPFR,
			AnnualFOMS:              m.AnnualFOMS,
			AnnualFSS:               m.AnnualFSS,
		})
	}

	return out
}
