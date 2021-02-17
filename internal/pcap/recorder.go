package pcap

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	_ "github.com/google/gopacket/layers"
)

const (
	DefaultReadTimeout = 30 * time.Second
)

var (
	ErrConsumerAlreadyRegistered = errors.New("consumer with the given name is already registered")
	DefaultRecordingOptions      = RecordingOptions{
		ReadTimeout: DefaultReadTimeout,
		Promiscuous: false,
	}
)

func NewRecorder() Recorder {
	return &recorder{
		locker:      new(sync.Mutex),
		openDevices: make(map[string]deviceConsumer),
	}
}

type recorder struct {
	locker      sync.Locker
	openDevices map[string]deviceConsumer
}

func (r recorder) Subscriptions() (subscriptions []Subscription) {
	r.locker.Lock()
	r.locker.Unlock()

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

func (r recorder) StartRecordingWithOptions(ctx context.Context, device string, consumer Consumer, opts RecordingOptions) (err error) {
	r.locker.Lock()
	defer r.locker.Unlock()

	consumerKey := fmt.Sprintf("%s:%s", device, consumer.Name())
	var openDev deviceConsumer
	var alreadyOpened bool
	if openDev, alreadyOpened = r.openDevices[consumerKey]; alreadyOpened {
		err = ErrConsumerAlreadyRegistered
		return
	}

	if openDev, err = openDeviceForConsumers(ctx, device, consumer, opts); err != nil {
		return
	}
	openDev.StartTransport()
	r.openDevices[consumerKey] = openDev

	return
}

func (r *recorder) StartRecording(ctx context.Context, device string, consumer Consumer) (err error) {
	return r.StartRecordingWithOptions(ctx, device, consumer, DefaultRecordingOptions)
}

func (r *recorder) StopRecording(consumerKey string) (err error) {
	r.locker.Lock()
	r.locker.Unlock()

	var dev deviceConsumer
	var known bool
	if dev, known = r.openDevices[consumerKey]; !known {
		return nil
	}

	delete(r.openDevices, consumerKey)

	err = dev.Close()

	return
}
