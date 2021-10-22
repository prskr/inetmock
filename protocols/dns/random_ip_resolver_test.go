package dns_test

import (
	"net"
	"testing"

	"gitlab.com/inetmock/inetmock/protocols/dns"
)

func TestRandomIPResolver_Lookup(t *testing.T) {
	t.Parallel()
	type fields struct {
		CIDR *net.IPNet
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Random IP from 24 bit CIDR",
			fields: fields{
				CIDR: mustParseCIDR("192.168.0.0/24"),
			},
		},
		{
			name: "Random IP from 29 bit CIDR",
			fields: fields{
				CIDR: mustParseCIDR("192.168.0.0/29"),
			},
		},
		{
			name: "Random IP from 16 bit CIDR",
			fields: fields{
				CIDR: mustParseCIDR("10.5.0.0/16"),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := dns.NewRandomIPResolver(tt.fields.CIDR)
			got := r.Lookup("")
			if !tt.fields.CIDR.Contains(got) {
				t.Errorf("Lookup() = %v, not in expected CIDR", got)
			}
		})
	}
}
