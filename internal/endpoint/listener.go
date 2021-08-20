//go:generate go-enum -f $GOFILE --lower --marshal --names

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

/* ENUM(
UDP,
TCP
)
*/
type NetProto int

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
	Uplink    *Uplink `mapstructure:"-"`
}

type Spec struct {
	HandlerRef HandlerReference `mapstructure:"protocols"`
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

	epNames := make([]string, len(l.Endpoints))
	multiplexEndpoints := make(map[string]MultiplexHandler)
	var idx int
	for name, spec := range l.Endpoints {
		epNames[idx] = name
		idx++
		if ep, ok := spec.Handler.(MultiplexHandler); !ok {
			return nil, nil, fmt.Errorf("protocols %s %w", spec.HandlerRef, ErrMultiplexingNotSupported)
		} else {
			multiplexEndpoints[name] = ep
		}
	}

	sort.Strings(epNames)

	plainMux := cmux.New(l.Uplink.Listener)
	tlsListener := plainMux.Match(cmux.TLS())
	tlsListener = tls.NewListener(tlsListener, tlsConfig)
	tlsMux := cmux.New(tlsListener)

	tlsRequired := false

	endpoints := make([]Endpoint, len(epNames))
	for i, epName := range epNames {
		epSpec := l.Endpoints[epName]
		epMux := plainMux
		if epSpec.TLS {
			epMux = tlsMux
			tlsRequired = true
		}
		epListener := Endpoint{
			name: fmt.Sprintf("%s:%s", l.Name, epName),
			uplink: Uplink{
				Proto:    NetProtoTCP,
				Listener: epMux.Match(multiplexEndpoints[epName].Matchers()...),
			},
			Spec: epSpec,
		}

		endpoints[i] = epListener
	}

	var muxes []cmux.CMux
	muxes = append(muxes, plainMux)

	if tlsRequired {
		muxes = append(muxes, tlsMux)
	} else {
		_ = tlsListener.Close()
	}

	return endpoints, muxes, nil
}

func (l *ListenerSpec) setupUplink() (err error) {
	l.Uplink = new(Uplink)
	switch l.Protocol {
	case "udp", "udp4", "udp6":
		l.Uplink.Proto = NetProtoUDP
		l.Uplink.PacketConn, err = net.ListenUDP(l.Protocol, &net.UDPAddr{
			IP:   net.ParseIP(l.Address),
			Port: int(l.Port),
		})
	case "tcp", "tcp4", "tcp6":
		l.Uplink.Proto = NetProtoTCP
		l.Uplink.Listener, err = net.ListenTCP(l.Protocol, &net.TCPAddr{
			IP:   net.ParseIP(l.Address),
			Port: int(l.Port),
		})
	default:
		err = errors.New("protocol not supported")
	}
	return
}
