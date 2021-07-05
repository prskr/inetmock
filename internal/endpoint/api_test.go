package endpoint_test

import (
	"context"
	"sync/atomic"
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
	var called int32
	tb.Cleanup(func() {
		if atomic.LoadInt32(&called) < 1 {
			tb.Error("ProtocolHandler got not called")
		}
	})
	return ProtocolHandlerDelegate(func(ctx context.Context, lifecycle endpoint.Lifecycle) error {
		atomic.AddInt32(&called, 1)
		return delegate(ctx, lifecycle)
	})
}
