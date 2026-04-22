package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	pb "github.com/kiselevos/new_tax/gen/grpc/api"
	"github.com/kiselevos/new_tax/web/testutils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Универсальный helper для POST запросов
func doPostPrivate(handler http.HandlerFunc, body string, headers map[string]string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/private-calc", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	w := httptest.NewRecorder()
	handler(w, req)
	return w
}

func TestPrivateCalc_OK(t *testing.T) {
	fake := &testutils.FakeTaxClient{
		PrivateResp: &pb.CalculatePrivateResponse{
			AnnualTaxAmount: 55555,
			GrossSalary:     120000,
		},
	}

	h := NewPrivateHandler(fake)

	w := doPostPrivate(h.HandlePrivateCalc,
		`{"gross_salary": 120000}`,
		map[string]string{"x-api-key": "valid"},
	)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp PrivateCalcResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("invalid json response: %v", err)
	}

	if resp.AnnualTaxAmount != 55555 {
		t.Errorf("expected 55555 got %d", resp.AnnualTaxAmount)
	}
}

func TestPrivateCalc_NoKey(t *testing.T) {
	fake := &testutils.FakeTaxClient{
		PrivateErr: status.Error(codes.PermissionDenied, "missing key"),
	}

	h := NewPrivateHandler(fake)

	w := doPostPrivate(h.HandlePrivateCalc, `{"gross_salary":120000}`, nil)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

func TestPrivateCalc_WrongKey(t *testing.T) {
	fake := &testutils.FakeTaxClient{
		PrivateErr: status.Error(codes.PermissionDenied, "invalid key"),
	}

	h := NewPrivateHandler(fake)

	w := doPostPrivate(h.HandlePrivateCalc,
		`{"gross_salary": 120000}`,
		map[string]string{"x-api-key": "WRONG"},
	)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

func TestPrivateCalc_InvalidJSON(t *testing.T) {
	fake := &testutils.FakeTaxClient{}

	h := NewPrivateHandler(fake)

	w := doPostPrivate(h.HandlePrivateCalc, `{"gross_salary": "oops"}`, map[string]string{
		"x-api-key": "valid",
	})

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestPrivateCalc_InvalidMethod(t *testing.T) {
	fake := &testutils.FakeTaxClient{}
	h := NewPrivateHandler(fake)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/private-calc", nil)
	w := httptest.NewRecorder()

	h.HandlePrivateCalc(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", w.Code)
	}
}

func TestPrivateCalc_BackendError(t *testing.T) {
	fake := &testutils.FakeTaxClient{
		PrivateErr: assertError(),
	}

	h := NewPrivateHandler(fake)

	w := doPostPrivate(h.HandlePrivateCalc,
		`{"gross_salary": 120000}`,
		map[string]string{"x-api-key": "valid"},
	)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if body["error"] != "internal server error" {
		t.Errorf("5xx must not leak details, got %q", body["error"])
	}
}

func TestPrivateCalc_BackendValidationError(t *testing.T) {
	fake := &testutils.FakeTaxClient{
		PrivateErr: status.Error(codes.InvalidArgument, "salary must be > 0"),
	}

	h := NewPrivateHandler(fake)

	w := doPostPrivate(h.HandlePrivateCalc,
		`{"gross_salary": 120000}`,
		map[string]string{"x-api-key": "valid"},
	)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	// Должна быть только description, без "rpc error: code = X desc = ..."
	if body["error"] != "salary must be > 0" {
		t.Errorf("expected clean message, got %q", body["error"])
	}
}

func TestPrivateCalc_InvalidEmploymentType(t *testing.T) {
	fake := &testutils.FakeTaxClient{}
	h := NewPrivateHandler(fake)

	w := doPostPrivate(h.HandlePrivateCalc,
		`{"gross_salary": 120000, "employment_type": "INVALID"}`,
		map[string]string{"x-api-key": "valid"},
	)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestPrivateCalc_WrongBonusesLength(t *testing.T) {
	fake := &testutils.FakeTaxClient{}
	h := NewPrivateHandler(fake)

	w := doPostPrivate(h.HandlePrivateCalc,
		`{"gross_salary": 120000, "monthly_bonuses": [0, 0, 0]}`,
		map[string]string{"x-api-key": "valid"},
	)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestPrivateCalc_GPH(t *testing.T) {
	fake := &testutils.FakeTaxClient{
		PrivateResp: &pb.CalculatePrivateResponse{
			AnnualTaxAmount: 23400000,
			GrossSalary:     15000000,
			AnnualFSS:       0,
		},
	}
	h := NewPrivateHandler(fake)

	w := doPostPrivate(h.HandlePrivateCalc,
		`{"gross_salary": 15000000, "employment_type": "GPH"}`,
		map[string]string{"x-api-key": "valid"},
	)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp PrivateCalcResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if resp.AnnualFSS != 0 {
		t.Errorf("GPH should have zero FSS, got %d", resp.AnnualFSS)
	}
}

func TestPrivateCalc_WithDeductionResult(t *testing.T) {
	fake := &testutils.FakeTaxClient{
		PrivateResp: &pb.CalculatePrivateResponse{
			AnnualTaxAmount: 10000000,
			GrossSalary:     20000000,
			DeductionResult: &pb.DeductionResult{
				ChildrenMonthlyDeduction: 280000,
				ChildrenMonths:           12,
				ChildrenReturn:           436800,
				TotalReturn:              436800,
			},
		},
	}
	h := NewPrivateHandler(fake)

	w := doPostPrivate(h.HandlePrivateCalc,
		`{"gross_salary": 20000000, "children_count": 2}`,
		map[string]string{"x-api-key": "valid"},
	)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp PrivateCalcResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if resp.DeductionResult == nil {
		t.Fatal("expected deduction_result in response")
	}
	if resp.DeductionResult.TotalReturn != 436800 {
		t.Errorf("expected total_return 436800, got %d", resp.DeductionResult.TotalReturn)
	}
}
