//go:generate mockgen -source=$GOFILE -destination=./../../internal/mock/audit/audit.mock.go -package=audit_mock

package audit

import (
	"errors"
	"io"
)

var (
	ErrSinkAlreadyRegistered = errors.New("sink with same name already registered")
)

type Emitter interface {
	Emit(ev Event)
}

type CloseHandle func()

type Sink interface {
	Name() string
	OnSubscribe(evs <-chan Event, close CloseHandle)
}

type EventStream interface {
	io.Closer
	Emitter
	RegisterSink(s Sink) error
	Sinks() []string
	RemoveSink(name string)
}
