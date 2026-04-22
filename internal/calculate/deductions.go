package calculate

import pb "github.com/kiselevos/new_tax/gen/grpc/api"

// Константы налоговых вычетов (актуальные значения: project-docs/TAX_CONSTANTS.md)
const (
	// Ст. 218 — стандартный вычет на детей
	childrenIncomeLimit    = 45_000_000  // лимит 450 000 ₽ в копейках
	child1Deduction        = 140_000     // 1 400 ₽/мес на первого ребёнка
	child2Deduction        = 280_000     // 2 800 ₽/мес на второго ребёнка
	childOtherDeduction    = 600_000     // 6 000 ₽/мес на третьего и последующих
	disabledChildDeduction = 1_200_000   // 12 000 ₽/мес за каждого ребёнка-инвалида

	// Ст. 220 — имущественный вычет
	housingDeductionLimit  = 200_000_000 // 2 000 000 ₽ в копейках
	mortgageDeductionLimit = 300_000_000 // 3 000 000 ₽ в копейках

	// Ст. 219 — социальный вычет
	socialDeductionLimit   = 15_000_000 // 150 000 ₽ в копейках
	childEduDeductionLimit = 11_000_000 // 110 000 ₽ в копейках
)

// DeductionInput содержит параметры для расчёта налоговых вычетов.
type DeductionInput struct {
	ChildrenCount         uint32
	DisabledChildrenCount uint32
	HousingExpense        uint64 // копейки, ст. 220
	MortgageExpense       uint64 // копейки, ст. 220
	SocialExpense         uint64 // копейки, ст. 219 (лечение + собственное обучение)
	ChildEduExpense       uint64 // копейки, ст. 219 (обучение ребёнка)
}

// HasDeductions возвращает true если заданы параметры для расчёта хотя бы одного вычета.
func (d DeductionInput) HasDeductions() bool {
	return d.ChildrenCount > 0 || d.DisabledChildrenCount > 0 ||
		d.HousingExpense > 0 || d.MortgageExpense > 0 ||
		d.SocialExpense > 0 || d.ChildEduExpense > 0
}

// CalcDeductions рассчитывает налоговые вычеты по ст. 218, 219, 220 НК РФ.
// Вычеты рассчитываются на основе фактических помесячных данных расчёта.
// Возвращает nil если вычеты не заданы или отсутствуют месяцы расчёта.
func CalcDeductions(input DeductionInput, months []MonthlyTax) *pb.DeductionResult {
	if !input.HasDeductions() || len(months) == 0 {
		return nil
	}

	annualTaxPaid := months[len(months)-1].AnnualTaxAmount
	marginalRate := months[len(months)-1].TaxRate

	// 1. Стандартный вычет на детей (ст. 218)
	// Действует по ставке 13%: лимит дохода (450 000 ₽) значительно ниже порога 15% (2 400 000 ₽).
	monthlyDeduction := childrenMonthlyDeduction(input.ChildrenCount, input.DisabledChildrenCount)
	var childrenMonths uint32
	for _, m := range months {
		incomeAtMonthStart := m.AnnualGrossIncome - m.MonthlyGrossIncome
		if incomeAtMonthStart < childrenIncomeLimit {
			childrenMonths++
		}
	}
	var childrenReturn uint64
	if childrenMonths > 0 && monthlyDeduction > 0 {
		totalDeductionBase := uint64(childrenMonths) * monthlyDeduction
		childrenReturn = RoundTaxAmount(totalDeductionBase * 13 / 100)
		childrenReturn = min(childrenReturn, annualTaxPaid)
	}

	// 2. Имущественный вычет (ст. 220)
	cappedHousing := min(input.HousingExpense, uint64(housingDeductionLimit))
	cappedMortgage := min(input.MortgageExpense, uint64(mortgageDeductionLimit))
	propertyReturnTotal := RoundTaxAmount((cappedHousing + cappedMortgage) * marginalRate / 100)
	propertyReturnThisYear := min(propertyReturnTotal, annualTaxPaid)

	// 3. Социальный вычет (ст. 219)
	cappedSocial := min(input.SocialExpense, uint64(socialDeductionLimit))
	cappedChildEdu := min(input.ChildEduExpense, uint64(childEduDeductionLimit))
	socialReturn := RoundTaxAmount((cappedSocial + cappedChildEdu) * marginalRate / 100)
	socialReturn = min(socialReturn, annualTaxPaid)

	// 4. Итог: сумма возвратов не может превышать уплаченный НДФЛ
	totalReturn := min(childrenReturn+propertyReturnThisYear+socialReturn, annualTaxPaid)

	return &pb.DeductionResult{
		ChildrenMonthlyDeduction: monthlyDeduction,
		ChildrenMonths:           childrenMonths,
		ChildrenReturn:           childrenReturn,
		PropertyReturnThisYear:   propertyReturnThisYear,
		PropertyReturnTotal:      propertyReturnTotal,
		SocialReturn:             socialReturn,
		TotalReturn:              totalReturn,
	}
}

// childrenMonthlyDeduction возвращает суммарный ежемесячный вычет на детей в копейках.
// Ставки за каждого ребёнка: 1-й — 1 400 ₽, 2-й — 2 800 ₽, 3-й и далее — 6 000 ₽.
// Дополнительно: 12 000 ₽ за каждого ребёнка-инвалида.
func childrenMonthlyDeduction(count uint32, disabledCount uint32) uint64 {
	var total uint64
	for i := uint32(1); i <= count; i++ {
		switch i {
		case 1:
			total += child1Deduction
		case 2:
			total += child2Deduction
		default:
			total += childOtherDeduction
		}
	}
	total += uint64(disabledCount) * disabledChildDeduction
	return total
}
