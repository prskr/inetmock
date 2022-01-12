package dhcp

import (
	"github.com/insomniacslk/dhcp/dhcpv4"

	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/audit/details"
	auditv1 "gitlab.com/inetmock/inetmock/pkg/audit/v1"
)

type EmittingHandler struct {
	Upstream DHCPv4MessageHandler
	Emitter  audit.Emitter
}

func (h *EmittingHandler) Handle(req, resp *dhcpv4.DHCPv4) error {
	ev := audit.Event{
		Application:   auditv1.AppProtocol_APP_PROTOCOL_DHCP,
		Transport:     auditv1.TransportProtocol_TRANSPORT_PROTOCOL_UDP,
		DestinationIP: req.ClientIPAddr,
		ProtocolDetails: details.DHCP{
			HopCount: req.HopCount,
			HWType:   auditv1.DHCPHwType(req.HWType),
			OpCode:   auditv1.DHCPOpCode(req.OpCode),
		},
	}

	ev.SourceIP = req.ClientIPAddr

	h.Emitter.Emit(ev)
	return h.Upstream.Handle(req, resp)
}
