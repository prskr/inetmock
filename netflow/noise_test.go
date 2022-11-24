package netflow_test

import (
	"context"
	"math"
	"math/rand"
	"net"
	"strconv"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"
)

func MakeSomeNoise(tb testing.TB, interval time.Duration) {
	tb.Helper()
	ticker := time.NewTicker(interval)
	ctx, cancel := context.WithCancel(context.Background())
	tb.Cleanup(func() {
		cancel()
		ticker.Stop()
	})
	grp, grpCtx := errgroup.WithContext(ctx)
	dialer := net.Dialer{
		Timeout: interval,
	}
	for range ticker.C {
		grp.Go(func() error {
			//nolint:gosec // close enough here
			if conn, err := dialer.DialContext(grpCtx, "tcp", net.JoinHostPort("localhost", strconv.Itoa(rand.Intn(math.MaxUint16)))); err == nil {
				_ = conn.Close()
			}
			return nil
		})
	}
}
