package mock_test

import (
	"net"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/maxatome/go-testdeep/td"
	mdns "github.com/miekg/dns"

	"gitlab.com/inetmock/inetmock/internal/endpoint/handler/dns"
	"gitlab.com/inetmock/inetmock/internal/endpoint/handler/dns/mock"
	auditmock "gitlab.com/inetmock/inetmock/internal/mock/audit"
	dnsmock "gitlab.com/inetmock/inetmock/internal/mock/dns"
	"gitlab.com/inetmock/inetmock/internal/test"
	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/logging"
)

var recordsMap = map[string]net.IP{
	"google.com.": net.IPv4(142, 250, 185, 174).To4(),
}

func TestRuleHandler_RegisterRule(t *testing.T) {
	t.Parallel()
	type args struct {
		rawRule string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Add rule without filters",
			args: args{
				rawRule: `=> IP(1.1.1.1)`,
			},
			wantErr: false,
		},
		{
			name: "Add rule with A filters",
			args: args{
				rawRule: `A(".*google\\.com$") => IP(1.1.1.1)`,
			},
			wantErr: false,
		},
		{
			name: "Add rule with AAAA filters",
			args: args{
				rawRule: `AAAA(".*google\\.com$") => IP(1.1.1.1)`,
			},
			wantErr: false,
		},
		{
			name: "Add invalid rule",
			args: args{
				rawRule: `=> IP(1.1.1.)`,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := &mock.RuleHandler{
				Logger: logging.CreateTestLogger(t),
			}
			if err := r.RegisterRule(tt.args.rawRule); (err != nil) != tt.wantErr {
				t.Errorf("RegisterRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRuleHandler_ServeDNS(t *testing.T) {
	const defaultTTL = 10 * time.Second

	t.Parallel()
	type fields struct {
		Fallback dns.IPResolver
		Cache    mock.Cache
	}
	type args struct {
		req *mdns.Msg
	}
	tests := []struct {
		name          string
		fields        fields
		rules         []string
		expectedEvent interface{}
		expectedMsg   interface{}
		args          args
	}{
		{
			name: "Resolve anything to ",
			fields: fields{
				Cache: new(mock.DelegateCache),
			},
			args: args{
				req: new(mdns.Msg).SetQuestion("google.com.", mdns.TypeA),
			},
			rules: []string{
				`=> IP(1.1.1.1)`,
			},
			expectedEvent: td.NotNil(),
			expectedMsg: td.Struct(&mdns.Msg{
				Answer: []mdns.RR{
					&mdns.A{
						Hdr: mdns.RR_Header{
							Class:  mdns.ClassINET,
							Ttl:    uint32(defaultTTL.Seconds()),
							Name:   "google.com.",
							Rrtype: mdns.TypeA,
						},
						A: net.IPv4(1, 1, 1, 1).To4(),
					},
				},
			}, td.StructFields{}),
		},
		{
			name: "Resolve with fallback handler",
			fields: fields{
				Fallback: dns.IPResolverFunc(func(string) net.IP {
					return net.IPv4(192, 168, 0, 1)
				}),
				Cache: new(mock.DelegateCache),
			},
			args: args{
				req: new(mdns.Msg).SetQuestion("google.com.", mdns.TypeA),
			},
			expectedEvent: td.NotNil(),
			expectedMsg: td.Struct(&mdns.Msg{
				Answer: []mdns.RR{
					&mdns.A{
						Hdr: mdns.RR_Header{
							Class:  mdns.ClassINET,
							Ttl:    uint32(defaultTTL.Seconds()),
							Name:   "google.com.",
							Rrtype: mdns.TypeA,
						},
						A: net.IPv4(192, 168, 0, 1).To4(),
					},
				},
			}, td.StructFields{}),
		},
		{
			name: "Resolve A request with Cache",
			fields: fields{
				Cache: &mock.DelegateCache{
					OnForwardLookup: func(host string) net.IP {
						if ip, ok := recordsMap[host]; ok {
							return ip
						}
						return nil
					},
				},
			},
			args: args{
				req: new(mdns.Msg).SetQuestion("google.com.", mdns.TypeA),
			},
			expectedEvent: td.NotNil(),
			expectedMsg: td.Struct(&mdns.Msg{
				Answer: []mdns.RR{
					&mdns.A{
						Hdr: mdns.RR_Header{
							Class:  mdns.ClassINET,
							Ttl:    uint32(defaultTTL.Seconds()),
							Name:   "google.com.",
							Rrtype: mdns.TypeA,
						},
						A: recordsMap["google.com."],
					},
				},
			}, td.StructFields{}),
		},
		{
			name: "Resolve PTR request with Cache",
			fields: fields{
				Cache: &mock.DelegateCache{
					OnReverseLookup: func(address net.IP) (host string, miss bool) {
						for h, i := range recordsMap {
							if address.Equal(i) {
								return h, false
							}
						}
						return "", true
					},
				},
			},
			args: args{
				req: new(mdns.Msg).SetQuestion("174.185.250.142.in-addr.arpa.", mdns.TypePTR),
			},
			expectedEvent: td.NotNil(),
			expectedMsg: td.Struct(&mdns.Msg{
				Answer: []mdns.RR{
					&mdns.PTR{
						Hdr: mdns.RR_Header{
							Class:  mdns.ClassINET,
							Ttl:    uint32(defaultTTL.Seconds()),
							Name:   "174.185.250.142.in-addr.arpa.",
							Rrtype: mdns.TypePTR,
						},
						Ptr: "google.com.",
					},
				},
			}, td.StructFields{}),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &mock.RuleHandler{
				TTL:         defaultTTL,
				Logger:      logging.CreateTestLogger(t),
				Cache:       tt.fields.Cache,
				HandlerName: t.Name(),
				Emitter:     mockEmitter(t, tt.expectedEvent),
				Fallback:    tt.fields.Fallback,
			}

			for idx := range tt.rules {
				rawRule := tt.rules[idx]
				if err := r.RegisterRule(rawRule); err != nil {
					t.Errorf("RegisterRule(%s) err = %v", rawRule, err)
				}
			}

			r.ServeDNS(mockResponseWriter(t, tt.expectedMsg), tt.args.req)
		})
	}
}

func mockEmitter(tb testing.TB, expectedEvent interface{}) audit.Emitter {
	tb.Helper()
	ctrl := gomock.NewController(tb)

	emitter := auditmock.NewMockEmitter(ctrl)
	emitter.
		EXPECT().
		Emit(test.GenericMatcher(tb, expectedEvent)).
		Times(1)

	return emitter
}

func mockResponseWriter(tb testing.TB, expectedMsg interface{}) mdns.ResponseWriter {
	tb.Helper()
	ctrl := gomock.NewController(tb)
	resolver := mock.NewRandomIPResolver(mustParseCIDR("10.10.0.0/8"))

	w := dnsmock.NewMockResponseWriter(ctrl)
	w.EXPECT().
		WriteMsg(test.GenericMatcher(tb, expectedMsg)).
		Times(1)

	w.EXPECT().
		LocalAddr().
		Return(&net.UDPAddr{IP: resolver.Lookup("")})

	w.EXPECT().
		RemoteAddr().
		Return(&net.UDPAddr{IP: resolver.Lookup("")})

	return w
}
