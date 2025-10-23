package handlers

import (
	"net/http"
	"strconv"
	"time"

	pb "github.com/kiselevos/new_tax/gen/grpc/api"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ParseFormToRequest превращает значения формы в gRPC-запрос
func ParseFormToRequest(r *http.Request) (*pb.CalculatePrivateRequest, error) {
	r.ParseForm()

	// --- базовые значения ---
	grossSalary := parseUint(r.FormValue("grossSalary"))
	monthStr := r.FormValue("startDate")
	territorialStr := r.FormValue("territorialMultiplier")
	northernStr := r.FormValue("northernCoefficient")
	hasTaxPrivilege := r.FormValue("hasTaxPrivilege") != ""
	isNotResident := r.FormValue("isNotResident") != ""

	// --- даты ---
	monthNum, _ := strconv.Atoi(monthStr)
	if monthNum == 0 {
		monthNum = 1 // дефолт — январь
	}
	startDate := time.Date(time.Now().Year(), time.Month(monthNum), 1, 0, 0, 0, 0, time.UTC)
	startTS := timestamppb.New(startDate)

	// --- коэффициенты ---
	territorial, _ := strconv.ParseFloat(territorialStr, 64)
	northern, _ := strconv.ParseFloat(northernStr, 64)

	if territorial == 0 {
		territorial = 1.0
	}
	if northern == 0 {
		northern = 1.0
	}

	territorialUint := uint64(territorial * 100)
	northernUint := uint64(100 + northern)

	return &pb.CalculatePrivateRequest{
		GrossSalary:           grossSalary,
		StartDate:             startTS,
		TerritorialMultiplier: &territorialUint,
		NorthernCoefficient:   &northernUint,
		HasTaxPrivilege:       &hasTaxPrivilege,
		IsNotResident:         &isNotResident,
	}, nil
}
