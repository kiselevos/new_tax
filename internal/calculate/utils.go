package calculate

import (
	"errors"
	"fmt"
	"time"

	"github.com/oapi-codegen/runtime/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ValidateCalculateInput проверяет валидность входных данных перед расчетом налога.
func ValidateCalculateInput(input CalculateInput) error {
	if input.GrossSalary < 1 || input.GrossSalary > 10_000_000_000 {
		return errors.New("gross_salary must be between 1 and 100,000,000,00 cents")
	}

	if !validateCoefficient(input.TerritorialMultiplier) {
		return fmt.Errorf("territorial_multiplier must be between 100 and 200 with step 10")
	}

	if !validateCoefficient(input.NorthernCoefficient) {
		return fmt.Errorf("northern_coefficient must be between 100 and 200 with step 10")
	}

	return nil
}

// validateCoefficient проверяет, что коэффициент от 100 до 200 с шагом 10.
func validateCoefficient(coeff uint64) bool {
	return coeff >= 100 && coeff <= 200 && (coeff-100)%10 == 0
}

// ptr возвращает указатель на переданное значение.
func ptr[T any](v T) *T {
	return &v
}

// GetStartMonth возвращает номер месяца из переданной даты.
func GetStartMonth(date time.Time) int {
	return int(date.Month())
}

// IntMonthFromDate возвращает дату первого числа месяца по номеру месяца и году из date.
func IntMonthFromDate(m int, date time.Time) time.Time {
	return time.Date(date.Year(), time.Month(m), 1, 0, 0, 0, 0, date.Location())
}

// ToOAPIDate преобразует time.Time в формат oapi-codegen Date.
func ToOAPIDate(t time.Time) *types.Date {
	return &types.Date{Time: t}
}

// ToProtoTimestamp преобразует time.Time в protobuf Timestamp.
func ToProtoTimestamp(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}
