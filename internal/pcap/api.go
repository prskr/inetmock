package pcap

import (
	"context"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type Device struct {
	Name        string
	IPAddresses []net.IP
}

type Subscription struct {
	ConsumerKey  string
	ConsumerName string
}

type CaptureParameters struct {
	LinkType layers.LinkType
}

type RecordingOptions struct {
	Promiscuous bool
	ReadTimeout time.Duration
}

type Recorder interface {
	AvailableDevices() (devices []Device, err error)
	Subscriptions() (subscriptions []Subscription)
	StartRecording(ctx context.Context, device string, consumer Consumer) (err error)
	StartRecordingWithOptions(ctx context.Context, device string, consumer Consumer, opts RecordingOptions) (err error)
	StopRecording(consumerKey string) (err error)
}

type Consumer interface {
	Name() string
	Observe(pkg gopacket.Packet)
	Init(params CaptureParameters)
}