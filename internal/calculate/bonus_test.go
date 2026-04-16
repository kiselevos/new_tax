package calculate

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var jan = time.Date(time.Now().UTC().Year(), time.January, 1, 0, 0, 0, 0, time.UTC)

// TestGetBonusForMonth проверяет вспомогательную функцию getBonusForMonth.
func TestGetBonusForMonth(t *testing.T) {
	bonuses := []uint64{0, 0, 50_000_00, 0, 0, 0, 0, 0, 0, 0, 0, 100_000_00}

	assert.Equal(t, uint64(0), getBonusForMonth(bonuses, 1))         // январь — 0
	assert.Equal(t, uint64(50_000_00), getBonusForMonth(bonuses, 3)) // март — 50k
	assert.Equal(t, uint64(100_000_00), getBonusForMonth(bonuses, 12)) // декабрь — 100k
	assert.Equal(t, uint64(0), getBonusForMonth(nil, 1))             // nil slice — 0
	assert.Equal(t, uint64(0), getBonusForMonth([]uint64{}, 1))      // пустой slice — 0
	assert.Equal(t, uint64(0), getBonusForMonth(bonuses, 13))        // за границей — 0
}

// TestBonus_NoBonus_MatchesBaselineCalculation проверяет что nil-бонусы не меняют результат.
func TestBonus_NoBonus_MatchesBaselineCalculation(t *testing.T) {
	salary := uint64(100_000_00)

	withoutBonus := TaxCalculateOnlySalary(salary, jan, 1, nil)
	withZeroBonuses := TaxCalculateOnlySalary(salary, jan, 1, make([]uint64, 12))

	require.Len(t, withoutBonus, 12)
	require.Len(t, withZeroBonuses, 12)

	for i := range withoutBonus {
		assert.Equal(t, withoutBonus[i].MonthlyTaxAmount, withZeroBonuses[i].MonthlyTaxAmount,
			"месяц %d: налог не должен отличаться", i+1)
		assert.Equal(t, withoutBonus[i].AnnualGrossIncome, withZeroBonuses[i].AnnualGrossIncome,
			"месяц %d: YTD gross не должен отличаться", i+1)
	}
}

// TestBonus_IncreasesGrossAndTax проверяет что бонус увеличивает gross и налог в нужном месяце.
func TestBonus_IncreasesGrossAndTax(t *testing.T) {
	salary := uint64(100_000_00) // 100k в месяц
	bonus := uint64(50_000_00)   // бонус 50k в марте

	bonuses := make([]uint64, 12)
	bonuses[2] = bonus // март = индекс 2

	result := TaxCalculateOnlySalary(salary, jan, 1, bonuses)
	require.Len(t, result, 12)

	// Январь и февраль — без изменений
	assert.Equal(t, salary, result[0].MonthlyGrossIncome, "январь: gross должен быть только оклад")
	assert.Equal(t, salary, result[1].MonthlyGrossIncome, "февраль: gross должен быть только оклад")

	// Март — gross включает бонус
	assert.Equal(t, salary+bonus, result[2].MonthlyGrossIncome, "март: gross должен включать бонус")
	assert.Equal(t, bonus, result[2].MonthlyBonus, "март: MonthlyBonus должен быть равен бонусу")

	// Апрель — снова только оклад
	assert.Equal(t, salary, result[3].MonthlyGrossIncome, "апрель: gross должен быть только оклад")
	assert.Equal(t, uint64(0), result[3].MonthlyBonus, "апрель: MonthlyBonus должен быть 0")

	// YTD в марте должен быть больше чем без бонуса
	resultNoBonus := TaxCalculateOnlySalary(salary, jan, 1, nil)
	assert.Greater(t, result[2].AnnualGrossIncome, resultNoBonus[2].AnnualGrossIncome,
		"март: YTD gross с бонусом должен быть больше")
	assert.Greater(t, result[2].AnnualTaxAmount, resultNoBonus[2].AnnualTaxAmount,
		"март: YTD налог с бонусом должен быть больше")
}

