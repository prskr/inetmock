package dhcp_test

import (
	"encoding/base64"
	"hash/fnv"
	"net"
	"path"
	"testing"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/maxatome/go-testdeep/td"

	"gitlab.com/inetmock/inetmock/internal/netutils"
	"gitlab.com/inetmock/inetmock/internal/state"
	"gitlab.com/inetmock/inetmock/internal/state/statetest"
	"gitlab.com/inetmock/inetmock/internal/test"
	"gitlab.com/inetmock/inetmock/protocols/dhcp"
)

const defaultTTL = 5 * time.Minute

var (
	defaultStartIP = net.IPv4(172, 20, 10, 100)
	defaultEndIP   = net.IPv4(172, 20, 10, 150)
)

func TestRangeMessageHandler_Handle(t *testing.T) {
	t.Parallel()
	type fields struct {
		TTL     time.Duration
		StartIP net.IP
		EndIP   net.IP
	}
	tests := []struct {
		name       string
		fields     fields
		req        *dhcpv4.DHCPv4
		storeSetup statetest.StoreSetup
		want       interface{}
		wantErr    bool
	}{
		{
			name: "No lease yet",
			fields: fields{
				TTL:     defaultTTL,
				StartIP: defaultStartIP,
				EndIP:   defaultEndIP,
			},
			req: &dhcpv4.DHCPv4{
				ClientHWAddr: netutils.MustParseMAC("54:df:83:56:2c:f3"),
			},
			want: td.Struct(new(dhcpv4.DHCPv4), td.StructFields{
				"YourIPAddr": test.IP("172.20.10.100"),
			}),
			wantErr: false,
		},
		{
			name: "Single lease present",
			fields: fields{
				TTL:     defaultTTL,
				StartIP: defaultStartIP,
				EndIP:   defaultEndIP,
			},
			storeSetup: statetest.StoreSetupFunc(func(_ testing.TB, store state.KVStore) error {
				rangeKey := calcRangeKey(defaultStartIP, defaultEndIP)
				return store.Set(path.Join("range", rangeKey, "172.20.10.100"), &dhcp.RangeLease{
					IP:  defaultStartIP,
					MAC: netutils.MustParseMAC("54:df:83:56:2d:f4"),
				})
			}),
			req: &dhcpv4.DHCPv4{
				ClientHWAddr: netutils.MustParseMAC("54:df:83:56:2c:f3"),
			},
			want: td.Struct(new(dhcpv4.DHCPv4), td.StructFields{
				"YourIPAddr": test.IP("172.20.10.101"),
			}),
			wantErr: false,
		},
		{
			name: "Lease for MAC already present",
			fields: fields{
				TTL:     defaultTTL,
				StartIP: defaultStartIP,
				EndIP:   defaultEndIP,
			},
			storeSetup: statetest.StoreSetupFunc(func(_ testing.TB, store state.KVStore) error {
				rangeKey := calcRangeKey(defaultStartIP, defaultEndIP)
				return store.Set(path.Join("range", rangeKey, "54:df:83:56:2c:f3"), &dhcp.RangeLease{
					IP:  net.IPv4(172, 20, 10, 115),
					MAC: netutils.MustParseMAC("54:df:83:56:2c:f3"),
				})
			}),
			req: &dhcpv4.DHCPv4{
				ClientHWAddr: netutils.MustParseMAC("54:df:83:56:2c:f3"),
			},
			want: td.Struct(new(dhcpv4.DHCPv4), td.StructFields{
				"YourIPAddr": test.IP("172.20.10.115"),
			}),
			wantErr: false,
		},
		{
			name: "Range full",
			fields: fields{
				TTL:     defaultTTL,
				StartIP: defaultStartIP,
				EndIP:   net.IPv4(172, 20, 10, 101),
			},
			storeSetup: statetest.StoreSetupFunc(func(_ testing.TB, store state.KVStore) error {
				rangeKey := calcRangeKey(defaultStartIP, net.IPv4(172, 20, 10, 101))
				return store.Set(path.Join("range", rangeKey, "172.20.10.100"), &dhcp.RangeLease{
					IP:  defaultStartIP,
					MAC: netutils.MustParseMAC("54:df:83:56:2d:f4"),
				})
			}),
			req: &dhcpv4.DHCPv4{
				ClientHWAddr: netutils.MustParseMAC("54:df:83:56:2c:f3"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			store := statetest.NewTestStore(t, tt.storeSetup)

			var resp *dhcpv4.DHCPv4
			if r, err := dhcpv4.NewReplyFromRequest(tt.req); err != nil {
				t.Errorf("dhcpv4.NewReplyFromRequest() error = %v", err)
				return
			} else {
				resp = r
			}

			h := &dhcp.RangeMessageHandler{
				Store:   store,
				TTL:     tt.fields.TTL,
				StartIP: tt.fields.StartIP,
				EndIP:   tt.fields.EndIP,
			}
			if err := h.Handle(tt.req, resp); err != nil {
				if !tt.wantErr {
					t.Errorf("Handle() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			td.Cmp(t, resp, tt.want)
		})
	}
}

func calcRangeKey(startIP, endIP net.IP) string {
	hash := fnv.New32a()
	return base64.URLEncoding.EncodeToString(hash.Sum(append(startIP, endIP...)))
}
