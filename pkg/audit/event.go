package audit

import (
	"net"
	"reflect"
	"sync"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"gitlab.com/inetmock/inetmock/internal/netutils"
	auditv1 "gitlab.com/inetmock/inetmock/pkg/audit/v1"
)

var (
	mappingLock     sync.Mutex
	wire2AppMapping = make(map[reflect.Type]func(msg *auditv1.EventEntity) Details)
)

func AddMapping(t reflect.Type, mapper func(msg *auditv1.EventEntity) Details) {
	mappingLock.Lock()
	defer mappingLock.Unlock()

	wire2AppMapping[t] = mapper
}

type Details interface {
	AddToMsg(msg *auditv1.EventEntity)
}

type Event struct {
	ID              int64
	Timestamp       time.Time
	Transport       auditv1.TransportProtocol
	Application     auditv1.AppProtocol
	SourceIP        net.IP
	DestinationIP   net.IP
	SourcePort      uint16
	DestinationPort uint16
	ProtocolDetails Details
	TLS             *TLSDetails
}

func (e *Event) ProtoMessage() *auditv1.EventEntity {
	var tlsDetails *auditv1.TLSDetailsEntity = nil
	if e.TLS != nil {
		tlsDetails = e.TLS.ProtoMessage()
	}

	msg := &auditv1.EventEntity{
		Id:              e.ID,
		Timestamp:       timestamppb.New(e.Timestamp),
		Transport:       e.Transport,
		Application:     e.Application,
		SourceIp:        e.SourceIP,
		DestinationIp:   e.DestinationIP,
		SourcePort:      uint32(e.SourcePort),
		DestinationPort: uint32(e.DestinationPort),
		Tls:             tlsDetails,
	}

	if e.ProtocolDetails != nil {
		e.ProtocolDetails.AddToMsg(msg)
	}

	return msg
}

func (e *Event) ApplyDefaults(id int64) {
	e.ID = id
	emptyTime := time.Time{}
	if e.Timestamp == emptyTime {
		e.Timestamp = time.Now().UTC()
	}
}

func (e *Event) SetSourceIPFromAddr(remoteAddr net.Addr) error {
	if ipPort, err := netutils.IPPortFromAddress(remoteAddr); err != nil {
		return err
	} else {
		e.SourceIP = ipPort.IP
		e.SourcePort = uint16(ipPort.Port)
	}
	return nil
}

func (e *Event) SetDestinationIPFromAddr(localAddr net.Addr) error {
	if ipPort, err := netutils.IPPortFromAddress(localAddr); err != nil {
		return err
	} else {
		e.DestinationIP = ipPort.IP
		e.DestinationPort = uint16(ipPort.Port)
	}
	return nil
}

func (e *Event) Dispose() {
	eventPool.Put(e)
}

func NewEventFromProto(msg *auditv1.EventEntity) (ev Event) {
	ev = Event{
		ID:              msg.GetId(),
		Timestamp:       msg.GetTimestamp().AsTime(),
		Transport:       msg.GetTransport(),
		Application:     msg.GetApplication(),
		SourceIP:        msg.SourceIp,
		DestinationIP:   msg.DestinationIp,
		SourcePort:      uint16(msg.GetSourcePort()),
		DestinationPort: uint16(msg.GetDestinationPort()),
		TLS:             NewTLSDetailsFromProto(msg.GetTls()),
		ProtocolDetails: unwrapDetails(msg),
	}

	return
}

func unwrapDetails(msg *auditv1.EventEntity) Details {
	if mapping, ok := wire2AppMapping[reflect.TypeOf(msg.ProtocolDetails)]; ok {
		return mapping(msg)
	}
	return nil
}
