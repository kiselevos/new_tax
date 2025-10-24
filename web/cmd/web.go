package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/kiselevos/new_tax/web"
	"github.com/kiselevos/new_tax/web/handlers"
)

func main() {

	tmpls, err := template.New("").Funcs(web.Funcs).ParseGlob("templates/*.tmpl")
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
