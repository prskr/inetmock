package client_test

import (
	"context"
	"net"
	"testing"

	"github.com/maxatome/go-testdeep/td"

	"inetmock.icb4dc0.de/inetmock/protocols/dns/client"
)

const resolveRetries = 10

func TestResolver_LookupA(t *testing.T) {
	t.Parallel()
	type args struct {
		host string
	}
	tests := []struct {
		name    string
		args    args
		wantRes any
		wantErr bool
	}{
		{
			name: "Resolve dns.google - incomplete FQDN",
			args: args{
				host: "dns.google",
			},
			wantRes: td.SuperBagOf(td.Code(func(ip net.IP) bool {
				return ip.Equal(net.IPv4(8, 8, 8, 8))
			})),
			wantErr: false,
		},
		{
			name: "Resolve dns.google - complete FQDN",
			args: args{
				host: "dns.google.",
			},
			wantRes: td.SuperBagOf(td.Code(func(ip net.IP) bool {
				return ip.Equal(net.IPv4(8, 8, 8, 8))
			})),
			wantErr: false,
		},
		{
			name: "Resolve one.one.one.one",
			args: args{
				host: "one.one.one.one",
			},
			wantRes: td.SuperBagOf(td.Code(func(ip net.IP) bool {
				return ip.Equal(net.IPv4(1, 1, 1, 1))
			})),
			wantErr: false,
		},
		{
			name: "Resolve dns9.quad9.net",
			args: args{
				host: "dns9.quad9.net",
			},
			wantRes: td.SuperBagOf(td.Code(func(ip net.IP) bool {
				return ip.Equal(net.IPv4(9, 9, 9, 9))
			})),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := client.Resolver{
				Transport: &client.TraditionalTransport{
					Network: "tcp",
					Address: "9.9.9.9:53",
				},
			}
			ctx, cancel := context.WithCancel(context.Background())
			t.Cleanup(cancel)
			var gotRes []net.IP
			err := retry(resolveRetries, func() (err error) {
				gotRes, err = r.LookupA(ctx, tt.args.host)
				return
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("LookupA() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			td.Cmp(t, gotRes, tt.wantRes)
		})
	}
}

func TestResolver_LookupAAAA(t *testing.T) {
	t.Parallel()
	type args struct {
		host string
	}
	tests := []struct {
		name    string
		args    args
		wantRes any
		wantErr bool
	}{
		{
			name: "Resolve dns.google - incomplete FQDN",
			args: args{
				host: "dns.google",
			},
			wantRes: td.SuperBagOf(td.Code(func(ip net.IP) bool {
				return ip.Equal(net.ParseIP("2001:4860:4860::8888"))
			})),
			wantErr: false,
		},
		{
			name: "Resolve dns.google - complete FQDN",
			args: args{
				host: "dns.google.",
			},
			wantRes: td.SuperBagOf(td.Code(func(ip net.IP) bool {
				return ip.Equal(net.ParseIP("2001:4860:4860::8888"))
			})),
			wantErr: false,
		},
		{
			name: "Resolve one.one.one.one",
			args: args{
				host: "one.one.one.one",
			},
			wantRes: td.SuperBagOf(td.Code(func(ip net.IP) bool {
				return ip.Equal(net.ParseIP("2606:4700:4700::1111"))
			})),
			wantErr: false,
		},
		{
			name: "Resolve dns9.quad9.net",
			args: args{
				host: "dns9.quad9.net",
			},
			wantRes: td.SuperBagOf(td.Code(func(ip net.IP) bool {
				return ip.Equal(net.ParseIP("2620:fe::fe:9"))
			})),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := client.Resolver{
				Transport: &client.TraditionalTransport{
					Network: "tcp",
					Address: "9.9.9.9:53",
				},
			}
			ctx, cancel := context.WithCancel(context.Background())
			t.Cleanup(cancel)
			var gotResp []net.IP
			err := retry(resolveRetries, func() (err error) {
				gotResp, err = r.LookupAAAA(ctx, tt.args.host)
				return
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("LookupAAAA() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			td.Cmp(t, gotResp, tt.wantRes)
		})
	}
}

func TestResolver_LookupPTR(t *testing.T) {
	t.Parallel()
	type args struct {
		inAddrArpa string
	}
	tests := []struct {
		name    string
		args    args
		wantRes any
		wantErr bool
	}{
		{
			name: "Resolve PTR 8.8.8.8 - invalid PTR syntax",
			args: args{
				inAddrArpa: "8.8.8.8",
			},
			wantRes: td.SuperBagOf("dns.google."),
			wantErr: false,
		},
		{
			name: "Resolve PTR 8.8.8.8 - valid PTR syntax",
			args: args{
				inAddrArpa: "8.8.8.8.in-addr.arpa",
			},
			wantRes: td.SuperBagOf("dns.google."),
			wantErr: false,
		},
		{
			name: "Resolve PTR 9.9.9.9",
			args: args{
				inAddrArpa: "9.9.9.9",
			},
			wantRes: td.SuperBagOf("dns9.quad9.net."),
			wantErr: false,
		},
		{
			name: "Resolve PTR 1.1.1.1",
			args: args{
				inAddrArpa: "1.1.1.1",
			},
			wantRes: td.SuperBagOf("one.one.one.one."),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := client.Resolver{
				Transport: &client.TraditionalTransport{
					Network: "tcp",
					Address: "9.9.9.9:53",
				},
			}
			ctx, cancel := context.WithCancel(context.Background())
			t.Cleanup(cancel)
			var gotRes []string
			err := retry(resolveRetries, func() (err error) {
				gotRes, err = r.LookupPTR(ctx, tt.args.inAddrArpa)
				return
			})

			if (err != nil) != tt.wantErr {
				t.Errorf("LookupPTR() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			td.Cmp(t, gotRes, tt.wantRes)
		})
	}
}

func retry(maxRetries int, retryTarget func() error) (err error) {
	for i := 0; i < maxRetries; i++ {
		if err = retryTarget(); err == nil {
			return nil
		}
	}
	return err
}
