package calculate

import (
	"time"

	pb "github.com/kiselevos/new_tax/gen/grpc/api"
)

// CalculateInput содержит входные данные для расчёта налога.
// Используется для преобразованных gRPC запросов.
type CalculateInput struct {
	GrossSalary           uint64
	TerritorialMultiplier uint64
	NorthernCoefficient   uint64
	StartDate             time.Time
	HasTaxPrivilege       bool
	IsNotResident         bool
}

// MonthlyTax содержит данные по налогам за конкретный месяц.
type MonthlyTax struct {
	// Название месяца
	Month time.Time

	// Месячные показатели
	MonthlyGrossIncome uint64
	MonthlyNetIncome   uint64
	MonthlyTaxAmount   uint64
	TaxRate            uint64

	// Годовые на текущий месяц
	AnnualGrossIncome uint64
	AnnualNetIncome   uint64
	AnnualTaxAmount   uint64

	// При наличии северной надбавки
	MonthlyNorthGrossIncome uint64
	MonthlyNorthTaxAmount   uint64
	MonthlyBaseGrossIncome  uint64
	MonthlyBaseTaxAmount    uint64

	AnnualNorthGrossIncome uint64
	AnnualNorthTaxAmount   uint64
	AnnualBaseGrossIncome  uint64
	AnnualBaseTaxAmount    uint64
}

// FromPrivateRequest преобразует gRPC-запрос от зарегистрированного пользователя в структуру CalculateInput.
// Устанавливает значения по умолчанию, если поля не указаны.
func FromPrivateRequest(req *pb.CalculatePrivateRequest) CalculateInput {
	var startDate time.Time

	if ts := req.GetStartDate(); ts != nil {
		startDate = ts.AsTime()
	} else {
		startDate = time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	}

	terr := req.GetTerritorialMultiplier()
	if terr == 0 {
		terr = 100
	}

	north := req.GetNorthernCoefficient()
	if north == 0 {
		north = 100
	}

	return CalculateInput{
		GrossSalary:           req.GetGrossSalary(),
		TerritorialMultiplier: terr,
		NorthernCoefficient:   north,
		StartDate:             startDate,
		HasTaxPrivilege:       req.GetHasTaxPrivilege(),
		IsNotResident:         req.GetIsNotResident(),
	}
}

// FromPublicRequest преобразует gRPC-запрос от гостя в структуру CalculateInput.
// Стартовая дата и льготы не указываются, значения по умолчанию.
func FromPublicRequest(req *pb.CalculatePublicRequest) CalculateInput {
	startDate := time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.UTC)

	terr := req.GetTerritorialMultiplier()
	if terr == 0 {
		terr = 100
	}

	north := req.GetNorthernCoefficient()
	if north == 0 {
		north = 100
	}

	return CalculateInput{
		GrossSalary:           req.GetGrossSalary(),
		TerritorialMultiplier: terr,
		NorthernCoefficient:   north,
		StartDate:             startDate,
		HasTaxPrivilege:       false,
		IsNotResident:         false,
	}
}

// ToGRPCPrivateResponse преобразует расчетные данные в gRPC-ответ для зарегестрированного юзера.
func ToGRPCPrivateResponse(monthDetails []MonthlyTax) []*pb.MonthlyPrivateTax {

	var result []*pb.MonthlyPrivateTax

	for _, m := range monthDetails {
		result = append(result, &pb.MonthlyPrivateTax{
			Month:                   ToProtoTimestamp(m.Month),
			MonthlyGrossIncome:      m.MonthlyGrossIncome,
			MonthlyNetIncome:        m.MonthlyNetIncome,
			MonthlyTaxAmount:        m.MonthlyTaxAmount,
			AnnualGrossIncome:       m.AnnualGrossIncome,
			AnnualNetIncome:         m.AnnualNetIncome,
			AnnualTaxAmount:         m.AnnualTaxAmount,
			TaxRate:                 m.TaxRate,
			AnnualBaseGrossIncome:   m.AnnualBaseGrossIncome,
			AnnualNorthGrossIncome:  m.AnnualNorthGrossIncome,
			AnnualBaseTaxAmount:     m.AnnualBaseTaxAmount,
			AnnualNorthTaxAmount:    m.AnnualNorthTaxAmount,
			MonthlyNorthTaxAmount:   m.MonthlyNorthTaxAmount,
			MonthlyBaseGrossIncome:  m.MonthlyBaseGrossIncome,
			MonthlyBaseTaxAmount:    m.MonthlyBaseTaxAmount,
			MonthlyNorthGrossIncome: m.MonthlyNorthGrossIncome,
		})
	}

	return result
}

// ToGRPCPublicResponse преобразует расчетные данные в gRPC-ответ для guest.
func ToGRPCPublicResponse(monthDetails []MonthlyTax) []*pb.MonthlyPublicTax {

	var result []*pb.MonthlyPublicTax

	for _, m := range monthDetails {
		result = append(result, &pb.MonthlyPublicTax{
			Month:              ToProtoTimestamp(m.Month),
			MonthlyGrossIncome: m.MonthlyGrossIncome,
			MonthlyNetIncome:   m.MonthlyNetIncome,
			MonthlyTaxAmount:   m.MonthlyTaxAmount,
			AnnualGrossIncome:  m.AnnualGrossIncome,
			AnnualNetIncome:    m.AnnualNetIncome,
			AnnualTaxAmount:    m.AnnualTaxAmount,
		})
	}

	return result
}
