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

	return &pb.CalculatePrivateRequest{
		GrossSalary:           grossSalary,
		StartDate:             startTS,
		TerritorialMultiplier: uint64Ptr(uint64(territorial)),
		NorthernCoefficient:   uint64Ptr(uint64(northern)),
		HasTaxPrivilege:       boolPtr(hasTaxPrivilege),
		IsNotResident:         boolPtr(isNotResident),
	}, nil
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
