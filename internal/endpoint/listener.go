package endpoint

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"sort"
	"strings"

	"github.com/soheilhy/cmux"
)

var (
	ErrUDPMultiplexer           = errors.New("UDP listeners don't support multiplexing")
	ErrMultiplexingNotSupported = errors.New("not all handlers do support multiplexing")
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
	Uplink    *Uplink `mapstructure:"-"`
}

type Spec struct {
	HandlerRef HandlerReference `mapstructure:"handler"`
	TLS        bool
	Handler    ProtocolHandler `mapstructure:"-"`
	Options    map[string]interface{}
}

func (l *ListenerSpec) ConfigureMultiplexing(tlsConfig *tls.Config) ([]Endpoint, []cmux.CMux, error) {
	if l.Uplink == nil {
		if err := l.setupUplink(); err != nil {
			return nil, nil, err
		}
	}

	if len(l.Endpoints) <= 1 {
		for name, s := range l.Endpoints {
			if s.TLS {
				l.Uplink.Listener = tls.NewListener(l.Uplink.Listener, tlsConfig)
			}
			endpoints := []Endpoint{
				{
					name:   fmt.Sprintf("%s:%s", l.Name, name),
					uplink: *l.Uplink,
					Spec:   s,
				},
			}
			return endpoints, nil, nil
		}
	}

	if l.Uplink.Proto == NetProtoUDP {
		return nil, nil, ErrUDPMultiplexer
	}

	plainGrp, tlsGrp, err := l.groupByTLS()
	if err != nil {
		return nil, nil, err
	}

	var muxes []cmux.CMux
	endpoints := make([]Endpoint, 0, len(l.Endpoints))
	lis := l.Uplink.Listener

	if len(plainGrp.Names) > 0 {
		plainMux := cmux.New(lis)
		endpoints = append(endpoints, l.setupMux(plainMux, plainGrp)...)
		muxes = append(muxes, plainMux)
		lis = plainMux.Match(cmux.Any())
	}

	if len(tlsGrp.Names) > 0 {
		tlsMux := cmux.New(tls.NewListener(lis, tlsConfig))
		endpoints = append(endpoints, l.setupMux(tlsMux, tlsGrp)...)
		muxes = append(muxes, tlsMux)
	}

	return endpoints, muxes, nil
}

func (l *ListenerSpec) groupByTLS() (plainGrp, tlsGrp *endpointGroup, err error) {
	if plainGrp, err = groupEndpoints(l.Endpoints, func(s Spec) bool { return !s.TLS }); err != nil {
		return nil, nil, err
	}

	if tlsGrp, err = groupEndpoints(l.Endpoints, func(s Spec) bool { return s.TLS }); err != nil {
		return nil, nil, err
	}

	return
}

func (l *ListenerSpec) setupMux(mux cmux.CMux, grp *endpointGroup) (endpoints []Endpoint) {
	for idx := range grp.Names {
		name := grp.Names[idx]
		epSpec := l.Endpoints[name]
		endpoints = append(endpoints, Endpoint{
			name: fmt.Sprintf("%s:%s", l.Name, name),
			uplink: Uplink{
				Proto:    NetProtoTCP,
				Addr:     l.Uplink.Addr,
				Listener: mux.Match(grp.Handlers[name].Matchers()...),
			},
			Spec: epSpec,
		})
	}
	return
}

func (l *ListenerSpec) setupUplink() (err error) {
	l.Uplink = &Uplink{
		Unmanaged: l.Unmanaged,
	}
	switch l.Protocol {
	case "udp", "udp4", "udp6":
		l.Uplink.Proto = NetProtoUDP
		addr := &net.UDPAddr{
			IP:   net.ParseIP(l.Address),
			Port: int(l.Port),
		}
		l.Uplink.Addr = addr
		if !l.Unmanaged {
			l.Uplink.PacketConn, err = net.ListenUDP(l.Protocol, addr)
		}
	case "tcp", "tcp4", "tcp6":
		l.Uplink.Proto = NetProtoTCP
		addr := &net.TCPAddr{
			IP:   net.ParseIP(l.Address),
			Port: int(l.Port),
		}
		l.Uplink.Addr = addr
		if !l.Unmanaged {
			l.Uplink.Listener, err = net.ListenTCP(l.Protocol, addr)
		}
	default:
		err = errors.New("protocol not supported")
	}
	return
}

type endpointGroup struct {
	Names    []string
	Handlers map[string]MultiplexHandler
}

func groupEndpoints(endpoints map[string]Spec, predicate func(s Spec) bool) (*endpointGroup, error) {
	grp := &endpointGroup{
		Names:    make([]string, 0, len(endpoints)),
		Handlers: make(map[string]MultiplexHandler),
	}

	for name, spec := range endpoints {
		var e MultiplexHandler
		if ep, ok := spec.Handler.(MultiplexHandler); !ok {
			return nil, fmt.Errorf("handler %s %w", spec.HandlerRef, ErrMultiplexingNotSupported)
		} else {
			e = ep
		}

		if predicate(spec) {
			grp.Names = append(grp.Names, name)
			grp.Handlers[name] = e
		}
	}
	sort.Strings(grp.Names)
	return grp, nil
}
