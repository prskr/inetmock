package pcap

import (
	"context"
	"io"
	"sync"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
)

func openDeviceForConsumers(ctx context.Context, device string, consumer Consumer, opts RecordingOptions) (dev deviceConsumer, err error) {
	var handle *pcapgo.EthernetHandle
	if handle, err = pcapgo.NewEthernetHandle(device); err != nil {
		return
	}

	if err = handle.SetPromiscuous(opts.Promiscuous); err != nil {
		return
	}

	consumerCtx, cancel := context.WithCancel(ctx)

	consumer.Init(CaptureParameters{
		LinkType: layers.LinkTypeEthernet,
	})

	dev = deviceConsumer{
		locker: new(sync.Mutex),
		ctx:    consumerCtx,
		cancel: cancel,
		captureParameters: CaptureParameters{
			LinkType: layers.LinkTypeEthernet,
		},
		handle:       handle,
		packetSource: gopacket.NewPacketSource(handle, layers.LinkTypeEthernet),
		consumer:     consumer,
	}

	go dev.removeConsumerOnContextEnd()

	return
}

type deviceConsumer struct {
	ctx               context.Context
	cancel            context.CancelFunc
	consumer          Consumer
	locker            sync.Locker
	captureParameters CaptureParameters
	handle            *pcapgo.EthernetHandle
	packetSource      *gopacket.PacketSource
}

func (o *deviceConsumer) removeConsumerOnContextEnd() {
	<-o.ctx.Done()

	o.locker.Lock()
	defer o.locker.Unlock()

	_ = o.Close()
}

func (o *deviceConsumer) Close() error {
	o.cancel()
	o.handle.Close()
	if closer, ok := o.consumer.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

func (o *deviceConsumer) StartTransport() {
	go o.transportToConsumers()
}

func (o *deviceConsumer) transportToConsumers() {
	for {
		select {
		case pkg, more := <-o.packetSource.Packets():
			if !more {
				return
			}
			o.locker.Lock()
			o.consumer.Observe(pkg)
			o.locker.Unlock()
		case <-o.ctx.Done():
			return
		}
	}
}