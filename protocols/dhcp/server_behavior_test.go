package dhcp_test

import (
	"net"
	"testing"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/maxatome/go-testdeep/td"

	"gitlab.com/inetmock/inetmock/internal/netutils"
	"gitlab.com/inetmock/inetmock/internal/rules"
	"gitlab.com/inetmock/inetmock/internal/state/statetest"
	"gitlab.com/inetmock/inetmock/internal/test"
	"gitlab.com/inetmock/inetmock/pkg/logging"
	"gitlab.com/inetmock/inetmock/protocols/dhcp"
)

func TestHandlerForRoutingRule(t *testing.T) {
	t.Parallel()
	type args struct {
		opts    dhcp.ProtocolOptions
		rawRule string
		req     *dhcpv4.DHCPv4
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "Static IP handler",
			args: args{
				rawRule: `=> IP(1.3.3.7)`,
				req: &dhcpv4.DHCPv4{
					ClientHWAddr: netutils.MustParseMAC("54:df:83:56:2c:f3"),
				},
			},
			want: td.Struct(new(dhcpv4.DHCPv4), td.StructFields{
				"YourIPAddr": test.IP("1.3.3.7"),
			}),
			wantErr: false,
		},
		{
			name: "Range IP handler",
			args: args{
				rawRule: `=> Range(3.3.6.110, 3.3.6.200)`,
				req: &dhcpv4.DHCPv4{
					ClientHWAddr: netutils.MustParseMAC("54:df:83:56:2c:f3"),
				},
			},
			want: td.Struct(new(dhcpv4.DHCPv4), td.StructFields{
				"YourIPAddr": test.IP("3.3.6.110"),
			}),
			wantErr: false,
		},
		{
			name: "Router option handler",
			args: args{
				rawRule: `=> Router(1.3.3.7)`,
				req: &dhcpv4.DHCPv4{
					ClientHWAddr: netutils.MustParseMAC("54:df:83:56:2c:f3"),
				},
			},
			want:    WantOption(dhcpv4.OptionRouter, net.IPv4(1, 3, 3, 7)),
			wantErr: false,
		},
		{
			name: "Netmask option handler",
			args: args{
				rawRule: `=> Netmask(255.255.255.0)`,
				req: &dhcpv4.DHCPv4{
					ClientHWAddr: netutils.MustParseMAC("54:df:83:56:2c:f3"),
				},
			},
			want:    WantOption(dhcpv4.OptionSubnetMask, net.IPv4(255, 255, 255, 0)),
			wantErr: false,
		},
		{
			name: "Single DNS option handler",
			args: args{
				rawRule: `=> DNS(1.1.1.1)`,
				req: &dhcpv4.DHCPv4{
					ClientHWAddr: netutils.MustParseMAC("54:df:83:56:2c:f3"),
				},
			},
			want:    WantDNS(net.IPv4(1, 1, 1, 1)),
			wantErr: false,
		},
		{
			name: "2 DNS option handler",
			args: args{
				rawRule: `=> DNS(1.1.1.1, 9.9.9.9)`,
				req: &dhcpv4.DHCPv4{
					ClientHWAddr: netutils.MustParseMAC("54:df:83:56:2c:f3"),
				},
			},
			want: WantDNS(
				net.IPv4(1, 1, 1, 1),
				net.IPv4(9, 9, 9, 9),
			),
			wantErr: false,
		},
		{
			name: "3 DNS option handler",
			args: args{
				rawRule: `=> DNS(1.1.1.1, 9.9.9.9, 8.8.8.8)`,
				req: &dhcpv4.DHCPv4{
					ClientHWAddr: netutils.MustParseMAC("54:df:83:56:2c:f3"),
				},
			},
			want: WantDNS(
				net.IPv4(1, 1, 1, 1),
				net.IPv4(9, 9, 9, 9),
				net.IPv4(8, 8, 8, 8),
			),
			wantErr: false,
		},
		{
			name: "Static IP & Router handler",
			args: args{
				rawRule: `=> IP(1.3.3.7) => Router(1.3.3.1)`,
				req: &dhcpv4.DHCPv4{
					ClientHWAddr: netutils.MustParseMAC("54:df:83:56:2c:f3"),
				},
			},
			want: td.All(
				td.Struct(new(dhcpv4.DHCPv4), td.StructFields{
					"YourIPAddr": test.IP("1.3.3.7"),
				}),
				WantOption(dhcpv4.OptionRouter, net.IPv4(1, 3, 3, 1)),
			),
			wantErr: false,
		},
		{
			name: "Static IP & router & netmask handler",
			args: args{
				rawRule: `=> IP(1.3.3.7) => Netmask(255.255.255.0) => Router(1.3.3.1)`,
				req: &dhcpv4.DHCPv4{
					ClientHWAddr: netutils.MustParseMAC("54:df:83:56:2c:f3"),
				},
			},
			want: td.All(
				td.Struct(new(dhcpv4.DHCPv4), td.StructFields{
					"YourIPAddr": test.IP("1.3.3.7"),
				}),
				WantOption(dhcpv4.OptionSubnetMask, net.IPv4(255, 255, 255, 0)),
				WantOption(dhcpv4.OptionRouter, net.IPv4(1, 3, 3, 1)),
			),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			store := statetest.NewTestStore(t)

			opts := dhcp.HandlerOptions{
				Logger:          logging.CreateTestLogger(t),
				StateStore:      store,
				ProtocolOptions: tt.args.opts,
			}

			var (
				rule *rules.ChainedResponsePipeline
				err  error
			)
			if rule, err = rules.Parse[rules.ChainedResponsePipeline](tt.args.rawRule); err != nil {
				t.Errorf("rules.Parse() error = %v", err)
				return
			}

			handlerChain, err := dhcp.HandlerForRoutingRule(rule, opts)
			if err != nil {
				t.Errorf("HandlerForRoutingRule() error = %v", err)
				return
			}

			var resp *dhcpv4.DHCPv4
			if r, err := dhcpv4.NewReplyFromRequest(tt.args.req); err != nil {
				t.Errorf("dhcpv4.NewReplyFromRequest() error = %v", err)
				return
			} else {
				resp = r
			}

			if err := handlerChain.Apply(tt.args.req, resp); err != nil {
				if !tt.wantErr {
					t.Errorf("handlerChain.Apply() error = %v", err)
				}
				return
			}

			td.Cmp(t, resp, tt.want)
		})
	}
}
