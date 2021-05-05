package rules

import (
	"errors"
	"fmt"
)

var (
	ErrAmbiguousParamCount = errors.New("the supplied number of arguments does not match the expected one")
)

func StringP(value string) *string {
	return &value
}

func IntP(value int) *int {
	return &value
}

func FloatP(value float64) *float64 {
	return &value
}

func ValidateParameterCount(params []Param, expected int) error {
	if len(params) < expected {
		return fmt.Errorf("%w: expected %d got %d", ErrAmbiguousParamCount, expected, len(params))
	}
	return nil
}