// TestBonus_ThresholdCrossing проверяет что бонус может перевести в более высокую ставку.
// Оклад 190k/мес = 2.28M/год — до порога 2.4M.
// Бонус 200k в декабре итого 2.48M — должна смениться ставка на 15%.
func TestBonus_ThresholdCrossing(t *testing.T) {
	salary := uint64(190_000_00)  // 190k/мес, за 12 мес = 2.28M < 2.4M порога
	bonus := uint64(200_000_00)   // бонус в декабре: 2.28M + 200k = 2.48M > 2.4M

	bonuses := make([]uint64, 12)
	bonuses[11] = bonus // декабрь = индекс 11

	withBonus := TaxCalculateOnlySalary(salary, jan, 1, bonuses)
	withoutBonus := TaxCalculateOnlySalary(salary, jan, 1, nil)

	require.Len(t, withBonus, 12)

	// Без бонуса все месяцы должны быть в ставке 13%
	for i, m := range withoutBonus {
		assert.Equal(t, uint64(13), m.TaxRate, "без бонуса месяц %d должен быть 13%%", i+1)
	}

	// С бонусом декабрь должен перейти в 15%
	assert.Equal(t, uint64(15), withBonus[11].TaxRate, "декабрь с бонусом должен быть 15%%")

	// Ноябрь остаётся в 13%
	assert.Equal(t, uint64(13), withBonus[10].TaxRate, "ноябрь должен оставаться в 13%%")
}

// TestBonus_WithNorth_BonusGoesToBaseA проверяет что бонус при северной надбавке идёт в базу A.
func TestBonus_WithNorth_BonusGoesToBaseA(t *testing.T) {
	salary := uint64(100_000_00)        // оклад
	northernAddition := uint64(30_000_00) // северная надбавка 30%
	bonus := uint64(50_000_00)          // бонус в марте

	bonuses := make([]uint64, 12)
	bonuses[2] = bonus // март

	result := TaxCalculateWithNorth(salary, northernAddition, jan, 1, bonuses)
	require.Len(t, result, 12)

	// В марте MonthlyBonus должен быть установлен
	assert.Equal(t, bonus, result[2].MonthlyBonus, "март: MonthlyBonus должен быть равен бонусу")

	// MonthlyBaseGrossIncome в марте должен включать бонус (база A)
	assert.Equal(t, salary+bonus, result[2].MonthlyBaseGrossIncome,
		"март: база A должна включать бонус")

	// Северная надбавка не меняется
	assert.Equal(t, northernAddition, result[2].MonthlyNorthGrossIncome,
		"март: северная надбавка не должна меняться от бонуса")

	// Общий gross в марте = оклад + надбавка + бонус
	assert.Equal(t, salary+northernAddition+bonus, result[2].MonthlyGrossIncome,
		"март: итоговый gross должен включать все три компонента")
}

// TestBonus_AnnualSum проверяет что сумма всех бонусов корректно отражается в YTD.
func TestBonus_AnnualSum(t *testing.T) {
	salary := uint64(100_000_00)
	quarterlyBonus := uint64(50_000_00)

	// Квартальные бонусы: март, июнь, сентябрь, декабрь
	bonuses := make([]uint64, 12)
	bonuses[2] = quarterlyBonus
	bonuses[5] = quarterlyBonus
	bonuses[8] = quarterlyBonus
	bonuses[11] = quarterlyBonus

	result := TaxCalculateOnlySalary(salary, jan, 1, bonuses)
	require.Len(t, result, 12)

	expectedAnnualBonus := quarterlyBonus * 4 // 200k итого
	expectedAnnualGross := salary*12 + expectedAnnualBonus

	assert.Equal(t, expectedAnnualGross, result[11].AnnualGrossIncome,
		"декабрь: годовой gross должен включать все квартальные бонусы")
}
