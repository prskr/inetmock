//go:build sudo

package netflow_test

import (
	"math"
	"math/rand"
	"net/netip"
	"testing"

	"github.com/cilium/ebpf"

	"inetmock.icb4dc0.de/inetmock/netflow"
)

func TestConnTrackCleaner_Cleanup(t *testing.T) {
	t.Parallel()
	RemoveMemlock(t)

	type args struct {
		highWaterMark float64
		interfaceName string
		mapSpec       ebpf.MapSpec
		batchOptions  []netflow.BatchOption
	}
	tests := []struct {
		name     string
		args     args
		mapSetup func(tb testing.TB, m *netflow.Map[netflow.ConnIdent, netflow.ConnMeta])
	}{
		{
			name: "No cleanup necessary",
			args: args{
				highWaterMark: 0.5,
				interfaceName: "lo",
				mapSpec: ebpf.MapSpec{
					Name:       "no_cleanup_necessary",
					Type:       ebpf.Hash,
					KeySize:    12,
					ValueSize:  16,
					MaxEntries: 100,
				},
			},
			mapSetup: randomEvents(1),
		},
		{
			name: "Cleanup some events",
			args: args{
				highWaterMark: 0.5,
				interfaceName: "lo",
				mapSpec: ebpf.MapSpec{
					Name:       "cleanup_some_evs",
					Type:       ebpf.Hash,
					KeySize:    12,
					ValueSize:  16,
					MaxEntries: 100,
				},
			},
			mapSetup: randomEvents(60),
		},
		{
			name: "Cleanup some events - always use fallback",
			args: args{
				highWaterMark: 0.5,
				interfaceName: "lo",
				mapSpec: ebpf.MapSpec{
					Name:       "cleanup_some_evs_with_fallback",
					Type:       ebpf.Hash,
					KeySize:    12,
					ValueSize:  16,
					MaxEntries: 100,
				},
				batchOptions: []netflow.BatchOption{netflow.WithUseFallback(true)},
			},
			mapSetup: randomEvents(60),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			m, err := ebpf.NewMap(&tt.args.mapSpec)
			if err != nil {
				t.Fatalf("Failed to setup cleanup map: %v", err)
			}

			t.Cleanup(func() {
				if err := m.Close(); err != nil {
					t.Errorf("Failed to close test map: %v", err)
				}
			})

			connTrackMap := netflow.MapOf[netflow.ConnIdent, netflow.ConnMeta](m)

			errSink := netflow.ErrorSinkFunc(func(err error) {
				t.Errorf("Error occurred: %v", err)
			})

			tt.mapSetup(t, connTrackMap)

			cleaner := netflow.NewConnTrackCleaner(connTrackMap, errSink, tt.args.highWaterMark, tt.args.interfaceName)
			cleaner.Cleanup()

			if all, err := connTrackMap.GetAll(tt.args.batchOptions...); err != nil {
				t.Errorf("Failed to get all entries: %v", err)
			} else if float64(len(all))/float64(tt.args.mapSpec.MaxEntries) > tt.args.highWaterMark {
				t.Errorf("Still too many entries in map")
			}
		})
	}
}

func randomEvents(eventCount int) func(tb testing.TB, m *netflow.Map[netflow.ConnIdent, netflow.ConnMeta]) {
	return func(tb testing.TB, m *netflow.Map[netflow.ConnIdent, netflow.ConnMeta]) {
		tb.Helper()

		for i := 0; i < eventCount; i++ {
			//nolint:gosec // close enough here
			id := netflow.ConnIdent{Addr: netip.MustParseAddr("10.10.1.1"), Port: uint16(rand.Intn(math.MaxUint16)), Transport: netflow.ProtocolTCP}
			meta := netflow.ConnMeta{Addr: netip.MustParseAddr("10.10.1.42"), LastObserved: uint32(i)}
			if err := m.Put(id, meta); err != nil {
				tb.Fatalf("Failed to setup map: %v", err)
			}
		}
	}
}
