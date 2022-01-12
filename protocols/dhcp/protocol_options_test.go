package dhcp_test

import (
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/td"

	"gitlab.com/inetmock/inetmock/internal/endpoint"
	"gitlab.com/inetmock/inetmock/internal/test"
	"gitlab.com/inetmock/inetmock/protocols/dhcp"
)

func TestLoadFromConfig(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		args    map[string]interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name:    "Empty args",
			args:    make(map[string]interface{}),
			want:    td.Struct(dhcp.ProtocolOptions{}, td.StructFields{}),
			wantErr: false,
		},
		{
			name: "Single rule",
			args: map[string]interface{}{
				"rules": []string{"some rule"},
			},
			want: td.Struct(dhcp.ProtocolOptions{
				Rules: []string{"some rule"},
			}, td.StructFields{}),
			wantErr: false,
		},
		{
			name: "Multiple rules",
			args: map[string]interface{}{
				"rules": []string{
					"some rule",
					"some other rule",
				},
			},
			want: td.Struct(dhcp.ProtocolOptions{
				Rules: []string{
					"some rule",
					"some other rule",
				},
			}, td.StructFields{}),
			wantErr: false,
		},
		{
			name: "Default ServerID",
			args: map[string]interface{}{
				"default": map[string]interface{}{
					"serverID": "1.2.3.4",
				},
			},
			want: td.Struct(dhcp.ProtocolOptions{}, td.StructFields{
				"Default": td.Struct(dhcp.DefaultOptions{}, td.StructFields{
					"ServerID": test.IP("1.2.3.4"),
				}),
			}),
			wantErr: false,
		},
		{
			name: "Default single DNS",
			args: map[string]interface{}{
				"default": map[string]interface{}{
					"dns": []string{
						"1.2.3.4",
					},
				},
			},
			want: td.Struct(dhcp.ProtocolOptions{}, td.StructFields{
				"Default": td.Struct(dhcp.DefaultOptions{}, td.StructFields{
					"DNS": td.Bag(test.IP("1.2.3.4")),
				}),
			}),
			wantErr: false,
		},
		{
			name: "Default multiple DNS",
			args: map[string]interface{}{
				"default": map[string]interface{}{
					"dns": []string{
						"1.2.3.4",
						"1.2.3.5",
					},
				},
			},
			want: td.Struct(dhcp.ProtocolOptions{}, td.StructFields{
				"Default": td.Struct(dhcp.DefaultOptions{}, td.StructFields{
					"DNS": td.Bag(
						test.IP("1.2.3.4"),
						test.IP("1.2.3.5"),
					),
				}),
			}),
			wantErr: false,
		},
		{
			name: "Default netmask",
			args: map[string]interface{}{
				"default": map[string]interface{}{
					"netmask": "255.255.255.0",
				},
			},
			want: td.Struct(dhcp.ProtocolOptions{}, td.StructFields{
				"Default": td.Struct(dhcp.DefaultOptions{}, td.StructFields{
					"Netmask": test.IP("255.255.255.0"),
				}),
			}),
			wantErr: false,
		},
		{
			name: "Default lease time",
			args: map[string]interface{}{
				"default": map[string]interface{}{
					"leaseTime": "1h",
				},
			},
			want: td.Struct(dhcp.ProtocolOptions{}, td.StructFields{
				"Default": td.Struct(dhcp.DefaultOptions{}, td.StructFields{
					"LeaseTime": 1 * time.Hour,
				}),
			}),
			wantErr: false,
		},
		{
			name: "Range fallback handler",
			args: map[string]interface{}{
				"fallback": map[string]interface{}{
					"type":    "range",
					"ttl":     "1h",
					"startIP": "172.20.0.100",
					"endIP":   "172.20.0.150",
				},
			},
			want: td.Struct(dhcp.ProtocolOptions{}, td.StructFields{
				"Fallback": td.Struct(&dhcp.RangeMessageHandler{
					TTL: 1 * time.Hour,
				}, td.StructFields{
					"StartIP": test.IP("172.20.0.100"),
					"EndIP":   test.IP("172.20.0.150"),
				}),
			}),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			lifecycle := endpoint.NewEndpointLifecycle(tt.name, endpoint.NewUplink(nil), tt.args)
			gotOpts, err := dhcp.LoadFromConfig(lifecycle, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadFromConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			td.Cmp(t, gotOpts, tt.want)
		})
	}
}
