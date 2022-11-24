package netflow_test

import (
	"testing"

	"github.com/maxatome/go-testdeep/td"
	"github.com/spf13/viper"

	"inetmock.icb4dc0.de/inetmock/netflow"
)

func TestPacketPolicyDecodeHook(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		config string
		want   netflow.PacketPolicy
	}{
		{
			name: "Test drop",
			config: `
Policy: drop
`,
			want: netflow.PacketPolicyDrop,
		},
		{
			name: "Test drop - case insensitive",
			config: `
Policy: Drop
`,
			want: netflow.PacketPolicyDrop,
		},
		{
			name: "Test pass",
			config: `
Policy: pass
`,
			want: netflow.PacketPolicyPass,
		},
		{
			name: "Test pass - case insensitive",
			config: `
Policy: PaSs
`,
			want: netflow.PacketPolicyPass,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cfg := PrepareViper(t, tt.config)
			hook := netflow.PacketPolicyDecodeHook()

			tmp := struct {
				Policy netflow.PacketPolicy
			}{}

			if err := cfg.Unmarshal(&tmp, viper.DecodeHook(hook)); err != nil {
				t.Errorf("Failed to unmarshal: %v", err)
				return
			}

			td.Cmp(t, tmp.Policy, tt.want)
		})
	}
}
