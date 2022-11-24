package dhcp_test

import (
	"testing"

	"github.com/insomniacslk/dhcp/dhcpv4"

	"inetmock.icb4dc0.de/inetmock/internal/netutils"
	"inetmock.icb4dc0.de/inetmock/internal/rules"
	"inetmock.icb4dc0.de/inetmock/protocols/dhcp"
)

func TestRequestFiltersForRoutingRule(t *testing.T) {
	t.Parallel()
	type args struct {
		rule string
		msg  *dhcpv4.DHCPv4
	}
	tests := []struct {
		name      string
		args      args
		wantMatch bool
	}{
		{
			name: "ExactMAC rule - match",
			args: args{
				rule: `ExactMAC("54:df:83:56:2c:f3") => IP(1.3.3.7)`,
				msg: &dhcpv4.DHCPv4{
					ClientHWAddr: netutils.MustParseMAC("54:df:83:56:2c:f3"),
				},
			},
			wantMatch: true,
		},
		{
			name: "ExactMAC rule - no match",
			args: args{
				rule: `ExactMAC("54:df:83:56:2c:f3") => IP(1.3.3.7)`,
				msg: &dhcpv4.DHCPv4{
					ClientHWAddr: netutils.MustParseMAC("54:df:83:56:2c:f4"),
				},
			},
			wantMatch: false,
		},
		{
			name: "MatchMAC rule - match",
			args: args{
				rule: `MatchMAC("(?i)00:06:7C:.*") => Range(3.3.6.110, 3.3.6.200)`,
				msg: &dhcpv4.DHCPv4{
					ClientHWAddr: netutils.MustParseMAC("00:06:7C:56:2c:f4"),
				},
			},
			wantMatch: true,
		},
		{
			name: "MatchMAC rule - no match",
			args: args{
				rule: `MatchMAC("00:06:7C:.*") => Range(3.3.6.110, 3.3.6.200)`,
				msg: &dhcpv4.DHCPv4{
					ClientHWAddr: netutils.MustParseMAC("00:06:8C:56:2c:f4"),
				},
			},
			wantMatch: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var (
				chain *rules.ChainedResponsePipeline
				err   error
			)
			if chain, err = rules.Parse[rules.ChainedResponsePipeline](tt.args.rule); err != nil {
				t.Errorf("rules.Parse() error = %v", err)
				return
			}
			gotFilters, err := dhcp.RequestFiltersForRoutingRule(chain)
			if err != nil {
				t.Errorf("RequestFiltersForRoutingRule() error = %v", err)
				return
			}

			if got := gotFilters.Matches(tt.args.msg); got != tt.wantMatch {
				t.Errorf("gotFilters.Matches() = %t, want = %t", got, tt.wantMatch)
			}
		})
	}
}
