package health

import (
	"context"
	"errors"
)

var (
	ErrAmbiguousCheckName = errors.New("a check with the same name is already registered")
)

func New() Checker {
	return &checker{}
}

type Checker interface {
	AddCheck(check Check) error
	Status(ctx context.Context) Result
}

type Check interface {
	Name() string
	Status(ctx context.Context) error
}
