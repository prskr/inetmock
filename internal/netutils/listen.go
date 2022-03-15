package netutils

import (
	"net"
	"time"

	"github.com/valyala/tcplisten"
	"go.uber.org/multierr"
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

func WrapToManaged(listener net.Listener) net.Listener {
	return &managedListener{Listener: listener}
}

func ListenTCP(addr *net.TCPAddr, opts ...TCPListenOption) (listener net.Listener, err error) {
	const network = "tcp4"
	l := new(managedListener)
	listenerCfg := new(tcplisten.Config)

	for i := range opts {
		opts[i].apply(listenerCfg)
	}

	l.Listener, err = listenerCfg.NewListener(network, addr.String())
	if err != nil {
		return nil, err
	}

	return l, err
}

type managedListener struct {
	net.Listener
}

func (l *managedListener) Accept() (c net.Conn, err error) {
	const defaultTCPKeepAlivePeriod = 30 * time.Second

	c, err = l.Listener.Accept()
	if err != nil {
		return nil, err
	}

	if tcpConn, ok := c.(*net.TCPConn); ok {
		err = multierr.Combine(
			tcpConn.SetLinger(0),
			tcpConn.SetKeepAlive(true),
			tcpConn.SetKeepAlivePeriod(defaultTCPKeepAlivePeriod),
		)
	}
	return
}

func (l *managedListener) Close() (err error) {
	if deadline, ok := l.Listener.(interface{ SetDeadline(t time.Time) error }); ok {
		err = deadline.SetDeadline(time.Now())
	}
	return multierr.Append(err, l.Listener.Close())
}
