package calculate

import "time"

// CalculateMonthlyTax рассчитывает помесячные значения налога на основе входных данных.
// Учитывает районный (территориальный) коэффициент и северную надбавку, если они заданы.
func CalculateMonthlyTax(input CalculateInput) []MonthlyTax {
	grossSalary := input.GrossSalary
	territorialMultiplier := input.TerritorialMultiplier
	northernCoefficient := input.NorthernCoefficient
	startDate := input.StartDate
	startMonth := GetStartMonth(startDate)

	noTerritorial := territorialMultiplier == 100
	noNorthern := northernCoefficient == 100

	// Оклад с учётом районного коэффициента (РК)
	monthlyGrossWithTerritorial := grossSalary * territorialMultiplier / 100

	var months []MonthlyTax

	bonuses := input.MonthlyBonuses
	if len(bonuses) < 12 {
		normalized := make([]uint64, 12)
		copy(normalized, bonuses)
		bonuses = normalized
	}

	switch {
	case input.IsNotResident: // Для нерезидентов: 30% со всего дохода
		totalMonthlyGross := monthlyGrossWithTerritorial * northernCoefficient / 100
		months = TaxCalculateForNotResident(totalMonthlyGross, startDate, startMonth, bonuses)

	case input.HasTaxPrivilege: // Для льготников (силовые ведомства): единая шкала 13%/15%
		totalMonthlyGross := monthlyGrossWithTerritorial * northernCoefficient / 100
		months = TaxCalculateWithPrivilege(totalMonthlyGross, startDate, startMonth, bonuses)

	case noTerritorial && noNorthern: // Без РК и СН
		months = TaxCalculateOnlySalary(grossSalary, startDate, startMonth, bonuses)

	case noNorthern: // Есть РК, нет СН
		months = TaxCalculateOnlySalary(monthlyGrossWithTerritorial, startDate, startMonth, bonuses)

	case noTerritorial: // Нет РК, есть СН
		northernAddition := grossSalary * (northernCoefficient - 100) / 100
		months = TaxCalculateWithNorth(grossSalary, northernAddition, startDate, startMonth, bonuses)

	default: // Есть и РК, и СН (СН начисляется на оклад с РК)
		northernAddition := monthlyGrossWithTerritorial * (northernCoefficient - 100) / 100
		months = TaxCalculateWithNorth(monthlyGrossWithTerritorial, northernAddition, startDate, startMonth, bonuses)
	}

	// Добавляем взносы от работодателя.
	var annualPFR, annualFOMS, annualFSS uint64

	for i := range months {
		gross := months[i].MonthlyGrossIncome
		incomeYTD := months[i].AnnualGrossIncome

		pfr, foms, fss := calcEmployerContributions(incomeYTD-gross, gross)

		months[i].MonthlyPFR = pfr
		months[i].MonthlyFOMS = foms
		months[i].MonthlyFSS = fss

		annualPFR += pfr
		annualFOMS += foms
		annualFSS += fss

		months[i].AnnualPFR = annualPFR
		months[i].AnnualFOMS = annualFOMS
		months[i].AnnualFSS = annualFSS
	}

	return months
}

