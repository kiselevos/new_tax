package web

import (
	"fmt"
	"html/template"
	"strings"
)

var Funcs = template.FuncMap{
	"fmtMoney": formatMoney,
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
