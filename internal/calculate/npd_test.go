package calculate

import (
	"testing"
	"time"

	pb "github.com/kiselevos/new_tax/gen/grpc/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func npdJan() time.Time {
	return time.Date(time.Now().UTC().Year(), time.January, 1, 0, 0, 0, 0, time.UTC)
}

// TestNPD_Rate4_Individual проверяет ставку 4% при доходе от физлица, без регистрационного вычета.
func TestNPD_Rate4_Individual(t *testing.T) {
	income := uint64(50_000_00) // 50 000 ₽ в копейках
	months := TaxCalculateForSelfEmployed(income, NpdRateIndividual, false, npdJan(), 1, make([]uint64, 12))

	require.Len(t, months, 12)
	for i, m := range months {
		expectedTax := RoundTaxAmount(income * NpdRateIndividual / 100) // 2 000 ₽
		assert.Equalf(t, expectedTax, m.MonthlyTaxAmount, "месяц %d: неверная сумма НПД", i+1)
		assert.Equal(t, uint64(NpdRateIndividual), m.TaxRate, "ставка должна быть 4%%")
		assert.Zero(t, m.NpdDeductionUsed, "вычет не применялся")
	}
}

// TestNPD_Rate6_Legal проверяет ставку 6% при доходе от юрлица, без регистрационного вычета.
func TestNPD_Rate6_Legal(t *testing.T) {
	income := uint64(50_000_00) // 50 000 ₽
	months := TaxCalculateForSelfEmployed(income, NpdRateLegal, false, npdJan(), 1, make([]uint64, 12))

	require.Len(t, months, 12)
	for i, m := range months {
		expectedTax := RoundTaxAmount(income * NpdRateLegal / 100) // 3 000 ₽
		assert.Equalf(t, expectedTax, m.MonthlyTaxAmount, "месяц %d: неверная сумма НПД", i+1)
		assert.Equal(t, uint64(NpdRateLegal), m.TaxRate, "ставка должна быть 6%%")
	}
}

// TestNPD_RegistrationDeduction_Applied проверяет применение регистрационного вычета 10 000 ₽.
// При малом доходе вычет не исчерпывается за год: эффективная ставка 3% вместо 4%.
func TestNPD_RegistrationDeduction_Applied(t *testing.T) {
	income := uint64(10_000_00) // 10 000 ₽ — потенциальная экономия 100 ₽/мес (1%)
	months := TaxCalculateForSelfEmployed(income, NpdRateIndividual, true, npdJan(), 1, make([]uint64, 12))

	require.Len(t, months, 12)

	// Суммарный использованный вычет за год: 100 ₽ * 12 = 1 200 ₽ < 10 000 ₽ → вычет всегда применяется
	totalDeductionUsed := uint64(0)
	for i, m := range months {
		expectedRaw := income * NpdRateIndividual / 100   // 400_00
		expectedDeduction := income * 1 / 100             // 100_00 (bonusRate = 1%)
		expectedTax := RoundTaxAmount(expectedRaw - expectedDeduction) // 300_00 = 300 ₽

		assert.Equalf(t, expectedTax, m.MonthlyTaxAmount, "месяц %d: неверная сумма НПД с вычетом", i+1)
		assert.Equalf(t, expectedDeduction, m.NpdDeductionUsed, "месяц %d: неверный использованный вычет", i+1)
		totalDeductionUsed += m.NpdDeductionUsed
	}
	// 100 ₽ * 12 месяцев = 1 200 ₽ в копейках = 1 200_00
	assert.Equal(t, uint64(1_200_00), totalDeductionUsed)
}

// TestNPD_RegistrationDeduction_Exhausted_4pct проверяет исчерпание вычета при ставке 4%.
// При доходе 500 000 ₽/мес потенциальная экономия — 5 000 ₽/мес (1%).
// Вычет 10 000 ₽ исчерпывается ровно за 2 месяца.
func TestNPD_RegistrationDeduction_Exhausted_4pct(t *testing.T) {
	income := uint64(500_000_00) // 500 000 ₽
	// bonusRate = 1%, potential = 500_000_00 * 1 / 100 = 500_000 kopecks = 5 000 ₽
	// Месяц 1: используем 500_000, остаток 500_000
	// Месяц 2: используем 500_000, остаток 0
	// Месяц 3+: вычет не применяется

	months := TaxCalculateForSelfEmployed(income, NpdRateIndividual, true, npdJan(), 1, make([]uint64, 12))
	require.Len(t, months, 12)

	potential := income * 1 / 100 // 500_000 kopecks

	// Месяц 1 и 2: вычет применяется
	for i := 0; i < 2; i++ {
		assert.Equalf(t, potential, months[i].NpdDeductionUsed, "месяц %d: вычет должен применяться", i+1) //nolint:govet
	}
	// Месяц 3 и далее: вычет исчерпан
	for i := 2; i < 12; i++ {
		assert.Zerof(t, months[i].NpdDeductionUsed, "месяц %d: вычет должен быть исчерпан", i+1)
		expectedTax := RoundTaxAmount(income * NpdRateIndividual / 100)
		assert.Equalf(t, expectedTax, months[i].MonthlyTaxAmount, "месяц %d: полная ставка после исчерпания", i+1)
	}
	// Суммарный вычет = NpdRegistrationDeduction
	totalUsed := uint64(0)
	for _, m := range months {
		totalUsed += m.NpdDeductionUsed
	}
	assert.Equal(t, uint64(NpdRegistrationDeduction), totalUsed)
}

