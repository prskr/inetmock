package queue_test

import (
	"math/rand"
	"net"
	"sync/atomic"
	"testing"
	"time"

	"gitlab.com/inetmock/inetmock/internal/queue"
	"gitlab.com/inetmock/inetmock/internal/test"
	"gitlab.com/inetmock/inetmock/protocols/dns"
)

const (
	entryTTL        = 20 * time.Millisecond
	initialCapacity = 500
)

type testData struct {
	domain string
	ip     net.IP
}

func Benchmark_DefaultQueue(b *testing.B) {
	ttl := queue.NewTTL(initialCapacity)
	ttl.OnEvicted(testCallback(b))
	data := generateTestData(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ttl.Push(data[i].domain, data[i].ip, entryTTL)
	}
}

func Benchmark_DefaultQueueParallel(b *testing.B) {
	ttl := queue.NewTTL(initialCapacity)
	ttl.OnEvicted(testCallback(b))
	data := generateTestData(b)
	b.ResetTimer()
	var idx int32 = -1
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			i := atomic.AddInt32(&idx, 1)
			ttl.Push(data[i].domain, data[i].ip, entryTTL)
		}
	})
}

func Benchmark_AutoEvictingQueue(b *testing.B) {
	ttl := queue.WrapToAutoEvict(queue.NewTTL(initialCapacity))
	ttl.OnEvicted(testCallback(b))
	data := generateTestData(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ttl.Push(data[i].domain, data[i].ip, entryTTL)
	}
}

func Benchmark_AutoEvictingQueueParallel(b *testing.B) {
	ttl := queue.WrapToAutoEvict(queue.NewTTL(initialCapacity))
	ttl.OnEvicted(testCallback(b))
	data := generateTestData(b)
	b.ResetTimer()
	var idx int32 = -1
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			i := atomic.AddInt32(&idx, 1)
			ttl.Push(data[i].domain, data[i].ip, entryTTL)
		}
	})
}

func generateTestData(b *testing.B) []testData {
	b.Helper()
	data := make([]testData, 0, b.N)
	//nolint:gosec
	for i := 0; i < b.N; i++ {
		data = append(data, testData{
			domain: test.GenerateDomain(),
			ip:     dns.Uint32ToIP(rand.Uint32()),
		})
	}
	return data
}

func testCallback(tb testing.TB) queue.EvictionCallback {
	tb.Helper()
	return queue.EvictionCallbackFunc(func(evictedEntries []*queue.Entry) {
		if len(evictedEntries) == 0 {
			tb.Logf("Evicted %d entries", len(evictedEntries))
		}
	})
}
