// +build linux

package consumers

import (
	"github.com/google/gopacket"
	"github.com/google/uuid"

	"gitlab.com/inetmock/inetmock/internal/pcap"
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