// TaxCalculateForNotResident - расчёт налога для налоговых нерезидентов РФ (30% со всех доходов).
// Округление выполняется на месячной дельте.
func TaxCalculateForNotResident(salary uint64, startDate time.Time, startMonth int, bonuses []uint64) []MonthlyTax {
	var (
		result            []MonthlyTax
		annualGrossIncome uint64
		annualTaxRaw      uint64 // накопленный «сырой» налог (в копейках)
		annualTaxAmount   uint64 // сумма округлённых месячных налогов
	)

	for m := startMonth; m <= 12; m++ {
		bonus := bonuses[m-1]
		monthlyGross := salary + bonus
		annualGrossIncome += monthlyGross

		newTaxRaw := CalculateNotResidentTax(annualGrossIncome) // YTD без округления
		monthlyRaw := newTaxRaw - annualTaxRaw                  // «сырая» месячная дельта
		monthlyRounded := RoundTaxAmount(monthlyRaw)            // округление по правилу 50 копеек

		annualTaxRaw = newTaxRaw
		annualTaxAmount += monthlyRounded

		monthlyNetIncome := monthlyGross - monthlyRounded
		annualNetIncome := annualGrossIncome - annualTaxAmount
		month := IntMonthFromDate(m, startDate)

		result = append(result, MonthlyTax{
			Month:              month,
			MonthlyGrossIncome: monthlyGross,
			MonthlyNetIncome:   monthlyNetIncome,
			MonthlyTaxAmount:   monthlyRounded,
			AnnualGrossIncome:  annualGrossIncome,
			AnnualNetIncome:    annualNetIncome,
			AnnualTaxAmount:    annualTaxAmount,
			TaxRate:            NotResident.Rate,
			MonthlyBonus:       bonus,
		})
	}
	return result
}

// TaxCalculateWithPrivilege - расчёт для льготников (силовые ведомства) по упрощённой шкале 13%/15%.
// Округление выполняется на месячной дельте.
func TaxCalculateWithPrivilege(salary uint64, startDate time.Time, startMonth int, bonuses []uint64) []MonthlyTax {
	var (
		result            []MonthlyTax
		annualGrossIncome uint64
		annualTaxRaw      uint64 // YTD без округления
		annualTaxAmount   uint64 // сумма округлённых месячных
	)

	for m := startMonth; m <= 12; m++ {
		bonus := bonuses[m-1]
		monthlyGross := salary + bonus
		annualGrossIncome += monthlyGross

		newTaxRaw := CalculateSimpleProgressiveTax(annualGrossIncome) // RAW 13/15
		monthlyRaw := newTaxRaw - annualTaxRaw
		monthlyRounded := RoundTaxAmount(monthlyRaw)

		annualTaxRaw = newTaxRaw
		annualTaxAmount += monthlyRounded

		monthlyNetIncome := monthlyGross - monthlyRounded
		annualNetIncome := annualGrossIncome - annualTaxAmount
		currentRate := findSimpleCurrentRate(annualGrossIncome)
		month := IntMonthFromDate(m, startDate)

		result = append(result, MonthlyTax{
			Month:              month,
			MonthlyGrossIncome: monthlyGross,
			MonthlyNetIncome:   monthlyNetIncome,
			MonthlyTaxAmount:   monthlyRounded,
			AnnualGrossIncome:  annualGrossIncome,
			AnnualNetIncome:    annualNetIncome,
			AnnualTaxAmount:    annualTaxAmount,
			TaxRate:            currentRate,
			MonthlyBonus:       bonus,
		})
	}
	return result
}

// TaxCalculateOnlySalary - расчёт по общей пятиступенчатой шкале (без учёта северной надбавки).
// Округление выполняется на месячной дельте.
func TaxCalculateOnlySalary(salary uint64, startDate time.Time, startMonth int, bonuses []uint64) []MonthlyTax {
	var (
		result            []MonthlyTax
		annualGrossIncome uint64
		annualTaxRaw      uint64 // YTD без округления
		annualTaxAmount   uint64 // сумма округлённых месячных
	)

	for m := startMonth; m <= 12; m++ {
		bonus := bonuses[m-1]
		monthlyGross := salary + bonus
		annualGrossIncome += monthlyGross

		newTaxRaw := CalculateProgressiveTax(annualGrossIncome) // RAW по 5-ступ. шкале
		monthlyRaw := newTaxRaw - annualTaxRaw
		monthlyRounded := RoundTaxAmount(monthlyRaw)

		annualTaxRaw = newTaxRaw
		annualTaxAmount += monthlyRounded

		monthlyNetIncome := monthlyGross - monthlyRounded
		annualNetIncome := annualGrossIncome - annualTaxAmount
		currentRate := findCurrentRate(annualGrossIncome)
		month := IntMonthFromDate(m, startDate)

		result = append(result, MonthlyTax{
			Month:              month,
			MonthlyGrossIncome: monthlyGross,
			MonthlyNetIncome:   monthlyNetIncome,
			MonthlyTaxAmount:   monthlyRounded,
			AnnualGrossIncome:  annualGrossIncome,
			AnnualNetIncome:    annualNetIncome,
			AnnualTaxAmount:    annualTaxAmount,
			TaxRate:            currentRate,
			MonthlyBonus:       bonus,
		})
	}
	return result
}

