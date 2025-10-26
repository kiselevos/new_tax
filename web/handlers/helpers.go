package handlers

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	pb "github.com/kiselevos/new_tax/gen/grpc/api"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CoefficientOption struct {
	Value int
	Label string
}

type BonusOption struct {
	Value int
	Label string
}

type Month struct {
	Value string
	Label string
}

type IndexData struct {
	CurrentYear int
	Months      []Month
	Territorial []CoefficientOption
	Northern    []BonusOption
	MinSalary   int
	LivingWage  int
}

func PrepareMonths() []Month {
	return []Month{
		{Value: "01", Label: "Январь"},
		{Value: "02", Label: "Февраль"},
		{Value: "03", Label: "Март"},
		{Value: "04", Label: "Апрель"},
		{Value: "05", Label: "Май"},
		{Value: "06", Label: "Июнь"},
		{Value: "07", Label: "Июль"},
		{Value: "08", Label: "Август"},
		{Value: "09", Label: "Сентябрь"},
		{Value: "10", Label: "Октябрь"},
		{Value: "11", Label: "Ноябрь"},
		{Value: "12", Label: "Декабрь"},
	}
}

func PrepareIndexData() IndexData {
	var territorial []CoefficientOption
	for i := 110; i <= 200; i += 5 {
		territorial = append(territorial, CoefficientOption{i, fmt.Sprintf("x%.2f", float64(i)/100)})
	}

	var northern []BonusOption
	for i := 10; i <= 100; i += 10 {
		northern = append(northern, BonusOption{100 + i, fmt.Sprintf("%d%%", i)})
	}

	minSalaryStr := os.Getenv("MIN_ALLOWED_SALARY")
	livingWageStr := os.Getenv("MIN_LIVING_WAGE")

	minSalary, err := strconv.Atoi(minSalaryStr)
	if err != nil || minSalary == 0 {
		minSalary = 4000
	}

	livingWage, err := strconv.Atoi(livingWageStr)
	if err != nil || livingWage == 0 {
		livingWage = 265000
	}

	return IndexData{
		CurrentYear: time.Now().Year(),
		Months:      PrepareMonths(),
		Territorial: territorial,
		Northern:    northern,
		MinSalary:   minSalary,
		LivingWage:  livingWage,
	}
}

