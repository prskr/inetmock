package endpoint

import (
	"net"

	"go.uber.org/multierr"
)

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
		multierr.AppendInvoke(&err, multierr.Close(u.PacketConn))
		u.PacketConn = nil
	}
	return
}
