package doh_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/td"
	mdns "github.com/miekg/dns"

	"inetmock.icb4dc0.de/inetmock/internal/endpoint"
	auditmock "inetmock.icb4dc0.de/inetmock/internal/mock/audit"
	"inetmock.icb4dc0.de/inetmock/internal/test"
	"inetmock.icb4dc0.de/inetmock/pkg/logging"
	"inetmock.icb4dc0.de/inetmock/protocols/dns/client"
	"inetmock.icb4dc0.de/inetmock/protocols/dns/doh"
)

func Test_dohHandler_Start(t *testing.T) {
	t.Parallel()
	type args struct {
		opts      map[string]any
		query     string
		queryType uint16
	}
	tests := []struct {
		name         string
		args         args
		want         any
		wantErr      bool
		wantQueryErr bool
	}{
		{
			name: "Resolve fake dns.google",
			args: args{
				opts: map[string]any{
					"ttl": 30 * time.Second,
					"cache": map[string]any{
						"type": "none",
					},
					"default": map[string]any{
						"type": "incremental",
						"cidr": "10.10.0.0/16",
					},
					"rules": []string{
						`A(".*\\.google\\.") => IP(1.1.1.1)`,
					},
				},
				query:     "dns.google.",
				queryType: mdns.TypeA,
			},
			want: td.Struct(new(mdns.Msg), td.StructFields{
				"Answer": td.SuperBagOf(td.Struct(new(mdns.A), td.StructFields{
					"A": td.Code(func(ip net.IP) bool {
						return ip.Equal(net.IPv4(1, 1, 1, 1))
					}),
				})),
			}),
			wantErr:      false,
			wantQueryErr: false,
		},
		{
			name: "Resolve fake reddit",
			args: args{
				opts: map[string]any{
					"ttl": 30 * time.Second,
					"cache": map[string]any{
						"type": "none",
					},
					"default": map[string]any{
						"type": "incremental",
						"cidr": "10.10.0.0/16",
					},
					"rules": []string{
						`A('.*\\.reddit\\.com') => IP(2.2.2.2)`,
					},
				},
				query:     "www.reddit.com.",
				queryType: mdns.TypeA,
			},
			want: td.Struct(new(mdns.Msg), td.StructFields{
				"Answer": td.SuperBagOf(td.Struct(new(mdns.A), td.StructFields{
					"A": td.Code(func(ip net.IP) bool {
						return ip.Equal(net.IPv4(2, 2, 2, 2))
					}),
				})),
			}),
			wantErr:      false,
			wantQueryErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			emitterMock := new(auditmock.EmitterMock)
			ctx, cancel := context.WithCancel(context.Background())
			t.Cleanup(cancel)
			d := doh.New(logging.CreateTestLogger(t), emitterMock)
			listener := test.NewInMemoryListener(t)
			resolver := client.Resolver{
				Transport: client.HTTPTransport{
					Client: test.HTTPClientForInMemListener(listener),
					Scheme: "https",
					Server: "one.one.one.one",
				},
			}
			lifecycle := endpoint.NewStartupSpec(t.Name(), endpoint.NewUplink(listener), tt.args.opts)
			if err := d.Start(ctx, lifecycle); (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			question := new(mdns.Msg).SetQuestion(tt.args.query, tt.args.queryType)
			got, err := resolver.Do(ctx, question)
			if err != nil {
				if !tt.wantQueryErr {
					t.Errorf("resolver.Do() err = %v", err)
				}
			}
			td.Cmp(t, got, tt.want)
		})
	}
}
