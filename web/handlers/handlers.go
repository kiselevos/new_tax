package handlers

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	pb "github.com/kiselevos/new_tax/gen/grpc/api"
	"github.com/kiselevos/new_tax/web/internal/client"
)

// Server — основной HTTP-сервер, хранящий шаблоны.
type Server struct {
	Tmpl *template.Template
}

// Month — структура для списка месяцев в форме.
type Month struct {
	Value string
	Label string
}

// loggingMiddleware — базовое логирование всех запросов.
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("🛠️  %s %s", r.Method, r.URL.Path)
		next(w, r)
		log.Printf("✅ Завершено за %v\n", time.Since(start))
	}
}

// Routes — регистрация всех маршрутов приложения.
func (s *Server) Routes() {
	http.HandleFunc("/", loggingMiddleware(s.Index))
	http.HandleFunc("/calculate", loggingMiddleware(s.Calculate))
	http.HandleFunc("/how-it-works", loggingMiddleware(s.HowItWorks))
	http.HandleFunc("/regional-info", loggingMiddleware(s.RegionalInfo))
	http.HandleFunc("/special-tax-modes", loggingMiddleware(s.SpecialTaxModes))
}

// Index — стартовая страница с формой расчёта.
func (s *Server) Index(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"CurrentYear": time.Now().Year(),
		"Months": []Month{
			{Value: "01", Label: "Январь"}, {Value: "02", Label: "Февраль"},
			{Value: "03", Label: "Март"}, {Value: "04", Label: "Апрель"},
			{Value: "05", Label: "Май"}, {Value: "06", Label: "Июнь"},
			{Value: "07", Label: "Июль"}, {Value: "08", Label: "Август"},
			{Value: "09", Label: "Сентябрь"}, {Value: "10", Label: "Октябрь"},
			{Value: "11", Label: "Ноябрь"}, {Value: "12", Label: "Декабрь"},
		},
	}

	if err := s.Tmpl.ExecuteTemplate(w, "index", data); err != nil {
		log.Printf("❌ Ошибка рендеринга index: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Calculate — обработка формы и запрос к gRPC-бэкенду.
func (s *Server) Calculate(w http.ResponseWriter, r *http.Request) {
	req, err := ParseFormToRequest(r)
	if err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	clientGRPC, conn, err := client.NewTaxClient()
	if err != nil {
		http.Error(w, "can't connect to backend", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	log.Println("→ sending request to backend")

	res, err := clientGRPC.CalculatePrivate(ctx, req)
	if err != nil {
		log.Printf("backend error: %v", err)
		http.Error(w, "backend error", http.StatusInternalServerError)
		return
	}

	data := struct {
		AnnualTaxAmount   uint64
		AnnualGrossIncome uint64
		AnnualNetIncome   uint64
		GrossSalary       uint64
		TerritorialMult   uint64
		NorthernCoeff     uint64
		MonthlyDetails    []*pb.MonthlyPrivateTax
	}{
		AnnualTaxAmount:   res.AnnualTaxAmount,
		AnnualGrossIncome: res.AnnualGrossIncome,
		AnnualNetIncome:   res.AnnualNetIncome,
		GrossSalary:       res.GrossSalary,
		TerritorialMult:   deref(res.TerritorialMultiplier),
		NorthernCoeff:     deref(res.NorthernCoefficient),
		MonthlyDetails:    res.MonthlyDetails,
	}

	s.Tmpl.ExecuteTemplate(w, "result", data)
}

// HowItWorks — страница с объяснением логики расчёта.
func (s *Server) HowItWorks(w http.ResponseWriter, r *http.Request) {
	log.Println("📄 Рендеринг страницы: how_it_works")
	if err := s.Tmpl.ExecuteTemplate(w, "how_it_works", nil); err != nil {
		log.Printf("❌ Ошибка рендеринга how_it_works: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// RegionalInfo — страница с информацией о коэффициентах и надбавках.
func (s *Server) RegionalInfo(w http.ResponseWriter, r *http.Request) {
	log.Println("📄 Рендеринг страницы: regional_info")
	if err := s.Tmpl.ExecuteTemplate(w, "regional_info", nil); err != nil {
		log.Printf("❌ Ошибка рендеринга regional_info: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// SpecialTaxModes — страница об особых налоговых режимах.
func (s *Server) SpecialTaxModes(w http.ResponseWriter, r *http.Request) {
	log.Println("📄 Рендеринг страницы: special_tax_modes")
	if err := s.Tmpl.ExecuteTemplate(w, "special_tax_modes", nil); err != nil {
		log.Printf("❌ Ошибка рендеринга special_tax_modes: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Вспомогательные функции
// ------------------------------------------------------------

func parseUint(s string) uint64 {
	val, _ := strconv.ParseUint(s, 10, 64)
	return val
}

func deref(p *uint64) uint64 {
	if p == nil {
		return 0
	}
	return *p
}
