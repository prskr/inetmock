package audit

import (
	"net"
	"time"

	"google.golang.org/protobuf/proto"
)

type EventDetails interface {
	ProtoMessage() proto.Message
}

type Event struct {
	ID              int64
	Timestamp       time.Time
	Transport       TransportProtocol
	Application     AppProtocol
	SourceIP        net.IP
	DestinationIP   net.IP
	SourcePort      uint16
	DestinationPort uint16
	ProtocolDetails EventDetails
	TLS             *TLSDetails
}
