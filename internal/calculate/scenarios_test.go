package calculate

import "time"

type Scenario struct {
	Name       string
	StartMonth time.Month

	Salary                uint64 // оклад в копейках
	TerritorialMultiplier uint64 // территориальный коэффициент, 20 => +20%
	NorthernCoefficient   uint64 // север, 50 => +50%
	Privilege             bool   // льготы силовых (для соответствующей функции)
	NonResident           bool   // нерезидент
}

func (s Scenario) StartDate() time.Time {
	y := time.Now().UTC().Year() // год берём текущий, т.к. на расчёт не влияет
	return time.Date(y, s.StartMonth, 1, 0, 0, 0, 0, time.UTC)
}

func (s Scenario) Bases() (baseA, baseB uint64) {
	baseA = s.Salary * (100 + s.TerritorialMultiplier) / 100
	baseB = baseA * s.NorthernCoefficient / 100
	return
}

var Scenarios = []Scenario{
	{
		Name:                  "A: Only salary usual person (100k, terr 20%, jan start)",
		StartMonth:            time.January,
		Salary:                100_000_00, // 100 000 ₽
		TerritorialMultiplier: 20,         // база A = 120 000 ₽
		NorthernCoefficient:   0,
		Privilege:             false,
		NonResident:           false,
	},
	{
		Name:                  "B: 200k salary, no coeffs, at 2.4M threshold (jan start)",
		StartMonth:            time.January,
		Salary:                200_000_00, // 200 000 ₽ → годом ровно 2.4M → 13% весь год
		TerritorialMultiplier: 0,
		NorthernCoefficient:   0,
		Privilege:             false,
		NonResident:           false,
	},
	{
		Name:                  "C: 300k salary, cross 2.4M in Sep (jan start)",
		StartMonth:            time.January,
		Salary:                300_000_00, // 300 000 ₽ → с сент. часть по 15%
		TerritorialMultiplier: 0,
		NorthernCoefficient:   0,
		Privilege:             false,
		NonResident:           false,
	},
	{
		Name:                  "D: 200k salary + terr 20% + north 50% (jan start)",
		StartMonth:            time.January,
		Salary:                200_000_00, // база A = 200k * 1.20 = 240k; север B = 240k * 0.50 = 120k
		TerritorialMultiplier: 20,
		NorthernCoefficient:   50,
		Privilege:             false,
		NonResident:           false,
	},
	{
		Name:                  "E: 120k salary + terr 20% + north 50% (start in Jun, projection 7m)",
		StartMonth:            time.June,
		Salary:                120_000_00, // база A = 144k; север B = 72k; 13% весь период
		TerritorialMultiplier: 20,
		NorthernCoefficient:   50,
		Privilege:             false,
		NonResident:           false,
	},
	{
		Name:                  "NR-1: Non-resident, tiny salary → monthly rounding to 0 (jan start)",
		StartMonth:            time.January,
		Salary:                1_00, // 1 ₽ → 30% = 30 коп → округляется до 0 руб помесячно
		TerritorialMultiplier: 0,
		NorthernCoefficient:   0,
		Privilege:             false,
		NonResident:           true,
	},
	{
		Name:                  "NR-2: Non-resident, 120k salary (start in Jun, 7 months)",
		StartMonth:            time.June,
		Salary:                120_000_00, // 120 000 ₽ → 36 000 ₽/мес налог (до округл. это кратно 100)
		TerritorialMultiplier: 0,
		NorthernCoefficient:   0,
		Privilege:             false,
		NonResident:           true,
	},
}
