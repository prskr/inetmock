package netutils

import (
	"net"

	"github.com/valyala/tcplisten"
)

type (
	TCPListenOption interface {
		apply(cfg *tcplisten.Config)
	}
	TCPListenOptionFunc func(cfg *tcplisten.Config)
)

func (f TCPListenOptionFunc) apply(cfg *tcplisten.Config) {
	f(cfg)
}

var (
	WithDeferAccept = func(deferAccept bool) TCPListenOption {
		return TCPListenOptionFunc(func(cfg *tcplisten.Config) {
			cfg.DeferAccept = deferAccept
		})
	}

	WithReusePort = func(reusePort bool) TCPListenOption {
		return TCPListenOptionFunc(func(cfg *tcplisten.Config) {
			cfg.ReusePort = reusePort
		})
	}

	WithFastOpen = func(fastOpen bool) TCPListenOption {
		return TCPListenOptionFunc(func(cfg *tcplisten.Config) {
			cfg.FastOpen = fastOpen
		})
	}
)

func ListenTCP(addr *net.TCPAddr, opts ...TCPListenOption) (listener net.Listener, err error) {
	const network = "tcp4"

	listenerCfg := new(tcplisten.Config)

	for i := range opts {
		opts[i].apply(listenerCfg)
	}

	return listenerCfg.NewListener(network, addr.String())
}
