package endpoint_test

import (
	"net"
	"reflect"
	"testing"

	"gitlab.com/inetmock/inetmock/internal/endpoint"
)

func TestListenerSpec_Addr(t *testing.T) {
	t.Parallel()
	type fields struct {
		Protocol string
		Address  string
		Port     uint16
	}
	tests := []struct {
		name    string
		fields  fields
		want    any
		wantErr bool
	}{
		{
			name: "TCP4 address",
			fields: fields{
				Protocol: "tcp4",
				Port:     1234,
			},
			want:    &net.TCPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 1234},
			wantErr: false,
		},
		{
			name: "TCP6 address",
			fields: fields{
				Protocol: "tcp6",
				Port:     1234,
			},
			want: &net.TCPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 1234},
		},
		{
			name: "TCP address",
			fields: fields{
				Protocol: "tcp",
				Port:     1234,
			},
			want: &net.TCPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 1234},
		},
		{
			name: "UDP4 address",
			fields: fields{
				Protocol: "udp4",
				Port:     1234,
			},
			want: &net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 1234},
		},
		{
			name: "UDP6 address",
			fields: fields{
				Protocol: "udp6",
				Port:     1234,
			},
			want: &net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 1234},
		},
		{
			name: "UDP address",
			fields: fields{
				Protocol: "udp",
				Port:     1234,
			},
			want: &net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 1234},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			l := endpoint.ListenerSpec{
				Protocol: tt.fields.Protocol,
				Address:  tt.fields.Address,
				Port:     tt.fields.Port,
			}
			got, err := l.Addr()
			if (err != nil) != tt.wantErr {
				t.Errorf("Addr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Addr() got = %v, want %v", got, tt.want)
			}
		})
	}
}
