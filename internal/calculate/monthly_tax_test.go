package calculate

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// -------------------------------------------------------------------
// CalculateMonthlyTax — маршрутизация по ветвям
// -------------------------------------------------------------------

func TestCalculateMonthlyTax_Routing(t *testing.T) {
	jan := time.Date(time.Now().UTC().Year(), time.January, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		input         CalculateInput
		wantMonths    int
		wantTaxRate   uint64 // ожидаемая ставка первого месяца
		wantMonthlyTax uint64 // ожидаемый налог первого месяца
	}{
		{
			name: "нерезидент без коэффициентов (ветвь 1)",
			input: CalculateInput{
				GrossSalary:           120_000_00,
				TerritorialMultiplier: 100,
				NorthernCoefficient:   100,
				StartDate:             jan,
				IsNotResident:         true,
			},
			wantMonths:     12,
			wantTaxRate:    30,
			wantMonthlyTax: 36_000_00, // 120 000 * 30%
		},
		{
			name: "нерезидент + северная надбавка — север применяется как коэффициент (ветвь 1)",
			input: CalculateInput{
				GrossSalary:           100_000_00,
				TerritorialMultiplier: 100,
				NorthernCoefficient:   150, // +50%
				StartDate:             jan,
				IsNotResident:         true,
			},
			wantMonths:     12,
			wantTaxRate:    30,
			wantMonthlyTax: 45_000_00, // (100k * 1.50) * 30%
		},
		{
			name: "льготник силовых ведомств (ветвь 2)",
			input: CalculateInput{
				GrossSalary:           200_000_00,
				TerritorialMultiplier: 100,
				NorthernCoefficient:   100,
				StartDate:             jan,
				HasTaxPrivilege:       true,
			},
			wantMonths:     12,
			wantTaxRate:    13,
			wantMonthlyTax: 26_000_00, // 200 000 * 13%
		},
		{
			name: "только оклад без коэффициентов (ветвь 3)",
			input: CalculateInput{
				GrossSalary:           100_000_00,
				TerritorialMultiplier: 100,
				NorthernCoefficient:   100,
				StartDate:             jan,
			},
			wantMonths:     12,
			wantTaxRate:    13,
			wantMonthlyTax: 13_000_00, // 100 000 * 13%
		},
		{
			name: "оклад с РК, без севера (ветвь 4)",
			input: CalculateInput{
				GrossSalary:           100_000_00,
				TerritorialMultiplier: 120, // +20%
				NorthernCoefficient:   100,
				StartDate:             jan,
			},
			wantMonths:     12,
			wantTaxRate:    13,
			wantMonthlyTax: 15_600_00, // 120 000 * 13%
		},
		{
			name: "оклад без РК, с севером (ветвь 5)",
			input: CalculateInput{
				GrossSalary:           100_000_00,
				TerritorialMultiplier: 100,
				NorthernCoefficient:   150, // +50%
				StartDate:             jan,
			},
			wantMonths:     12,
			wantTaxRate:    13,
			wantMonthlyTax: 19_500_00, // оклад 100k * 13% + север 50k * 13%
		},
		{
			name: "оклад с РК и севером (ветвь 6 default)",
			input: CalculateInput{
				GrossSalary:           100_000_00,
				TerritorialMultiplier: 120, // база A = 120k
				NorthernCoefficient:   150, // север = 60k (120k * 50%)
				StartDate:             jan,
			},
			wantMonths:     12,
			wantTaxRate:    13,
			wantMonthlyTax: 23_400_00, // 120k*13% + 60k*13%
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateMonthlyTax(tt.input)
			require.Len(t, result, tt.wantMonths)
			assert.Equal(t, tt.wantTaxRate, result[0].TaxRate)
			assert.Equal(t, tt.wantMonthlyTax, result[0].MonthlyTaxAmount)
		})
	}
}

