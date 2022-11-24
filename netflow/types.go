package netflow

import (
	"encoding"
	"encoding/binary"
	"fmt"
	"net"
	"strings"
)

var _ encoding.TextUnmarshaler = (*Protocol)(nil)

type Protocol uint32

const (
	ProtocolUnspecified Protocol = iota
	ProtocolTCP
	ProtocolUDP
)

func (p *Protocol) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		*p = ProtocolUnspecified
		return nil
	}

	switch strings.ToLower(string(text)) {
	case "tcp":
		*p = ProtocolTCP
	case "udp":
		*p = ProtocolUDP
	default:
		*p = ProtocolUnspecified
	}

	return nil
}

type (
	PacketSink interface {
		OnObservedPacket(pkt *Packet)
	}
	ErrorSink interface {
		OnError(err error)
	}

	PacketSinkFunc func(pkt *Packet)
	ErrorSinkFunc  func(err error)
)

var noOpErrorSink ErrorSink = ErrorSinkFunc(func(error) {
})

func (f PacketSinkFunc) OnObservedPacket(pkt *Packet) {
	f(pkt)
}

func (f ErrorSinkFunc) OnError(err error) {
	f(err)
}

type (
	MonitorMode  uint32
	PacketReader interface {
		Read() (*Packet, error)
		Close() error
	}
)

const (
	MonitorModeUnspecified MonitorMode = iota
	MonitorModePerfEvent
	MonitorModeRingBuf
	monitorModeMock
)

func (mm MonitorMode) Function() string {
	switch mm {
	case MonitorModeRingBuf:
		return "xdp_ingress_ring"
	case MonitorModePerfEvent:
		return "xdp_ingress_perf"
	case monitorModeMock:
		return "xdp_mock"
	case MonitorModeUnspecified:
		fallthrough
	default:
		return ""
	}
}

func (mm MonitorMode) Section() string {
	switch mm {
	case MonitorModeRingBuf:
		return "xdp/ring"
	case MonitorModePerfEvent:
		return "xdp/perf"
	case monitorModeMock:
		return "xdp/mock"
	case MonitorModeUnspecified:
		fallthrough
	default:
		return ""
	}
}

var _ encoding.BinaryUnmarshaler = (*Packet)(nil)

const (
	packetBinarySize = 16
)

type Packet struct {
	SourceIP   net.IP
	DestIP     net.IP
	SourcePort uint16
	DestPort   uint16
	Transport  Protocol
}

func (p *Packet) UnmarshalBinary(data []byte) error {
	if dataLen := len(data); dataLen < packetBinarySize {
		return fmt.Errorf("required 16 bytes bot got %d", dataLen)
	}

	p.SourceIP = make(net.IP, net.IPv4len)
	copy(p.SourceIP, data[:4])

	p.DestIP = make(net.IP, net.IPv4len)
	copy(p.DestIP, data[4:8])

	p.SourcePort = binary.LittleEndian.Uint16(data[8:10])
	p.DestPort = binary.LittleEndian.Uint16(data[10:12])
	p.Transport = Protocol(binary.LittleEndian.Uint32(data[12:]))

	return nil
}

func (p *Packet) Dispose() {
	p.SourceIP = nil
	p.DestIP = nil
	p.SourcePort = 0
	p.DestPort = 0
	p.Transport = ProtocolUnspecified
	packetPool.Put(p)
}
