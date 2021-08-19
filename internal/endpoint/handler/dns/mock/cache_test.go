package mock_test

import (
	"net"
	"testing"

	"gitlab.com/inetmock/inetmock/internal/endpoint/handler/dns/mock"
)

func Test_noOpCache_ForwardLookup(t *testing.T) {
	t.Parallel()
	type args struct {
		host string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Lookup google.com expect nil",
			args: args{
				host: "google.com",
			},
		},
		{
			name: "Lookup stackoverflow.com expect nil",
			args: args{
				host: "stackoverflow.com",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			n := mock.DelegateCache{}
			if got := n.ForwardLookup(tt.args.host); got != nil {
				t.Errorf("ForwardLookup() = %v, want nil", got)
			}
		})
	}
}

func Test_noOpCache_ReverseLookup(t *testing.T) {
	t.Parallel()
	type args struct {
		ip net.IP
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Reverse lookup 192.168.0.1 want miss",
			args: args{
				ip: net.ParseIP("192.168.0.1"),
			},
		},
		{
			name: "Reverse lookup 9.9.9.9 want miss",
			args: args{
				ip: net.ParseIP("9.9.9.9"),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			n := mock.DelegateCache{}
			gotHost, gotMiss := n.ReverseLookup(tt.args.ip)
			if gotHost != "" {
				t.Errorf("ReverseLookup() gotHost = %v, want ''", gotHost)
			}
			if !gotMiss {
				t.Errorf("ReverseLookup() gotMiss = %v, want true", gotMiss)
			}
		})
	}
}
