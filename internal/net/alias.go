package net

import "net"

type (
	UDPAddr    = net.UDPAddr
	TCPAddr    = net.TCPAddr
	Listener   = net.Listener
	PacketConn = net.PacketConn
)

var ListenUDP = net.ListenUDP
