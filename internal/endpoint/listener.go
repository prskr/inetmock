package endpoint

import (
	"errors"
	"net"
	"strings"
)

var (
	ErrUDPMultiplexer           = errors.New("UDP listeners don't support multiplexing")
	ErrMultiplexingNotSupported = errors.New("not all handlers do support multiplexing")
	ErrUnsupportedProtocol      = errors.New("protocol not supported")
)

type HandlerReference string

func (h HandlerReference) ToLower() HandlerReference {
	return HandlerReference(strings.ToLower(string(h)))
}

type ListenerSpec struct {
	Name      string
	Protocol  string
	Address   string `mapstructure:"listenAddress"`
	Port      uint16
	Endpoints map[string]Spec
	Unmanaged bool
}

func (l ListenerSpec) Addr() (net.Addr, error) {
	ip := net.IPv4(0, 0, 0, 0)
	if l.Address != "" {
		ip = net.ParseIP(l.Address)
	}

	switch strings.ToLower(strings.TrimRight(l.Protocol, "46")) {
	case "tcp":
		return &net.TCPAddr{
			IP:   ip,
			Port: int(l.Port),
		}, nil
	case "udp":
		return &net.UDPAddr{
			IP:   ip,
			Port: int(l.Port),
		}, nil
	default:
		return nil, ErrUnsupportedProtocol
	}
}

type Spec struct {
	HandlerRef HandlerReference `mapstructure:"handler"`
	TLS        bool
	Handler    ProtocolHandler `mapstructure:"-"`
	Options    map[string]any
}
