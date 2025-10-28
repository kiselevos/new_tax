package handlers

import (
	"context"
	"html/template"
	"net/http"
	"strings"
	"time"

	pb "github.com/kiselevos/new_tax/gen/grpc/api"
	"github.com/kiselevos/new_tax/pkg/logx"
	"github.com/kiselevos/new_tax/web"
	"github.com/kiselevos/new_tax/web/internal/client"
	"github.com/kiselevos/new_tax/web/internal/middleware"
	"google.golang.org/grpc/metadata"
)

// Server — основной HTTP-сервер, хранящий шаблоны.
type Server struct {
	Tmpl *template.Template
}

// Routes — регистрация всех маршрутов приложения.
func (s *Server) Routes(mux *http.ServeMux) {
	mux.HandleFunc("/", s.Index)
	mux.HandleFunc("/calculate", s.Calculate)
	mux.HandleFunc("/how-it-works", s.HowItWorks)
	mux.HandleFunc("/regional-info", s.RegionalInfo)
	mux.HandleFunc("/special-tax-modes", s.SpecialTaxModes)
}

// Index — главная страница
func (s *Server) Index(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logx.From(ctx).With("page", "index")

	if err := s.Tmpl.ExecuteTemplate(w, "index", PrepareIndexData()); err != nil {
		log.Error("template_render_failed", "err", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	log.Info("page_rendered")
}

// Calculate — обработка формы и запрос к gRPC-бэкенду
func (s *Server) Calculate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logx.From(ctx).With("path", r.URL.Path, "method", r.Method)

	ct := r.Header.Get("Content-Type")
	if !strings.Contains(ct, "application/x-www-form-urlencoded") &&
		!strings.Contains(ct, "multipart/form-data") {
		log.Warn("unsupported_content_type", "content_type", ct)
		http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Warn("form_parse_failed", "err", err)
		http.Error(w, "invalid form data", http.StatusBadRequest)
		return
	}

	req, err := ParseFormToRequest(r)
	if err != nil {
		log.Warn("form_validation_failed", "err", err)
		http.Error(w, "invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	clientGRPC, conn, err := client.NewTaxClient()
	if err != nil {
		log.Error("grpc_dial_failed", "err", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	rid := middleware.GetRID(ctx)

	md := metadata.New(map[string]string{"x-request-id": rid})
	rpcCtx, cancel := context.WithTimeout(metadata.NewOutgoingContext(ctx, md), 3*time.Second)

	defer cancel()

	res, err := clientGRPC.CalculatePrivate(rpcCtx, req)
	if err != nil {
		log.Error("grpc_call_failed", "method", "CalculatePrivate", "err", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	minWage := web.GetMinLivingWage()
	showWarning := req.GrossSalary < minWage

	data := struct {
		AnnualTaxAmount   uint64
		AnnualGrossIncome uint64
		AnnualNetIncome   uint64
		GrossSalary       uint64
		TerritorialMult   uint64
		NorthernCoeff     uint64
		MonthlyDetails    []*pb.MonthlyPrivateTax
		ShowWarning       bool
		HasTaxPrivilege   bool
		IsNotResident     bool
	}{
		AnnualTaxAmount:   res.AnnualTaxAmount,
		AnnualGrossIncome: res.AnnualGrossIncome,
		AnnualNetIncome:   res.AnnualNetIncome,
		GrossSalary:       res.GrossSalary,
		TerritorialMult:   deref(res.TerritorialMultiplier),
		NorthernCoeff:     deref(res.NorthernCoefficient),
		MonthlyDetails:    res.MonthlyDetails,
		ShowWarning:       showWarning,
		HasTaxPrivilege:   getBool(req.HasTaxPrivilege),
		IsNotResident:     getBool(req.IsNotResident),
	}

	if err := s.Tmpl.ExecuteTemplate(w, "result", data); err != nil {
		log.Error("template_render_failed", "page", "result", "err", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	log.Info("calculation_completed",
		"gross_salary", req.GrossSalary,
		"annual_tax", res.AnnualTaxAmount,
		"warning", showWarning,
	)
}

// HowItWorks — страница с описанием расчёта
func (s *Server) HowItWorks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logx.From(ctx).With("page", "how_it_works")

	if err := s.Tmpl.ExecuteTemplate(w, "how_it_works", nil); err != nil {
		log.Error("template_render_failed", "err", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	log.Info("page_rendered")
}

// RegionalInfo — страница с региональными коэффициентами
func (s *Server) RegionalInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logx.From(ctx).With("page", "regional_info")

	if err := s.Tmpl.ExecuteTemplate(w, "regional_info", nil); err != nil {
		log.Error("template_render_failed", "err", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	log.Info("page_rendered")
}

// SpecialTaxModes — страница о льготах и нерезидентстве
func (s *Server) SpecialTaxModes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logx.From(ctx).With("page", "special_tax_modes")

	if err := s.Tmpl.ExecuteTemplate(w, "special_tax_modes", nil); err != nil {
		log.Error("template_render_failed", "err", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	log.Info("page_rendered")
}

// ----- helpers -----

func deref(p *uint64) uint64 {
	if p == nil {
		return 0
	}
	return *p
}

func getBool(p *bool) bool {
	if p == nil {
		return false
	}
	return *p
}
