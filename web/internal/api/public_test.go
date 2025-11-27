package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	pb "github.com/kiselevos/new_tax/gen/grpc/api"
	"github.com/kiselevos/new_tax/web/testutils"
)

// helper для запроса
func doPost(handler http.HandlerFunc, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/calc", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler(w, req)
	return w
}

func TestPublicCalc_OK(t *testing.T) {
	fake := &testutils.FakeTaxClient{
		PublicResp: &pb.CalculatePublicResponse{
			AnnualTaxAmount:   12000,
			AnnualGrossIncome: 600000,
			AnnualNetIncome:   588000,
			GrossSalary:       50000,
		},
	}

	h := NewPublicHandler(fake)

	w := doPost(h.HandlePublicCalc, `{"gross_salary": 50000}`)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp PublicCalcResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}

	if resp.AnnualTaxAmount != 12000 {
		t.Errorf("expected annual_tax_amount=12000, got %d", resp.AnnualTaxAmount)
	}
	if resp.GrossSalary != 50000 {
		t.Errorf("gross salary mismatch: %d", resp.GrossSalary)
	}
}

func TestPublicCalc_InvalidMethod(t *testing.T) {
	fake := &testutils.FakeTaxClient{}
	h := NewPublicHandler(fake)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/calc", nil)
	w := httptest.NewRecorder()

	h.HandlePublicCalc(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", w.Code)
	}
}

func TestPublicCalc_InvalidJSON(t *testing.T) {
	fake := &testutils.FakeTaxClient{}
	h := NewPublicHandler(fake)

	w := doPost(h.HandlePublicCalc, `{"gross_salary": "abc"}`)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestPublicCalc_BackendError(t *testing.T) {
	fake := &testutils.FakeTaxClient{
		PublicErr: assertError(),
	}
	h := NewPublicHandler(fake)

	w := doPost(h.HandlePublicCalc, `{"gross_salary": 50000}`)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// helper that returns a standard error
func assertError() error {
	return http.ErrHandlerTimeout
}
