package dhcp

import (
	"fmt"
	"net"
	"strings"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"go.uber.org/zap"

	"gitlab.com/inetmock/inetmock/internal/rules"
	"gitlab.com/inetmock/inetmock/internal/state"
	"gitlab.com/inetmock/inetmock/pkg/logging"
)

var knownResponseHandlers = map[string]func(opts HandlerOptions, args ...rules.Param) (DHCPv4MessageHandler, error){
	"ip":      StaticIPHandler,
	"range":   IPRangeHandler,
	"router":  RouterIPHandler,
	"dns":     DNSHandler,
	"netmask": NetmaskHandler,
}

type HandlerOptions struct {
	ProtocolOptions
	Logger     logging.Logger
	StateStore state.KVStore
}

func HandlerForRoutingRule(rule *rules.ChainedResponsePipeline, opts HandlerOptions) (HandlerChain, error) {
	if rule.Response == nil || len(rule.Response) == 0 {
		return nil, rules.ErrNoTerminatorDefined
	}

	chain := make(HandlerChain, 0, len(rule.Response))
	for idx := range rule.Response {
		if constructor, ok := knownResponseHandlers[strings.ToLower(rule.Response[idx].Name)]; !ok {
			return nil, fmt.Errorf("%w %s", rules.ErrUnknownTerminator, rule.Response[idx].Name)
		} else if handler, err := constructor(opts, rule.Response[idx].Params...); err != nil {
			return nil, err
		} else {
			chain = append(chain, handler)
		}
	}

	return chain, nil
}

func StaticIPHandler(opts HandlerOptions, args ...rules.Param) (DHCPv4MessageHandler, error) {
	return singleIPModifier("static_ip", opts.Logger, args, func(ip net.IP, resp *dhcpv4.DHCPv4) {
		resp.YourIPAddr = ip
	})
}

func IPRangeHandler(opts HandlerOptions, args ...rules.Param) (DHCPv4MessageHandler, error) {
	const expectedParamsStartAndEnd = 2
	if err := rules.ValidateParameterCount(args, expectedParamsStartAndEnd); err != nil {
		return nil, err
	}

	var startIP, endIP net.IP

	if ip, err := args[0].AsIP(); err != nil {
		return nil, err
	} else {
		startIP = ip
	}

	if ip, err := args[1].AsIP(); err != nil {
		return nil, err
	} else {
		endIP = ip
	}

	return &RangeMessageHandler{
		Store:   opts.StateStore,
		TTL:     opts.Default.LeaseTime,
		StartIP: startIP,
		EndIP:   endIP,
	}, nil
}

func RouterIPHandler(opts HandlerOptions, args ...rules.Param) (DHCPv4MessageHandler, error) {
	return singleIPModifier("router", opts.Logger, args, func(ip net.IP, resp *dhcpv4.DHCPv4) {
		resp.Options.Update(dhcpv4.OptRouter(ip))
	})
}

func NetmaskHandler(opts HandlerOptions, args ...rules.Param) (DHCPv4MessageHandler, error) {
	return singleIPModifier("netmask", opts.Logger, args, func(ip net.IP, resp *dhcpv4.DHCPv4) {
		if ip = ip.To4(); ip == nil || len(ip) < 4 {
			return
		}
		mask := net.IPv4Mask(ip[0], ip[1], ip[2], ip[3])
		resp.Options.Update(dhcpv4.OptSubnetMask(mask))
	})
}

func DNSHandler(opts HandlerOptions, args ...rules.Param) (DHCPv4MessageHandler, error) {
	dnsIPs, err := multiIPArguments(args)
	if err != nil {
		return nil, err
	}

	handlerLogger := opts.Logger.With(zap.String("handler_type", "dns_handler"), logging.IPs("dns_ips", dnsIPs))
	return DHCPv4MessageHandlerFunc(func(req, resp *dhcpv4.DHCPv4) error {
		handlerLogger.Info("Set DNS servers", zap.Stringer("client_mac", req.ClientHWAddr))
		resp.Options.Update(dhcpv4.OptDNS(dnsIPs...))
		return nil
	}), nil
}

func singleIPModifier(
	name string,
	logger logging.Logger,
	args []rules.Param,
	modifier func(ip net.IP, resp *dhcpv4.DHCPv4),
) (DHCPv4MessageHandler, error) {
	ip, err := singleIPArgument(args)
	if err != nil {
		return nil, err
	}

	logger = logger.With(zap.String("handler_type", name), logging.IP("ip_value", ip))
	return DHCPv4MessageHandlerFunc(func(req, resp *dhcpv4.DHCPv4) error {
		logger.Info("Set IP value", zap.Stringer("client_mac", req.ClientHWAddr))
		modifier(ip, resp)
		return nil
	}), nil
}

func singleIPArgument(args []rules.Param) (ip net.IP, err error) {
	if err := rules.ValidateParameterCount(args, 1); err != nil {
		return nil, err
	}

	return args[0].AsIP()
}

func multiIPArguments(args []rules.Param) (ips []net.IP, err error) {
	if err := rules.ValidateParameterCount(args, 1); err != nil {
		return nil, err
	}

	ips = make([]net.IP, 0, len(args))
	for _, p := range args {
		if ip, err := p.AsIP(); err != nil {
			return nil, err
		} else {
			ips = append(ips, ip)
		}
	}
	return ips, nil
}
