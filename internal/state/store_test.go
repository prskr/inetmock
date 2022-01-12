package state_test

import (
	"errors"
	"path"
	"reflect"
	"testing"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/maxatome/go-testdeep/td"
	"go.uber.org/multierr"

	"gitlab.com/inetmock/inetmock/internal/state"
	"gitlab.com/inetmock/inetmock/internal/state/statetest"
)

type sampleStruct struct {
	FirstName string
	LastName  string
	Age       uint
}

var (
	ted = sampleStruct{
		FirstName: "Ted",
		LastName:  "Tester",
		Age:       42,
	}
	simon = sampleStruct{
		FirstName: "Simon",
		LastName:  "Sample",
		Age:       21,
	}

	suffixes = []string{
		"",
		"subsystem1",
	}
)

func TestStore_Get(t *testing.T) {
	t.Parallel()

	type args struct {
		key string
		v   interface{}
	}
	tests := []struct {
		name    string
		args    args
		setup   statetest.StoreSetup
		want    interface{}
		wantErr error
	}{
		{
			name: "Get a sample value",
			args: args{
				key: "ted.tester",
				v:   new(sampleStruct),
			},
			setup: statetest.StoreSetupFunc(func(tb testing.TB, store state.KVStore) error {
				tb.Helper()
				return store.Set("ted.tester", &ted)
			}),
			want:    td.Struct(&ted, td.StructFields{}),
			wantErr: nil,
		},
		{
			name: "Error receiver not a pointer",
			args: args{
				key: "ted.tester",
				v:   sampleStruct{},
			},
			wantErr: state.ErrReceiverNotAPointer,
		},
		{
			name: "Error - missing value",
			args: args{
				key: "ted.tester",
				v:   new(sampleStruct),
			},
			wantErr: badger.ErrKeyNotFound,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			for sfxIdx := range suffixes {
				store := statetest.NewTestStore(t).WithSuffixes(suffixes[sfxIdx])
				statetest.SetupStore(t, store, tt.setup)
				if err := store.Get(tt.args.key, tt.args.v); err != nil {
					if !errors.Is(err, tt.wantErr) {
						t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
					}
					continue
				}
				td.Cmp(t, tt.args.v, tt.want)
			}
		})
	}
}

