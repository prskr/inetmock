package audit

import (
	"encoding/binary"
	"math/big"
	"net"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Details interface {
	MarshalToWireFormat() (*anypb.Any, error)
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
	ProtocolDetails *anypb.Any
	TLS             *TLSDetails
}

func (e *Event) ProtoMessage() proto.Message {
	var sourceIP isEventEntity_SourceIP
	if ipv4 := e.SourceIP.To4(); ipv4 != nil {
		if len(ipv4) == 16 {
			sourceIP = &EventEntity_SourceIPv4{SourceIPv4: binary.BigEndian.Uint32(ipv4[12:16])}
		} else {
			sourceIP = &EventEntity_SourceIPv4{SourceIPv4: binary.BigEndian.Uint32(ipv4)}
		}
	} else {
		ipv6 := big.NewInt(0)
		ipv6.SetBytes(e.SourceIP)
		sourceIP = &EventEntity_SourceIPv6{SourceIPv6: ipv6.Uint64()}
	}

	var destinationIP isEventEntity_DestinationIP
	if ipv4 := e.DestinationIP.To4(); ipv4 != nil {
		if len(ipv4) == 16 {
			destinationIP = &EventEntity_DestinationIPv4{DestinationIPv4: binary.BigEndian.Uint32(ipv4[12:16])}
		} else {
			destinationIP = &EventEntity_DestinationIPv4{DestinationIPv4: binary.BigEndian.Uint32(ipv4)}
		}
	} else {
		ipv6 := big.NewInt(0)
		ipv6.SetBytes(e.SourceIP)
		destinationIP = &EventEntity_DestinationIPv6{DestinationIPv6: ipv6.Uint64()}
	}

	var tlsDetails *TLSDetailsEntity = nil
	if e.TLS != nil {
		tlsDetails = e.TLS.ProtoMessage()
	}

	return &EventEntity{
		Id:              e.ID,
		Timestamp:       timestamppb.New(e.Timestamp),
		Transport:       e.Transport,
		Application:     e.Application,
		SourceIP:        sourceIP,
		DestinationIP:   destinationIP,
		SourcePort:      uint32(e.SourcePort),
		DestinationPort: uint32(e.DestinationPort),
		Tls:             tlsDetails,
		ProtocolDetails: e.ProtocolDetails,
	}
}

func (e *Event) ApplyDefaults(id int64) {
	e.ID = id
	emptyTime := time.Time{}
	if e.Timestamp == emptyTime {
		e.Timestamp = time.Now().UTC()
	}
}

func NewEventFromProto(msg *EventEntity) (ev Event) {
	var sourceIP net.IP
	switch ip := msg.GetSourceIP().(type) {
	case *EventEntity_SourceIPv4:
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, ip.SourceIPv4)
		sourceIP = buf
		sourceIP = sourceIP.To4()
	case *EventEntity_SourceIPv6:
		sourceIP = big.NewInt(int64(ip.SourceIPv6)).Bytes()
	}

	var destinationIP net.IP
	switch ip := msg.GetDestinationIP().(type) {
	case *EventEntity_DestinationIPv4:
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, ip.DestinationIPv4)
		destinationIP = buf
		destinationIP = destinationIP.To4()
	case *EventEntity_DestinationIPv6:
		destinationIP = big.NewInt(int64(ip.DestinationIPv6)).Bytes()
	}

	ev = Event{
		ID:              msg.GetId(),
		Timestamp:       msg.GetTimestamp().AsTime(),
		Transport:       msg.GetTransport(),
		Application:     msg.GetApplication(),
		SourceIP:        sourceIP,
		DestinationIP:   destinationIP,
		SourcePort:      uint16(msg.GetSourcePort()),
		DestinationPort: uint16(msg.GetDestinationPort()),
		ProtocolDetails: msg.GetProtocolDetails(),
		TLS:             NewTLSDetailsFromProto(msg.GetTls()),
	}
	return
}
