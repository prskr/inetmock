//go:build sudo

package netflow_test

import (
	"testing"
	"time"

	"github.com/cilium/ebpf"

	"inetmock.icb4dc0.de/inetmock/netflow"
)

func TestEpoch_Sync(t *testing.T) {
	t.Parallel()
	RemoveMemlock(t)

	m, err := ebpf.NewMap(&ebpf.MapSpec{
		Name:       "sync_epoch_0",
		Type:       ebpf.Array,
		KeySize:    4,
		ValueSize:  4,
		MaxEntries: 2,
	})
	if err != nil {
		t.Fatalf("Failed to setup cleanup map: %v", err)
	}

	t.Cleanup(func() {
		if err := m.Close(); err != nil {
			t.Errorf("Failed to close test map: %v", err)
		}
	})

	configMap := netflow.MapOf[netflow.NATConfigKey, uint32](m)
	epoch := netflow.NewEpoch(configMap, netflow.WithStart(time.Now().Add(-5*time.Second)))

	if err := epoch.Sync(); err != nil {
		t.Errorf("Failed to sync epoch: %v", err)
		return
	}

	if val, err := configMap.Get(netflow.NATConfigKey(0)); err != nil {
		t.Errorf("Failed to get config key value: %v", err)
	} else if val < 1 {
		t.Errorf("Expected value greater 0 but got %d", val)
	}
}

func TestEpoch_StartStopSync(t *testing.T) {
	t.Parallel()
	RemoveMemlock(t)

	m, err := ebpf.NewMap(&ebpf.MapSpec{
		Name:       "start_stop_sync",
		Type:       ebpf.Array,
		KeySize:    4,
		ValueSize:  4,
		MaxEntries: 2,
	})
	if err != nil {
		t.Fatalf("Failed to setup cleanup map: %v", err)
	}

	t.Cleanup(func() {
		if err := m.Close(); err != nil {
			t.Errorf("Failed to close test map: %v", err)
		}
	})

	configMap := netflow.MapOf[netflow.NATConfigKey, uint32](m)
	epoch := netflow.NewEpoch(configMap, netflow.WithStart(time.Now().Add(-5*time.Second)))

	syncWindowSize := 50 * time.Millisecond
	if err := epoch.StartSync(syncWindowSize); err != nil {
		t.Errorf("Failed to start sync: %v", err)
	}

	time.Sleep(2 * syncWindowSize)

	if val, err := configMap.Get(netflow.NATConfigKey(0)); err != nil {
		t.Errorf("Failed to get config key value: %v", err)
	} else if val < 1 {
		t.Errorf("Expected value greater 0 but got %d", val)
	}

	epoch.Stop()
}
