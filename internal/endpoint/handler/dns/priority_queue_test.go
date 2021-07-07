package dns_test

import (
	"net"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/td"

	"gitlab.com/inetmock/inetmock/internal/endpoint/handler/dns"
)

func Test_ttlQueue_Evict(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name             string
		t                dns.TTLQueue
		wantEvicted      int
		expectedLength   int
		expectedCapacity int
	}{
		{
			name:           "Evict empty",
			t:              dns.NewQueue(0),
			wantEvicted:    0,
			expectedLength: 0,
		},
		{
			name: "Evict nothing",
			t: dns.NewFromSeed([]*dns.Entry{
				dns.Entry{
					Name:    "gogle.ru",
					Address: net.ParseIP("127.0.0.1"),
				}.WithTTL(60 * time.Second),
			}),
			wantEvicted:    0,
			expectedLength: 1,
		},
		{
			name: "Evict one",
			t: dns.NewFromSeed([]*dns.Entry{
				dns.Entry{
					Name:    "gogle.ru",
					Address: net.ParseIP("127.0.0.1"),
				}.WithTTL(-60 * time.Second),
				dns.Entry{
					Name:    "google.com",
					Address: net.ParseIP("127.0.0.2"),
				}.WithTTL(60 * time.Second),
			}),
			wantEvicted:    1,
			expectedLength: 1,
		},
		{
			name: "Evict multiple",
			t: dns.NewFromSeed([]*dns.Entry{
				dns.Entry{
					Name:    "gogle.ru",
					Address: net.ParseIP("127.0.0.1"),
				}.WithTTL(-60 * time.Second),
				dns.Entry{
					Name:    "gugle.ru",
					Address: net.ParseIP("127.0.0.3"),
				}.WithTTL(-30 * time.Second),
				dns.Entry{
					Name:    "google.com",
					Address: net.ParseIP("127.0.0.2"),
				}.WithTTL(60 * time.Second),
			}),
			wantEvicted:    2,
			expectedLength: 1,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var (
				gotEvicted []*dns.Entry
				wait       = make(chan struct{})
			)

			tt.t.OnEvicted(dns.EvictionCallbackFunc(func(evictedEntries []*dns.Entry) {
				t.Log("Evicted elements")
				gotEvicted = evictedEntries
				close(wait)
			}))

			tt.t.Evict()
			if tt.wantEvicted > 0 {
				t.Log("waiting to get callback")
				for range wait {

				}
			}
			td.CmpLen(t, gotEvicted, tt.wantEvicted)
			td.Cmp(t, tt.t.Len(), tt.expectedLength)
		})
	}
}

func Test_ttlQueue_Push(t *testing.T) {
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
			q := dns.NewQueue(tt.fields.initialCapacity)
			for _, a := range tt.args {
				q.Push(a.name, a.address, a.ttl)
			}
			t.Cmp(q.Len(), tt.expectedLength)
			t.Cmp(q.Cap(), tt.expectedCapacity)
		})
	}
}

func Test_ttlQueue_UpdateTTL(t1 *testing.T) {
	t1.Parallel()
	type seedEntry struct {
		host string
		ip   net.IP
		ttl  time.Duration
	}
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
						host: "mail.gogle.ru",
						ip:   net.ParseIP("1.2.3.4"),
						ttl:  100 * time.Millisecond,
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
						host: "mail.gogle.ru",
						ip:   net.ParseIP("1.2.3.4"),
						ttl:  100 * time.Millisecond,
					},
					{
						host: "asdf.gogle.ru",
						ip:   net.ParseIP("1.2.3.5"),
						ttl:  120 * time.Millisecond,
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
						host: "mail.gogle.ru",
						ip:   net.ParseIP("1.2.3.4"),
						ttl:  50 * time.Millisecond,
					},
					{
						host: "asdf.gogle.ru",
						ip:   net.ParseIP("1.2.3.5"),
						ttl:  150 * time.Millisecond,
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
						host: "mail.gogle.ru",
						ip:   net.ParseIP("1.2.3.4"),
						ttl:  50 * time.Millisecond,
					},
					{
						host: "honey.gogle.ru",
						ip:   net.ParseIP("1.2.3.4"),
						ttl:  100 * time.Millisecond,
					},
					{
						host: "asdf.gogle.ru",
						ip:   net.ParseIP("1.2.3.5"),
						ttl:  150 * time.Millisecond,
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
			var queue = dns.NewQueue(tt.fields.initialCapacity)
			for i := range tt.fields.seeds {
				seed := tt.fields.seeds[i]
				_ = queue.Push(seed.host, seed.ip, seed.ttl)
			}
			var entry = queue.Get(tt.args.idxToUpdate)
			queue.UpdateTTL(entry, tt.args.newTTL)

			var current = queue.PeekFront().TTL()
			for i := 0; i < queue.Len(); i++ {
				if current.After(queue.Get(i).TTL()) {
					t1.Errorf("TTLs in wrong order got = %v, current = %v", queue.Get(i).TTL(), current)
				}
			}
		})
	}
}
