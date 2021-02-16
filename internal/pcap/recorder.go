package pcap

import (
	"context"
	"errors"
	"net"
	"sync"
	"time"
)

const (
	defaultSnapshotLength int32 = 1600
	defaultReadTimeout          = 30 * time.Second
)

var (
	ErrConsumerAlreadyRegistered = errors.New("consumer with the given name is already registered")
	WithSnapshotLength           = func(snapshotLength int32) RecorderOption {
		return func(opt recorderOptions) recorderOptions {
			opt.snapshotLength = snapshotLength
			return opt
		}
	}

	WithPromiscuous = func(promiscuous bool) RecorderOption {
		return func(opt recorderOptions) recorderOptions {
			opt.promiscuous = promiscuous
			return opt
		}
	}
	WithReadTimeout = func(readTimeout time.Duration) RecorderOption {
		return func(opt recorderOptions) recorderOptions {
			opt.readTimeout = readTimeout
			return opt
		}
	}
)

type recorderOptions struct {
	snapshotLength int32
	promiscuous    bool
	readTimeout    time.Duration
}

func NewRecorder(options ...RecorderOption) Recorder {
	opts := recorderOptions{
		snapshotLength: defaultSnapshotLength,
		promiscuous:    true,
		readTimeout:    defaultReadTimeout,
	}

	for i := range options {
		opts = options[i](opts)
	}

	return &recorder{
		locker:      new(sync.Mutex),
		opts:        opts,
		openDevices: make(map[string]deviceConsumer),
	}
}

type recorder struct {
	locker      sync.Locker
	opts        recorderOptions
	openDevices map[string]deviceConsumer
}

func (r recorder) Subscriptions() (subscriptions []Subscription) {
	r.locker.Lock()
	r.locker.Unlock()

	for devName, dev := range r.openDevices {
		sub := Subscription{
			Device: devName,
		}
		for name := range dev.consumers {
			sub.Consumers = append(sub.Consumers, name)
		}
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

func (r *recorder) Subscribe(ctx context.Context, device string, consumer Consumer) (err error) {
	r.locker.Lock()
	defer r.locker.Unlock()

	var openDev deviceConsumer
	var alreadyOpened bool
	if openDev, alreadyOpened = r.openDevices[device]; !alreadyOpened || openDev.CleanupOrphaned() {
		if openDev, err = openDeviceForConsumers(device, r.opts); err != nil {
			return
		}
		openDev.StartTransport()
		r.openDevices[device] = openDev
	}

	if err = openDev.AddConsumer(ctx, consumer); err != nil {
		return
	}

	return
}

func (r *recorder) RemoveSubscriptions(device, consumerName string) (removed bool) {
	r.locker.Lock()
	r.locker.Unlock()

	var dev deviceConsumer
	var known bool
	if dev, known = r.openDevices[device]; !known {
		return false
	}

	removed, _ = dev.RemoveConsumer(consumerName)
	return
}
