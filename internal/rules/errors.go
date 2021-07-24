package rules

import "errors"

var (
	ErrNoTerminatorDefined = errors.New("no terminator defined")
	ErrUnknownTerminator   = errors.New("no terminator with the given name is known")
	ErrUnknownFilterMethod = errors.New("no filter with the given name is known")
)
