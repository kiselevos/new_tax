package calculate

import "time"

type MonthlyTaxLite struct {
	Month            time.Time // 1-е число месяца (UTC)
	MonthlyTaxAmount uint64
	MonthlyBaseTax   uint64 // опционально: если в сценарии есть север, фиксируем разложение A/B
	MonthlyNorthTax  uint64
	TaxRate          uint64 // маржинальная по базе A (13/15/18/20/22) - полезно на переходах
}

type ResultTest struct {
	Name    string
	Monthly []MonthlyTaxLite
}

func firstOf(m time.Month) time.Time {
	y := time.Now().UTC().Year() // год берём текущий, т.к. на расчёт не влияет
	return time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)
}

// быстрая генерация одинаковых месяцев
func repeatMonthly(start time.Month, n int, mk func(i int) MonthlyTaxLite) []MonthlyTaxLite {
	out := make([]MonthlyTaxLite, n)
	for i := 0; i < n; i++ {
		mt := mk(i)
		mt.Month = firstOf(time.Month(int(start) + i))
		out[i] = mt
	}
	return out
}

var Results = []ResultTest{
	// A: 100k, terr 20% => база A=120k; север=0; 13% весь год
	{
		Name: "A: Only salary usual person (100k, terr 20%, jan start)",
		Monthly: repeatMonthly(time.January, 12, func(i int) MonthlyTaxLite {
			return MonthlyTaxLite{
				MonthlyTaxAmount: 15_600_00, // 120 000 * 13%
				MonthlyBaseTax:   15_600_00,
				MonthlyNorthTax:  0,
				TaxRate:          13,
			}
		}),
	},

	// B: 200k, no coeffs; годом ровно 2.4М → 13% весь год
	{
		Name: "B: 200k salary, no coeffs, at 2.4M threshold (jan start)",
		Monthly: repeatMonthly(time.January, 12, func(i int) MonthlyTaxLite {
			return MonthlyTaxLite{
				MonthlyTaxAmount: 26_000_00, // 200 000 * 13%
				MonthlyBaseTax:   26_000_00,
				MonthlyNorthTax:  0,
				TaxRate:          13,
			}
		}),
	},

	// C: 300k, порог 2.4М после 8 месяцев:
	// Jan–Aug: 39 000; Sep–Dec: 45 000; ставки 13% → 15%
	{
		Name: "C: 300k salary, cross 2.4M in Sep (jan start)",
		Monthly: func() []MonthlyTaxLite {
			months := make([]MonthlyTaxLite, 12)
			for i := 0; i < 12; i++ {
				months[i].Month = firstOf(time.Month(int(time.January) + i))
				if i < 8 { // Jan..Aug
					months[i].MonthlyTaxAmount = 39_000_00
					months[i].MonthlyBaseTax = 39_000_00
					months[i].TaxRate = 13
				} else { // Sep..Dec
					months[i].MonthlyTaxAmount = 45_000_00
					months[i].MonthlyBaseTax = 45_000_00
					months[i].TaxRate = 15
				}
			}
			return months
		}(),
	},

	// D: 200k + terr 20% + north 50%
	// База A=240k; Север B=120k
	// A: Jan–Oct 31 200, Nov–Dec 36 000
	// B: весь год 15 600
	// Итого: Jan–Oct 46 800; Nov–Dec 51 600
	{
		Name: "D: 200k salary + terr 20% + north 50% (jan start)",
		Monthly: func() []MonthlyTaxLite {
			months := make([]MonthlyTaxLite, 12)
			for i := 0; i < 12; i++ {
				months[i].Month = firstOf(time.Month(int(time.January) + i))
				if i < 10 {
					months[i].MonthlyBaseTax = 31_200_00
					months[i].MonthlyNorthTax = 15_600_00
					months[i].MonthlyTaxAmount = 46_800_00
					months[i].TaxRate = 13
				} else {
					months[i].MonthlyBaseTax = 36_000_00
					months[i].MonthlyNorthTax = 15_600_00
					months[i].MonthlyTaxAmount = 51_600_00
					months[i].TaxRate = 15
				}
			}
			return months
		}(),
	},

	// E: 120k + terr 20% + north 50%, старт июнь (7 мес)
	// База A=144k → 18 720; Север B=72k → 9 360; Итого 28 080 все 7 мес; ставка 13%
	{
		Name: "E: 120k salary + terr 20% + north 50% (start in Jun, projection 7m)",
		Monthly: repeatMonthly(time.June, 7, func(i int) MonthlyTaxLite {
			return MonthlyTaxLite{
				MonthlyTaxAmount: 28_080_00,
				MonthlyBaseTax:   18_720_00,
				MonthlyNorthTax:  9_360_00,
				TaxRate:          13,
			}
		}),
	},

	// NR-1: нерезидент, 1 ₽/мес → 30% = 30 коп → округляется до 0 каждый месяц
	{
		Name: "NR-1: Non-resident, tiny salary → monthly rounding to 0 (jan start)",
		Monthly: repeatMonthly(time.January, 12, func(i int) MonthlyTaxLite {
			return MonthlyTaxLite{
				MonthlyTaxAmount: 0,
				MonthlyBaseTax:   0,
				MonthlyNorthTax:  0,
				TaxRate:          30, // плоская ставка NR
			}
		}),
	},

	// NR-2: нерезидент, 120k, старт июнь (7 мес) → 36 000/мес
	{
		Name: "NR-2: Non-resident, 120k salary (start in Jun, 7 months)",
		Monthly: repeatMonthly(time.June, 7, func(i int) MonthlyTaxLite {
			return MonthlyTaxLite{
				MonthlyTaxAmount: 36_000_00,
				MonthlyBaseTax:   36_000_00,
				MonthlyNorthTax:  0,
				TaxRate:          30,
			}
		}),
	},
}
