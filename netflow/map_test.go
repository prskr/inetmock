//go:build sudo

package netflow_test

import (
	"net/netip"
	"testing"

	"github.com/cilium/ebpf"
	"github.com/maxatome/go-testdeep/td"

	"inetmock.icb4dc0.de/inetmock/netflow"
)

type mapPutAllTestCase[K comparable, V any] struct {
	name     string
	toInsert map[K]V
	opts     []netflow.BatchOption
	mapSpec  ebpf.MapSpec
}

func (tt mapPutAllTestCase[K, V]) Name() string {
	return tt.name
}

//nolint:thelper // is not really a helper function
func (tt mapPutAllTestCase[K, V]) Run(t *testing.T) {
	t.Parallel()

	underlying, err := ebpf.NewMap(&tt.mapSpec)
	if err != nil {
		t.Fatalf("Failed to setup cleanup map: %v", err)
	}

	t.Cleanup(func() {
		if err := underlying.Close(); err != nil {
			t.Errorf("Failed to close test map: %v", err)
		}
	})

	testMap := netflow.MapOf[K, V](underlying)

	if err := testMap.PutAll(tt.toInsert, tt.opts...); err != nil {
		t.Errorf("PutAll() error = %v", err)
		return
	}

	if got, err := testMap.GetAll(tt.opts...); err != nil {
		t.Errorf("GetAll() error = %v", err)
	} else {
		td.Cmp(t, got, tt.toInsert)
	}
}

func TestMap_PutAll(t *testing.T) {
	t.Parallel()
	RemoveMemlock(t)

	tests := []testCase{
		mapPutAllTestCase[uint32, uint32]{
			name:     "Put empty map",
			toInsert: make(map[uint32]uint32),
			mapSpec: ebpf.MapSpec{
				Name:       "map_test_empty",
				Type:       ebpf.Hash,
				KeySize:    4,
				ValueSize:  4,
				MaxEntries: 1,
			},
		},
		mapPutAllTestCase[uint32, uint32]{
			name: "Put single element",
			toInsert: map[uint32]uint32{
				13: 37,
			},
			mapSpec: ebpf.MapSpec{
				Name:       "map_test_single_elem",
				Type:       ebpf.Hash,
				KeySize:    4,
				ValueSize:  4,
				MaxEntries: 1,
			},
		},
		mapPutAllTestCase[uint32, uint32]{
			name: "Put single element - force fallback",
			toInsert: map[uint32]uint32{
				13: 37,
			},
			mapSpec: ebpf.MapSpec{
				Name:       "map_test_single_elem_fb",
				Type:       ebpf.Hash,
				KeySize:    4,
				ValueSize:  4,
				MaxEntries: 1,
			},
			opts: []netflow.BatchOption{netflow.WithUseFallback(true)},
		},
		mapPutAllTestCase[uint32, netflow.ConnIdent]{
			name: "Put single complex element",
			toInsert: map[uint32]netflow.ConnIdent{
				13: {
					Addr:      netip.MustParseAddr("1.2.3.4"),
					Port:      2345,
					Transport: netflow.ProtocolTCP,
				},
			},
			mapSpec: ebpf.MapSpec{
				Name:       "map_test_single_complex_elem",
				Type:       ebpf.Hash,
				KeySize:    4,
				ValueSize:  12,
				MaxEntries: 1,
			},
		},
		mapPutAllTestCase[uint32, netflow.ConnIdent]{
			name: "Put single complex element - force fallback",
			toInsert: map[uint32]netflow.ConnIdent{
				13: {
					Addr:      netip.MustParseAddr("1.2.3.4"),
					Port:      2345,
					Transport: netflow.ProtocolTCP,
				},
			},
			mapSpec: ebpf.MapSpec{
				Name:       "map_test_single_complex_elem_fb",
				Type:       ebpf.Hash,
				KeySize:    4,
				ValueSize:  12,
				MaxEntries: 1,
			},
			opts: []netflow.BatchOption{netflow.WithUseFallback(true)},
		},
	}

	//nolint:paralleltest // done in Run function
	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name(), tt.Run)
	}
}

type mapSyncTestCase[K comparable, V any] struct {
	name         string
	initialState map[K]V
	newState     map[K]V
	opts         []netflow.BatchOption
	mapSpec      ebpf.MapSpec
}

func (tt mapSyncTestCase[K, V]) Name() string {
	return tt.name
}

