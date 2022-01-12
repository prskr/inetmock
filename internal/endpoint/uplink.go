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
		u.Proto = NetProtoTCP
	case net.PacketConn:
		u.PacketConn = c
		u.Addr = c.LocalAddr()
		u.Proto = NetProtoUDP
	case net.Addr:
		u.Unmanaged = true
		u.Addr = c
		switch c.(type) {
		case *net.TCPAddr:
			u.Proto = NetProtoTCP
		case *net.UDPAddr:
			u.Proto = NetProtoUDP
		}
	}

	return
}

type Uplink struct {
	Addr       net.Addr
	Unmanaged  bool
	Proto      NetProto
	Listener   net.Listener
	PacketConn net.PacketConn
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
