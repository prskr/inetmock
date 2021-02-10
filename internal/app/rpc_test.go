package app

import (
	"net/url"
	"reflect"
	"testing"
)

func TestRPC_ListenURL(t *testing.T) {
	type fields struct {
		Listen string
	}
	tests := []struct {
		name   string
		fields fields
		wantU  *url.URL
	}{
		{
			name: "Parse valid TCP URL",
			fields: fields{
				Listen: "tcp://localhost:8080",
			},
			wantU: func() *url.URL {
				if u, e := url.Parse("tcp://localhost:8080"); e != nil {
					t.Errorf("Error during URL parsing: %v", e)
					return nil
				} else {
					return u
				}
			}(),
		},
		{
			name: "Parse valid unix socket url",
			fields: fields{
				Listen: "unix:///var/run/inetmock.sock",
			},
			wantU: func() *url.URL {
				if u, e := url.Parse("unix:///var/run/inetmock.sock"); e != nil {
					t.Errorf("Error during URL parsing: %v", e)
					return nil
				} else {
					return u
				}
			}(),
		},
		{
			name: "Expect fallback value due to parse error",
			fields: fields{
				Listen: `"tcp;\\asdf:234sedf`,
			},
			wantU: func() *url.URL {
				if u, e := url.Parse("tcp://:0"); e != nil {
					t.Errorf("Error during URL parsing: %v", e)
					return nil
				} else {
					return u
				}
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RPC{
				Listen: tt.fields.Listen,
			}
			if gotU := r.ListenURL(); !reflect.DeepEqual(gotU, tt.wantU) {
				t.Errorf("ListenURL() = %v, want %v", gotU, tt.wantU)
			}
		})
	}
}
