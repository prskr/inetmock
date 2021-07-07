package rules_test

import (
	"net"
	"testing"

	"gitlab.com/inetmock/inetmock/internal/rules"
)

func TestParseCIDR(t *testing.T) {
	t.Parallel()
	type args struct {
		cidr string
	}
	tests := []struct {
		name    string
		args    args
		wantC   *rules.CIDR
		wantErr bool
	}{
		{
			name: "Parse valid CIDR",
			args: args{
				cidr: "8.8.8.8/32",
			},
			wantC: &rules.CIDR{
				IPNet: &net.IPNet{
					IP:   net.ParseIP("8.8.8.8"),
					Mask: net.CIDRMask(32, 32),
				},
			},
			wantErr: false,
		},
		{
			name: "Parse invalid CIDR",
			args: args{
				cidr: "8.8.8.8/33",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotC, err := rules.ParseCIDR(tt.args.cidr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCIDR() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotC != nil {
				if !gotC.IP.Equal(tt.wantC.IP) {
					t.Errorf("ParseCIDR() IP got = %v, want %v", gotC.IP.String(), tt.wantC.IP.String())
					return
				}
				if gotC.Mask.String() != tt.wantC.Mask.String() {
					t.Errorf("ParseCIDR() Mask got = %v, want %v", gotC.Mask.String(), tt.wantC.Mask.String())
					return
				}
			}
		})
	}
}
