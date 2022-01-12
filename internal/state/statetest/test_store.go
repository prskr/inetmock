package statetest

import (
	"testing"

	"github.com/maxatome/go-testdeep/td"

	"gitlab.com/inetmock/inetmock/internal/state"
)

type (
	StoreSetup interface {
		Setup(tb testing.TB, store state.KVStore) error
	}
	StoreSetupFunc func(tb testing.TB, store state.KVStore) error
)

func (s StoreSetupFunc) Setup(tb testing.TB, store state.KVStore) error {
	tb.Helper()
	return s(tb, store)
}

func NewTestStore(tb testing.TB, setups ...StoreSetup) *state.Store {
	tb.Helper()
	s, err := state.NewDefault(state.WithInMemory(), state.WithLogger(state.TestLogger{TB: tb}))
	if err != nil {
		tb.Fatalf("NewDefault() error = %v", err)
		return nil
	}

	tb.Cleanup(func() {
		td.CmpNoError(tb, s.Close())
	})

	SetupStore(tb, s, setups...)

	return s
}

func SetupStore(tb testing.TB, store state.KVStore, setups ...StoreSetup) {
	tb.Helper()
	for idx := range setups {
		setup := setups[idx]
		if setup == nil {
			continue
		}
		if err := setup.Setup(tb, store); err != nil {
			tb.Fatalf("setup failed: %s", err)
		}
	}
}
