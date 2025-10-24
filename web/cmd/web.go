package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/kiselevos/new_tax/web/handlers"
)

func formatMoney(amount uint64) string {
	rubles := float64(amount) / 100
	s := fmt.Sprintf("%.2f", rubles)

	parts := strings.Split(s, ".")
	intPart := parts[0]
	fracPart := parts[1]

	// добавляем пробелы каждые 3 цифры
	for i := len(intPart) - 3; i > 0; i -= 3 {
		intPart = intPart[:i] + " " + intPart[i:] // тонкий пробел (U+202F)
	}

	return intPart + "," + fracPart + " ₽"
}

func main() {
	funcs := template.FuncMap{

		"fmtMoney": formatMoney,
	}

	tmpls, err := template.New("").Funcs(funcs).ParseGlob("templates/*.tmpl")
	if err != nil {
		log.Fatalf("ошибка загрузки шаблонов: %v", err)
	}

	s := &handlers.Server{Tmpl: tmpls}
	s.Routes()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("🌐 Web server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
