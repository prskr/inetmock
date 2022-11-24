package netutils_test

import (
	"net"
	"testing"

	"inetmock.icb4dc0.de/inetmock/internal/netutils"
)

type myVerySpecialAddress string

func (m myVerySpecialAddress) Network() string {
	return string(m)
}

func (m myVerySpecialAddress) String() string {
	return string(m)
}

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
			name:    "Get address for IPv4 TCP address",
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
			name:    "Get address for IPv6 TCP address",
			want:    "::1",
			wantErr: false,
			args: args{
				addr: &net.TCPAddr{
					IP:   net.ParseIP("::1"),
					Port: 23494,
				},
			},
		},
		{
			name:    "Get address for IPv4 UDP address",
			want:    "127.0.0.1",
			wantErr: false,
			args: args{
				addr: &net.UDPAddr{
					IP:   net.ParseIP("127.0.0.1"),
					Port: 23494,
				},
			},
		},
		{
			name:    "Get address for IPv6 UDP address",
			want:    "::1",
			wantErr: false,
			args: args{
				addr: &net.UDPAddr{
					IP:   net.ParseIP("::1"),
					Port: 23494,
				},
			},
		},
		{
			name:    "Error due to unknown address type",
			wantErr: true,
			args: args{
				addr: myVerySpecialAddress("127.0.0.1:1234"),
			},
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ip, _, err := netutils.IPPortFromAddress(tt.args.addr)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("extractIPFromAddress() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if !ip.Equal(net.ParseIP(tt.want)) {
				t.Errorf("extractIPFromAddress() got = %v, want %v", ip, tt.want)
			}
		})
	}
}
