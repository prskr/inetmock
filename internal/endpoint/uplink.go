package endpoint

import (
	"net"

	"go.uber.org/multierr"
)

type Uplink struct {
	Proto      NetProto
	Listener   net.Listener
	PacketConn net.PacketConn
}

func (u Uplink) Addr() net.Addr {
	if u.Listener != nil {
		return u.Listener.Addr()
	}
	if u.PacketConn != nil {
		return u.PacketConn.LocalAddr()
	}
	return nil
}

func (u Uplink) Close() (err error) {
	if u.Listener != nil {
		multierr.AppendInvoke(&err, multierr.Close(u.Listener))
	}
	if u.PacketConn != nil {
		multierr.AppendInvoke(&err, multierr.Close(u.PacketConn))
	}
	return
}
