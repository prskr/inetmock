package audit

import (
	"net"
	"reflect"
	"sync"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"inetmock.icb4dc0.de/inetmock/internal/netutils"
	auditv1 "inetmock.icb4dc0.de/inetmock/pkg/audit/v1"
)

var (
	mappingLock     sync.Mutex
	wire2AppMapping = make(map[reflect.Type]func(msg *auditv1.EventEntity) Details)
	eventPool       = sync.Pool{
		New: func() any {
			return new(Event)
		},
	}
	eventBuilderPool = sync.Pool{
		New: func() any {
			return new(eventBuilder)
		},
	}
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
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
}

func (e *Event) SetSourceIPFromAddr(remoteAddr net.Addr) error {
	if ip, port, err := netutils.IPPortFromAddress(remoteAddr); err != nil {
		return err
	} else {
		e.SourceIP = ip
		e.SourcePort = uint16(port)
	}
	return nil
}

func (e *Event) SetDestinationIPFromAddr(localAddr net.Addr) error {
	if ip, port, err := netutils.IPPortFromAddress(localAddr); err != nil {
		return err
	} else {
		e.DestinationIP = ip
		e.DestinationPort = uint16(port)
	}
	return nil
}

func (e *Event) Dispose() {
	eventPool.Put(e)
}

func (e *Event) Reset() {
	e.ID = 0
	e.Timestamp = time.Time{}
	e.SourceIP = nil
	e.DestinationIP = nil
	e.SourcePort = 0
	e.DestinationPort = 0
	e.TLS = nil
	e.Transport = auditv1.TransportProtocol_TRANSPORT_PROTOCOL_UNSPECIFIED
	e.ProtocolDetails = nil
}

func NewEventFromProto(msg *auditv1.EventEntity) (ev *Event) {
	ev = eventFromPool()
	ev.ID = msg.GetId()
	ev.Timestamp = msg.GetTimestamp().AsTime()
	ev.Transport = msg.GetTransport()
	ev.Application = msg.GetApplication()
	ev.SourceIP = msg.SourceIp
	ev.DestinationIP = msg.DestinationIp
	ev.SourcePort = uint16(msg.GetSourcePort())
	ev.DestinationPort = uint16(msg.GetDestinationPort())
	ev.TLS = NewTLSDetailsFromProto(msg.GetTls())
	ev.ProtocolDetails = unwrapDetails(msg)

	return
}

func unwrapDetails(msg *auditv1.EventEntity) Details {
	if mapping, ok := wire2AppMapping[reflect.TypeOf(msg.ProtocolDetails)]; ok {
		return mapping(msg)
	}
	return nil
}

func eventFromPool() *Event {
	ev := eventPool.Get().(*Event)
	ev.Reset()
	return ev
}