func ParseFormToRequest(r *http.Request) (*pb.CalculatePrivateRequest, error) {
	// Логируем все полученные значения формы для отладки
	log.Printf("🔍 All form values: %+v", r.Form)

	// Получаем значение зарплаты
	rawSalary := r.FormValue("grossSalary")
	log.Printf("💰 Raw salary value: '%s'", rawSalary)

	if rawSalary == "" {
		return nil, fmt.Errorf("gross salary is required")
	}

	// Убираем ВСЕ пробелы (включая неразрывные) и заменяем запятые на точки
	rawSalary = strings.ReplaceAll(rawSalary, "\u00A0", "") // неразрывный пробел
	rawSalary = strings.ReplaceAll(rawSalary, " ", "")      // обычный пробел
	rawSalary = strings.ReplaceAll(rawSalary, ",", ".")     // запятая на точку

	log.Printf("💰 Cleaned salary value: '%s'", rawSalary)

	// Парсим зарплату
	salaryFloat, err := strconv.ParseFloat(rawSalary, 64)
	if err != nil {
		log.Printf("❌ Salary parsing error: %v, raw value: '%s'", err, rawSalary)
		return nil, fmt.Errorf("invalid gross salary format: '%s'. Use numbers only (e.g., 50000 or 50000.50)", rawSalary)
	}

	grossSalary := uint64(math.Round(salaryFloat * 100))
	log.Printf("✅ Parsed salary: %.2f -> %d", salaryFloat, grossSalary)

	// Валидация зарплаты
	check := ValidateSalary(grossSalary)
	if !check.Valid {
		return nil, fmt.Errorf(check.Message)
	}
	if check.ShowWarning {
		log.Println(check.Message)
	}

	// Получаем остальные значения с значениями по умолчанию
	monthStr := r.FormValue("startDate")
	territorialStr := r.FormValue("territorialMultiplier")
	northernStr := r.FormValue("northernCoefficient")
	hasTaxPrivilege := r.FormValue("hasTaxPrivilege") != ""
	isNotResident := r.FormValue("isNotResident") != ""

	log.Printf("📋 Other form values: startDate=%s, territorial=%s, northern=%s, hasTaxPrivilege=%t, isNotResident=%t",
		monthStr, territorialStr, northernStr, hasTaxPrivilege, isNotResident)

	// Обработка месяца
	monthNum, err := strconv.Atoi(monthStr)
	if err != nil || monthNum < 1 || monthNum > 12 {
		monthNum = 1 // значение по умолчанию
		log.Printf("⚠️  Invalid month, using default: 1")
	}

	startDate := time.Date(time.Now().Year(), time.Month(monthNum), 1, 0, 0, 0, 0, time.UTC)
	startTS := timestamppb.New(startDate)

	// Обработка коэффициентов с валидацией
	territorial := 100
	if territorialStr != "" {
		if v, err := strconv.Atoi(territorialStr); err == nil && v >= 100 && v <= 200 {
			territorial = v
		} else {
			log.Printf("⚠️  Invalid territorial multiplier: %s, using default: 100", territorialStr)
		}
	}

	northern := 100
	if northernStr != "" {
		if v, err := strconv.Atoi(northernStr); err == nil && v >= 100 && v <= 200 {
			northern = v
		} else {
			log.Printf("⚠️  Invalid northern coefficient: %s, using default: 100", northernStr)
		}
	}

	log.Printf("📄 Form parsed successfully: GrossSalary=%d, Territorial=%d, Northern=%d, HasTaxPrivilege=%t, IsNotResident=%t, StartDate=%s",
		grossSalary, territorial, northern, hasTaxPrivilege, isNotResident,
		startDate.Format("2006-01-02"))

	return &pb.CalculatePrivateRequest{
		GrossSalary:           grossSalary,
		StartDate:             startTS,
		TerritorialMultiplier: uint64Ptr(uint64(territorial)),
		NorthernCoefficient:   uint64Ptr(uint64(northern)),
		HasTaxPrivilege:       boolPtr(hasTaxPrivilege),
		IsNotResident:         boolPtr(isNotResident),
	}, nil
}

// SalaryValidationResult — структура результата проверки
type SalaryValidationResult struct {
	Valid       bool
	ShowWarning bool
	Message     string
}

// ValidateSalary — проверяет оклад по бизнес-правилам
func ValidateSalary(grossSalary uint64) SalaryValidationResult {
	minWageStr := os.Getenv("MIN_LIVING_WAGE")
	minAllowedSalary := os.Getenv("MIN_ALLOWED_SALARY")
	minWage, err := strconv.ParseUint(minWageStr, 10, 64)
	if err != nil || minWage == 0 {
		minWage = 2244000
	}
	minSalary, err := strconv.ParseUint(minAllowedSalary, 10, 64)
	if err != nil || minSalary == 0 {
		minSalary = 500000
	}

	log.Println("FFFFFF")

	// 1. Ошибка, если меньше 5000 ₽
	if grossSalary < minSalary {
		return SalaryValidationResult{
			Valid:       false,
			ShowWarning: false,
			Message:     fmt.Sprintf("❌ Минимальная сумма оклада — %d ₽", minSalary/100),
		}
	}

	// 2. Предупреждение, если меньше прожиточного минимума
	if grossSalary < minWage {
		log.Printf("⚠️  Предупреждение: сумма %d меньше прожиточного минимума (%d)", grossSalary/100, minWage)
		return SalaryValidationResult{
			Valid:       true,
			ShowWarning: true,
			Message:     "⚠️ Сумма меньше прожиточного минимума, что неприемлемо при полной занятости.",
		}
	}

	return SalaryValidationResult{Valid: true}
}

// Вспомогательные функции:
func uint64Ptr(v uint64) *uint64 { return &v }
func boolPtr(v bool) *bool       { return &v }
