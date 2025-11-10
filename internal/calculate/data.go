package calculate

import "math"

const (
	PfrLimit  = 275900000 // Ставка 22% до лимита в 2 759 000 рублей, после 10%
	FssLimit  = 132500000 // Ставка 2,9% до лимита в 1 325 000 рублей, после 0%
	PfrRate   = 220       // 22% (Делим на 1000)
	PfrRateHi = 100       // 10% (Делим на 1000)
	FomsRate  = 51        // 5.1% (Делим на 1000)
	FssRate   = 29        // 2.9% (Делим на 1000)
)

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
