package queue_test

import (
	"math/rand"
	"testing"
	"time"

	"gitlab.com/inetmock/inetmock/internal/endpoint/handler/dns"
	"gitlab.com/inetmock/inetmock/internal/queue"
	"gitlab.com/inetmock/inetmock/internal/test"
)

const (
	entryTTL        = 20 * time.Millisecond
	initialCapacity = 500
)

func testCallback(tb testing.TB) queue.EvictionCallback {
	tb.Helper()
	return queue.EvictionCallbackFunc(func(evictedEntries []*queue.Entry) {
		if len(evictedEntries) == 0 {
			tb.Logf("Evicted %d entries", len(evictedEntries))
		}
	})
}

func Benchmark_DefaultQueue(b *testing.B) {
	ttl := queue.NewTTL(initialCapacity)
	ttl.OnEvicted(testCallback(b))
	for i := 0; i < b.N; i++ {
		//nolint:gosec
		ttl.Push(test.GenerateDomain(), dns.Uint32ToIP(rand.Uint32()), entryTTL)
	}
}

func Benchmark_DefaultQueueParallel(b *testing.B) {
	ttl := queue.NewTTL(initialCapacity)
	ttl.OnEvicted(testCallback(b))
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			//nolint:gosec
			ttl.Push(test.GenerateDomain(), dns.Uint32ToIP(rand.Uint32()), entryTTL)
		}
	})
}

func Benchmark_AutoEvictingQueue(b *testing.B) {
	ttl := queue.WrapToAutoEvict(queue.NewTTL(initialCapacity))
	ttl.OnEvicted(testCallback(b))
	for i := 0; i < b.N; i++ {
		//nolint:gosec
		ttl.Push(test.GenerateDomain(), dns.Uint32ToIP(rand.Uint32()), entryTTL)
	}
}

func Benchmark_AutoEvictingQueueParallel(b *testing.B) {
	ttl := queue.WrapToAutoEvict(queue.NewTTL(initialCapacity))
	ttl.OnEvicted(testCallback(b))
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			//nolint:gosec
			ttl.Push(test.GenerateDomain(), dns.Uint32ToIP(rand.Uint32()), entryTTL)
		}
	})
}
