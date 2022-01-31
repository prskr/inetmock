package netutils

import (
	"errors"
	"net"
	"time"

	"go.uber.org/multierr"
)

var ErrNotATCPListener = errors.New("is not a TCP listener")

func WrapToManaged(listener net.Listener) (net.Listener, error) {
	if tcpListener, ok := listener.(*net.TCPListener); ok {
		return &managedListener{TCPListener: tcpListener}, nil
	}
	return nil, ErrNotATCPListener
}

func ListenTCP(addr *net.TCPAddr) (listener net.Listener, err error) {
	const network = "tcp"
	l := new(managedListener)
	l.TCPListener, err = net.ListenTCP(network, addr)
	if err != nil {
		return nil, err
	}

	return l, err
}

type managedListener struct {
	*net.TCPListener
}

func (l *managedListener) Accept() (c net.Conn, err error) {
	const defaultTCPKeepAlivePeriod = 30 * time.Second

	c, err = l.TCPListener.Accept()
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

func (l *managedListener) Close() error {
	return multierr.Append(l.TCPListener.SetDeadline(time.Now()), l.TCPListener.Close())
}