func TestCalculateMonthlyTax_ZeroSalary(t *testing.T) {
	jan := time.Date(time.Now().UTC().Year(), time.January, 1, 0, 0, 0, 0, time.UTC)

	input := CalculateInput{
		GrossSalary:           0,
		TerritorialMultiplier: 100,
		NorthernCoefficient:   100,
		StartDate:             jan,
	}
	result := CalculateMonthlyTax(input)
	require.Len(t, result, 12)
	for _, m := range result {
		assert.Equal(t, uint64(0), m.MonthlyTaxAmount, "нулевой оклад — нулевой налог")
		assert.Equal(t, uint64(0), m.MonthlyNetIncome, "нулевой оклад — нулевой net")
		assert.Equal(t, uint64(0), m.MonthlyPFR)
		assert.Equal(t, uint64(0), m.MonthlyFOMS)
		assert.Equal(t, uint64(0), m.MonthlyFSS)
	}
}

func TestCalculateMonthlyTax_ContributionsPopulated(t *testing.T) {
	// Проверяем что взносы работодателя заполняются
	jan := time.Date(time.Now().UTC().Year(), time.January, 1, 0, 0, 0, 0, time.UTC)

	input := CalculateInput{
		GrossSalary:           100_000_00,
		TerritorialMultiplier: 100,
		NorthernCoefficient:   100,
		StartDate:             jan,
	}
	result := CalculateMonthlyTax(input)
	require.NotEmpty(t, result)
	// Взносы должны быть ненулевыми при ненулевом окладе
	assert.Greater(t, result[0].MonthlyPFR, uint64(0))
	assert.Greater(t, result[0].MonthlyFOMS, uint64(0))
	assert.Greater(t, result[0].MonthlyFSS, uint64(0))
	// Годовые — это накопленная сумма, к декабрю должно быть больше чем в январе
	assert.Greater(t, result[11].AnnualPFR, result[0].AnnualPFR)
}

// -------------------------------------------------------------------
// TaxCalculateWithPrivilege — льготники (13%/15%)
// -------------------------------------------------------------------

func TestTaxCalculateWithPrivilege_BasicAndBoundaries(t *testing.T) {
	jan := time.Date(time.Now().UTC().Year(), time.January, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		salary  uint64
		start   time.Time
		wantLen int
		monthly []uint64 // ожидаемый налог по месяцам
	}{
		{
			name:    "300k/мес, весь год ≤5М → 13%",
			salary:  300_000_00, // 12 мес = 3.6М < 5М
			start:   jan,
			wantLen: 12,
			// 300 000 * 13% = 39 000 каждый месяц
			monthly: repeatAmount(39_000_00, 12),
		},
		{
			name:    "500k/мес, переход 5М в ноябре (ставка 13%→15%)",
			salary:  500_000_00, // 10 мес = 5.0М ровно, потом 15%
			start:   jan,
			wantLen: 12,
			// Jan–Oct: 500k * 13% = 65 000
			// Nov–Dec: дельта = 75 000 (пересечение 5М)
			monthly: append(repeatAmount(65_000_00, 10), repeatAmount(75_000_00, 2)...),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TaxCalculateWithPrivilege(tt.salary, tt.start, 1)
			require.Len(t, result, tt.wantLen)
			for i, m := range result {
				assert.Equal(t, tt.monthly[i], m.MonthlyTaxAmount,
					"месяц %d: ожидалось %d, получили %d", i+1, tt.monthly[i], m.MonthlyTaxAmount)
			}
		})
	}
}

func TestTaxCalculateWithPrivilege_TaxRateTransition(t *testing.T) {
	jan := time.Date(time.Now().UTC().Year(), time.January, 1, 0, 0, 0, 0, time.UTC)

	// 500k * 10 = 5M — ровно на границе порога в конце октября
	result := TaxCalculateWithPrivilege(500_000_00, jan, 1)
	require.Len(t, result, 12)

	for i := 0; i < 10; i++ {
		assert.Equal(t, uint64(13), result[i].TaxRate, "месяц %d должен быть 13%%", i+1)
	}
	for i := 10; i < 12; i++ {
		assert.Equal(t, uint64(15), result[i].TaxRate, "месяц %d должен быть 15%%", i+1)
	}
}

