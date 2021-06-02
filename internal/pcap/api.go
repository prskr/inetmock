// +build linux
//go:build linux

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

type StartRecordingResult struct {
	ConsumerKey string
}

type Recorder interface {
	io.Closer
	AvailableDevices() (devices []Device, err error)
	Subscriptions() (subscriptions []Subscription)
	StartRecording(ctx context.Context, device string, consumer Consumer) (result *StartRecordingResult, err error)
	StartRecordingWithOptions(
		ctx context.Context,
		device string,
		consumer Consumer,
		opts RecordingOptions,
	) (result *StartRecordingResult, err error)
	StopRecording(consumerKey string) (err error)
}
