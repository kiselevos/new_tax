package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	pb "github.com/kiselevos/new_tax/gen/grpc/api"
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
	CurrentYear int
	Months      []Month
	Territorial []CoefficientOption
	Northern    []BonusOption
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

	return IndexData{
		CurrentYear: time.Now().Year(),
		Months:      PrepareMonths(),
		Territorial: territorial,
		Northern:    northern,
	}
}

func ParseFormToRequest(r *http.Request) (*pb.CalculatePrivateRequest, error) {
	r.ParseForm()

	grossSalary := parseUint(r.FormValue("grossSalary"))
	monthStr := r.FormValue("startDate")
	territorialStr := r.FormValue("territorialMultiplier")
	northernStr := r.FormValue("northernCoefficient")
	hasTaxPrivilege := r.FormValue("hasTaxPrivilege") != ""
	isNotResident := r.FormValue("isNotResident") != ""

	monthNum, _ := strconv.Atoi(monthStr)
	if monthNum == 0 {
		monthNum = 1
	}
	startDate := time.Date(time.Now().Year(), time.Month(monthNum), 1, 0, 0, 0, 0, time.UTC)
	startTS := timestamppb.New(startDate)

	territorial, _ := strconv.Atoi(territorialStr)
	northern, _ := strconv.Atoi(northernStr)

	log.Printf("📄 Form parsed: GrossSalary=%d, Territorial=%d, Northern=%d, HasTaxPrivilege=%t, IsNotResident=%t, StartDate=%s",
		grossSalary, territorial, northern, hasTaxPrivilege, isNotResident,
		startDate.Format("2006-01-02"))

	return &pb.CalculatePrivateRequest{
		GrossSalary:           grossSalary,
		StartDate:             startTS,
		TerritorialMultiplier: uint64Ptr(uint64(territorial)),
		NorthernCoefficient:   uint64Ptr(uint64(northern)),
		HasTaxPrivilege:       boolPtr(hasTaxPrivilege),
		IsNotResident:         boolPtr(isNotResident),
	}, nil
}

// Вспомогательные функции:
func uint64Ptr(v uint64) *uint64 { return &v }
func boolPtr(v bool) *bool       { return &v }
