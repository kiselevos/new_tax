package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/kiselevos/new_tax/web"
	"github.com/kiselevos/new_tax/web/handlers"
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Не найден или не читается web/.env: %v", err)
	}

	tmpls, err := template.New("").Funcs(web.Funcs).ParseGlob("templates/*.tmpl")
	if err != nil {
		log.Fatalf("ошибка загрузки шаблонов: %v", err)
	}

	s := &handlers.Server{Tmpl: tmpls}
	s.Routes()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Fatal(http.ListenAndServe(os.Getenv("WEB_PORT"), nil))
}
