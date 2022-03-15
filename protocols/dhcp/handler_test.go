//go:build sudo
// +build sudo

package dhcp_test

import (
	"context"
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/nclient4"
	"github.com/maxatome/go-testdeep/td"

	"gitlab.com/inetmock/inetmock/internal/endpoint"
	audit_mock "gitlab.com/inetmock/inetmock/internal/mock/audit"
	"gitlab.com/inetmock/inetmock/internal/state/statetest"
	"gitlab.com/inetmock/inetmock/internal/test"
	"gitlab.com/inetmock/inetmock/pkg/logging"
	"gitlab.com/inetmock/inetmock/protocols/dhcp"
)

func Test_dhcpHandler_Start(t *testing.T) {
	t.Parallel()
	anyInterface, ifAddr := anyPhysicalInterface(t)
	type args struct {
		opts map[string]any
		mac  string
	}
	tests := []struct {
		name      string
		args      args
		wantEvent any
		want      any
		wantErr   bool
	}{
		{
			name: "Exact match test",
			args: args{
				mac: "d1:15:b8:0c:0c:9a",
				opts: map[string]any{
					"default": map[string]any{
						"serverID": "1.2.3.4",
						"netmask":  "255.255.255.0",
						"dns": []string{
							"1.2.3.4",
						},
						"router": "1.2.3.4",
					},
					"rules": []string{
						fmt.Sprintf(`ExactMAC(%q) => IP(1.3.3.7)`, anyInterface.HardwareAddr.String()),
					},
				},
			},
			wantEvent: td.NotNil(),
			want: td.Struct(new(nclient4.Lease), td.StructFields{
				"Offer": td.Struct(new(dhcpv4.DHCPv4), td.StructFields{
					"YourIPAddr": test.IP("1.3.3.7"),
				}),
				"ACK": td.NotNil(),
			}),
			wantErr: false,
		},
		{
			name: "MatchMAC test",
			args: args{
				mac: "db:2d:f0:0f:6e:aa",
				opts: map[string]any{
					"default": map[string]any{
						"serverID": "1.2.3.4",
						"netmask":  "255.255.255.0",
						"dns": []string{
							"1.2.3.4",
						},
						"router": "1.2.3.4",
					},
					"rules": []string{
						fmt.Sprintf(`MatchMAC("%s.*") => IP(1.2.3.4) => Netmask(255.255.255.0)`, anyInterface.HardwareAddr.String()[:5]),
					},
				},
			},
			wantEvent: td.NotNil(),
			want: td.Struct(new(nclient4.Lease), td.StructFields{
				"Offer": td.Struct(new(dhcpv4.DHCPv4), td.StructFields{
					"YourIPAddr": test.IP("1.2.3.4"),
					"Options": td.Code(func(opts dhcpv4.Options) error {
						if mask := net.IPMask(opts.Get(dhcpv4.OptionSubnetMask)); mask == nil {
							return errors.New("missing subnet mask option")
						} else if ones, _ := mask.Size(); ones != 24 {
							return errors.New("subnet mask has not expected value of ones")
						}
						return nil
					}),
				}),
				"ACK": td.NotNil(),
			}),
			wantErr: false,
		},
		{
			name: "Range test",
			args: args{
				mac: "db:2d:f0:0f:6e:aa",
				opts: map[string]any{
					"default": map[string]any{
						"serverID": "1.2.4.5",
						"netmask":  "255.255.255.0",
						"dns": []string{
							"1.2.4.5",
						},
						"router": "1.2.4.5",
					},
					"rules": []string{
						`=> Range(1.2.4.100, 1.2.4.150) => Netmask(255.255.255.0)`,
					},
				},
			},
			wantEvent: td.NotNil(),
			want: td.Struct(new(nclient4.Lease), td.StructFields{
				"Offer": td.Struct(new(dhcpv4.DHCPv4), td.StructFields{
					"YourIPAddr": test.IP("1.2.4.100"),
					"Options": td.Code(func(opts dhcpv4.Options) error {
						if mask := net.IPMask(opts.Get(dhcpv4.OptionSubnetMask)); mask == nil {
							return errors.New("missing subnet mask option")
						} else if ones, _ := mask.Size(); ones != 24 {
							return errors.New("subnet mask has not expected value of ones")
						}
						return nil
					}),
				}),
				"ACK": td.NotNil(),
			}),
			wantErr: false,
		},
		{
			name: "Range fallback test",
			args: args{
				mac: "db:2d:f0:0f:6e:aa",
				opts: map[string]any{
					"default": map[string]any{
						"serverID": "1.2.4.5",
						"netmask":  "255.255.255.0",
						"dns": []string{
							"1.2.4.5",
						},
						"router": "1.2.4.5",
					},
					"fallback": map[string]any{
						"type":    "range",
						"ttl":     "1h",
						"startIP": "1.2.4.100",
						"endIP":   "1.2.4.150",
					},
					"rules": []string{},
				},
			},
			wantEvent: td.NotNil(),
			want: td.Struct(new(nclient4.Lease), td.StructFields{
				"Offer": td.Struct(new(dhcpv4.DHCPv4), td.StructFields{
					"YourIPAddr": test.IP("1.2.4.100"),
					"Options": td.Code(func(opts dhcpv4.Options) error {
						if mask := net.IPMask(opts.Get(dhcpv4.OptionSubnetMask)); mask == nil {
							return errors.New("missing subnet mask option")
						} else if ones, _ := mask.Size(); ones != 24 {
							return errors.New("subnet mask has not expected value of ones")
						}
						return nil
					}),
				}),
				"ACK": td.NotNil(),
			}),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			ctx, cancel := context.WithCancel(test.Context(t))
			t.Cleanup(cancel)
			logger := logging.CreateTestLogger(t)

			emitterMock := audit_mock.NewMockEmitter(ctrl)
			if !tt.wantErr {
				emitterMock.EXPECT().Emit(test.GenericMatcher(t, tt.wantEvent)).MinTimes(1)
			}

			srvAddr := randomUDPAddr(t)
			lifecycle := endpoint.NewStartupSpec(t.Name(), endpoint.NewUplink(srvAddr), tt.args.opts)
			handler := dhcp.New(logger, emitterMock, statetest.NewTestStore(t))
			if err := handler.Start(ctx, lifecycle); err != nil {
				if !tt.wantErr {
					t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
				}
			}

			client := setupClient(t, srvAddr, ifAddr, anyInterface.Name)
			var got *nclient4.Lease
			if l, err := client.Request(ctx); (err != nil) != tt.wantErr {
				t.Errorf("client.Request() error = %v", err)
				return
			} else {
				got = l
			}

			td.Cmp(t, got, tt.want)
		})
	}
}

