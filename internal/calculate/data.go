package calculate

import "math"

// TaxLimit описывает одну ступень прогрессивной налоговой шкалы.
// Доход до UpperLimit облагается по ставке Rate (в процентах).
type TaxLimit struct {
	UpperLimit uint64
	Rate       uint64
}

// Limits содержит прогрессивную шкалу налогообложения в РФ с 2025 года.
// Применяется поэтапно: каждая часть дохода до указанного порога облагается по своей ставке.
var Limits = []TaxLimit{
	{UpperLimit: 2_400_000_00, Rate: 13},
	{UpperLimit: 5_000_000_00, Rate: 15},
	{UpperLimit: 20_000_000_00, Rate: 18},
	{UpperLimit: 50_000_000_00, Rate: 20},
	{UpperLimit: math.MaxUint64, Rate: 22},
}

// SimpleLimits упращенная шкала 13% до 5 млн рублей, 15% после
var SimpleLimits = TaxLimit{
	UpperLimit: 5_000_000_00,
	Rate:       13,
}

// NotResident налог для нерезедентов 30%
var NotResident = TaxLimit{
	UpperLimit: math.MaxUint64,
	Rate:       30,
}
