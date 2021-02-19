// +build linux

package pcap

import (
	"context"
	"io"
	"net"
	"time"

	"github.com/google/gopacket"
)

type Consumer interface {
	Name() string
	Observe(pkg gopacket.Packet)
	Init() error
}

type Device struct {
	Name        string
	IPAddresses []net.IP
}

type Subscription struct {
	ConsumerKey  string
	ConsumerName string
}

type RecordingOptions struct {
	Promiscuous bool
	ReadTimeout time.Duration
}

type Recorder interface {
	io.Closer
	AvailableDevices() (devices []Device, err error)
	Subscriptions() (subscriptions []Subscription)
	StartRecording(ctx context.Context, device string, consumer Consumer) (err error)
	StartRecordingWithOptions(ctx context.Context, device string, consumer Consumer, opts RecordingOptions) (err error)
	StopRecording(consumerKey string) (err error)
}
