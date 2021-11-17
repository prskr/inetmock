package test

import (
	"context"
	"testing"
)

func Context(t *testing.T) context.Context {
	t.Helper()
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	if deadline, ok := t.Deadline(); ok {
		ctx, cancel = context.WithDeadline(context.Background(), deadline)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	t.Cleanup(cancel)
	return ctx
}
