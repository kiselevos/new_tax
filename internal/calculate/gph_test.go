package calculate

import (
	"testing"
	"time"

	pb "github.com/kiselevos/new_tax/gen/grpc/api"
	"github.com/stretchr/testify/assert"
)

func gpbJan() time.Time {
	return time.Date(time.Now().UTC().Year(), time.January, 1, 0, 0, 0, 0, time.UTC)
}

// TestGPH_FSS_IsZero проверяет, что при ГПХ ФСС равен нулю во всех месяцах.
func TestGPH_FSS_IsZero(t *testing.T) {
	months := CalculateMonthlyTax(CalculateInput{
		GrossSalary:           100_000_00, // 100 000 ₽
		TerritorialMultiplier: 100,
		NorthernCoefficient:   100,
		StartDate:             gpbJan(),
		EmploymentType:        pb.EmploymentType_GPH,
	})

	assert.Len(t, months, 12)
	for _, m := range months {
		assert.Zero(t, m.MonthlyFSS, "ФСС должен быть 0 при ГПХ, месяц %v", m.Month)
	}
	// Накопленный ФСС за год тоже должен быть 0
	assert.Zero(t, months[len(months)-1].AnnualFSS)
}

// TestTD_FSS_IsPositive проверяет, что при ТД ФСС начисляется (> 0).
func TestTD_FSS_IsPositive(t *testing.T) {
	months := CalculateMonthlyTax(CalculateInput{
		GrossSalary:           100_000_00,
		TerritorialMultiplier: 100,
		NorthernCoefficient:   100,
		StartDate:             gpbJan(),
		EmploymentType:        pb.EmploymentType_TD,
	})

	assert.Len(t, months, 12)
	// Хотя бы в первых месяцах ФСС должен быть положительным (до лимита базы ФСС)
	assert.Positive(t, months[0].MonthlyFSS, "ФСС должен быть > 0 при ТД")
}

// TestGPH_NDFLSameAsTD проверяет, что НДФЛ при ГПХ совпадает с ТД (та же шкала).
func TestGPH_NDFLSameAsTD(t *testing.T) {
	input := CalculateInput{
		GrossSalary:           200_000_00, // 200 000 ₽
		TerritorialMultiplier: 100,
		NorthernCoefficient:   100,
		StartDate:             gpbJan(),
	}

	inputGPH := input
	inputGPH.EmploymentType = pb.EmploymentType_GPH

	inputTD := input
	inputTD.EmploymentType = pb.EmploymentType_TD

	monthsGPH := CalculateMonthlyTax(inputGPH)
	monthsTD := CalculateMonthlyTax(inputTD)

	require := assert.New(t)
	require.Len(monthsGPH, 12)
	require.Len(monthsTD, 12)

	for i := range monthsGPH {
		assert.Equalf(t, monthsTD[i].MonthlyTaxAmount, monthsGPH[i].MonthlyTaxAmount,
			"НДФЛ должен совпадать у ТД и ГПХ в месяц %d", i+1)
		assert.Equalf(t, monthsTD[i].MonthlyNetIncome, monthsGPH[i].MonthlyNetIncome,
			"Чистый доход должен совпадать у ТД и ГПХ в месяц %d", i+1)
	}
}

// TestGPH_PFR_FOMS_Charged проверяет, что ПФР и ФОМС при ГПХ всё равно начисляются.
func TestGPH_PFR_FOMS_Charged(t *testing.T) {
	months := CalculateMonthlyTax(CalculateInput{
		GrossSalary:           100_000_00,
		TerritorialMultiplier: 100,
		NorthernCoefficient:   100,
		StartDate:             gpbJan(),
		EmploymentType:        pb.EmploymentType_GPH,
	})

	assert.Positive(t, months[0].MonthlyPFR, "ПФР должен начисляться при ГПХ")
	assert.Positive(t, months[0].MonthlyFOMS, "ФОМС должен начисляться при ГПХ")
	assert.Zero(t, months[0].MonthlyFSS, "ФСС не должен начисляться при ГПХ")
}
