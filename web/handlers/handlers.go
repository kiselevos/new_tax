package handlers

import (
	"context"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"time"

	pb "github.com/kiselevos/new_tax/gen/grpc/api"
	"github.com/kiselevos/new_tax/pkg/logx"
	"github.com/kiselevos/new_tax/web"
	"github.com/kiselevos/new_tax/web/internal/metrics"
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
	mux.HandleFunc("/tax-deductions", s.TaxDeductions)
	mux.HandleFunc("/employment-types", s.EmploymentTypes)
	mux.HandleFunc("/api-docs", s.HandleApiDocs)
	mux.HandleFunc("/robots.txt", s.GetRobots)
	mux.HandleFunc("/sitemap.xml", s.GetSitemap)

	// Метрики
	metrics.Mount(mux)

	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/favicon.ico")
	})
}

func (s *Server) Index(w http.ResponseWriter, r *http.Request) {
	data := PrepareIndexData()
	data.FormError = r.URL.Query().Get("error")

	if err := s.Tmpl.ExecuteTemplate(w, "index", data); err != nil {
		logx.From(r.Context()).Error("template_render_failed", "page", "index", "err", err)
		http.Error(w, "internal server error", 500)
		return
	}
}

func (s *Server) Calculate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logx.From(ctx)

	start := time.Now()
	region := middleware.GetRegion(ctx)
	client := "ui"

	metrics.M.Calculator.Attempts.
		WithLabelValues(client, region.Label).
		Inc()

	defer func() {
		metrics.M.Calculator.Duration.
			WithLabelValues(client, region.Label).
			Observe(time.Since(start).Seconds())
	}()

	// Validate content type
	ct := r.Header.Get("Content-Type")
	if !strings.Contains(ct, "application/x-www-form-urlencoded") &&
		!strings.Contains(ct, "multipart/form-data") {
		log.Warn("unsupported_content_type", "content_type", ct)
		metrics.M.ErrorTypes.WithLabelValues(client, "content_type").Inc()
		metrics.M.Calculator.Failed.WithLabelValues(client, region.Label).Inc()
		http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	// Parse form
	if err := r.ParseForm(); err != nil {
		log.Warn("form_parse_failed", "err", err)
		metrics.M.ErrorTypes.WithLabelValues(client, "parse_form").Inc()
		metrics.M.Calculator.Failed.WithLabelValues(client, region.Label).Inc()
		http.Error(w, "invalid form data", http.StatusBadRequest)
		return
	}

	// Validate business input
	req, err := ParseFormToRequest(r)
	if err != nil {
		log.Warn("form_validation_failed", "err", err)
		metrics.M.ErrorTypes.WithLabelValues(client, "validate").Inc()
		metrics.M.Calculator.Failed.WithLabelValues(client, region.Label).Inc()
		http.Redirect(w, r, "/?error="+url.QueryEscape("Введите корректный оклад (например, 50 000 ₽)"), http.StatusSeeOther)
		return
	}

	// gRPC request
	rid := middleware.GetRID(ctx)
	md := metadata.New(map[string]string{"x-request-id": rid, "x-internal": "true"})
	rpcCtx, cancel := context.WithTimeout(metadata.NewOutgoingContext(ctx, md), 3*time.Second)
	defer cancel()

	res, err := s.TaxClient.CalculatePrivate(rpcCtx, req)
	if err != nil {
		log.Error("grpc_call_failed", "method", "CalculatePrivate", "err", err)
		metrics.M.ErrorTypes.WithLabelValues(client, "grpc").Inc()
		metrics.M.Calculator.Failed.WithLabelValues(client, region.Label).Inc()
		http.Redirect(w, r, "/?error="+url.QueryEscape("Ошибка при расчёте. Попробуйте ещё раз"), http.StatusSeeOther)
		return
	}

	minWage := web.GetMinLivingWage()
	showWarning := req.GrossSalary < minWage

	// Суммируем премии за год и определяем месяц начала расчёта
	var annualBonus uint64
	for _, m := range res.MonthlyDetails {
		annualBonus += m.MonthlyBonus
	}
	startMonthNum := 1
	if len(res.MonthlyDetails) > 0 && res.MonthlyDetails[0].Month != nil {
		startMonthNum = int(res.MonthlyDetails[0].Month.AsTime().Month())
	}

	// Нормализуем бонусы в срез из 12 элементов для шаблона
	bonuses := make([]uint64, 12)
	for _, m := range res.MonthlyDetails {
		if m.Month != nil {
			idx := int(m.Month.AsTime().Month()) - 1
			bonuses[idx] = m.MonthlyBonus
		}
	}

	// Первый месяц без премии — для виджета «На руки в месяц»
	var baseMonth *pb.MonthlyPrivateTax
	for _, m := range res.MonthlyDetails {
		if m.MonthlyBonus == 0 {
			baseMonth = m
			break
		}
	}
	if baseMonth == nil && len(res.MonthlyDetails) > 0 {
		baseMonth = res.MonthlyDetails[0]
	}

	indexData := PrepareIndexData()

	// Prepare data
	data := ResultPayload{
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
		IsGPH:             req.GetEmploymentType() == pb.EmploymentType_GPH,
		EmploymentTypeStr: req.GetEmploymentType().String(),
		AnnualPFR:         res.AnnualPFR,
		AnnualFOMS:        res.AnnualFOMS,
		AnnualFSS:         res.AnnualFSS,
		MonthlyBonuses:    bonuses,
		AnnualBonus:       annualBonus,
		StartMonthNum:     startMonthNum,
		BaseMonth:         baseMonth,
		Months:            indexData.Months,
		Territorial:       indexData.Territorial,
		Northern:          indexData.Northern,
		DeductionResult:            res.DeductionResult,
		ChildrenCountInput:         req.GetChildrenCount(),
		DisabledChildrenCountInput: req.GetDisabledChildrenCount(),
		HousingExpenseInput:        req.GetHousingExpense() / 100,
		MortgageExpenseInput:       req.GetMortgageExpense() / 100,
		SocialExpenseInput:         req.GetSocialExpense() / 100,
		ChildEduExpenseInput:       req.GetChildEduExpense() / 100,
	}

	// render template
	if err := s.Tmpl.ExecuteTemplate(w, "result", data); err != nil {
		log.Error("template_render_failed", "page", "result", "err", err)
		metrics.M.ErrorTypes.WithLabelValues(client, "template").Inc()
		metrics.M.Calculator.Failed.WithLabelValues(client, region.Label).Inc()
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	gross := float64(data.GrossSalary) / 100.0

	log.Info("business_calc",
		"client", client,
		"region", region.Name,
		"rid", middleware.GetRID(ctx),

		"gross_salary_rub", gross,
		"territorial_multiplier", data.TerritorialMult,
		"northern_coefficient", data.NorthernCoeff,
		"has_tax_privilege", data.HasTaxPrivilege,
		"is_not_resident", data.IsNotResident,

		"annual_tax", data.AnnualTaxAmount,
	)

	metrics.M.Calculator.GrossSalary.
		WithLabelValues(client, region.Label).
		Observe(gross)

	// UI completed successfully
	metrics.M.Calculator.Success.WithLabelValues(client, region.Label).Inc()
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

func (s *Server) TaxDeductions(w http.ResponseWriter, r *http.Request) {
	if err := s.Tmpl.ExecuteTemplate(w, "tax_deductions", nil); err != nil {
		logx.From(r.Context()).Error("template_render_failed", "page", "tax_deductions", "err", err)
		http.Error(w, "internal server error", 500)
	}
}

func (s *Server) EmploymentTypes(w http.ResponseWriter, r *http.Request) {
	if err := s.Tmpl.ExecuteTemplate(w, "employment_types", nil); err != nil {
		logx.From(r.Context()).Error("template_render_failed", "page", "employment_types", "err", err)
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
