package queue_test

import (
	"net"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/td"

	"gitlab.com/inetmock/inetmock/internal/endpoint/handler/dns"
	"gitlab.com/inetmock/inetmock/internal/queue"
)

type seedEntry struct {
	value *dns.Record
	ttl   time.Duration
}

func seed(vals ...seedEntry) []*queue.Entry {
	s := make([]*queue.Entry, 0, len(vals))
	for idx := range vals {
		s = append(s, (&queue.Entry{
			Key:   vals[idx].value.Name,
			Value: vals[idx].value,
		}).WithTTL(vals[idx].ttl))
	}
	return s
}

func Test_TTL_Evict(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name             string
		t                queue.TTL
		wantEvicted      int
		expectedLength   int
		expectedCapacity int
	}{
		{
			name:           "Evict empty",
			t:              queue.NewTTL(0),
			wantEvicted:    0,
			expectedLength: 0,
		},
		{
			name: "Evict nothing",
			t: queue.NewTTLFromSeed(
				seed(
					seedEntry{
						value: &dns.Record{
							Name:    "gogle.ru",
							Address: net.ParseIP("127.0.0.1"),
						},
						ttl: 60 * time.Second,
					},
				)),
			wantEvicted:    0,
			expectedLength: 1,
		},
		{
			name: "Evict one",
			t: queue.NewTTLFromSeed(
				seed(
					seedEntry{
						value: &dns.Record{
							Name:    "gogle.ru",
							Address: net.ParseIP("127.0.0.1"),
						},
						ttl: -60 * time.Second,
					},
					seedEntry{
						value: &dns.Record{
							Name:    "google.com",
							Address: net.ParseIP("127.0.0.2"),
						},
						ttl: 60 * time.Second,
					},
				)),
			wantEvicted:    1,
			expectedLength: 1,
		},
		{
			name: "Evict multiple",
			t: queue.NewTTLFromSeed(
				seed(
					seedEntry{
						value: &dns.Record{
							Name:    "gogle.ru",
							Address: net.ParseIP("127.0.0.1"),
						},
						ttl: -60 * time.Second,
					},
					seedEntry{
						value: &dns.Record{
							Name:    "gogle.ru",
							Address: net.ParseIP("127.0.0.3"),
						},
						ttl: -30 * time.Second,
					},
					seedEntry{
						value: &dns.Record{
							Name:    "google.com",
							Address: net.ParseIP("127.0.0.2"),
						},
						ttl: 60 * time.Second,
					},
				)),
			wantEvicted:    2,
			expectedLength: 1,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t1 *testing.T) {
			t1.Parallel()
			var (
				gotEvicted []*queue.Entry
				wait       = make(chan struct{})
				t          = td.NewT(t1)
			)

			tt.t.OnEvicted(queue.EvictionCallbackFunc(func(evictedEntries []*queue.Entry) {
				t1.Log("Evicted elements")
				gotEvicted = evictedEntries
				close(wait)
			}))

			tt.t.Evict()
			if tt.wantEvicted > 0 {
				t1.Log("waiting to get callback")
				for range wait {

				}
			}
			td.CmpLen(t1, gotEvicted, tt.wantEvicted)
			t.Cmp(tt.t.Len(), tt.expectedLength)
			validateQueue(t, tt.t)
		})
	}
}

func Test_TTL_Push(t *testing.T) {
	t.Parallel()
	type fields struct {
		initialCapacity int
	}
	type args struct {
		name    string
		address net.IP
		ttl     time.Duration
	}
	tests := []struct {
		name             string
		fields           fields
		args             []args
		expectedLength   int
		expectedCapacity int
	}{
		{
			name: "Push multiple elements",
			fields: fields{
				initialCapacity: 10,
			},
			args: []args{
				{
					name:    "mail.google.ru",
					address: net.ParseIP("192.168.199.10"),
					ttl:     10 * time.Millisecond,
				},
				{
					name:    "www.google.ru",
					address: net.ParseIP("192.168.199.11"),
					ttl:     10 * time.Millisecond,
				},
			},
			expectedLength:   2,
			expectedCapacity: 10,
		},
		{
			name: "Push multiple elements with unsorted TTLs",
			fields: fields{
				initialCapacity: 10,
			},
			args: []args{
				{
					name:    "mail.google.ru",
					address: net.ParseIP("192.168.199.10"),
					ttl:     100 * time.Millisecond,
				},
				{
					name:    "first.google.ru",
					address: net.ParseIP("192.168.199.10"),
					ttl:     10 * time.Millisecond,
				},
				{
					name:    "www.google.ru",
					address: net.ParseIP("192.168.199.11"),
					ttl:     200 * time.Millisecond,
				},
			},
			expectedLength:   3,
			expectedCapacity: 10,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(tb *testing.T) {
			tb.Parallel()
			t := td.NewT(tb)
			q := queue.NewTTL(tt.fields.initialCapacity)
			for _, a := range tt.args {
				q.Push(a.name, a.address, a.ttl)
			}
			t.Cmp(q.Len(), tt.expectedLength)
			t.Cmp(q.Cap(), tt.expectedCapacity)
			validateQueue(t, q)
		})
	}
}

