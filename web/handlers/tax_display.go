package handlers

import "time"

// TaxConstants содержит налоговые константы для отображения на информационных страницах.
// Актуальные значения и источники: project-docs/TAX_CONSTANTS.md
// При изменении налогового законодательства обновлять этот файл вместе с internal/calculate/data.go
type TaxConstants struct {
	CurrentYear int

	// НДФЛ — прогрессивная шкала 2025 (ст. 224 НК РФ, ФЗ-176 от 12.07.2024)
	NdflRate1          int    // 13% — до 2 400 000 ₽
	NdflRate2          int    // 15% — 2 400 001 — 5 000 000 ₽
	NdflRate3          int    // 18% — 5 000 001 — 20 000 000 ₽
	NdflRate4          int    // 20% — 20 000 001 — 50 000 000 ₽
	NdflRate5          int    // 22% — свыше 50 000 000 ₽
	NdflThreshold1     string // "2 400 000"
	NdflThreshold2     string // "5 000 000"
	NdflThreshold3     string // "20 000 000"
	NdflThreshold4     string // "50 000 000"
	NdflRateNonResident int   // 30%

	// Страховые взносы работодателя (ст. 425 НК РФ)
	PfrRate     int    // 22%
	PfrRateHi   int    // 10% — сверх лимита
	PfrLimit    string // "2 979 000"
	FomsRateStr string // "5,1"
	FssRateStr  string // "2,9"

	// НПД — самозанятые (Федеральный закон № 422-ФЗ)
	NpdRateIndividual    int    // 4%
	NpdRateLegal         int    // 6%
	NpdLimit             string // "2 400 000"
	NpdRegistrationBonus string // "10 000"

	// Фиксированные взносы ИП «за себя» (ст. 430 НК РФ, ФЗ-389 от 31.07.2023)
	// Обновлять ежегодно: на 2026 год — 57 390 ₽ (ФЗ-425 от 28.11.2025)
	IpFixedContrib        string // текущий год
	IpAdditionalRate      int    // 1% с дохода свыше 300 000 ₽
	IpAdditionalThreshold string // "300 000"

	// Имущественные вычеты (ст. 220 НК РФ)
	PropertyDeductionMax string // "2 000 000"
	MortgageDeductionMax string // "3 000 000"
	PropertyMaxReturn    string // "260 000"  (2 000 000 × 13%)
	MortgageMaxReturn    string // "390 000"  (3 000 000 × 13%)

	// Социальные вычеты (ст. 219 НК РФ)
	SocialDeductionMax   string // "150 000"
	SocialMaxReturn      string // "19 500"   (150 000 × 13%)
	ChildEduDeductionMax string // "110 000"

	// Стандартные вычеты на детей (ст. 218 НК РФ)
	ChildDeduction1        string // "1 400" — первый ребёнок
	ChildDeduction2        string // "2 800" — второй ребёнок
	ChildDeduction3        string // "6 000" — третий и последующие (с 2025 года)
	ChildDisabledDeduction string // "12 000" — ребёнок-инвалид
	ChildDeductionLimit    string // "450 000" — лимит дохода для вычета (с 2025 года)
}

// PrepareTaxConstants возвращает актуальные налоговые константы для шаблонов.
func PrepareTaxConstants() TaxConstants {
	return TaxConstants{
		CurrentYear: time.Now().Year(),

		NdflRate1: 13, NdflRate2: 15, NdflRate3: 18, NdflRate4: 20, NdflRate5: 22,
		NdflThreshold1:      "2 400 000",
		NdflThreshold2:      "5 000 000",
		NdflThreshold3:      "20 000 000",
		NdflThreshold4:      "50 000 000",
		NdflRateNonResident: 30,

		PfrRate: 22, PfrRateHi: 10, PfrLimit: "2 979 000",
		FomsRateStr: "5,1", FssRateStr: "2,9",

		NpdRateIndividual:    4,
		NpdRateLegal:         6,
		NpdLimit:             "2 400 000",
		NpdRegistrationBonus: "10 000",

		IpFixedContrib:        "57 390",
		IpAdditionalRate:      1,
		IpAdditionalThreshold: "300 000",

		PropertyDeductionMax: "2 000 000",
		MortgageDeductionMax: "3 000 000",
		PropertyMaxReturn:    "260 000",
		MortgageMaxReturn:    "390 000",

		SocialDeductionMax:   "150 000",
		SocialMaxReturn:      "19 500",
		ChildEduDeductionMax: "110 000",

		ChildDeduction1:        "1 400",
		ChildDeduction2:        "2 800",
		ChildDeduction3:        "6 000",
		ChildDisabledDeduction: "12 000",
		ChildDeductionLimit:    "450 000",
	}
}
