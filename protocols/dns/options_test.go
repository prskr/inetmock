package dns_test

import (
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/td"

	"inetmock.icb4dc0.de/inetmock/internal/endpoint"
	dnsmock "inetmock.icb4dc0.de/inetmock/internal/mock/dns"
	"inetmock.icb4dc0.de/inetmock/protocols/dns"
)

func TestOptionsFromLifecycle(t *testing.T) {
	t.Parallel()
	type args struct {
		opts map[string]any
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		{
			name: "None cache config",
			args: args{
				opts: map[string]any{
					"cache": map[string]any{
						"type": "none",
					},
					"default": map[string]any{
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
				opts: map[string]any{
					"cache": map[string]any{
						"type": "inMemory",
						"ttl":  30 * time.Second,
					},
					"default": map[string]any{
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
				opts: map[string]any{
					"cache": map[string]any{
						"type": "none",
					},
					"default": map[string]any{
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
				opts: map[string]any{
					"cache": map[string]any{
						"type": "none",
					},
					"default": map[string]any{
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
			lifecycle := endpoint.NewStartupSpec("", endpoint.NewUplink(nil), tt.args.opts)
			got, err := dns.OptionsFromLifecycle(lifecycle)
			if (err != nil) != tt.wantErr {
				t.Errorf("OptionsFromLifecycle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			td.Cmp(t, got, tt.want)
		})
	}
}