func Test_TTL_UpdateTTL(t1 *testing.T) {
	t1.Parallel()
	type fields struct {
		initialCapacity int
		seeds           []seedEntry
	}
	type args struct {
		idxToUpdate int
		newTTL      time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Single element",
			fields: fields{
				initialCapacity: 10,
				seeds: []seedEntry{
					{
						value: &dns.Record{
							Name:    "mail.gogle.ru",
							Address: net.ParseIP("1.2.3.4"),
						},
						ttl: 100 * time.Millisecond,
					},
				},
			},
			args: args{
				idxToUpdate: 0,
				newTTL:      200 * time.Millisecond,
			},
		},
		{
			name: "Last element",
			fields: fields{
				initialCapacity: 10,
				seeds: []seedEntry{
					{
						value: &dns.Record{
							Name:    "mail.gogle.ru",
							Address: net.ParseIP("1.2.3.4"),
						},
						ttl: 100 * time.Millisecond,
					},
					{
						value: &dns.Record{
							Name:    "asdf.gogle.ru",
							Address: net.ParseIP("1.2.3.5"),
						},
						ttl: 120 * time.Millisecond,
					},
				},
			},
			args: args{
				idxToUpdate: 1,
				newTTL:      200 * time.Millisecond,
			},
		},
		{
			name: "Switch elements",
			fields: fields{
				initialCapacity: 10,
				seeds: []seedEntry{
					{
						value: &dns.Record{
							Name:    "mail.gogle.ru",
							Address: net.ParseIP("1.2.3.4"),
						},
						ttl: 50 * time.Millisecond,
					},
					{
						value: &dns.Record{
							Name:    "asdf.gogle.ru",
							Address: net.ParseIP("1.2.3.5"),
						},
						ttl: 150 * time.Millisecond,
					},
				},
			},
			args: args{
				idxToUpdate: 0,
				newTTL:      500 * time.Millisecond,
			},
		},
		{
			name: "Switch elements",
			fields: fields{
				initialCapacity: 10,
				seeds: []seedEntry{
					{
						value: &dns.Record{
							Name:    "mail.gogle.ru",
							Address: net.ParseIP("1.2.3.4"),
						},
						ttl: 50 * time.Millisecond,
					},
					{
						value: &dns.Record{
							Name:    "honey.gogle.ru",
							Address: net.ParseIP("1.2.3.6"),
						},
						ttl: 100 * time.Millisecond,
					},
					{
						value: &dns.Record{
							Name:    "asdf.gogle.ru",
							Address: net.ParseIP("1.2.3.5"),
						},
						ttl: 150 * time.Millisecond,
					},
				},
			},
			args: args{
				idxToUpdate: 1,
				newTTL:      500 * time.Millisecond,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t1.Run(tt.name, func(t1 *testing.T) {
			t1.Parallel()
			var ttlQueue = queue.NewTTL(tt.fields.initialCapacity)
			for i := range tt.fields.seeds {
				seed := tt.fields.seeds[i]
				_ = ttlQueue.Push(seed.value.Name, seed.value, seed.ttl)
			}
			var entry = ttlQueue.Get(tt.args.idxToUpdate)
			ttlQueue.UpdateTTL(entry, tt.args.newTTL)

			validateQueue(td.NewT(t1), ttlQueue)
		})
	}
}

var (
	baseSeeds = []seedEntry{
		{
			value: &dns.Record{
				Name:    "a.gogle.ru",
				Address: net.ParseIP("1.2.3.4"),
			},
			ttl: -150 * time.Millisecond,
		},
		{
			value: &dns.Record{
				Name:    "b.gogle.ru",
				Address: net.ParseIP("1.2.3.5"),
			},
			ttl: -140 * time.Millisecond,
		},
		{
			value: &dns.Record{
				Name:    "c.gogle.ru",
				Address: net.ParseIP("1.2.3.6"),
			},
			ttl: 450 * time.Millisecond,
		},
		{
			value: &dns.Record{
				Name:    "d.gogle.ru",
				Address: net.ParseIP("1.2.3.7"),
			},
			ttl: 500 * time.Millisecond,
		},
		{
			value: &dns.Record{
				Name:    "e.gogle.ru",
				Address: net.ParseIP("1.2.3.8"),
			},
			ttl: 600 * time.Millisecond,
		},
	}
)

func Test_TTL_PushAfterEviction(t1 *testing.T) {
	t1.Parallel()
	tests := []struct {
		name                   string
		seeds                  []seedEntry
		idxToUpdate            int
		newTTL                 time.Duration
		wantItemsAfterEviction int
	}{
		{
			name:                   "Update element at index 0 of 3",
			seeds:                  baseSeeds,
			wantItemsAfterEviction: 3,
			idxToUpdate:            0,
			newTTL:                 800 * time.Millisecond,
		},
		{
			name:                   "Update element at index 1 of 3",
			seeds:                  baseSeeds,
			wantItemsAfterEviction: 3,
			idxToUpdate:            1,
			newTTL:                 800 * time.Millisecond,
		},
		{
			name:                   "Update element at index 2 of 3",
			seeds:                  baseSeeds,
			wantItemsAfterEviction: 3,
			idxToUpdate:            2,
			newTTL:                 800 * time.Millisecond,
		},
	}
	for _, tt := range tests {
		tt := tt
		t1.Run(tt.name, func(t1 *testing.T) {
			t1.Parallel()
			ttlQueue := queue.NewTTL(20)
			t := td.NewT(t1)
			for _, s := range tt.seeds {
				ttlQueue.Push(s.value.Name, s.value, s.ttl)
			}

			ttlQueue.Evict()
			t.Cmp(ttlQueue.Len(), tt.wantItemsAfterEviction)
			e := ttlQueue.Get(tt.idxToUpdate)
			ttlQueue.UpdateTTL(e, tt.newTTL)
			validateQueue(t, ttlQueue)
		})
	}
}

func validateQueue(t *td.T, q queue.TTL) {
	t.Helper()
	if q.Len() < 1 {
		return
	}
	var current = q.PeekFront().TTL()
	for i := 0; i < q.Len(); i++ {
		entry := q.Get(i)
		t.Cmp(q.IndexOf(entry), i)
		if current.After(entry.TTL()) {
			t.Errorf("TTLs in wrong order got = %v, current = %v", entry.TTL(), current)
		}
	}
}
