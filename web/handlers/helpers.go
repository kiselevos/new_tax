package handlers

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	pb "github.com/kiselevos/new_tax/gen/grpc/api"
	"github.com/kiselevos/new_tax/pkg/logx"
	"github.com/kiselevos/new_tax/web"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CoefficientOption struct {
	Value int
	Label string
}

type BonusOption struct {
	Value int
	Label string
}

type Month struct {
	Value string
	Label string
}

type IndexData struct {
	CurrentYear   int
	FeedbackEmail string
	Months        []Month
	Territorial   []CoefficientOption
	Northern      []BonusOption
	FormError     string
}

type ResultPayload struct {
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
	MonthlyBonuses    []uint64              // 12 элементов, индекс 0 = январь (копейки)
	AnnualBonus       uint64                // сумма всех премий за год (копейки)
	StartMonthNum     int                   // номер месяца начала расчёта (1-12)
	BaseMonth         *pb.MonthlyPrivateTax // первый месяц без премии (для виджета «На руки»)
	Months            []Month               // опции для select месяца в панели редактирования
	Territorial       []CoefficientOption   // опции РК
	Northern          []BonusOption         // опции СН
	IsGPH             bool                  // договор ГПХ (нет ФСС, нет трудовых гарантий)
	EmploymentTypeStr string                // "TD" | "GPH" | "SELF_EMPLOYED" — для pre-select в форме
	DeductionResult   *pb.DeductionResult   // результат расчёта налоговых вычетов (если переданы параметры)

	// Значения из последнего запроса вычетов (для предзаполнения формы после пересчёта)
	ChildrenCountInput         uint32
	DisabledChildrenCountInput uint32
	HousingExpenseInput        uint64 // рубли (для input[type=number])
	MortgageExpenseInput       uint64
	SocialExpenseInput         uint64
	ChildEduExpenseInput       uint64
}

func PrepareMonths() []Month {
	return []Month{
		{Value: "01", Label: "Январь"},
		{Value: "02", Label: "Февраль"},
		{Value: "03", Label: "Март"},
		{Value: "04", Label: "Апрель"},
		{Value: "05", Label: "Май"},
		{Value: "06", Label: "Июнь"},
		{Value: "07", Label: "Июль"},
		{Value: "08", Label: "Август"},
		{Value: "09", Label: "Сентябрь"},
		{Value: "10", Label: "Октябрь"},
		{Value: "11", Label: "Ноябрь"},
		{Value: "12", Label: "Декабрь"},
	}
}

func PrepareIndexData() IndexData {
	var territorial []CoefficientOption
	for i := 110; i <= 200; i += 5 {
		territorial = append(territorial, CoefficientOption{i, fmt.Sprintf("x%.2f", float64(i)/100)})
	}

	var northern []BonusOption
	for i := 10; i <= 100; i += 10 {
		northern = append(northern, BonusOption{100 + i, fmt.Sprintf("%d%%", i)})
	}

	feedbackEmail := web.GetFeedbackEmail()

	return IndexData{
		CurrentYear:   time.Now().Year(),
		FeedbackEmail: feedbackEmail,
		Months:        PrepareMonths(),
		Territorial:   territorial,
		Northern:      northern,
	}
}

