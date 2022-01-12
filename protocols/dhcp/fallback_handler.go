package dhcp

import (
	"net"

	"github.com/insomniacslk/dhcp/dhcpv4"

	"gitlab.com/inetmock/inetmock/pkg/logging"
)

type FallbackHandler struct {
	Previous DHCPv4MessageHandler
	Logger   logging.Logger
	DefaultOptions
}

func (h *FallbackHandler) Handle(req, resp *dhcpv4.DHCPv4) error {
	if err := h.Previous.Handle(req, resp); err != nil {
		return err
	}

	internalHandlers := []DHCPv4MessageHandler{
		DHCPv4MessageHandlerFunc(h.handleServerID),
		DHCPv4MessageHandlerFunc(h.handleRouter),
		DHCPv4MessageHandlerFunc(h.handleNetmask),
		DHCPv4MessageHandlerFunc(h.handleDNS),
	}

	for idx := range internalHandlers {
		if err := internalHandlers[idx].Handle(req, resp); err != nil {
			return err
		}
	}

	return nil
}

func (h *FallbackHandler) handleRouter(_, resp *dhcpv4.DHCPv4) error {
	if !resp.Options.Has(dhcpv4.OptionRouter) {
		h.Logger.Info("Set fallback router", logging.IP("ip_value", h.Router))
		resp.Options.Update(dhcpv4.OptRouter(h.Router))
	}
	return nil
}

func (h *FallbackHandler) handleNetmask(_, resp *dhcpv4.DHCPv4) error {
	if !resp.Options.Has(dhcpv4.OptionSubnetMask) {
		if ip := h.Netmask.To4(); ip != nil && len(ip) > 3 {
			h.Logger.Info("Set fallback netmask", logging.IP("ip_value", h.Netmask))
			resp.Options.Update(dhcpv4.OptSubnetMask(net.IPv4Mask(ip[0], ip[1], ip[2], ip[3])))
		}
	}
	return nil
}

func (h *FallbackHandler) handleDNS(_, resp *dhcpv4.DHCPv4) error {
	if !resp.Options.Has(dhcpv4.OptionDomainNameServer) {
		h.Logger.Info("Set fallback DNS servers", logging.IPs("ip_value", h.DNS))
		resp.Options.Update(dhcpv4.OptDNS(h.DNS...))
	}
	return nil
}

func (h *FallbackHandler) handleServerID(req, resp *dhcpv4.DHCPv4) error {
	if req.OpCode != dhcpv4.OpcodeBootRequest {
		return nil
	}

	if req.ServerIPAddr != nil && !req.ServerIPAddr.Equal(net.IPv4zero) && !req.ServerIPAddr.Equal(h.ServerID) {
		return ErrDropRequest
	}

	h.Logger.Info("Set server_id", logging.IP("ip_value", h.ServerID))
	resp.ServerIPAddr = make(net.IP, len(h.ServerID))
	copy(resp.ServerIPAddr, h.ServerID)
	resp.UpdateOption(dhcpv4.OptServerIdentifier(h.ServerID))

	return nil
}
