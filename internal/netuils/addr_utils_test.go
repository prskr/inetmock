package netuils_test

import (
	"net"
	"testing"

	"gitlab.com/inetmock/inetmock/internal/netuils"
)

func Test_extractIPFromAddress(t *testing.T) {
	t.Parallel()
	type args struct {
		addr net.Addr
	}
	type testCase struct {
		name    string
		args    args
		want    string
		wantErr bool
	}
	tests := []testCase{
		{
			name:    "Get address for IPv4 address",
			want:    "127.0.0.1",
			wantErr: false,
			args: args{
				addr: &net.TCPAddr{
					IP:   net.ParseIP("127.0.0.1"),
					Port: 23494,
				},
			},
		},
		{
			name:    "Get address for IPv6 address",
			want:    "::1",
			wantErr: false,
			args: args{
				addr: &net.TCPAddr{
					IP:   net.ParseIP("::1"),
					Port: 23494,
				},
			},
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := netuils.IPPortFromAddress(tt.args.addr)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractIPFromAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.IP.Equal(net.ParseIP(tt.want)) {
				t.Errorf("extractIPFromAddress() got = %v, want %v", got, tt.want)
			}
		})
	}
}
