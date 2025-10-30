package server

import (
	"net/http"
	"testing"
)

func TestNewServer_DoesNotPanic(t *testing.T) {
	mux := http.NewServeMux()
	srv := New(":0", mux)

	if srv == nil {
		t.Fatal("expected non-nil *http.Server")
	}
	if srv.Addr != ":0" {
		t.Errorf("expected addr ':0', got %q", srv.Addr)
	}
	if srv.Handler != mux {
		t.Error("handler not set correctly")
	}
}
