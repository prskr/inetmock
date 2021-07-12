package mock

import (
	"net"
	"testing"
)

func mustParseCIDR(cidr string) *net.IPNet {
	_, n, err := net.ParseCIDR(cidr)
	if err != nil {
		panic(err)
	}
	return n
}

func TestIncrementalIPResolver_Lookup(t *testing.T) {
	t.Parallel()
	type fields struct {
		cidr   *net.IPNet
		offset uint32
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Default offset",
			fields: fields{
				cidr: mustParseCIDR("192.168.0.0/24"),
			},
		},
		{
			name: "offset at max address",
			fields: fields{
				cidr:   mustParseCIDR("192.168.0.0/24"),
				offset: 255,
			},
		},
		{
			name: "offset at max address",
			fields: fields{
				cidr:   mustParseCIDR("192.168.0.0/23"),
				offset: 511,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			i := IncrementalIPResolver{
				cidr:   tt.fields.cidr,
				offset: tt.fields.offset,
			}

			got := i.Lookup("")
			if !tt.fields.cidr.Contains(got) {
				t.Errorf("Lookup() = %v, is not in expected CIDR", got)
			}
		})
	}
}
