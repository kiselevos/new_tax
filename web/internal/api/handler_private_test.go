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
}