// TaxCalculateWithNorth - расчёт при наличии северной надбавки.
// База A (оклад с РК) облагается по общей шкале; база B (северная надбавка) - по упрощённой 13%/15%.
// Бонус относится к базе A (общая прогрессивная шкала).
// Округление выполняется на месячной дельте по КАЖДОЙ базе отдельно.
func TaxCalculateWithNorth(salary, northernAddition uint64, startDate time.Time, startMonth int, bonuses []uint64) []MonthlyTax {
	var (
		result []MonthlyTax

		// Доходы YTD по базам
		annualBaseGrossIncome  uint64
		annualNorthGrossIncome uint64

		// Накопленные «сырые» налоги (в копейках) по базам
		baseTaxRawYTD  uint64
		northTaxRawYTD uint64

		// Накопленные округлённые суммы налога по базам
		annualBaseTaxAmount  uint64
		annualNorthTaxAmount uint64
	)

	for m := startMonth; m <= 12; m++ {
		bonus := bonuses[m-1] // бонус идёт в базу A (общая шкала)

		// Доходы YTD по базам
		annualBaseGrossIncome += salary + bonus
		annualNorthGrossIncome += northernAddition

		// «Сырые» YTD налоги по базам
		newBaseRaw := CalculateProgressiveTax(annualBaseGrossIncome)         // RAW общая шкала
		newNorthRaw := CalculateSimpleProgressiveTax(annualNorthGrossIncome) // RAW 13/15

		// Месячные «сырые» дельты и их округление
		monthlyBaseRaw := newBaseRaw - baseTaxRawYTD
		monthlyNorthRaw := newNorthRaw - northTaxRawYTD
		monthlyBaseRounded := RoundTaxAmount(monthlyBaseRaw)
		monthlyNorthRounded := RoundTaxAmount(monthlyNorthRaw)

		// Обновляем YTD
		baseTaxRawYTD = newBaseRaw
		northTaxRawYTD = newNorthRaw
		annualBaseTaxAmount += monthlyBaseRounded
		annualNorthTaxAmount += monthlyNorthRounded

		// Итоги месяца/года
		monthlyTaxAmount := monthlyBaseRounded + monthlyNorthRounded
		monthlyGrossIncome := salary + northernAddition + bonus
		monthlyNetIncome := monthlyGrossIncome - monthlyTaxAmount

		annualGrossIncome := annualBaseGrossIncome + annualNorthGrossIncome
		annualTaxAmount := annualBaseTaxAmount + annualNorthTaxAmount
		annualNetIncome := annualGrossIncome - annualTaxAmount

		currentRate := findCurrentRate(annualBaseGrossIncome) // маржинальная по базе A
		month := IntMonthFromDate(m, startDate)

		result = append(result, MonthlyTax{
			Month: month,

			MonthlyGrossIncome: monthlyGrossIncome,
			MonthlyNetIncome:   monthlyNetIncome,
			MonthlyTaxAmount:   monthlyTaxAmount,
			TaxRate:            currentRate,
			MonthlyBonus:       bonus,

			AnnualGrossIncome: annualGrossIncome,
			AnnualNetIncome:   annualNetIncome,
			AnnualTaxAmount:   annualTaxAmount,

			MonthlyNorthGrossIncome: northernAddition,
			MonthlyNorthTaxAmount:   monthlyNorthRounded,
			MonthlyBaseGrossIncome:  salary + bonus,
			MonthlyBaseTaxAmount:    monthlyBaseRounded,

			AnnualNorthGrossIncome: annualNorthGrossIncome,
			AnnualNorthTaxAmount:   annualNorthTaxAmount,
			AnnualBaseGrossIncome:  annualBaseGrossIncome,
			AnnualBaseTaxAmount:    annualBaseTaxAmount,
		})
	}
	return result
}

