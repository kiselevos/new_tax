package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/kiselevos/new_tax/web/handlers"
)

func main() {
	// --- 1. Шаблонные функции для HTML ---
	funcs := template.FuncMap{
		// Генерация последовательности чисел (например: 1.10, 1.15, 1.20 ...)
		"seq": func(start, end, step float64) []float64 {
			var result []float64
			for v := start; v <= end+0.001; v += step {
				result = append(result, v)
			}
			return result
		},
		// Генерация целых чисел от 0 до n-1
		"until": func(n int) []int {
			arr := make([]int, n)
			for i := range arr {
				arr[i] = i
			}
			return arr
		},
		// Простейшие арифметические операции
		"add": func(a, b int) int { return a + b },
		"mul": func(a, b int) int { return a * b },
	}

	tmpl := template.Must(template.New("").Funcs(funcs).ParseFiles(
		"templates/index.tmpl",
		"templates/result.tmpl",
		"templates/how_it_works.tmpl",
		"templates/regional_info.tmpl",
		"templates/special_tax_modes.tmpl",
	))

	// --- 3. Инициализация сервера и маршрутов ---
	s := &handlers.Server{Tmpl: tmpl}
	s.Routes()

	// --- 4. Статика (CSS, изображения и т.п.) ---
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// --- 5. Запуск сервера ---
	log.Println("🌐 Web server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
