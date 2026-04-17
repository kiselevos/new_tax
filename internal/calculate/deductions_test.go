package calculate

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// helpers

func makeMonths(grossMonthly uint64, startMonth int) []MonthlyTax {
	input := CalculateInput{
		GrossSalary:           grossMonthly,
		TerritorialMultiplier: 100,
		NorthernCoefficient:   100,
		StartDate:             time.Date(time.Now().Year(), time.Month(startMonth), 1, 0, 0, 0, 0, time.UTC),
	}
	return CalculateMonthlyTax(input)
}

// ===== childrenMonthlyDeduction =====

func TestChildrenMonthlyDeduction_Zero(t *testing.T) {
	assert.Equal(t, uint64(0), childrenMonthlyDeduction(0, 0))
}

func TestChildrenMonthlyDeduction_OneChild(t *testing.T) {
	// 1 ребёнок = 1 400 ₽ = 140 000 копеек
	assert.Equal(t, uint64(140_000), childrenMonthlyDeduction(1, 0))
}

func TestChildrenMonthlyDeduction_TwoChildren(t *testing.T) {
	// 1 400 + 2 800 = 4 200 ₽
	assert.Equal(t, uint64(420_000), childrenMonthlyDeduction(2, 0))
}

func TestChildrenMonthlyDeduction_ThreeChildren(t *testing.T) {
	// 1 400 + 2 800 + 6 000 = 10 200 ₽
	assert.Equal(t, uint64(1_020_000), childrenMonthlyDeduction(3, 0))
}

func TestChildrenMonthlyDeduction_WithDisabled(t *testing.T) {
	// 1 ребёнок + 1 инвалид = 1 400 + 12 000 = 13 400 ₽
	assert.Equal(t, uint64(1_340_000), childrenMonthlyDeduction(1, 1))
}

func TestChildrenMonthlyDeduction_TwoDisabled(t *testing.T) {
	// 2 ребёнка: 1 400 + 2 800 = 4 200. Оба инвалиды: +12 000×2 = 24 000. Итого 28 200 ₽
	assert.Equal(t, uint64(2_820_000), childrenMonthlyDeduction(2, 2))
}

func TestChildrenMonthlyDeduction_DisabledOnly(t *testing.T) {
	// disabled без children — inвалид без указания обычных детей
	// disabledCount > 0 без count = 0 → только надбавка
	assert.Equal(t, uint64(1_200_000), childrenMonthlyDeduction(0, 1))
}

// ===== CalcDeductions: nil cases =====

func TestCalcDeductions_NilWhenNoInput(t *testing.T) {
	months := makeMonths(100_000_00, 1) // 100 000 ₽/мес
	result := CalcDeductions(DeductionInput{}, months)
	assert.Nil(t, result)
}

func TestCalcDeductions_NilWhenNoMonths(t *testing.T) {
	input := DeductionInput{ChildrenCount: 1}
	result := CalcDeductions(input, nil)
	assert.Nil(t, result)
}

// ===== CalcDeductions: дети (ст. 218) =====

func TestCalcDeductions_Children_FullYear(t *testing.T) {
	// 30 000 ₽/мес → годовой доход 360 000 ₽ < 450 000 ₽ → вычет действует 12 мес.
	months := makeMonths(3_000_000, 1) // 30 000 ₽/мес с января
	input := DeductionInput{ChildrenCount: 1}

	result := CalcDeductions(input, months)
	require.NotNil(t, result)

	// Ежемесячный вычет = 1 400 ₽
	assert.Equal(t, uint64(140_000), result.ChildrenMonthlyDeduction)
	// 12 месяцев
	assert.Equal(t, uint32(12), result.ChildrenMonths)
	// Возврат = 12 × 1 400 × 13% = 2 184 ₽ = 218 400 копеек
	assert.Equal(t, uint64(218_400), result.ChildrenReturn)
}

func TestCalcDeductions_Children_IncomeExceedsLimit(t *testing.T) {
	// 200 000 ₽/мес → к марту доход превысит 450 000 ₽
	// Янв: 200К (итого 200К < 450К → применяется)
	// Фев: 200К (итого 400К < 450К → применяется)
	// Мар: 200К (итого 600К > 450К, начало месяца 400К < 450К → применяется)
	// Апр: начало месяца 600К > 450К → НЕ применяется
	// Итого 3 месяца
	months := makeMonths(20_000_000, 1) // 200 000 ₽/мес
	input := DeductionInput{ChildrenCount: 2}

	result := CalcDeductions(input, months)
	require.NotNil(t, result)
	assert.Equal(t, uint32(3), result.ChildrenMonths)
}

// ===== CalcDeductions: имущество (ст. 220) =====

func TestCalcDeductions_Property_Housing(t *testing.T) {
	// Оклад 200 000 ₽/мес → годовой налог ≈ 312 000 ₽ (13%)
	months := makeMonths(20_000_000, 1)
	input := DeductionInput{
		HousingExpense: 150_000_000, // 1 500 000 ₽ (ниже лимита 2 000 000)
	}

	result := CalcDeductions(input, months)
	require.NotNil(t, result)

	// Возврат = 1 500 000 × 13% = 195 000 ₽ = 19 500 000 копеек
	assert.Equal(t, uint64(19_500_000), result.PropertyReturnTotal)
	// Этот год = min(19 500 000, annualTax)
	assert.LessOrEqual(t, result.PropertyReturnThisYear, result.PropertyReturnTotal)
}

