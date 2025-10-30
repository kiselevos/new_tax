package client

import (
	"testing"
)

func TestNewTaxClient_DefaultAddr(t *testing.T) {
	client, conn, err := NewTaxClient()
	if err == nil {
		t.Log("as expected: connection refused (backend not running)")
	}
	if client == nil || conn == nil {
		t.Errorf("expected non-nil client and conn even if backend unavailable")
	}
}
