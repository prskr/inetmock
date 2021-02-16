package pcap

import (
	"context"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type Device struct {
	Name        string
	IPAddresses []net.IP
}

type Subscription struct {
	Device    string
	Consumers []string
}

type CaptureParameters struct {
	SnapshotLength int32
	LinkType       layers.LinkType
}

type RecorderOption func(opt recorderOptions) recorderOptions

type Recorder interface {
	AvailableDevices() (devices []Device, err error)
	Subscriptions() (subscriptions []Subscription)
	Subscribe(ctx context.Context, device string, consumer Consumer) (err error)
	RemoveSubscriptions(device, consumerName string) (removed bool)
}

type Consumer interface {
	Name() string
	Observe(pkg gopacket.Packet)
	Init(params CaptureParameters)
}
