// +build linux
//go:build linux

package pcap

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	_ "github.com/google/gopacket/layers"
	"go.uber.org/multierr"
)

const (
	DefaultReadTimeout = 30 * time.Second
)

var (
	ErrConsumerAlreadyRegistered    = errors.New("consumer with the given name is already registered")
	ErrNoMatchingConsumerRegistered = errors.New("no consumer with given key is registered")
	DefaultRecordingOptions         = RecordingOptions{
		ReadTimeout: DefaultReadTimeout,
		Promiscuous: false,
	}

	_ Recorder = (*recorder)(nil)
)

func NewRecorder() Recorder {
	return &recorder{
		locker:      new(sync.Mutex),
		openDevices: make(map[string]*deviceConsumer),
	}
}

type recorder struct {
	locker      sync.Locker
	openDevices map[string]*deviceConsumer
}

func (r recorder) Subscriptions() (subscriptions []Subscription) {
	r.locker.Lock()
	defer r.locker.Unlock()

	for devName, dev := range r.openDevices {
		sub := Subscription{
			ConsumerKey: devName,
		}
		sub.ConsumerName = dev.consumer.Name()
		subscriptions = append(subscriptions, sub)
	}
	return
}

func (recorder) AvailableDevices() (devices []Device, err error) {
	var ifs []net.Interface
	if ifs, err = net.Interfaces(); err != nil {
		return
	}

	for di := range ifs {
		d := Device{
			Name: ifs[di].Name,
		}

		if addrs, err := ifs[di].Addrs(); err == nil {
			for ai := range addrs {
				if ipAddr, ok := addrs[ai].(*net.IPNet); ok {
					d.IPAddresses = append(d.IPAddresses, ipAddr.IP)
				}
			}
		}

		devices = append(devices, d)
	}
	return
}

func (r recorder) StartRecordingWithOptions(
	ctx context.Context,
	device string,
	consumer Consumer,
	opts RecordingOptions,
) (result *StartRecordingResult, err error) {
	r.locker.Lock()
	defer r.locker.Unlock()

	result = &StartRecordingResult{
		ConsumerKey: fmt.Sprintf("%s:%s", device, consumer.Name()),
	}

	var (
		openDev       *deviceConsumer
		alreadyOpened bool
	)

	if _, alreadyOpened = r.openDevices[result.ConsumerKey]; alreadyOpened {
		return nil, ErrConsumerAlreadyRegistered
	}

	if openDev, err = openDeviceForConsumers(device, consumer, opts); err != nil {
		return nil, err
	}
	openDev.StartTransport(ctx)
	r.openDevices[result.ConsumerKey] = openDev
	go r.removeConsumerOnContextEnd(ctx, result.ConsumerKey)

	return result, nil
}

func (r *recorder) StartRecording(ctx context.Context, device string, consumer Consumer) (result *StartRecordingResult, err error) {
	return r.StartRecordingWithOptions(ctx, device, consumer, DefaultRecordingOptions)
}

func (r *recorder) StopRecording(consumerKey string) error {
	r.locker.Lock()
	defer r.locker.Unlock()

	var (
		dev   *deviceConsumer
		known bool
	)

	if dev, known = r.openDevices[consumerKey]; !known {
		return ErrNoMatchingConsumerRegistered
	}

	delete(r.openDevices, consumerKey)
	return dev.Close()
}

func (r recorder) Close() (err error) {
	r.locker.Lock()
	defer r.locker.Unlock()

	for _, consumer := range r.openDevices {
		err = multierr.Append(err, consumer.Close())
	}
	return
}

func (r *recorder) removeConsumerOnContextEnd(ctx context.Context, consumerKey string) {
	<-ctx.Done()
	_ = r.StopRecording(consumerKey)
}
