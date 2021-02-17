// +build linux

package pcap

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type CaptureParameters struct {
	LinkType layers.LinkType
}

type Consumer interface {
	Name() string
	Observe(pkg gopacket.Packet)
	Init(params CaptureParameters)
}
