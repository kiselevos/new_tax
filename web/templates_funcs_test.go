package web

import "testing"

func TestGetMinLivingWage(t *testing.T) {
	got := GetMinLivingWage()
	if got <= 0 {
		t.Errorf("expected positive wage, got %v", got)
	}
}

func TestFuncs_ContainsGetMinSalary(t *testing.T) {
	if _, ok := Funcs["getMinSalary"]; !ok {
		t.Error("expected getMinSalary to be registered in Funcs")
	}
}
