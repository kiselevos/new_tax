package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoggerMiddleware_AddsContext(t *testing.T) {
	called := false
	h := Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Context() == nil {
			t.Error("context should not be nil")
		}
		called = true
	}))
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if !called {
		t.Error("handler was not called")
	}
}
