//go:build linux

package consumers

import (
	"github.com/google/uuid"
	"github.com/gopacket/gopacket"

	"inetmock.icb4dc0.de/inetmock/internal/pcap"
)

var _ pcap.Consumer = (*noopConsumer)(nil)

type noopConsumer struct {
	name string
}

func NewNoOpConsumer() pcap.Consumer {
	return NewNoOpConsumerWithName(uuid.NewString())
}

func NewNoOpConsumerWithName(name string) pcap.Consumer {
	return &noopConsumer{
		name: name,
	}
}

func (n noopConsumer) Name() string {
	return n.name
}

func (n noopConsumer) Observe(gopacket.Packet) {
}

func (n noopConsumer) Init() error {
	return nil
}
