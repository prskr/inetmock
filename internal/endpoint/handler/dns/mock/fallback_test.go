package mock

import (
	"net"
	"reflect"
	"testing"

	"github.com/maxatome/go-testdeep/td"
)

func Test_randomIPFallback_GetIP(t *testing.T) {
	t.Parallel()
	ra := randomIPFallback{}
	for i := 0; i < 1000; i++ {
		if got := ra.GetIP(); reflect.DeepEqual(got, net.IP{}) {
			t.Errorf("GetIP() = %v", got)
		}
	}
}

func Test_incrementalIPFallback_GetIP(t *testing.T) {
	t.Parallel()
	type fields struct {
		latestIP uint32
	}
	type testCase struct {
		name   string
		fields fields
		want   []net.IP
	}
	tests := []testCase{
		{
			name: "Expect the next icremental IP",
			fields: fields{
				latestIP: 167772160,
			},
			want: []net.IP{
				net.IPv4(10, 0, 0, 1),
			},
		},
		{
			name: "Expect a sequence of 5",
			fields: fields{
				latestIP: 167772160,
			},
			want: []net.IP{
				net.IPv4(10, 0, 0, 1),
				net.IPv4(10, 0, 0, 2),
				net.IPv4(10, 0, 0, 3),
				net.IPv4(10, 0, 0, 4),
				net.IPv4(10, 0, 0, 5),
			},
		},
		{
			name: "Expect next block to be incremented",
			fields: fields{
				latestIP: 167772413,
			},
			want: []net.IP{
				net.IPv4(10, 0, 0, 254),
				net.IPv4(10, 0, 0, 255),
				net.IPv4(10, 0, 1, 0),
			},
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			i := &incrementalIPFallback{
				latestIP: tt.fields.latestIP,
			}
			for k := 0; k < len(tt.want); k++ {
				td.Cmp(t, i.GetIP(), tt.want[k])
			}
		})
	}
}

func Test_ipToInt32(t *testing.T) {
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
			td.Cmp(t, ipToInt32(tt.args.ip), tt.want)
		})
	}
}