//nolint:thelper // is not really a helper function
func (tt mapSyncTestCase[K, V]) Run(t *testing.T) {
	t.Parallel()

	underlying, err := ebpf.NewMap(&tt.mapSpec)
	if err != nil {
		t.Fatalf("Failed to setup cleanup map: %v", err)
	}

	t.Cleanup(func() {
		if err := underlying.Close(); err != nil {
			t.Errorf("Failed to close test map: %v", err)
		}
	})

	testMap := netflow.MapOf[K, V](underlying)

	if err := testMap.PutAll(tt.initialState, tt.opts...); err != nil {
		t.Errorf("PutAll() error = %v", err)
		return
	}

	if err := testMap.Sync(tt.newState, tt.opts...); err != nil {
		t.Errorf("Sync() error = %v", err)
		return
	}

	if got, err := testMap.GetAll(tt.opts...); err != nil {
		t.Errorf("GetAll() error = %v", err)
	} else {
		td.Cmp(t, got, tt.newState)
	}
}

func TestMap_Sync(t *testing.T) {
	t.Parallel()
	RemoveMemlock(t)

	tests := []testCase{
		mapSyncTestCase[uint32, uint32]{
			name:         "Initial state empty",
			initialState: make(map[uint32]uint32),
			newState: map[uint32]uint32{
				13: 37,
			},
			mapSpec: ebpf.MapSpec{
				Name:       "map_sync_empty",
				Type:       ebpf.Hash,
				KeySize:    4,
				ValueSize:  4,
				MaxEntries: 1,
			},
		},
		mapSyncTestCase[uint32, uint32]{
			name: "Add single element",
			initialState: map[uint32]uint32{
				13: 37,
			},
			newState: map[uint32]uint32{
				13: 37,
				9:  81,
			},
			mapSpec: ebpf.MapSpec{
				Name:       "map_sync_add_single",
				Type:       ebpf.Hash,
				KeySize:    4,
				ValueSize:  4,
				MaxEntries: 2,
			},
		},
		mapSyncTestCase[uint32, uint32]{
			name: "Delete single element",
			initialState: map[uint32]uint32{
				13: 37,
				9:  81,
			},
			newState: map[uint32]uint32{
				13: 37,
			},
			mapSpec: ebpf.MapSpec{
				Name:       "map_sync_del_single",
				Type:       ebpf.Hash,
				KeySize:    4,
				ValueSize:  4,
				MaxEntries: 2,
			},
		},
		mapSyncTestCase[uint32, uint32]{
			name: "Add and delete elements",
			initialState: map[uint32]uint32{
				13: 37,
				9:  81,
			},
			newState: map[uint32]uint32{
				13: 37,
				45: 89,
				34: 23,
			},
			mapSpec: ebpf.MapSpec{
				Name:       "map_sync_del_single",
				Type:       ebpf.Hash,
				KeySize:    4,
				ValueSize:  4,
				MaxEntries: 3,
			},
		},
		mapSyncTestCase[uint32, uint32]{
			name: "Add single element - force fallback",
			initialState: map[uint32]uint32{
				13: 37,
			},
			newState: map[uint32]uint32{
				13: 37,
				9:  81,
			},
			opts: []netflow.BatchOption{netflow.WithUseFallback(true)},
			mapSpec: ebpf.MapSpec{
				Name:       "map_sync_add_single_fb",
				Type:       ebpf.Hash,
				KeySize:    4,
				ValueSize:  4,
				MaxEntries: 2,
			},
		},
		mapSyncTestCase[uint32, uint32]{
			name: "Delete single element - force fallback",
			initialState: map[uint32]uint32{
				13: 37,
				9:  81,
			},
			newState: map[uint32]uint32{
				13: 37,
			},
			opts: []netflow.BatchOption{netflow.WithUseFallback(true)},
			mapSpec: ebpf.MapSpec{
				Name:       "map_sync_del_single_fb",
				Type:       ebpf.Hash,
				KeySize:    4,
				ValueSize:  4,
				MaxEntries: 2,
			},
		},
		mapSyncTestCase[uint32, uint32]{
			name: "Add and delete elements - force fallback",
			initialState: map[uint32]uint32{
				13: 37,
				9:  81,
			},
			newState: map[uint32]uint32{
				13: 37,
				45: 89,
				34: 23,
			},
			opts: []netflow.BatchOption{netflow.WithUseFallback(true)},
			mapSpec: ebpf.MapSpec{
				Name:       "map_sync_del_single_fb",
				Type:       ebpf.Hash,
				KeySize:    4,
				ValueSize:  4,
				MaxEntries: 3,
			},
		},
	}

	//nolint:paralleltest // done in Run function
	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name(), tt.Run)
	}
}
