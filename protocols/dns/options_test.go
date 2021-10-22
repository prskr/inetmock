package dns_test

import (
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/td"

	"gitlab.com/inetmock/inetmock/internal/endpoint"
	dnsmock "gitlab.com/inetmock/inetmock/internal/mock/dns"
	"gitlab.com/inetmock/inetmock/protocols/dns"
)

func TestOptionsFromLifecycle(t *testing.T) {
	t.Parallel()
	type args struct {
		opts map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "None cache config",
			args: args{
				opts: map[string]interface{}{
					"cache": map[string]interface{}{
						"type": "none",
					},
					"default": map[string]interface{}{
						"type": "random",
						"cidr": "192.168.0.1/24",
					},
				},
			},
			want: td.Struct(new(dns.Options), td.StructFields{
				"Cache": td.Isa(new(dnsmock.CacheMock)),
			}),
		},
		{
			name: "TTL cache config",
			args: args{
				opts: map[string]interface{}{
					"cache": map[string]interface{}{
						"type": "ttl",
						"ttl":  30 * time.Second,
					},
					"default": map[string]interface{}{
						"type": "random",
						"cidr": "192.168.0.1/24",
					},
				},
			},
			want: td.Struct(new(dns.Options), td.StructFields{
				"Cache": td.Isa(new(dns.Cache)),
			}),
		},
		{
			name: "Random IP resolver",
			args: args{
				opts: map[string]interface{}{
					"cache": map[string]interface{}{
						"type": "none",
					},
					"default": map[string]interface{}{
						"type": "random",
						"cidr": "192.168.0.1/24",
					},
				},
			},
			want: td.Struct(new(dns.Options), td.StructFields{
				"Default": td.Struct(new(dns.RandomIPResolver), td.StructFields{
					"CIDR":   td.NotNil(),
					"Random": td.NotNil(),
				}),
			}),
		},
		{
			name: "Incremental IP resolver",
			args: args{
				opts: map[string]interface{}{
					"cache": map[string]interface{}{
						"type": "none",
					},
					"default": map[string]interface{}{
						"type": "incremental",
						"cidr": "192.168.0.1/24",
					},
				},
			},
			want: td.Struct(new(dns.Options), td.StructFields{
				"Default": td.Struct(new(dns.IncrementalIPResolver), td.StructFields{
					"CIDR": td.NotNil(),
				}),
			}),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			lifecycle := endpoint.NewEndpointLifecycle("", endpoint.Uplink{}, tt.args.opts)
			got, err := dns.OptionsFromLifecycle(lifecycle)
			if (err != nil) != tt.wantErr {
				t.Errorf("OptionsFromLifecycle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			td.Cmp(t, got, tt.want)
		})
	}
}
