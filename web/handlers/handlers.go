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
	"github.com/kiselevos/new_tax/web/internal/middleware"
	"google.golang.org/grpc/metadata"
)

type Server struct {
	Tmpl      *template.Template
	TaxClient pb.TaxServiceClient
}

func NewServer(tmpl *template.Template, client pb.TaxServiceClient) *Server {
	return &Server{
		Tmpl:      tmpl,
		TaxClient: client,
	}
}

func (s *Server) Routes(mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			s.NotFound(w, r)
			return
		}
		s.Index(w, r)
	})
	mux.HandleFunc("/calculate", s.Calculate)
	mux.HandleFunc("/about", s.About)
	mux.HandleFunc("/regional-info", s.RegionalInfo)
	mux.HandleFunc("/special-tax-modes", s.SpecialTaxModes)
	mux.HandleFunc("/api-docs", s.HandleApiDocs)
	mux.HandleFunc("/robots.txt", s.GetRobots)
	mux.HandleFunc("/sitemap.xml", s.GetSitemap)

	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/favicon.ico")
	})
}

func (s *Server) Index(w http.ResponseWriter, r *http.Request) {
	data := PrepareIndexData()

	if err := s.Tmpl.ExecuteTemplate(w, "index", data); err != nil {
		logx.From(r.Context()).Error("template_render_failed", "page", "index", "err", err)
		http.Error(w, "internal server error", 500)
		return
	}
}

func (s *Server) Calculate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logx.From(ctx)

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

	rid := middleware.GetRID(ctx)
	md := metadata.New(map[string]string{"x-request-id": rid, "x-internal": "true"})
	rpcCtx, cancel := context.WithTimeout(metadata.NewOutgoingContext(ctx, md), 3*time.Second)
	defer cancel()

	res, err := s.TaxClient.CalculatePrivate(rpcCtx, req)
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
		AnnualPFR         uint64
		AnnualFOMS        uint64
		AnnualFSS         uint64
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
		AnnualPFR:         res.AnnualPFR,
		AnnualFOMS:        res.AnnualFOMS,
		AnnualFSS:         res.AnnualFSS,
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

func (s *Server) About(w http.ResponseWriter, r *http.Request) {
	if err := s.Tmpl.ExecuteTemplate(w, "about", nil); err != nil {
		logx.From(r.Context()).Error("template_render_failed", "page", "about", "err", err)
		http.Error(w, "internal server error", 500)
	}
}

func (s *Server) RegionalInfo(w http.ResponseWriter, r *http.Request) {
	if err := s.Tmpl.ExecuteTemplate(w, "regional_info", nil); err != nil {
		logx.From(r.Context()).Error("template_render_failed", "page", "regional_info", "err", err)
		http.Error(w, "internal server error", 500)
	}
}

func (s *Server) SpecialTaxModes(w http.ResponseWriter, r *http.Request) {
	if err := s.Tmpl.ExecuteTemplate(w, "special_tax_modes", nil); err != nil {
		logx.From(r.Context()).Error("template_render_failed", "page", "special_tax_modes", "err", err)
		http.Error(w, "internal server error", 500)
	}
}

func (s *Server) HandleApiDocs(w http.ResponseWriter, r *http.Request) {

	data, err := PrepareApiData()
	if err != nil {
		http.Error(w, "cannot load api docs", 500)
		return
	}

	if err := s.Tmpl.ExecuteTemplate(w, "api-docs", data); err != nil {
		logx.From(r.Context()).Error("template_render_failed", "err", err)
		http.Error(w, "internal server error", 500)
	}
}

func (s *Server) GetRobots(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	http.ServeFile(w, r, "static/robots.txt")
}

func (s *Server) GetSitemap(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	http.ServeFile(w, r, "static/sitemap.xml")
}

func (s *Server) NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	if err := s.Tmpl.ExecuteTemplate(w, "404", nil); err != nil {
		logx.From(r.Context()).Error("template_render_failed", "page", "404", "err", err)
		_, _ = w.Write([]byte("404 page not found"))
	}
}

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