func TestCalcDeductions_Property_HousingCappedAtLimit(t *testing.T) {
	// Расходы выше лимита 2 000 000 ₽ → лимит применяется.
	// 300 000 ₽/мес × 12 = 3 600 000 ₽ → маржинальная ставка 15%.
	// Возврат = 2 000 000 × 15% = 300 000 ₽ = 30 000 000 копеек.
	months := makeMonths(30_000_000, 1) // 300 000 ₽/мес
	input := DeductionInput{
		HousingExpense: 500_000_000, // 5 000 000 ₽ >> лимит
	}

	result := CalcDeductions(input, months)
	require.NotNil(t, result)

	assert.Equal(t, uint64(30_000_000), result.PropertyReturnTotal)
}

func TestCalcDeductions_Property_MortgageCap(t *testing.T) {
	// Ипотека выше лимита 3 000 000 ₽ → лимит применяется.
	// 400 000 ₽/мес × 12 = 4 800 000 ₽ → маржинальная ставка 15%.
	// Возврат = 3 000 000 × 15% = 450 000 ₽ = 45 000 000 копеек.
	months := makeMonths(40_000_000, 1) // 400 000 ₽/мес
	input := DeductionInput{
		MortgageExpense: 500_000_000, // 5 000 000 ₽ >> лимит
	}

	result := CalcDeductions(input, months)
	require.NotNil(t, result)

	assert.Equal(t, uint64(45_000_000), result.PropertyReturnTotal)
}

func TestCalcDeductions_Property_ThisYearLimitedByTax(t *testing.T) {
	// Маленький оклад → мало налога → thisYear < total
	months := makeMonths(5_000_000, 1) // 50 000 ₽/мес → ~78 000 ₽ налога за год
	input := DeductionInput{
		HousingExpense: 200_000_000, // 2 000 000 ₽ → total 260 000 ₽
	}

	result := CalcDeductions(input, months)
	require.NotNil(t, result)

	assert.Equal(t, uint64(26_000_000), result.PropertyReturnTotal)
	assert.Less(t, result.PropertyReturnThisYear, result.PropertyReturnTotal)
	assert.LessOrEqual(t, result.PropertyReturnThisYear, result.PropertyReturnTotal)
}

// ===== CalcDeductions: социальный (ст. 219) =====

func TestCalcDeductions_Social_WithinLimit(t *testing.T) {
	months := makeMonths(20_000_000, 1) // 200 000 ₽/мес
	input := DeductionInput{
		SocialExpense: 10_000_000, // 100 000 ₽ (ниже лимита 150 000)
	}

	result := CalcDeductions(input, months)
	require.NotNil(t, result)

	// 100 000 × 13% = 13 000 ₽ = 1 300 000 копеек
	assert.Equal(t, uint64(1_300_000), result.SocialReturn)
}

func TestCalcDeductions_Social_LimitCapped(t *testing.T) {
	months := makeMonths(20_000_000, 1)
	input := DeductionInput{
		SocialExpense: 30_000_000, // 300 000 ₽ >> лимит 150 000
	}

	result := CalcDeductions(input, months)
	require.NotNil(t, result)

	// Лимит 150 000 × 13% = 19 500 ₽ = 1 950 000 копеек
	assert.Equal(t, uint64(1_950_000), result.SocialReturn)
}

func TestCalcDeductions_Social_ChildEdu(t *testing.T) {
	months := makeMonths(20_000_000, 1)
	input := DeductionInput{
		ChildEduExpense: 8_000_000, // 80 000 ₽ (ниже лимита 110 000)
	}

	result := CalcDeductions(input, months)
	require.NotNil(t, result)

	// 80 000 × 13% = 10 400 ₽ = 1 040 000 копеек
	assert.Equal(t, uint64(1_040_000), result.SocialReturn)
}

// ===== CalcDeductions: total не превышает уплаченный НДФЛ =====

func TestCalcDeductions_TotalCappedByTax(t *testing.T) {
	// Небольшой оклад → мало налога. Большие вычеты → total > налога
	months := makeMonths(2_000_000, 1) // 20 000 ₽/мес → ~31 200 ₽ налога
	annualTax := months[len(months)-1].AnnualTaxAmount

	input := DeductionInput{
		HousingExpense:  200_000_000, // 2 000 000 ₽ → 260 000 ₽ возврат
		SocialExpense:   15_000_000,  // 150 000 ₽ → 19 500 ₽ возврат
		ChildrenCount:   3,
	}

	result := CalcDeductions(input, months)
	require.NotNil(t, result)

	// total не превышает уплаченный НДФЛ
	assert.LessOrEqual(t, result.TotalReturn, annualTax)
}

// ===== HasDeductions =====

func TestHasDeductions_False(t *testing.T) {
	assert.False(t, DeductionInput{}.HasDeductions())
}

func TestHasDeductions_True_Children(t *testing.T) {
	assert.True(t, DeductionInput{ChildrenCount: 1}.HasDeductions())
}

func TestHasDeductions_True_Disabled(t *testing.T) {
	assert.True(t, DeductionInput{DisabledChildrenCount: 1}.HasDeductions())
}

func TestHasDeductions_True_Housing(t *testing.T) {
	assert.True(t, DeductionInput{HousingExpense: 1_000_000}.HasDeductions())
}
