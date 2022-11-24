//go:build sudo

package netflow_test

import (
	"net/netip"
	"reflect"
	"testing"

	"github.com/cilium/ebpf"
	"github.com/maxatome/go-testdeep/td"

	"inetmock.icb4dc0.de/inetmock/netflow"
)

type testCase interface {
	Run(t *testing.T)
	Name() string
}

func TestEBPFTypes_MarshalUnmarshal(t *testing.T) {
	t.Parallel()
	RemoveMemlock(t)

	tests := []testCase{
		connTrackTypeTest[netflow.ConnMeta]{
			MapName:  "conn_meta_marshal",
			TypeSize: 16,
			Value: netflow.ConnMeta{
				Addr:         netip.MustParseAddr("1.2.3.4"),
				Port:         1234,
				Transport:    netflow.ProtocolTCP,
				LastObserved: 42,
			},
		},
		connTrackTypeTest[netflow.ConnIdent]{
			MapName:  "conn_ident_marshal",
			TypeSize: 12,
			Value: netflow.ConnIdent{
				Addr:      netip.MustParseAddr("5.6.7.8"),
				Port:      5678,
				Transport: netflow.ProtocolUDP,
			},
		},
		connTrackTypeTest[netflow.FirewallRule]{
			MapName:  "fw_rule_marshal",
			TypeSize: 8,
			Value: netflow.FirewallRule{
				Policy:         netflow.XDPActionDrop,
				MonitorTraffic: true,
			},
		},
	}

	//nolint:paralleltest // parallel call is done in struct function
	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name(), tt.Run)
	}
}

type connTrackTypeTest[T comparable] struct {
	MapName  string
	TypeSize uint32
	Value    T
}

func (tt connTrackTypeTest[T]) Name() string {
	return reflect.TypeOf(tt.Value).Name()
}

func (tt connTrackTypeTest[T]) Run(t *testing.T) {
	t.Helper()
	t.Parallel()

	m, err := ebpf.NewMap(&ebpf.MapSpec{
		Name:       tt.MapName,
		Type:       ebpf.Array,
		KeySize:    4,
		ValueSize:  tt.TypeSize,
		MaxEntries: 1,
	})
	if err != nil {
		t.Fatalf("Failed to setup cleanup map: %v", err)
	}

	t.Cleanup(func() {
		if err := m.Close(); err != nil {
			t.Errorf("Failed to close test map: %v", err)
		}
	})

	testMap := netflow.MapOf[uint32, T](m)

	if err := testMap.Put(0, tt.Value); err != nil {
		t.Errorf("Failed to put elem: %v", err)
	}

	if got, err := testMap.Get(0); err != nil {
		t.Errorf("Failed to get elem: %v", err)
	} else {
		td.Cmp(t, got, tt.Value)
	}
}