// TestNPD_RegistrationDeduction_Exhausted_6pct проверяет исчерпание вычета при ставке 6%.
// bonusRate = 2%, potential = 500 000 * 2% = 10 000 ₽ → вычет исчерпывается в первый же месяц.
func TestNPD_RegistrationDeduction_Exhausted_6pct(t *testing.T) {
	income := uint64(500_000_00)
	// bonusRate = 2%, potential = 500_000_00 * 2 / 100 = 1_000_000 kopecks = 10 000 ₽
	// Месяц 1: potential == deductionLeft → используем весь вычет, остаток 0
	// Месяц 2+: вычет не применяется

	months := TaxCalculateForSelfEmployed(income, NpdRateLegal, true, npdJan(), 1, make([]uint64, 12))
	require.Len(t, months, 12)

	// Месяц 1: весь вычет
	assert.Equal(t, uint64(NpdRegistrationDeduction), months[0].NpdDeductionUsed)

	// Месяц 2+: нет вычета, полная ставка
	for i := 1; i < 12; i++ {
		assert.Zerof(t, months[i].NpdDeductionUsed, "месяц %d: вычет исчерпан", i+1)
		expectedTax := RoundTaxAmount(income * NpdRateLegal / 100)
		assert.Equalf(t, expectedTax, months[i].MonthlyTaxAmount, "месяц %d: полная ставка", i+1)
	}
}

// TestNPD_NoContributions проверяет, что при НПД взносы ПФР/ФОМС/ФСС не начисляются.
func TestNPD_NoContributions(t *testing.T) {
	months := CalculateMonthlyTax(CalculateInput{
		GrossSalary:           100_000_00,
		TerritorialMultiplier: 100,
		NorthernCoefficient:   100,
		StartDate:             npdJan(),
		EmploymentType:        pb.EmploymentType_SELF_EMPLOYED,
		NpdIncomeSource:       pb.NpdIncomeSource_LEGAL_ENTITY,
	})

	require.Len(t, months, 12)
	for i, m := range months {
		assert.Zerof(t, m.MonthlyPFR, "месяц %d: ПФР должен быть 0 при НПД", i+1)
		assert.Zerof(t, m.MonthlyFOMS, "месяц %d: ФОМС должен быть 0 при НПД", i+1)
		assert.Zerof(t, m.MonthlyFSS, "месяц %d: ФСС должен быть 0 при НПД", i+1)
	}
	assert.Zero(t, months[11].AnnualPFR)
	assert.Zero(t, months[11].AnnualFOMS)
	assert.Zero(t, months[11].AnnualFSS)
}

// TestNPD_Via_CalculateMonthlyTax_Individual проверяет ставку 4% через полный пайплайн.
func TestNPD_Via_CalculateMonthlyTax_Individual(t *testing.T) {
	months := CalculateMonthlyTax(CalculateInput{
		GrossSalary:           50_000_00,
		TerritorialMultiplier: 100,
		NorthernCoefficient:   100,
		StartDate:             npdJan(),
		EmploymentType:        pb.EmploymentType_SELF_EMPLOYED,
		NpdIncomeSource:       pb.NpdIncomeSource_INDIVIDUAL,
	})

	require.Len(t, months, 12)
	expectedTax := RoundTaxAmount(uint64(50_000_00) * NpdRateIndividual / 100)
	for i, m := range months {
		assert.Equalf(t, uint64(NpdRateIndividual), m.TaxRate, "месяц %d: ставка должна быть 4%%", i+1)
		assert.Equalf(t, expectedTax, m.MonthlyTaxAmount, "месяц %d: сумма НПД", i+1)
	}
}

// TestNPD_AnnualAccumulation проверяет корректный накопительный итог за год.
func TestNPD_AnnualAccumulation(t *testing.T) {
	income := uint64(30_000_00) // 30 000 ₽
	months := TaxCalculateForSelfEmployed(income, NpdRateIndividual, false, npdJan(), 1, make([]uint64, 12))

	require.Len(t, months, 12)

	monthlyTax := RoundTaxAmount(income * NpdRateIndividual / 100)

	for i, m := range months {
		expectedAnnualGross := income * uint64(i+1)
		expectedAnnualTax := monthlyTax * uint64(i+1)
		expectedAnnualNet := expectedAnnualGross - expectedAnnualTax

		assert.Equalf(t, expectedAnnualGross, m.AnnualGrossIncome, "месяц %d: накопленный доход", i+1)
		assert.Equalf(t, expectedAnnualTax, m.AnnualTaxAmount, "месяц %d: накопленный налог", i+1)
		assert.Equalf(t, expectedAnnualNet, m.AnnualNetIncome, "месяц %d: накопленный чистый доход", i+1)
	}
}

// TestNPD_StartFromMarch проверяет расчёт с марта — должно быть 10 месяцев.
func TestNPD_StartFromMarch(t *testing.T) {
	startDate := time.Date(time.Now().UTC().Year(), time.March, 1, 0, 0, 0, 0, time.UTC)
	months := TaxCalculateForSelfEmployed(20_000_00, NpdRateIndividual, false, startDate, 3, make([]uint64, 12))

	assert.Len(t, months, 10) // март–декабрь
	assert.Equal(t, time.March, months[0].Month.Month())
	assert.Equal(t, time.December, months[9].Month.Month())
}
