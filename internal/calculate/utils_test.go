package calculate

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateCalculateInput_OK(t *testing.T) {
	in := CalculateInput{
		GrossSalary:           1_000_00, // 1 000 ₽ в копейках
		TerritorialMultiplier: 120,      // +20%
		NorthernCoefficient:   150,      // +50%
	}
	require.NoError(t, ValidateCalculateInput(in))
}

func TestValidateCalculateInput_Errors(t *testing.T) {
	tests := []struct {
		name string
		in   CalculateInput
		want string
	}{
		{
			name: "salary too small",
			in:   CalculateInput{GrossSalary: 0, TerritorialMultiplier: 100, NorthernCoefficient: 100},
			want: "gross_salary must be between 1",
		},
		{
			name: "salary too big",
			in:   CalculateInput{GrossSalary: 10_000_000_001, TerritorialMultiplier: 100, NorthernCoefficient: 100},
			want: "gross_salary must be between 1",
		},
		{
			name: "territorial below min",
			in:   CalculateInput{GrossSalary: 100, TerritorialMultiplier: 90, NorthernCoefficient: 100},
			want: "territorial_multiplier must be between 100 and 200",
		},
		{
			name: "territorial above max",
			in:   CalculateInput{GrossSalary: 100, TerritorialMultiplier: 210, NorthernCoefficient: 100},
			want: "territorial_multiplier must be between 100 and 200",
		},
		{
			name: "territorial wrong step",
			in:   CalculateInput{GrossSalary: 100, TerritorialMultiplier: 117, NorthernCoefficient: 100},
			want: "territorial_multiplier must be between 100 and 200",
		},
		{
			name: "northern below min",
			in:   CalculateInput{GrossSalary: 100, TerritorialMultiplier: 100, NorthernCoefficient: 90},
			want: "northern_coefficient must be between 100 and 200",
		},
		{
			name: "northern above max",
			in:   CalculateInput{GrossSalary: 100, TerritorialMultiplier: 100, NorthernCoefficient: 210},
			want: "northern_coefficient must be between 100 and 200",
		},
		{
			name: "northern wrong step",
			in:   CalculateInput{GrossSalary: 100, TerritorialMultiplier: 100, NorthernCoefficient: 133},
			want: "northern_coefficient must be between 100 and 200",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCalculateInput(tt.in)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.want)
		})
	}
}

func TestValidateCoefficient(t *testing.T) {
	type T = struct {
		in   uint64
		want bool
	}
	ok := []T{
		{100, true}, {110, true}, {115, true}, {120, true}, {200, true},
	}
	bad := []T{
		{99, false}, {201, false}, {116, false}, {0, false},
	}
	for _, tc := range ok {
		assert.Truef(t, validateCoefficient(tc.in), "should accept %d", tc.in)
	}
	for _, tc := range bad {
		assert.Falsef(t, validateCoefficient(tc.in), "should reject %d", tc.in)
	}
}

func TestGetStartMonth(t *testing.T) {
	d := time.Date(2025, time.June, 1, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, int(time.June), GetStartMonth(d))
}

func TestIntMonthFromDate(t *testing.T) {
	ref := time.Date(2025, time.September, 10, 12, 0, 0, 0, time.UTC)
	got := IntMonthFromDate(int(time.March), ref)
	assert.Equal(t, 2025, got.Year())
	assert.Equal(t, time.March, got.Month())
	assert.Equal(t, 1, got.Day())
	assert.Equal(t, ref.Location(), got.Location())
}