// ParseFormToRequest - парсит данные из формы в gRPC-запрос.
func ParseFormToRequest(r *http.Request) (*pb.CalculatePrivateRequest, error) {
	log := logx.From(r.Context()).With("component", "form_parser")

	rawSalary := r.FormValue("grossSalary")
	if rawSalary == "" {
		log.Warn("form_missing_field", "field", "grossSalary")
		return nil, fmt.Errorf("gross salary is required")
	}

	// Очистка от пробелов и замена запятых на точки
	rawSalary = strings.ReplaceAll(rawSalary, "\u00A0", "")
	rawSalary = strings.ReplaceAll(rawSalary, " ", "")
	rawSalary = strings.ReplaceAll(rawSalary, ",", ".")

	salaryFloat, err := strconv.ParseFloat(rawSalary, 64)
	if err != nil {
		log.Warn("form_invalid_salary", "raw", rawSalary, "err", err)
		return nil, fmt.Errorf("invalid gross salary format: '%s'. Use only numbers (e.g., 50000 or 50000.50)", rawSalary)
	}
	grossSalary := uint64(math.Round(salaryFloat * 100))

	// Извлекаем остальные поля
	monthStr := r.FormValue("startDate")
	territorialStr := r.FormValue("territorialMultiplier")
	northernStr := r.FormValue("northernCoefficient")
	hasTaxPrivilege := r.FormValue("hasTaxPrivilege") != ""
	isNotResident := r.FormValue("isNotResident") != ""
	employmentType := parseEmploymentType(r.FormValue("employmentType"))

	// Месяц
	monthNum, err := strconv.Atoi(monthStr)
	if err != nil || monthNum < 1 || monthNum > 12 {
		log.Warn("form_invalid_month", "input", monthStr)
		monthNum = 1
	}
	startDate := time.Date(time.Now().Year(), time.Month(monthNum), 1, 0, 0, 0, 0, time.UTC)
	startTS := timestamppb.New(startDate)

	// Коэффициенты
	territorial := 100
	if territorialStr != "" {
		if v, err := strconv.Atoi(territorialStr); err == nil && v >= 100 && v <= 200 {
			territorial = v
		} else {
			log.Warn("form_invalid_territorial", "input", territorialStr)
		}
	}

	northern := 100
	if northernStr != "" {
		if v, err := strconv.Atoi(northernStr); err == nil && v >= 100 && v <= 200 {
			northern = v
		} else {
			log.Warn("form_invalid_northern", "input", northernStr)
		}
	}

	// Бонусы по месяцам: поля bonus_1 … bonus_12
	bonuses := make([]uint64, 12)
	for i := 1; i <= 12; i++ {
		raw := strings.TrimSpace(r.FormValue(fmt.Sprintf("bonus_%d", i)))
		if raw == "" {
			continue
		}
		raw = strings.ReplaceAll(raw, "\u00A0", "")
		raw = strings.ReplaceAll(raw, " ", "")
		raw = strings.ReplaceAll(raw, ",", ".")
		if v, err := strconv.ParseFloat(raw, 64); err == nil && v >= 0 {
			bonuses[i-1] = uint64(math.Round(v * 100))
		}
	}

	// Налоговые вычеты (опциональные поля, передаются только при заполнении)
	childrenCount := parseUint32Form(r, "childrenCount")
	disabledChildrenCount := parseUint32Form(r, "disabledChildrenCount")
	housingExpense := parseKopecksForm(r, "housingExpense")
	mortgageExpense := parseKopecksForm(r, "mortgageExpense")
	socialExpense := parseKopecksForm(r, "socialExpense")
	childEduExpense := parseKopecksForm(r, "childEduExpense")

	req := &pb.CalculatePrivateRequest{
		GrossSalary:           grossSalary,
		StartDate:             startTS,
		TerritorialMultiplier: uint64Ptr(uint64(territorial)),
		NorthernCoefficient:   uint64Ptr(uint64(northern)),
		HasTaxPrivilege:       boolPtr(hasTaxPrivilege),
		IsNotResident:         boolPtr(isNotResident),
		EmploymentType:        employmentTypePtr(employmentType),
		MonthlyBonuses:        bonuses,
	}
	if childrenCount > 0 {
		req.ChildrenCount = &childrenCount
	}
	if disabledChildrenCount > 0 {
		req.DisabledChildrenCount = &disabledChildrenCount
	}
	if housingExpense > 0 {
		req.HousingExpense = &housingExpense
	}
	if mortgageExpense > 0 {
		req.MortgageExpense = &mortgageExpense
	}
	if socialExpense > 0 {
		req.SocialExpense = &socialExpense
	}
	if childEduExpense > 0 {
		req.ChildEduExpense = &childEduExpense
	}
	return req, nil
}

// parseUint32Form читает числовое поле формы как uint32.
func parseUint32Form(r *http.Request, field string) uint32 {
	s := strings.TrimSpace(r.FormValue(field))
	if s == "" {
		return 0
	}
	v, err := strconv.ParseUint(s, 10, 32)
	if err != nil || v == 0 {
		return 0
	}
	return uint32(v)
}

// parseKopecksForm читает поле суммы в рублях и конвертирует в копейки (uint64).
func parseKopecksForm(r *http.Request, field string) uint64 {
	s := strings.TrimSpace(r.FormValue(field))
	if s == "" {
		return 0
	}
	s = strings.ReplaceAll(s, "\u00A0", "")
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, ",", ".")
	v, err := strconv.ParseFloat(s, 64)
	if err != nil || v <= 0 {
		return 0
	}
	return uint64(math.Round(v * 100))
}

// parseEmploymentType конвертирует строку формы в proto-enum EmploymentType.
func parseEmploymentType(s string) pb.EmploymentType {
	switch s {
	case "GPH":
		return pb.EmploymentType_GPH
	case "SELF_EMPLOYED":
		return pb.EmploymentType_SELF_EMPLOYED
	default:
		return pb.EmploymentType_TD
	}
}

// employmentTypePtr возвращает указатель на pb.EmploymentType для proto optional-поля.
func employmentTypePtr(v pb.EmploymentType) *pb.EmploymentType {
	return &v
}

func PrepareApiData() (*ApiDocsData, error) {

	raw, err := web.ApiDocsFS.ReadFile("api_docs/swagger.json")
	if err != nil {
		return nil, err
	}

	var d ApiDocsData
	if err := json.Unmarshal(raw, &d); err != nil {
		return nil, err
	}

	v := web.GetApiVersion()
	d.ApiVers = v

	for i := range d.Endpoints {
		d.Endpoints[i].Path = strings.ReplaceAll(
			d.Endpoints[i].Path,
			"{version}",
			v,
		)

		if obj, ok := d.Endpoints[i].ExampleRequest.(map[string]interface{}); ok {
			pretty, err := json.MarshalIndent(obj, "", "  ")
			if err == nil {
				d.Endpoints[i].ExampleRequest = string(pretty)
			}
		}

		if obj, ok := d.Endpoints[i].ExampleResponse.(map[string]interface{}); ok {
			pretty, err := json.MarshalIndent(obj, "", "  ")
			if err == nil {
				d.Endpoints[i].ExampleResponse = string(pretty)
			}
		}
	}

	return &d, nil
}

// Вспомогательные функции:
func uint64Ptr(v uint64) *uint64 { return &v }
func boolPtr(v bool) *bool       { return &v }
