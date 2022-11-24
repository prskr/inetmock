package dns_test

import (
	"net"
	"reflect"
	"testing"

	mdns "github.com/miekg/dns"

	"inetmock.icb4dc0.de/inetmock/protocols/dns"
)

func mustReverseAddr(addr string) string {
	if arpa, err := mdns.ReverseAddr(addr); err != nil {
		panic(err)
	} else {
		return arpa
	}
}

func TestParseInAddrArpa(t *testing.T) {
	t.Parallel()
	type args struct {
		inAddrArpa string
	}
	tests := []struct {
		name string
		args args
		want net.IP
	}{
		{
			name: "Parse 1.1.1.1",
			args: args{
				inAddrArpa: mustReverseAddr("1.1.1.1"),
			},
			want: net.IPv4(1, 1, 1, 1),
		},
		{
			name: "Parse 192.168.0.1",
			args: args{
				inAddrArpa: mustReverseAddr("192.168.0.1"),
			},
			want: net.IPv4(192, 168, 0, 1),
		},
		{
			name: "Invalid reverse address without suffix",
			args: args{
				inAddrArpa: "1.0.168.192",
			},
			want: nil,
		},
		{
			name: "Invalid number of blocks",
			args: args{
				inAddrArpa: "12.1.0.168.192.in-addr.arpa.",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := dns.ParseInAddrArpa(tt.args.inAddrArpa); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseInAddrArpa() = %v, want %v", got, tt.want)
			}
		})
	}
}
