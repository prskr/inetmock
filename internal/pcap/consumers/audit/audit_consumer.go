package audit

import (
	"bytes"
	"fmt"
	"time"

	"github.com/google/gopacket"

	"gitlab.com/inetmock/inetmock/internal/pcap"
	"gitlab.com/inetmock/inetmock/pkg/audit"
)

type auditConsumer struct {
	name             string
	emitter          audit.Emitter
	knownConnections map[uint64]int64
}

func NewAuditConsumer(name string, emitter audit.Emitter) pcap.Consumer {
	return &auditConsumer{
		name:             name,
		emitter:          emitter,
		knownConnections: make(map[uint64]int64),
	}
}

func (a auditConsumer) Name() string {
	return a.name
}

func (a auditConsumer) Observe(pkg gopacket.Packet) {
	var appLayer gopacket.ApplicationLayer
	if appLayer = pkg.ApplicationLayer(); appLayer == nil {
		return
	}

	connHash := (37 * pkg.NetworkLayer().NetworkFlow().FastHash()) ^ pkg.TransportLayer().TransportFlow().FastHash()

	if _, known := a.knownConnections[connHash]; !known {
		fmt.Printf("new connection - network = %s, transport = %s \n", pkg.NetworkLayer().NetworkFlow().String(), pkg.TransportLayer().TransportFlow().String())
		payload := filterPayload(appLayer.Payload())
		if bytes.Contains(payload, []byte("HTTP")) {
			fmt.Println(string(payload))
			fmt.Println("found HTTP")
		}
	}
	a.knownConnections[connHash] = time.Now().Add(60 * time.Second).Unix()
}

func (a auditConsumer) Init(pcap.CaptureParameters) {
}

func filterPayload(payload []byte) []byte {
	if len(payload) < 4096 {
		return payload
	}
	return payload[:4096]
}
