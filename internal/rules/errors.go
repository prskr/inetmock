package rules

import "errors"

var (
	ErrNoTerminatorDefined = errors.New("no terminator defined")
	ErrUnknownTerminator   = errors.New("no terminator with the given name is known")
	ErrUnknownFilterMethod = errors.New("no filter with the given name is known")
	ErrNoInitiatorDefined  = errors.New("no initiator defined")
	ErrUnknownInitiator    = errors.New("no initiator with the given name is known")
)
