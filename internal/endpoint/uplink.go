package endpoint

import (
	"net"
	"time"

	"go.uber.org/multierr"
)

const defaultTCPKeepAlivePeriod = 30 * time.Second

type AutoLingeringListener struct {
	net.Listener
}

func (l AutoLingeringListener) Accept() (c net.Conn, err error) {
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

func NewUplink(conn interface{}) (u Uplink) {
	switch c := conn.(type) {
	case net.Listener:
		u.Listener = c
		u.Addr = c.Addr()
	case net.PacketConn:
		u.PacketConn = c
		u.Addr = c.LocalAddr()
	case net.Addr:
		u.Unmanaged = true
		u.Addr = c
	}

	return
}

type Uplink struct {
	Addr       net.Addr
	Unmanaged  bool
	Listener   net.Listener
	PacketConn net.PacketConn
}

func (u Uplink) IsUDP() bool {
	_, ok := u.Addr.(*net.UDPAddr)
	return ok
}

func (u Uplink) IsTCP() bool {
	_, ok := u.Addr.(*net.TCPAddr)
	return ok
}

func (u *Uplink) Close() (err error) {
	if u.Listener != nil {
		multierr.AppendInvoke(&err, multierr.Close(u.Listener))
		u.Listener = nil
	}
	if u.PacketConn != nil {
		err = multierr.Combine(err, u.PacketConn.SetDeadline(time.Now()), u.PacketConn.Close())
		u.PacketConn = nil
	}
	return
}
