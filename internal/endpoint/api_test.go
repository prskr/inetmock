package endpoint_test

import (
	"context"
	"testing"

	"gitlab.com/inetmock/inetmock/internal/endpoint"
)

type ProtocolHandlerDelegate func(ctx context.Context, lifecycle endpoint.Lifecycle) error

func (p ProtocolHandlerDelegate) Start(ctx context.Context, lifecycle endpoint.Lifecycle) error {
	return p(ctx, lifecycle)
}

func VerifiedProtocolHandler(
	tb testing.TB,
	delegate func(ctx context.Context, lifecycle endpoint.Lifecycle) error,
) endpoint.ProtocolHandler {
	tb.Helper()
	var called bool
	tb.Cleanup(func() {
		if !called {
			tb.Error("ProtocolHandler got not called")
		}
	})
	return ProtocolHandlerDelegate(func(ctx context.Context, lifecycle endpoint.Lifecycle) error {
		called = true
		return delegate(ctx, lifecycle)
	})
}
