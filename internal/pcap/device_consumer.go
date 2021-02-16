package pcap

import (
	"context"
	"io"
	"sync"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
)

func openDeviceForConsumers(device string, opts recorderOptions) (dev deviceConsumer, err error) {
	var handle *pcapgo.EthernetHandle
	if handle, err = pcapgo.NewEthernetHandle(device); err != nil {
		return
	}

	if err = handle.SetPromiscuous(opts.promiscuous); err != nil {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	dev = deviceConsumer{
		locker: new(sync.Mutex),
		ctx:    ctx,
		cancel: cancel,
		captureParameters: CaptureParameters{
			SnapshotLength: opts.snapshotLength,
			LinkType:       layers.LinkTypeEthernet,
		},
		handle:       handle,
		packetSource: gopacket.NewPacketSource(handle, layers.LinkTypeEthernet),
		consumers:    make(map[string]Consumer),
	}
	return
}

type deviceConsumer struct {
	ctx               context.Context
	cancel            context.CancelFunc
	locker            sync.Locker
	captureParameters CaptureParameters
	handle            *pcapgo.EthernetHandle
	packetSource      *gopacket.PacketSource
	consumers         map[string]Consumer
}

func (o deviceConsumer) CleanupOrphaned() bool {
	if o.ctx.Err() != nil {
		o.handle.Close()
		return true
	}
	return false
}

func (o *deviceConsumer) AddConsumer(ctx context.Context, consumer Consumer) error {
	o.locker.Lock()
	defer o.locker.Unlock()

	if _, alreadyPresent := o.consumers[consumer.Name()]; alreadyPresent {
		return ErrConsumerAlreadyRegistered
	}

	consumer.Init(o.captureParameters)
	o.consumers[consumer.Name()] = consumer

	go o.removeConsumerOnContextEnd(ctx, consumer)
	return nil
}

func (o *deviceConsumer) removeConsumerOnContextEnd(ctx context.Context, consumer Consumer) {
	<-ctx.Done()
	o.locker.Lock()
	defer o.locker.Unlock()
	if _, stillPresent := o.consumers[consumer.Name()]; !stillPresent {
		return
	}

	delete(o.consumers, consumer.Name())

	if len(o.consumers) == 0 {
		o.cancel()
	}
}

func (o *deviceConsumer) RemoveConsumer(name string) (existed bool, consumerClosed bool) {
	o.locker.Lock()
	defer o.locker.Unlock()

	var consumer Consumer
	consumer, existed = o.consumers[name]

	if existed {
		delete(o.consumers, name)
		if closer, ok := consumer.(io.Closer); ok {
			_ = closer.Close()
		}
	}

	if len(o.consumers) == 0 {
		o.cancel()
		o.handle.Close()
		consumerClosed = true
	}
	return
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
			for _, consumer := range o.consumers {
				consumer.Observe(pkg)
			}
			o.locker.Unlock()
		case <-o.ctx.Done():
			return
		}
	}
}
