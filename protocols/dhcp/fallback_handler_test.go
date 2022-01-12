package dhcp_test

import (
	"fmt"
	"net"
	"testing"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/maxatome/go-testdeep/td"

	"gitlab.com/inetmock/inetmock/pkg/logging"
	"gitlab.com/inetmock/inetmock/protocols/dhcp"
)

func TestFallbackHandler_Handle(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		DefaultOptions dhcp.DefaultOptions
		want           interface{}
		wantErr        bool
	}{
		{
			name: "Set ServerID if missing",
			DefaultOptions: dhcp.DefaultOptions{
				ServerID: net.IPv4(1, 2, 3, 4),
			},
			want: td.Code(func(resp *dhcpv4.DHCPv4) error {
				serverIP := net.IPv4(1, 2, 3, 4)
				if !resp.ServerIPAddr.Equal(serverIP) {
					return fmt.Errorf("ServerIP %v does not match", resp.ServerIPAddr)
				}

				if optVal := resp.Options.Get(dhcpv4.OptionServerIdentifier); !serverIP.Equal(optVal) {
					return fmt.Errorf("ServerIP %v does not match", optVal)
				}

				return nil
			}),
			wantErr: false,
		},
		{
			name: "Set router if missing",
			DefaultOptions: dhcp.DefaultOptions{
				Router: net.IPv4(1, 2, 3, 4),
			},
			want:    WantOption(dhcpv4.OptionRouter, net.IPv4(1, 2, 3, 4)),
			wantErr: false,
		},
		{
			name: "Set netmask if missing",
			DefaultOptions: dhcp.DefaultOptions{
				Netmask: net.IPv4(255, 255, 255, 255),
			},
			want:    WantOption(dhcpv4.OptionSubnetMask, net.IPv4(255, 255, 255, 255)),
			wantErr: false,
		},
		{
			name: "Set DNS if missing - single",
			DefaultOptions: dhcp.DefaultOptions{
				DNS: []net.IP{
					net.IPv4(1, 1, 1, 1),
				},
			},
			want:    WantDNS(net.IPv4(1, 1, 1, 1)),
			wantErr: false,
		},
		{
			name: "Set DNS if missing - multiple",
			DefaultOptions: dhcp.DefaultOptions{
				DNS: []net.IP{
					net.IPv4(1, 1, 1, 1),
					net.IPv4(9, 9, 9, 9),
				},
			},
			want:    WantDNS(net.IPv4(1, 1, 1, 1), net.IPv4(9, 9, 9, 9)),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := &dhcp.FallbackHandler{
				Previous:       dhcp.NoOpHandler,
				Logger:         logging.CreateTestLogger(t),
				DefaultOptions: tt.DefaultOptions,
			}
			var (
				req = &dhcpv4.DHCPv4{
					OpCode: dhcpv4.OpcodeBootRequest,
				}
				resp *dhcpv4.DHCPv4
			)

			if r, err := dhcpv4.NewReplyFromRequest(req); err != nil {
				t.Errorf("dhcpv4.NewReplyFromRequest() error = %v", err)
				return
			} else {
				resp = r
			}

			if err := h.Handle(req, resp); err != nil {
				if !tt.wantErr {
					t.Errorf("Handle() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			td.Cmp(t, resp, tt.want)
		})
	}
}

func WantDNS(expected ...net.IP) td.TestDeep {
	expectedMap := make(map[string]bool, len(expected))
	for idx := range expected {
		expectedMap[expected[idx].String()] = true
	}
	return td.Code(func(msg *dhcpv4.DHCPv4) error {
		dns := msg.DNS()
		if len(dns) != len(expectedMap) {
			return fmt.Errorf("expected %d servers but are %d", len(expectedMap), len(dns))
		}
		for idx := range dns {
			if _, ok := expectedMap[dns[idx].String()]; !ok {
				return fmt.Errorf("DNS server %v missing", dns[idx].String())
			}
		}
		return nil
	})
}

func WantOption(code dhcpv4.OptionCode, want net.IP) td.TestDeep {
	return td.Code(func(msg *dhcpv4.DHCPv4) error {
		if optVal := msg.Options.Get(code); !want.Equal(optVal) {
			return fmt.Errorf("option %v want %v got %v", code, want, optVal)
		}
		return nil
	})
}
