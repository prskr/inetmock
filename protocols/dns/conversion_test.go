package dns_test

import (
	"net"
	"testing"

	"github.com/maxatome/go-testdeep/td"

	"gitlab.com/inetmock/inetmock/protocols/dns"
)

func TestIPToInt32(t *testing.T) {
	t.Parallel()
	type args struct {
		ip net.IP
	}
	type testCase struct {
		name string
		args args
		want uint32
	}
	tests := []testCase{
		{
			name: "Convert 188.193.106.113 to int",
			args: args{
				ip: net.ParseIP("188.193.106.113"),
			},
			want: 3166792305,
		},
		{
			name: "Convert 192.168.178.10 to int",
			args: args{
				ip: net.ParseIP("192.168.178.10"),
			},
			want: 3232281098,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			td.Cmp(t, dns.IPToInt32(tt.args.ip), tt.want)
		})
	}
}
