package dns_test

import (
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/td"

	"gitlab.com/inetmock/inetmock/internal/netutils"
	dns2 "gitlab.com/inetmock/inetmock/protocols/dns"
)

func Test_cache_PutRecord(t *testing.T) {
	t.Parallel()
	type args struct {
		host    string
		address net.IP
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Put localhost",
			args: args{
				host:    "localhost",
				address: net.IPv4(127, 0, 0, 1),
			},
		},
		{
			name: "Put Quad9",
			args: args{
				host:    "dns9.quad9.net",
				address: net.IPv4(9, 9, 9, 9),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(tb *testing.T) {
			tb.Parallel()
			t := td.NewT(tb)
			c := dns2.NewCache(dns2.WithTTL(100*time.Millisecond), dns2.WithInitialSize(500))
			c.PutRecord(tt.args.host, tt.args.address)
			t.Cmp(c.ForwardLookup(tt.args.host), tt.args.address)
			host, miss := c.ReverseLookup(tt.args.address)
			t.Cmp(host, tt.args.host)
			t.Cmp(miss, false)
		})
	}
}

//nolint:gosec
func Test_cache_ForwardLookup(t *testing.T) {
	t.Parallel()
	type args struct {
		host     string
		resolver dns2.IPResolver
	}
	type seed struct {
		host    string
		address net.IP
	}
	tests := []struct {
		name  string
		args  args
		times int
		seeds []seed
		want  interface{}
	}{
		{
			name: "Lookup with known entry",
			args: args{
				host: "dns9.quad9.net",
			},
			seeds: []seed{
				{
					host:    "dns9.quad9.net",
					address: net.IPv4(9, 9, 9, 9),
				},
			},
			want: net.IPv4(9, 9, 9, 9),
		},
		{
			name: "Lookup with resolver",
			args: args{
				host: "mail.gogle.ru",
				resolver: dns2.IPResolverFunc(func(host string) net.IP {
					return netutils.Uint32ToIP(rand.Uint32())
				}),
			},
			want: td.NotNil(),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(tb *testing.T) {
			tb.Parallel()
			t := td.NewT(tb)
			c := dns2.NewCache(dns2.WithTTL(200*time.Millisecond), dns2.WithInitialSize(500))
			for _, s := range tt.seeds {
				c.PutRecord(s.host, s.address)
			}
			var resolved net.IP
			for i := 0; i < tt.times; i++ {
				got := c.ForwardLookup(tt.args.host)
				t.Cmp(got, tt.want)
				if resolved != nil {
					t.Cmp(got, resolved)
				}
				resolved = got
			}
		})
	}
}
