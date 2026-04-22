package calculate

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Ровно 50 копеек округляем ВВЕРХ (к ближайшим 100)
func TestRoundTaxAmount_Basic(t *testing.T) {
	tests := []struct {
		in   uint64
		want uint64
	}{
		{13200045, 13200000},
		{13200085, 13200100},
		{13200000, 13200000},

		{0, 0},
		{1, 0},
		{49, 0},
		{50, 100},
		{51, 100},
		{99, 100},
		{100, 100},
		{101, 100},
		{149, 100},
		{150, 200},

		{12_345_678_900, 12_345_678_900},
		{12_345_678_949, 12_345_678_900},
		{12_345_678_950, 12_345_679_000},
	}

	for _, test := range tests {
		got := RoundTaxAmount(test.in)
		assert.Equal(t, test.want, got, "in=%d", test.in)
	}
}

func TestRoundTaxAmount_IdempotentAndMultiplieOf100(t *testing.T) {

	tests := []uint64{0, 1, 49, 50, 51, 99, 100, 101, 13200045, 13200085}

	for _, i := range tests {
		r := RoundTaxAmount(i)

		if r%100 != 0 {
			t.Fatalf("not multiple of 100: in=%d r=%d", i, r)
		}

		if rr := RoundTaxAmount(r); r != rr {
			t.Fatalf("not idempotent: rr=%d r=%d", rr, r)
		}
	}
}

// До 5 млн включительно применяется 13%, ровно на границе остаёмся на 13%.
func TestFindSimpleCurrentRate_Basic(t *testing.T) {
	tests := []struct {
		name string
		in   uint64
		want uint64
	}{
		{"zero", 0, 13},
		{"well_below", 1_000_000_00, 13},
		{"just_below_by_99", 4_999_999_99, 13},
		{"at_threshold", 5_000_000_00, 13},
		{"just_above_by_1", 5_000_000_01, 15},
		{"well_above", 6_000_000_00, 15},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := findSimpleCurrentRate(test.in)
			assert.Equal(t, test.want, got, "want=%d got=%d", test.want, got)
		})
	}
}

func TestFindCurrentRate_Basic(t *testing.T) {
	tests := []struct {
		name string
		in   uint64
		want uint64
	}{
		{"zero", 0, 13},
		{"first_step", 1_000_000_00, 13},
		{"first_step_by_99", 2_399_999_99, 13},
		{"first_step_border", 2_400_000_00, 13},
		{"second_step_plus_1", 2_400_001_00, 15},
		{"second_step", 3_400_000_00, 15},
		{"second_step_by_99", 4_999_999_99, 15},
		{"second_step_border", 5_000_000_00, 15},
		{"third_step_plus_1", 5_000_001_00, 18},
		{"third_step", 10_000_000_00, 18},
		{"third_step_by_99", 19_999_999_99, 18},
		{"third_step_border", 20_000_000_00, 18},
		{"fourth_step_plus_1", 20_000_001_00, 20},
		{"fourth_step", 40_000_000_00, 20},
		{"fourth_step_by_99", 49_999_999_99, 20},
		{"fourth_step_border", 50_000_000_00, 20},
		{"five_step_plus_1", 50_000_001_00, 22},
		{"five_step_boarder", math.MaxUint64 - 1, 22},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := findCurrentRate(test.in)
			assert.Equal(t, test.want, got, "want=%d got=%d", test.want, got)
		})
	}
}

func TestCalculateNotResidentTax_Basic(t *testing.T) {
	// Зафиксируем ставку и вернём обратно после теста
	old := NotResident.Rate
	NotResident.Rate = 30
	t.Cleanup(func() { NotResident.Rate = old })

	tests := []struct {
		name   string
		income uint64
		want   uint64
	}{
		{"zero", 0, 0},
		{"one_ruble", 100, 30},
		{"two_rubles", 200, 60},
		{"odd_amount", 12_345_679, 3_703_703}, // 12_345_679 * 30 / 100 = 3_703_703 коп
		{"big_amount", 12_000_000, 3_600_000}, // 120 000 ₽ → 3 600 ₽ = 3 600 000 коп
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := CalculateNotResidentTax(tc.income)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestCalculateSimpleProgressiveTax_BasicAndBoundaries(t *testing.T) {
	limit := SimpleLimits.UpperLimit

	tests := []struct {
		name   string
		income uint64
		want   uint64
	}{
		{"zero", 0, 0},

		// ниже порога - 13%
		{"well_below", 1_000_000_00, 1_000_000_00 * 13 / 100},
		{"just_below_by_99kop", limit - 1, (limit - 1) * 13 / 100},

		// ровно на пороге - всё по 13%
		{"at_threshold", limit, limit * 13 / 100},

		// выше порога - кусочно: 13% на limit + 15% на остаток
		{"just_above_by_1kop", limit + 1, limit*13/100 + (1 * 15 / 100)}, // обычно даст 0 на «хвосте» из-за целочисленного деления - и это ок
		{"above_some_amount", 6_000_000_00, limit*13/100 + (6_000_000_00-limit)*15/100},

		// большая сумма, но в допустимых бизнес-границах
		{"big_amount", 40_000_000_00, limit*13/100 + (40_000_000_00-limit)*15/100},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := CalculateSimpleProgressiveTax(tc.income)
			assert.Equalf(t, tc.want, got, "income=%d", tc.income)
		})
	}
}

