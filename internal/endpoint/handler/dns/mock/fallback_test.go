package mock

import (
	"net"
	"testing"

	"github.com/maxatome/go-testdeep/td"
)

func Test_randomIPFallback_GetIP(t *testing.T) {
	t.Parallel()
	ra := randomIPFallback{}
	for i := 0; i < 1000; i++ {
		got := ra.GetIP()
		td.CmpNot(t, got, net.IP{})
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
