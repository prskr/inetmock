package netflow

import (
	"net"

	"inetmock.icb4dc0.de/inetmock/pkg/audit"
	auditv1 "inetmock.icb4dc0.de/inetmock/pkg/audit/v1"
)

var _ PacketSink = (*EmittingPacketSink)(nil)

type Lookup interface {
	ReverseLookup(address net.IP) (host string, miss bool)
}

type EmittingPacketSink struct {
	Lookup
	audit.Emitter
}

func (e EmittingPacketSink) OnObservedPacket(pkt *Packet) {
	var details audit.Details
	if host, miss := e.ReverseLookup(pkt.DestIP); !miss {
		details = &audit.NetMon{ReverseResolvedDomain: host}
	}

	e.Builder().
		WithApplication(auditv1.AppProtocol_APP_PROTOCOL_UNSPECIFIED).
		WithTransport(auditv1.TransportProtocol(pkt.Transport)).
		WithSource(pkt.SourceIP, pkt.SourcePort).
		WithDestination(pkt.DestIP, pkt.DestPort).
		WithProtocolDetails(details).
		Emit()
}