func TestCalculateProgressiveTax_BasicAndBoundaries(t *testing.T) {
	U := []struct{ U, R uint64 }{
		{2_400_000_00, 13},
		{5_000_000_00, 15},
		{20_000_000_00, 18},
		{50_000_000_00, 20},
		{math.MaxUint64, 22},
	}

	at := func(u []struct{ U, R uint64 }, upto int) uint64 {
		var tax, prev uint64
		for i := 0; i <= upto; i++ {
			tax += (u[i].U - prev) * u[i].R / 100
			prev = u[i].U
		}
		return tax
	}

	tests := []struct {
		name   string
		income uint64
		want   uint64
	}{
		{"zero", 0, 0},
		{"below_first", 1_000_000_00, 1_000_000_00 * 13 / 100},

		// Ровно на границах
		{"at_1st", U[0].U, at(U, 0)},
		{"at_2nd", U[1].U, at(U, 0) + (U[1].U-U[0].U)*15/100},
		{"at_3rd", U[2].U, at(U, 0) + (U[1].U-U[0].U)*15/100 + (U[2].U-U[1].U)*18/100},
		{"at_4th", U[3].U, at(U, 0) + (U[1].U-U[0].U)*15/100 + (U[2].U-U[1].U)*18/100 + (U[3].U-U[2].U)*20/100},

		// Чуть по обе стороны границ
		{"just_below_2nd", U[1].U - 1, at(U, 0) + (U[1].U-1-U[0].U)*15/100},
		{"just_above_2nd", U[1].U + 1, at(U, 0) + (U[1].U-U[0].U)*15/100 + (1*18)/100}, // 1 коп под 18% даёт 0 коп «сырых» - это ожидаемо

		{"mid_3rd", 10_000_000_00, // внутри 3-й ступени
			at(U, 0) + (U[1].U-U[0].U)*15/100 + (10_000_000_00-U[1].U)*18/100,
		},

		// Выше всех порогов
		{"above_all", 60_000_000_00,
			at(U, 0) + (U[1].U-U[0].U)*15/100 + (U[2].U-U[1].U)*18/100 +
				(U[3].U-U[2].U)*20/100 + (60_000_000_00-U[3].U)*22/100,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := CalculateProgressiveTax(tc.income)
			assert.Equalf(t, tc.want, got, "income=%d", tc.income)
		})
	}
}

func resultsByName(rs []ResultTest) map[string]ResultTest {
	m := make(map[string]ResultTest, len(rs))
	for _, r := range rs {
		m[r.Name] = r
	}
	return m
}

// базовые проверки инвариантов на месяце + накопление годового
func assertMonthInvariant(
	t *testing.T,
	i int,
	got MonthlyTax,
	expect MonthlyTaxLite,
	baseA, baseB uint64,
	acc *uint64,
	checkSplit bool,
) {
	t.Helper()

	// 1) месяц: сравниваем YYYY-MM (год для детерминизма)
	assert.Equalf(t, expect.Month.Format("2006-01"), got.Month.Format("2006-01"), "month[%d]", i)

	// 2) gross фиксирован входом
	assert.Equalf(t, baseA+baseB, got.MonthlyGrossIncome, "gross[%d]", i)

	// 3) если в expectations поле ненулевое - сверяем, иначе пропускаем (маска)
	if expect.MonthlyTaxAmount != 0 {
		assert.Equalf(t, expect.MonthlyTaxAmount, got.MonthlyTaxAmount, "tax[%d]", i)
	}

	if checkSplit {
		if expect.MonthlyBaseTax != 0 {
			assert.Equalf(t, expect.MonthlyBaseTax, got.MonthlyBaseTaxAmount, "baseTax[%d]", i)
		}
		// Если есть северная база (baseB>0), проверяем её даже если в ожидаемом 0 - чтобы не забыть про разложение.
		if expect.MonthlyNorthTax != 0 || baseB > 0 {
			assert.Equalf(t, expect.MonthlyNorthTax, got.MonthlyNorthTaxAmount, "northTax[%d]", i)
		}
	}

	if expect.TaxRate != 0 {
		assert.EqualValuesf(t, expect.TaxRate, got.TaxRate, "rate[%d]", i)
	}

	// 4) инварианты
	if checkSplit {
		assert.Equalf(t,
			got.MonthlyBaseTaxAmount+got.MonthlyNorthTaxAmount,
			got.MonthlyTaxAmount,
			"split==sum[%d]", i,
		)
	}
	assert.Equalf(t,
		got.MonthlyGrossIncome-got.MonthlyTaxAmount,
		got.MonthlyNetIncome,
		"net==gross-tax[%d]", i,
	)

	// 5) накопление годового налога = сумма месячных (округлённых)
	*acc += got.MonthlyTaxAmount
	assert.Equalf(t, *acc, got.AnnualTaxAmount, "annual sum[%d]", i)
}