func TestTaxCalculateWithPrivilege_YTDMonotone(t *testing.T) {
	jan := time.Date(time.Now().UTC().Year(), time.January, 1, 0, 0, 0, 0, time.UTC)
	result := TaxCalculateWithPrivilege(200_000_00, jan, 1)
	for i := 1; i < len(result); i++ {
		assert.GreaterOrEqual(t, result[i].AnnualTaxAmount, result[i-1].AnnualTaxAmount,
			"годовой налог должен монотонно расти")
		assert.GreaterOrEqual(t, result[i].AnnualGrossIncome, result[i-1].AnnualGrossIncome)
	}
}

// -------------------------------------------------------------------
// calcEmployerContributions — взносы ПФР/ФОМС/ФСС
// -------------------------------------------------------------------

func TestCalcEmployerContributions_BelowLimits(t *testing.T) {
	// Оба дохода ниже лимитов → полные ставки
	pfr, foms, fss := calcEmployerContributions(0, 100_000_00)

	assert.Equal(t, uint64(22_000_00), pfr, "ПФР 22%")
	assert.Equal(t, uint64(5_100_00), foms, "ФОМС 5.1%")
	assert.Equal(t, uint64(2_900_00), fss, "ФСС 2.9%")
}

func TestCalcEmployerContributions_ZeroGross(t *testing.T) {
	pfr, foms, fss := calcEmployerContributions(0, 0)
	assert.Equal(t, uint64(0), pfr)
	assert.Equal(t, uint64(0), foms)
	assert.Equal(t, uint64(0), fss)
}

func TestCalcEmployerContributions_CrossingPfrLimit(t *testing.T) {
	// Половина гросса ещё в лимите, половина — уже нет
	income := uint64(PfrLimit - 50_000_00)
	gross := uint64(100_000_00)

	pfr, foms, fss := calcEmployerContributions(income, gross)

	// 50k по 22% + 50k по 10%
	wantPFR := uint64(50_000_00*PfrRate/1000 + 50_000_00*PfrRateHi/1000)
	assert.Equal(t, wantPFR, pfr, "ПФР: частичное пересечение лимита")

	assert.Equal(t, uint64(gross*FomsRate/1000), foms, "ФОМС всегда 5.1%")

	// ФСС: такой же лимит как ПФР — тоже пересекаем
	wantFSS := uint64(50_000_00 * FssRate / 1000)
	assert.Equal(t, wantFSS, fss, "ФСС: только часть до лимита")
}

func TestCalcEmployerContributions_AboveLimits(t *testing.T) {
	// Доход уже за лимитом → ПФР по пониженной, ФСС = 0
	pfr, foms, fss := calcEmployerContributions(PfrLimit, 100_000_00)

	assert.Equal(t, uint64(100_000_00*PfrRateHi/1000), pfr, "ПФР снижен до 10%")
	assert.Equal(t, uint64(100_000_00*FomsRate/1000), foms, "ФОМС всегда 5.1%")
	assert.Equal(t, uint64(0), fss, "ФСС = 0 после лимита")
}

func TestCalcEmployerContributions_FomsUnlimited(t *testing.T) {
	// ФОМС не имеет лимита — всегда 5.1%
	grossValues := []uint64{1_000_00, 100_000_00, 500_000_00, 10_000_000_00}
	for _, gross := range grossValues {
		_, foms, _ := calcEmployerContributions(PfrLimit*10, gross) // income точно за лимитом
		assert.Equal(t, gross*FomsRate/1000, foms, "ФОМС 5.1%% при gross=%d", gross)
	}
}

// -------------------------------------------------------------------
// helpers
// -------------------------------------------------------------------

func repeatAmount(amount uint64, n int) []uint64 {
	out := make([]uint64, n)
	for i := range out {
		out[i] = amount
	}
	return out
}
