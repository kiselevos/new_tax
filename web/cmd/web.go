package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/kiselevos/new_tax/web/handlers"
)

func main() {
	funcs := template.FuncMap{

		"seq": func(start, end, step float64) []float64 {
			var result []float64
			for v := start; v <= end+0.001; v += step {
				result = append(result, v)
			}
			return result
		},

		"until": func(n int) []int {
			arr := make([]int, n)
			for i := range arr {
				arr[i] = i
			}
			return arr
		},

		"add": func(a, b int) int { return a + b },
		"mul": func(a, b int) int { return a * b },
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
