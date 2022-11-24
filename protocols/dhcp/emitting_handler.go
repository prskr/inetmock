package dhcp

import (
	"github.com/insomniacslk/dhcp/dhcpv4"

	"inetmock.icb4dc0.de/inetmock/pkg/audit"
	auditv1 "inetmock.icb4dc0.de/inetmock/pkg/audit/v1"
)

type EmittingHandler struct {
	Upstream DHCPv4MessageHandler
	Emitter  audit.Emitter
}

func (h *EmittingHandler) Handle(req, resp *dhcpv4.DHCPv4) error {
	const defaultDHCPPort = 67
	h.Emitter.Builder().
		WithApplication(auditv1.AppProtocol_APP_PROTOCOL_DHCP).
		WithTransport(auditv1.TransportProtocol_TRANSPORT_PROTOCOL_UDP).
		WithDestination(req.ClientIPAddr, defaultDHCPPort).
		WithSource(req.ClientIPAddr, 0).
		WithProtocolDetails(&audit.DHCP{
			HopCount: req.HopCount,
			HWType:   auditv1.DHCPHwType(req.HWType),
			OpCode:   auditv1.DHCPOpCode(req.OpCode),
		}).
		Emit()

	return h.Upstream.Handle(req, resp)
}
