package calculate

import "time"

// TaxCalculateForSelfEmployed рассчитывает помесячные значения налога для самозанятых (НПД).
//
// Ставки (ст. 10 Закона № 422-ФЗ):
//   - 4% — доход от физических лиц
//   - 6% — доход от юридических лиц и ИП
//
// Регистрационный вычет (ст. 12 Закона № 422-ФЗ):
//   - При регистрации самозанятый получает бонус 10 000 ₽ (NpdRegistrationDeduction)
//   - Ставка снижается: 4% → 3%, 6% → 4% — пока бонус не исчерпан
//   - Экономия за месяц = monthlyGross * bonusRate / 100, где bonusRate: 4% → 1%, 6% → 2%
//
// Страховые взносы при НПД не начисляются.
// Территориальный коэффициент и северная надбавка на ставку НПД не влияют.
func TaxCalculateForSelfEmployed(
	income uint64,
	rate uint64,
	hasDeduction bool,
	startDate time.Time,
	startMonth int,
	bonuses []uint64,
) []MonthlyTax {
	var (
		result            []MonthlyTax
		annualGrossIncome uint64
		annualTaxAmount   uint64
		deductionLeft     uint64
	)

	// Ставка снижения при вычете: 4%→1%, 6%→2%
	var bonusRate uint64
	if rate == NpdRateIndividual {
		bonusRate = 1
	} else {
		bonusRate = 2
	}

	if hasDeduction {
		deductionLeft = NpdRegistrationDeduction
	}

	for m := startMonth; m <= 12; m++ {
		bonus := bonuses[m-1]
		monthlyGross := income + bonus
		annualGrossIncome += monthlyGross

		// Базовый налог без вычета (в копейках, без округления)
		monthlyTaxRaw := monthlyGross * rate / 100

		// Применяем регистрационный вычет
		var deductionUsed uint64
		if deductionLeft > 0 {
			potential := monthlyGross * bonusRate / 100
			if potential <= deductionLeft {
				deductionUsed = potential
			} else {
				deductionUsed = deductionLeft
			}
			deductionLeft -= deductionUsed
		}

		monthlyTax := RoundTaxAmount(monthlyTaxRaw - deductionUsed)
		annualTaxAmount += monthlyTax
		monthlyNet := monthlyGross - monthlyTax
		annualNet := annualGrossIncome - annualTaxAmount

		result = append(result, MonthlyTax{
			Month:             IntMonthFromDate(m, startDate),
			MonthlyGrossIncome: monthlyGross,
			MonthlyNetIncome:  monthlyNet,
			MonthlyTaxAmount:  monthlyTax,
			TaxRate:           rate,
			AnnualGrossIncome: annualGrossIncome,
			AnnualNetIncome:   annualNet,
			AnnualTaxAmount:   annualTaxAmount,
			MonthlyBonus:      bonus,
			NpdDeductionUsed:  deductionUsed,
		})
	}

	return result
}
