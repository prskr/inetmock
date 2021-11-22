//go:build linux
// +build linux

package pcap

import (
	"context"
	"errors"
	"io"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
	"go.uber.org/multierr"
)

const (
	transportClosingTimeout = 100 * time.Millisecond
)

var ErrTransportStillRunning = errors.New("transport to consumers did not stop in time")

func openDeviceForConsumers(device string, consumer Consumer, opts RecordingOptions) (*deviceConsumer, error) {
	var (
		handle *pcapgo.EthernetHandle
		err    error
	)

	if handle, err = pcapgo.NewEthernetHandle(device); err != nil {
		return nil, err
	}

	if err := handle.SetPromiscuous(opts.Promiscuous); err != nil {
		return nil, err
	}

	if err := consumer.Init(); err != nil {
		return nil, err
	}

	packetSrc := gopacket.NewZeroCopyPacketSource(handle, layers.LinkTypeEthernet)
	packetSrc.Lazy = true
	packetSrc.NoCopy = true
	dev := &deviceConsumer{
		locker:        new(sync.Mutex),
		handle:        handle,
		packetSource:  packetSrc,
		consumer:      consumer,
		transportStat: make(chan struct{}),
	}

	return dev, nil
}

type deviceConsumer struct {
	locker        sync.Locker
	consumer      Consumer
	cancel        context.CancelFunc
	handle        *pcapgo.EthernetHandle
	packetSource  *gopacket.PacketSource
	transportStat chan struct{}
}

func (o *deviceConsumer) Close() error {
	o.locker.Lock()
	defer o.locker.Unlock()

	err := o.handle.Close()

	if o.cancel != nil {
		o.cancel()
	}
	select {
	case _, more := <-o.transportStat:
		if more {
			return ErrTransportStillRunning
		}
	case <-time.After(transportClosingTimeout):
	}

	if closer, ok := o.consumer.(io.Closer); ok {
		return multierr.Append(err, closer.Close())
	}
	return err
}

func (o *deviceConsumer) StartTransport(ctx context.Context) {
	o.locker.Lock()
	defer o.locker.Unlock()
	ctx, o.cancel = context.WithCancel(ctx)
	go o.transportToConsumers(ctx)
}

func (o *deviceConsumer) transportToConsumers(ctx context.Context) {
	defer close(o.transportStat)
	for {
		select {
		case pkg, more := <-o.packetSource.Packets(ctx):
			if !more {
				return
			}
			o.locker.Lock()
			o.consumer.Observe(pkg)
			o.locker.Unlock()
		case <-ctx.Done():
			return
		}
	}
}