// CalculateProgressiveTax - расчёт по пятиступенчатой шкале (без округления).
// Не учитывает северную надбавку, облагаемую по упрощённой системе.
func CalculateProgressiveTax(income uint64) uint64 {
	var tax uint64
	var prev uint64
	for _, limit := range Limits {
		if income > limit.UpperLimit {
			diff := limit.UpperLimit - prev
			tax += diff * limit.Rate / 100
			prev = limit.UpperLimit
		} else {
			if income > prev {
				diff := income - prev
				tax += diff * limit.Rate / 100
			}
			break
		}
	}
	return tax
}

// CalculateSimpleProgressiveTax - расчёт по упрощённой шкале 13%/15% (без округления).
// До 5 млн - 13%, всё сверх - 15%.
func CalculateSimpleProgressiveTax(income uint64) uint64 {
	limit := SimpleLimits.UpperLimit
	if income <= limit {
		return income * 13 / 100
	}
	thirteen := uint64(limit * 13 / 100)
	fifteen := (income - limit) * 15 / 100
	return thirteen + fifteen
}

// CalculateNotResidentTax - расчёт налога для нерезидентов (30% на все доходы).
func CalculateNotResidentTax(income uint64) uint64 {
	return income * NotResident.Rate / 100
}

// findCurrentRate - текущая маржинальная ставка по общей шкале в зависимости от годового дохода.
func findCurrentRate(income uint64) uint64 {
	for _, limit := range Limits {
		if income <= limit.UpperLimit {
			return limit.Rate
		}
	}
	return Limits[len(Limits)-1].Rate
}

// findSimpleCurrentRate - текущая ставка по упрощённой шкале (13% до 5 млн, далее 15%).
func findSimpleCurrentRate(income uint64) uint64 {
	if income <= SimpleLimits.UpperLimit {
		return SimpleLimits.Rate // 13
	}
	return 15
}

/*
Согласно п. 6 ст. 52 НК РФ сумма налога исчисляется в полных рублях:
сумма менее 50 копеек отбрасывается, 50 копеек и более - округляется до полного рубля.
*/

// RoundTaxAmount - округляет налог до полных рублей по правилу 50 копеек.
func RoundTaxAmount(value uint64) uint64 {
	remainder := value % 100
	if remainder < 50 {
		return value - remainder
	}
	return value + (100 - remainder)
}

// Расчет налоговой нагрузки на работодателя.
func calcEmployerContributions(income, gross uint64) (pfr, foms, fss uint64) {

	// Расчет ПФР
	if income < PfrLimit {
		remaining := PfrLimit - income // Вычисляем сколько осталось до перехода лимита
		if gross <= remaining {
			pfr = gross * PfrRate / 1000
		} else {
			pfr = remaining*PfrRate/1000 + (gross-remaining)*PfrRateHi/1000
		}
	} else {
		pfr = gross * PfrRateHi / 1000
	}

	// Расчет ФОМС
	foms = gross * FomsRate / 1000 // всегда 5.1%

	// Расчет ФСС
	if income < FssLimit {
		remaining := FssLimit - income
		if gross <= remaining {
			fss = gross * FssRate / 1000
		} else {
			fss = remaining * FssRate / 1000
		}
	}

	return
}