func TestStore_GetAll(t *testing.T) {
	t.Parallel()

	type args struct {
		prefix string
		into   func() interface{}
	}
	tests := []struct {
		name    string
		args    args
		setup   statetest.StoreSetup
		want    interface{}
		wantErr error
	}{
		{
			name: "Empty database - empty result",
			args: args{
				prefix: "",
				into: func() interface{} {
					return new([]sampleStruct)
				},
			},
			want: td.Empty(),
		},
		{
			name: "Empty database - non-existing prefix - empty result",
			args: args{
				prefix: "teachers",
				into: func() interface{} {
					return new([]sampleStruct)
				},
			},
			want: td.Empty(),
		},
		{
			name: "Single result",
			args: args{
				into: func() interface{} {
					return new([]sampleStruct)
				},
			},
			setup: statetest.StoreSetupFunc(func(tb testing.TB, store state.KVStore) error {
				tb.Helper()
				return store.Set("ted.tester", ted)
			}),
			want: td.Bag(ted),
		},
		{
			name: "Single result - as reference",
			args: args{
				into: func() interface{} {
					return new([]*sampleStruct)
				},
			},
			setup: statetest.StoreSetupFunc(func(tb testing.TB, store state.KVStore) error {
				tb.Helper()
				return store.Set("ted.tester", ted)
			}),
			want: td.Bag(&ted),
		},
		{
			name: "Multiple results",
			args: args{
				into: func() interface{} {
					return new([]sampleStruct)
				},
			},
			setup: statetest.StoreSetupFunc(func(tb testing.TB, store state.KVStore) error {
				tb.Helper()
				return multierr.Append(
					store.Set("ted.tester", ted),
					store.Set("simon.sample", simon),
				)
			}),
			want: td.Bag(ted, simon),
		},
		{
			name: "Single result with filter",
			args: args{
				prefix: "teachers",
				into: func() interface{} {
					return new([]sampleStruct)
				},
			},
			setup: statetest.StoreSetupFunc(func(tb testing.TB, store state.KVStore) error {
				tb.Helper()
				return multierr.Append(
					store.Set(path.Join("teachers", "ted.tester"), ted),
					store.Set(path.Join("pupils", "simon.sample"), simon),
				)
			}),
			want: td.Bag(ted),
		},
		{
			name: "Multiple results with filter",
			args: args{
				prefix: "teachers",
				into: func() interface{} {
					return new([]sampleStruct)
				},
			},
			setup: statetest.StoreSetupFunc(func(tb testing.TB, store state.KVStore) error {
				tb.Helper()
				return multierr.Combine(
					store.Set(path.Join("teachers", "ted.tester"), ted),
					store.Set(path.Join("teachers", "simon.sample"), simon),
					store.Set(path.Join("pupils", "simon.sample"), simon),
				)
			}),
			want: td.Bag(ted, simon),
		},
		{
			name: "Err not a pointer",
			args: args{
				into: func() interface{} {
					return make([]sampleStruct, 0)
				},
			},
			wantErr: state.ErrReceiverNotAPointer,
		},
		{
			name: "Err not a slice",
			args: args{
				into: func() interface{} {
					return new(sampleStruct)
				},
			},
			wantErr: state.ErrReceiverNotASlice,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			for sfxIdx := range suffixes {
				store := statetest.NewTestStore(t).WithSuffixes(suffixes[sfxIdx])
				statetest.SetupStore(t, store, tt.setup)
				into := tt.args.into()
				if err := store.GetAll(tt.args.prefix, into); err != nil {
					if !errors.Is(err, tt.wantErr) {
						t.Errorf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
					}
					continue
				}
				td.Cmp(t, into, tt.want)
			}
		})
	}
}

func TestStore_Set(t *testing.T) {
	t.Parallel()
	type args struct {
		key  string
		v    interface{}
		opts []state.SetOption
	}
	tests := []struct {
		name   string
		args   args
		setup  statetest.StoreSetup
		sleep  time.Duration
		setErr error
		getErr error
	}{
		{
			name: "Set sample struct",
			args: args{
				key: "ted.tester",
				v:   &ted,
			},
		},
		{
			name: "Override existing value",
			args: args{
				key: "ted.tester",
				v:   &ted,
			},
			setup: statetest.StoreSetupFunc(func(tb testing.TB, store state.KVStore) error {
				tb.Helper()
				return store.Set("ted.tester", nil)
			}),
		},
		{
			name: "Set with TTL",
			args: args{
				key:  "ted.tester",
				v:    &ted,
				opts: []state.SetOption{state.WithTTL(100 * time.Millisecond)},
			},
			sleep: 200 * time.Millisecond,
			setup: statetest.StoreSetupFunc(func(tb testing.TB, store state.KVStore) error {
				tb.Helper()
				return store.Set("ted.tester", nil)
			}),
			getErr: badger.ErrKeyNotFound,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			for sfxIdx := range suffixes {
				store := statetest.NewTestStore(t).WithSuffixes(suffixes[sfxIdx])
				statetest.SetupStore(t, store, tt.setup)
				if err := store.Set(tt.args.key, tt.args.v, tt.args.opts...); err != nil {
					if !errors.Is(err, tt.setErr) {
						t.Errorf("Set() error = %v, wantErr %v", err, tt.setErr)
					}
					continue
				}

				if tt.sleep > 0 {
					time.Sleep(tt.sleep)
				}

				got := reflect.New(reflect.TypeOf(tt.args.v).Elem()).Interface()
				if err := store.Get(tt.args.key, got); err != nil {
					if !errors.Is(err, tt.getErr) {
						t.Errorf("store.Get() error = %v", err)
					}
					continue
				}

				td.Cmp(t, got, tt.args.v)
			}
		})
	}
}