// Быстро получить базы A/B из сценария
func basesFromScenario(sc Scenario) (uint64, uint64) {
	baseA := sc.Salary * (100 + sc.TerritorialMultiplier) / 100 // оклад + теркоэф
	baseB := baseA * sc.NorthernCoefficient / 100               // север от базы A
	return baseA, baseB
}

func Test_TaxCalculateWithNorth_Scenarios(t *testing.T) {
	wantMap := resultsByName(Results)

	for _, sc := range Scenarios {
		// Нерезидентов здесь не тестируем: у них плоская ставка 30% и другая функция.
		if sc.NonResident {
			continue
		}
		want, ok := wantMap[sc.Name]
		if !ok || len(want.Monthly) == 0 {
			continue
		}

		t.Run("WithNorth/"+sc.Name, func(t *testing.T) {
			baseA, baseB := basesFromScenario(sc)
			got := TaxCalculateWithNorth(baseA, baseB, sc.StartDate(), int(sc.StartMonth), make([]uint64, 12))

			require.Equal(t, len(want.Monthly), len(got), "months length")

			var acc uint64
			for i := range got {
				assertMonthInvariant(t, i, got[i], want.Monthly[i], baseA, baseB, &acc, true)
			}
			assert.Equal(t, acc, got[len(got)-1].AnnualTaxAmount)
		})
	}
}

func Test_TaxCalculateOnlySalary_Scenarios(t *testing.T) {
	wantMap := resultsByName(Results)

	for _, sc := range Scenarios {
		// Только резиденты и без севера
		if sc.NonResident || sc.NorthernCoefficient != 0 {
			continue
		}
		want, ok := wantMap[sc.Name]
		if !ok || len(want.Monthly) == 0 {
			continue
		}

		t.Run("OnlySalary/"+sc.Name, func(t *testing.T) {
			baseA, _ := basesFromScenario(sc)
			got := TaxCalculateOnlySalary(baseA, sc.StartDate(), int(sc.StartMonth), make([]uint64, 12))

			require.Equal(t, len(want.Monthly), len(got), "months length")

			var acc uint64
			for i := range got {
				expect := want.Monthly[i]
				// не проверяем разложение в этом тесте
				expect.MonthlyBaseTax = 0
				expect.MonthlyNorthTax = 0
				assertMonthInvariant(t, i, got[i], expect, baseA, 0, &acc, false)
			}
			assert.Equal(t, acc, got[len(got)-1].AnnualTaxAmount)
		})
	}
}

func Test_TaxCalculateForNotResident_Scenarios(t *testing.T) {
	wantMap := resultsByName(Results)

	// Фиксируем ставку нерезидента на время теста (если глобальная).
	oldRate := NotResident.Rate
	if oldRate != 30 {
		NotResident.Rate = 30
		defer func() { NotResident.Rate = oldRate }()
	}

	for _, sc := range Scenarios {
		if !sc.NonResident {
			continue
		}
		want, ok := wantMap[sc.Name]
		if !ok || len(want.Monthly) == 0 {
			continue
		}

		t.Run("NonResident/"+sc.Name, func(t *testing.T) {
			baseA, _ := basesFromScenario(sc) // у NR считаем только по базе A
			got := TaxCalculateForNotResident(baseA, sc.StartDate(), int(sc.StartMonth), make([]uint64, 12))

			require.Equal(t, len(want.Monthly), len(got), "months length")

			var acc uint64
			for i := range got {
				expect := want.Monthly[i]

				assert.Equalf(t, expect.Month.Format("2006-01"), got[i].Month.Format("2006-01"), "month[%d]", i)

				// Gross: только база A (без севера).
				assert.Equalf(t, baseA, got[i].MonthlyGrossIncome, "gross[%d]", i)

				// Общий налог за месяц (после округления).
				if expect.MonthlyTaxAmount != 0 {
					assert.Equalf(t, expect.MonthlyTaxAmount, got[i].MonthlyTaxAmount, "tax[%d]", i)
				}

				// Net = Gross - Tax.
				assert.Equalf(t, got[i].MonthlyGrossIncome-got[i].MonthlyTaxAmount, got[i].MonthlyNetIncome, "net[%d]", i)

				// Плоская ставка (если ожидаем задаём).
				if expect.TaxRate != 0 {
					assert.EqualValuesf(t, expect.TaxRate, got[i].TaxRate, "rate[%d]", i)
				}

				// Накопление (сумма округлённых месячных).
				acc += got[i].MonthlyTaxAmount
				assert.Equalf(t, acc, got[i].AnnualTaxAmount, "annual sum[%d]", i)
			}

			// Финальная сверка годового итога.
			assert.Equal(t, acc, got[len(got)-1].AnnualTaxAmount)
		})
	}
}
