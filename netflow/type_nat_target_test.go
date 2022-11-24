package netflow_test

import (
	"strings"
	"testing"

	"github.com/maxatome/go-testdeep/td"
	"github.com/spf13/viper"

	"inetmock.icb4dc0.de/inetmock/netflow"
)

func TestNATTargetDecodingHook(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		config string
		want   netflow.NATTarget
	}{
		{
			name: "Empty value",
			config: `
Target: ""
`,
			want: netflow.NATTargetInterface,
		},
		{
			name: "NAT target interface",
			config: `
Target: interface
`,
			want: netflow.NATTargetInterface,
		},
		{
			name: "NAT target IP",
			config: `
Target: ip
`,
			want: netflow.NATTargetIP,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cfg := PrepareViper(t, tt.config)
			hook := netflow.NATTargetDecodingHook()

			tmp := struct {
				Target netflow.NATTarget
			}{}

			if err := cfg.Unmarshal(&tmp, viper.DecodeHook(hook)); err != nil {
				t.Errorf("Failed to unmarshal: %v", err)
				return
			}

			td.Cmp(t, tmp.Target, tt.want)
		})
	}
}

func PrepareViper(tb testing.TB, configText string) *viper.Viper {
	tb.Helper()

	v := viper.New()
	v.SetConfigType("yaml")
	if err := v.ReadConfig(strings.NewReader(configText)); err != nil {
		tb.Fatalf("Failed to read in config: %s err = %v", configText, err)
	}

	return v
}
