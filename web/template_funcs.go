package web

import (
	"fmt"
	"html/template"
	"os"
	"strconv"
	"strings"
)

var Funcs = template.FuncMap{
	"fmtMoney":         formatMoney,
	"getMinSalary":     GetMinSalary,
	"getMinLivingWage": GetMinLivingWage,
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
