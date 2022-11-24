package client_test

import (
	"context"
	"net"
	"testing"

	"github.com/maxatome/go-testdeep/td"
	mdns "github.com/miekg/dns"

	"inetmock.icb4dc0.de/inetmock/protocols/dns/client"
)

func TestTraditionalTransport_RoundTrip(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		question *mdns.Msg
		wantResp any
		wantErr  bool
	}{
		{
			name:     "Resolve one.one.one.one.",
			question: new(mdns.Msg).SetQuestion(mdns.Fqdn("one.one.one.one"), mdns.TypeA),
			wantResp: td.Struct(new(mdns.Msg), td.StructFields{
				"Answer": td.SuperBagOf(td.Struct(new(mdns.A), td.StructFields{
					"A": td.Code(func(ip net.IP) bool {
						return ip.Equal(net.IPv4(1, 1, 1, 1))
					}),
				})),
			}),
			wantErr: false,
		},
		{
			name:     "Resolve 1.1.1.1",
			question: new(mdns.Msg).SetQuestion("1.1.1.1.in-addr.arpa.", mdns.TypePTR),
			wantResp: td.Struct(new(mdns.Msg), td.StructFields{
				"Answer": td.SuperBagOf(td.Struct(new(mdns.PTR), td.StructFields{
					"Ptr": mdns.Fqdn("one.one.one.one"),
				})),
			}),
			wantErr: false,
		},
		{
			name:     "Resolve dns9.quad9.net.",
			question: new(mdns.Msg).SetQuestion(mdns.Fqdn("dns9.quad9.net"), mdns.TypeA),
			wantResp: td.Struct(new(mdns.Msg), td.StructFields{
				"Answer": td.SuperBagOf(td.Struct(new(mdns.A), td.StructFields{
					"A": td.Code(func(ip net.IP) bool {
						return ip.Equal(net.IPv4(9, 9, 9, 9))
					}),
				})),
			}),
			wantErr: false,
		},
		{
			name:     "Resolve 9.9.9.9",
			question: new(mdns.Msg).SetQuestion("9.9.9.9.in-addr.arpa.", mdns.TypePTR),
			wantResp: td.Struct(new(mdns.Msg), td.StructFields{
				"Answer": td.SuperBagOf(td.Struct(new(mdns.PTR), td.StructFields{
					"Ptr": mdns.Fqdn("dns9.quad9.net"),
				})),
			}),
			wantErr: false,
		},
		{
			name:     "Resolve dns.google.",
			question: new(mdns.Msg).SetQuestion(mdns.Fqdn("dns.google"), mdns.TypeA),
			wantResp: td.Struct(new(mdns.Msg), td.StructFields{
				"Answer": td.SuperBagOf(td.Struct(new(mdns.A), td.StructFields{
					"A": td.Code(func(ip net.IP) bool {
						return ip.Equal(net.IPv4(8, 8, 8, 8))
					}),
				})),
			}),
			wantErr: false,
		},
		{
			name:     "Resolve 8.8.8.8",
			question: new(mdns.Msg).SetQuestion("8.8.8.8.in-addr.arpa.", mdns.TypePTR),
			wantResp: td.Struct(new(mdns.Msg), td.StructFields{
				"Answer": td.SuperBagOf(td.Struct(new(mdns.PTR), td.StructFields{
					"Ptr": mdns.Fqdn("dns.google"),
				})),
			}),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t1 *testing.T) {
			t1.Parallel()
			t := td.NewT(t)

			ctx, cancel := context.WithCancel(context.Background())
			t.Cleanup(cancel)
			resolver := client.Resolver{Transport: &client.TraditionalTransport{
				Network: "tcp",
				Address: "9.9.9.9:53",
			}}
			gotResp, err := resolver.Do(ctx, tt.question)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("Resolver.Do() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			t.Cmp(gotResp, tt.wantResp)
		})
	}
}
