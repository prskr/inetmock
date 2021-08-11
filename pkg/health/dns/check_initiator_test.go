package dns_test

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/maxatome/go-testdeep/td"

	"gitlab.com/inetmock/inetmock/internal/rules"
	"gitlab.com/inetmock/inetmock/pkg/health/dns"
)

var _ dns.Resolver = (*mockResolver)(nil)

type mockResolver struct {
	LookupAddrDelegate func(ctx context.Context, addr string) (names []string, err error)
	LookupHostDelegate func(ctx context.Context, host string) (addrs []string, err error)
}

func (r *mockResolver) LookupHost(ctx context.Context, host string) (addrs []string, err error) {
	if r == nil || r.LookupHostDelegate == nil {
		return nil, nil
	}
	return r.LookupHostDelegate(ctx, host)
}

func (r *mockResolver) LookupAddr(ctx context.Context, addr string) (names []string, err error) {
	if r == nil || r.LookupAddrDelegate == nil {
		return nil, nil
	}
	return r.LookupAddrDelegate(ctx, addr)
}

func TestCheckForRule(t *testing.T) {
	t.Parallel()
	type args struct {
		rule     string
		resolver dns.Resolver
	}
	tests := []struct {
		name          string
		args          args
		wantResp      interface{}
		wantParseErr  bool
		wantErr       bool
		wantResolvErr bool
	}{
		{
			name: "A record lookup initiator",
			args: args{
				rule: `dns.A("gitlab.com")`,
				resolver: &mockResolver{
					LookupHostDelegate: func(context.Context, string) (addrs []string, err error) {
						return []string{"192.168.0.11"}, nil
					},
				},
			},
			wantResp: td.Struct(&dns.Response{
				Addresses: []net.IP{net.IPv4(192, 168, 0, 11)},
			}, td.StructFields{}),
			wantParseErr:  false,
			wantErr:       false,
			wantResolvErr: false,
		},
		{
			name: "PTR lookup initiator",
			args: args{
				rule: `dns.PTR(192.168.0.11)`,
				resolver: &mockResolver{
					LookupAddrDelegate: func(ctx context.Context, addr string) (names []string, err error) {
						return []string{"google.com"}, nil
					},
				},
			},
			wantResp: td.Struct(&dns.Response{
				Hosts: []string{"google.com"},
			}, td.StructFields{}),
			wantParseErr:  false,
			wantErr:       false,
			wantResolvErr: false,
		},
		{
			name: "Check misses 'dns.' module",
			args: args{
				rule: `A("gitlab.com")`,
			},
			wantParseErr:  false,
			wantErr:       true,
			wantResolvErr: false,
		},
		{
			name: "Check is not recognized",
			args: args{
				rule: `dns.SRV("smtp", "tcp", "inetmock.loc")`,
			},
			wantParseErr:  false,
			wantErr:       true,
			wantResolvErr: false,
		},
		{
			name: "Check is wrong formatted",
			args: args{
				rule: `dns.A("gitlab.com)`,
			},
			wantParseErr:  true,
			wantErr:       false,
			wantResolvErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var (
				parsedRule = new(rules.Check)
				initiator  dns.Initiator
				resp       *dns.Response
				err        error
			)
			if err = rules.Parse(tt.args.rule, parsedRule); err != nil {
				if !tt.wantParseErr {
					t.Errorf("rules.Parse() error = %v", err)
				}
				return
			}

			if initiator, err = dns.CheckForRule(parsedRule); err != nil {
				if !tt.wantErr {
					t.Errorf("CheckForRule() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			var ctx, cancel = context.WithCancel(context.Background())
			t.Cleanup(cancel)

			if resp, err = initiator.Do(ctx, tt.args.resolver); err != nil {
				if !tt.wantResolvErr {
					t.Errorf("initiator.Do() error = %v, wantResolvErr = %t", err, tt.wantResolvErr)
				}
				return
			}
			td.Cmp(t, resp, tt.wantResp)
		})
	}
}

func TestPTRInitiator(t *testing.T) {
	t.Parallel()
	type args struct {
		args     []rules.Param
		resolver dns.Resolver
	}
	tests := []struct {
		name          string
		args          args
		wantResp      interface{}
		wantErr       bool
		wantResolvErr bool
	}{
		{
			name: "Mocked resolver - expect empty list of addresses",
			args: args{
				args: []rules.Param{
					{
						IP: net.IPv4(192, 168, 0, 10),
					},
				},
				resolver: new(mockResolver),
			},
			wantResp: td.Struct(new(dns.Response), td.StructFields{
				"Hosts": td.Empty(),
			}),
			wantErr:       false,
			wantResolvErr: false,
		},
		{
			name: "Mocked resolver - expect a single host",
			args: args{
				args: []rules.Param{
					{
						IP: net.IPv4(192, 168, 0, 10),
					},
				},
				resolver: &mockResolver{
					LookupAddrDelegate: func(context.Context, string) (names []string, err error) {
						return []string{"my-laptop.fritz.box"}, nil
					},
				},
			},
			wantResp: td.Struct(&dns.Response{
				Hosts: []string{"my-laptop.fritz.box"},
			}, td.StructFields{}),
			wantErr:       false,
			wantResolvErr: false,
		},
		{
			name: "Missing param for PTR initiator",
			args: args{
				args: make([]rules.Param, 0),
			},
			wantErr: true,
		},
		{
			name: "Wrong param type for PTR initiator",
			args: args{
				args: []rules.Param{
					{
						Int: rules.IntP(42),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Return error from mocked resolver",
			args: args{
				args: []rules.Param{
					{
						IP: net.IPv4(192, 168, 0, 10),
					},
				},
				resolver: &mockResolver{
					LookupAddrDelegate: func(context.Context, string) (names []string, err error) {
						return nil, errors.New("some random error")
					},
				},
			},
			wantResp: td.Struct(&dns.Response{
				Hosts: []string{"my-laptop.fritz.box"},
			}, td.StructFields{}),
			wantErr:       false,
			wantResolvErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var (
				initiator dns.Initiator
				resp      *dns.Response
				err       error
			)
			initiator, err = dns.PTRInitiator(tt.args.args...)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("PTRInitiator() error = %v, wantErr %t", err, tt.wantErr)
				}
				return
			}
			var ctx, cancel = context.WithCancel(context.Background())
			t.Cleanup(cancel)
			resp, err = initiator.Do(ctx, tt.args.resolver)
			if err != nil {
				if !tt.wantResolvErr {
					t.Errorf("initiator.Do() error = %v, wantResolvErr %t", err, tt.wantResolvErr)
				}
				return
			}
			td.Cmp(t, resp, tt.wantResp)
		})
	}
}

func TestAorAAAAInitiator(t *testing.T) {
	t.Parallel()
	type args struct {
		args     []rules.Param
		resolver dns.Resolver
	}
	tests := []struct {
		name          string
		args          args
		wantResp      interface{}
		wantErr       bool
		wantResolvErr bool
	}{
		{
			name: "Mocked resolver - expect empty list of addresses",
			args: args{
				args: []rules.Param{
					{
						String: rules.StringP("gitlab.com"),
					},
				},
				resolver: new(mockResolver),
			},
			wantResp: td.Struct(new(dns.Response), td.StructFields{
				"Addresses": td.Empty(),
			}),
			wantErr:       false,
			wantResolvErr: false,
		},
		{
			name: "Mocked resolver - expect single result",
			args: args{
				args: []rules.Param{
					{
						String: rules.StringP("gitlab.com"),
					},
				},
				resolver: &mockResolver{
					LookupHostDelegate: func(ctx context.Context, host string) (addrs []string, err error) {
						return []string{"192.168.0.12"}, nil
					},
				},
			},
			wantResp: td.Struct(&dns.Response{
				Addresses: []net.IP{net.IPv4(192, 168, 0, 12)},
			}, td.StructFields{}),
			wantErr:       false,
			wantResolvErr: false,
		},
		{
			name: "Mocked resolver - expect multiple result",
			args: args{
				args: []rules.Param{
					{
						String: rules.StringP("gitlab.com"),
					},
				},
				resolver: &mockResolver{
					LookupHostDelegate: func(ctx context.Context, host string) (addrs []string, err error) {
						return []string{
							"192.168.0.12",
							"192.168.0.13",
						}, nil
					},
				},
			},
			wantResp: td.Struct(&dns.Response{
				Addresses: []net.IP{
					net.IPv4(192, 168, 0, 12),
					net.IPv4(192, 168, 0, 13),
				},
			}, td.StructFields{}),
			wantErr:       false,
			wantResolvErr: false,
		},
		{
			name: "Missing parameter for AorAAAA initiator",
			args: args{
				args: make([]rules.Param, 0),
			},
			wantErr: true,
		},
		{
			name: "Wrong parameter for AorAAAA initiator",
			args: args{
				args: []rules.Param{
					{
						Int: rules.IntP(42),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Mocked resolver returns error",
			args: args{
				args: []rules.Param{
					{
						String: rules.StringP("gitlab.com"),
					},
				},
				resolver: &mockResolver{
					LookupHostDelegate: func(ctx context.Context, host string) (addrs []string, err error) {
						if host == "gitlab.com" {
							return nil, errors.New("expected error")
						}
						return nil, nil
					},
				},
			},
			wantErr:       false,
			wantResolvErr: true,
			wantResp:      td.Nil(),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var (
				initiator dns.Initiator
				resp      *dns.Response
				err       error
			)
			initiator, err = dns.AorAAAAInitiator(tt.args.args...)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("AorAAAAInitiator() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			var ctx, cancel = context.WithCancel(context.Background())
			t.Cleanup(cancel)
			if resp, err = initiator.Do(ctx, tt.args.resolver); err != nil {
				if !tt.wantResolvErr {
					t.Errorf("initiator.Do() error = %v, wantResolvErr %t", err, tt.wantResolvErr)
				}
			}
			td.Cmp(t, resp, tt.wantResp)
		})
	}
}
