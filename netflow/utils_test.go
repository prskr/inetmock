//go:build sudo

package netflow_test

import (
	"sync"
	"testing"

	"github.com/cilium/ebpf/rlimit"
)

var memlockRemovalOnce sync.Once

func RemoveMemlock(tb testing.TB) {
	tb.Helper()
	memlockRemovalOnce.Do(func() {
		if err := rlimit.RemoveMemlock(); err != nil {
			tb.Fatalf("Failed to remove memlock: %v", err)
		}
	})
}
