package netflow_test

import (
	"net"
	"testing"

	"inetmock.icb4dc0.de/inetmock/netflow"
)

func TestIPPortProto_Protocol(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		p    netflow.IPPortProto
		want netflow.Protocol
	}{
		{
			name: "TCP for full specified IPPortProto",
			p:    netflow.IPPortProto("0.0.0.0:80/tcp"),
			want: netflow.ProtocolTCP,
		},
		{
			name: "UDP for full specified IPPortProto",
			p:    netflow.IPPortProto("0.0.0.0:53/udp"),
			want: netflow.ProtocolUDP,
		},
		{
			name: "Missing protocol",
			p:    netflow.IPPortProto("0.0.0.0:80"),
			want: netflow.ProtocolUnspecified,
		},
		{
			name: "TCP port - missing IP",
			p:    netflow.IPPortProto("80/tcp"),
			want: netflow.ProtocolTCP,
		},
		{
			name: "TCP port - missing IP, case insensitive",
			p:    netflow.IPPortProto("80/TCP"),
			want: netflow.ProtocolTCP,
		},
		{
			name: "UDP port - missing IP",
			p:    netflow.IPPortProto("80/udp"),
			want: netflow.ProtocolUDP,
		},
		{
			name: "UDP port - missing IP, case insensitive",
			p:    netflow.IPPortProto("80/UDP"),
			want: netflow.ProtocolUDP,
		},
		{
			name: "unspecified port - wrong syntax",
			p:    netflow.IPPortProto("80:udp"),
			want: netflow.ProtocolUnspecified,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.p.Protocol(); got != tt.want {
				t.Errorf("Protocol() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPortProto_Port(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		pp   netflow.IPPortProto
		want uint16
	}{
		{
			name: "Port 53 - valid IP",
			pp:   netflow.IPPortProto("1.2.3.4:53/udp"),
			want: 53,
		},
		{
			name: "Port 53 - zero IP",
			pp:   netflow.IPPortProto("0.0.0.0:53/udp"),
			want: 53,
		},
		{
			name: "Port 53 - missing IP",
			pp:   netflow.IPPortProto("53/udp"),
			want: 53,
		},
		{
			name: "Port 443 - missing IP",
			pp:   netflow.IPPortProto("443/tcp"),
			want: 443,
		},
		{
			name: "Port 31876 - missing IP",
			pp:   netflow.IPPortProto("31876/tcp"),
			want: 31876,
		},
		{
			name: "uint16 overflow - missing IP",
			pp:   netflow.IPPortProto("66789/tcp"),
			want: 0,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.pp.Port(); got != tt.want {
				t.Errorf("Port() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIPPortProto_IP(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		p    netflow.IPPortProto
		want net.IP
	}{
		{
			name: "Explicitly get zero val",
			p:    netflow.IPPortProto("0.0.0.0:80/tcp"),
			want: net.IPv4zero,
		},
		{
			name: "Implicitly get zero val",
			p:    netflow.IPPortProto("80/tcp"),
			want: net.IPv4zero,
		},
		{
			name: "Get nil value for empty string",
			p:    netflow.IPPortProto(""),
			want: nil,
		},
		{
			name: "Get actual IP",
			p:    netflow.IPPortProto("1.2.3.4:80/tcp"),
			want: net.IPv4(1, 2, 3, 4),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.p.IP(); !tt.want.Equal(got) {
				t.Errorf("IP() = %v, want %v", got, tt.want)
			}
		})
	}
}
