package web

import (
	"fmt"
	"html/template"
	"os"
	"strconv"
	"strings"

	"google.golang.org/protobuf/types/known/timestamppb"
)

var Funcs = template.FuncMap{
	"fmtMoney":         formatMoney,
	"fmtMoneyRaw":      formatMoneyRaw,
	"monthNum":         monthNum,
	"getMinSalary":     GetMinSalary,
	"getMinLivingWage": GetMinLivingWage,
	"getFeedbackEmail": GetFeedbackEmail,
	"russianMonth":     formatRussianMonth,
	"sub": func(a, b int) int {
		return a - b
	},
	"minus100": func(a uint64) uint64 { return a - 100 },
	"divf":     func(a uint64, b float64) float64 { return float64(a) / b },
	"toInt":    func(n uint64) int { return int(n) },
	"sum":      sum,
}

func sum(nums ...uint64) uint64 {
	var total uint64
	for _, n := range nums {
		total += n
	}
	return total
}

func formatMoney(amount uint64) string {
	rubles := float64(amount) / 100
	s := fmt.Sprintf("%.2f", rubles)

	parts := strings.Split(s, ".")
	intPart := parts[0]
	fracPart := parts[1]

	for i := len(intPart) - 3; i > 0; i -= 3 {
		intPart = intPart[:i] + " " + intPart[i:]
	}

	return intPart + "," + fracPart + " ₽"
}

func GetMinSalary() float64 {
	minAllowedSalaryStr := os.Getenv("MIN_ALLOWED_SALARY")
	minSalary, err := strconv.ParseFloat(minAllowedSalaryStr, 64)
	if err != nil || minSalary == 0 {
		minSalary = 1
	}
	return minSalary
}

func GetMinLivingWage() uint64 {
	minWageStr := os.Getenv("MIN_LIVING_WAGE")
	minWage, err := strconv.Atoi(minWageStr)
	if err != nil || minWage == 0 {
		minWage = 1
	}
	return uint64(minWage)
}

func GetFeedbackEmail() string {
	feedbackEmail := os.Getenv("FEEDBACK_EMAIL")
	if feedbackEmail == "" {
		return "okiselev421@gmail.com"
	}
	return feedbackEmail
}

func GetApiVersion() string {
	apiVersion := os.Getenv("API_VERSION")
	if apiVersion == "" {
		return "v1"
	}
	return apiVersion
}

// formatMoneyRaw переводит копейки в целые рубли для подстановки в поля формы.
// Например: 5000000 → "50000"
func formatMoneyRaw(kopecks uint64) string {
	return strconv.FormatUint(kopecks/100, 10)
}

// monthNum возвращает номер месяца (1-12) из timestamp.
func monthNum(ts *timestamppb.Timestamp) int {
	if ts == nil {
		return 0
	}
	return int(ts.AsTime().Month())
}

// Функция для форматирования месяца на русском из timestamp
func formatRussianMonth(ts *timestamppb.Timestamp) string {
	if ts == nil {
		return "Неизвестный месяц"
	}

	t := ts.AsTime()
	return getRussianMonthName(int(t.Month())) + " " + fmt.Sprintf("%d", t.Year())
}

// Функция для получения русского названия месяца по номеру
func getRussianMonthName(monthNumber int) string {
	months := map[int]string{
		1:  "Январь",
		2:  "Февраль",
		3:  "Март",
		4:  "Апрель",
		5:  "Май",
		6:  "Июнь",
		7:  "Июль",
		8:  "Август",
		9:  "Сентябрь",
		10: "Октябрь",
		11: "Ноябрь",
		12: "Декабрь",
	}

	if name, exists := months[monthNumber]; exists {
		return name
	}
	return "Неизвестный месяц"
}
