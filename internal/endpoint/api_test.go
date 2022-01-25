package endpoint_test

import (
	"context"

	"github.com/soheilhy/cmux"

	"gitlab.com/inetmock/inetmock/internal/endpoint"
)

type MultiplexHandlerMock struct {
	endpoint.ProtocolHandler
	MultiplexMatchers []cmux.Matcher
}

func (h MultiplexHandlerMock) Matchers() []cmux.Matcher {
	return h.MultiplexMatchers
}

type ProtocolHandlerFunc func(ctx context.Context, startupSpec *endpoint.StartupSpec) error

func (p ProtocolHandlerFunc) Start(ctx context.Context, startupSpec *endpoint.StartupSpec) error {
	return p(ctx, startupSpec)
}

type StoppableProtocolHandlerMock struct {
	OnStart func(ctx context.Context, startupSpec *endpoint.StartupSpec) error
	OnStop  func(ctx context.Context) error
}

func (p StoppableProtocolHandlerMock) Start(ctx context.Context, startupSpec *endpoint.StartupSpec) error {
	return p.OnStart(ctx, startupSpec)
}

func (p StoppableProtocolHandlerMock) Stop(ctx context.Context) error {
	return p.OnStop(ctx)
}
