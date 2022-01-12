package state_test

import (
	"net"
	"reflect"
	"testing"

	"github.com/maxatome/go-testdeep/td"

	"gitlab.com/inetmock/inetmock/internal/netutils"
	"gitlab.com/inetmock/inetmock/internal/state"
)

func TestMsgPackEncoding_Encode(t *testing.T) {
	t.Parallel()
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Encode & decode string",
			args: args{
				v: &struct {
					Greeting string
				}{"Hello, world"},
			},
		},
		{
			name: "Encode & decode integer",
			args: args{
				v: &struct {
					Nr int
				}{42},
			},
		},
		{
			name: "Encode & decode net.IP and net.HardwareAddress",
			args: args{
				v: &struct {
					IP  net.IP
					MAC net.HardwareAddr
				}{
					IP:  net.IPv4(1, 2, 3, 4),
					MAC: netutils.MustParseMAC("54:df:83:56:2c:f3"),
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			m := state.MsgPackEncoding{}
			into := reflect.New(reflect.TypeOf(tt.args.v).Elem()).Interface()
			gotData, err := m.Encode(tt.args.v)
			if err != nil {
				t.Errorf("Encode() error = %v", err)
				return
			}

			if err := m.Decode(gotData, into); err != nil {
				t.Errorf("Decode() error = %v", err)
				return
			}

			td.Cmp(t, into, tt.args.v)
		})
	}
}
