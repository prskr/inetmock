package mock

import (
	"net"
	"reflect"
	"testing"
)

func Test_randomIPFallback_GetIP(t *testing.T) {
	ra := randomIPFallback{}
	for i := 0; i < 1000; i++ {
		if got := ra.GetIP(); reflect.DeepEqual(got, net.IP{}) {
			t.Errorf("GetIP() = %v", got)
		}
	}
}

func Test_incrementalIPFallback_GetIP(t *testing.T) {
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
	scenario := func(tt testCase) func(t *testing.T) {
		return func(t *testing.T) {
			i := &incrementalIPFallback{
				latestIP: tt.fields.latestIP,
			}
			for k := 0; k < len(tt.want); k++ {
				if got := i.GetIP(); !reflect.DeepEqual(got, tt.want[k]) {
					t.Errorf("GetIP() = %v, want %v", got, tt.want[k])
				}
			}
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, scenario(tt))
	}
}

func Test_ipToInt32(t *testing.T) {
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
	scenario := func(tt testCase) func(t *testing.T) {
		return func(t *testing.T) {
			if got := ipToInt32(tt.args.ip); got != tt.want {
				t.Errorf("ipToInt32() = %v, want %v", got, tt.want)
			}
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, scenario(tt))
	}
}
