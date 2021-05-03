package health

import (
	"context"
	"errors"
	"fmt"
)

type Result map[string]error

func (r Result) IsHealthy() (healthy bool) {
	for _, e := range r {
		if e != nil {
			return false
		}
	}
	return true
}

func (r Result) CheckResult(name string) (knownCheck bool, result error) {
	result, knownCheck = r[name]
	return
}

type Checker interface {
	AddCheck(check Check) error
	Status(ctx context.Context) (Result, error)
}

type CheckError struct {
	Check   string
	Message string
	Orig    error
}

func (c CheckError) Error() string {
	return fmt.Sprintf("check %s failed: %s - %v", c.Check, c.Message, c.Orig)
}

func (c CheckError) Is(err error) bool {
	return errors.Is(c.Orig, err)
}

func (c CheckError) Unwrap() error {
	return c.Orig
}

type Check interface {
	Name() string
	Status(ctx context.Context) CheckError
}

var (
	ErrAmbiguousCheckName = errors.New("a check with the same name is already registered")
)

func New() Checker {
	return &checker{}
}
