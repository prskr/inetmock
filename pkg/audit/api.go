//go:generate mockgen -source=$GOFILE -destination=./../../internal/mock/audit/audit.mock.go -package=audit_mock

package audit

import (
	"context"
	"errors"
	"io"
)

var (
	ErrSinkAlreadyRegistered = errors.New("sink with same name already registered")
)

type Emitter interface {
	Emit(ev Event)
}

type Sink interface {
	Name() string
	OnSubscribe(evs <-chan Event)
}

type EventStream interface {
	io.Closer
	Emitter
	RegisterSink(ctx context.Context, s Sink) error
	Sinks() []string
	RemoveSink(name string) (exists bool)
}
