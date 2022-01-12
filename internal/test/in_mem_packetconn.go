package test

import (
	"net"
)

var _ net.PacketConn = (*pipePacketConn)(nil)

func NewInMemoryPacketConnPipe(clientAddr, srvAddr net.Addr) (upstream, downstream net.PacketConn) {
	conn1, conn2 := net.Pipe()
	upstream = &pipePacketConn{
		peer: clientAddr,
		Conn: conn1,
	}

	downstream = &pipePacketConn{
		peer: srvAddr,
		Conn: conn2,
	}

	return
}

type pipePacketConn struct {
	peer net.Addr
	net.Conn
}

func (pc pipePacketConn) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	n, err = pc.Read(p)
	addr = pc.peer
	return
}

func (pc pipePacketConn) WriteTo(p []byte, _ net.Addr) (n int, err error) {
	return pc.Write(p)
}