func setupClient(tb testing.TB, srvAddr *net.UDPAddr, ifAddr net.Addr, ifName string) *nclient4.Client {
	tb.Helper()
	dialAddr := &net.UDPAddr{
		Port: srvAddr.Port,
		IP:   ifAddr.(*net.IPNet).IP,
	}
	if cli, err := nclient4.New(
		ifName,
		nclient4.WithServerAddr(dialAddr),
		nclient4.WithUnicast(randomUDPAddr(tb)),
	); err != nil {
		tb.Fatalf("nclient4.NewWithConn() error = %v", err)
	} else {
		return cli
	}
	return nil
}

func randomUDPAddr(tb testing.TB) *net.UDPAddr {
	tb.Helper()
	var conn *net.UDPConn
	if c, err := net.ListenUDP("udp4", &net.UDPAddr{}); err != nil {
		tb.Fatalf("net.ListenUDP() error = %v", err)
	} else {
		conn = c
	}

	defer conn.Close()
	switch la := conn.LocalAddr(); addr := la.(type) {
	case *net.UDPAddr:
		return addr
	default:
		tb.Fatalf("Address %v is not an UDP address", la)
	}
	return nil
}

func anyPhysicalInterface(tb testing.TB) (net.Interface, net.Addr) {
	tb.Helper()
	var interfaces []net.Interface
	if ifs, err := net.Interfaces(); err != nil {
		tb.Fatalf("net.Interfaces() error = %v", err)
	} else {
		interfaces = ifs
	}

	for idx := range interfaces {
		cif := interfaces[idx]
		if len(cif.HardwareAddr.String()) == 0 {
			continue
		}
		if addr, err := cif.Addrs(); err != nil {
			continue
		} else if len(addr) == 0 {
			continue
		} else {
			return cif, addr[0]
		}
	}
	tb.Fatal("No physical interface found")
	return net.Interface{}, nil
}
