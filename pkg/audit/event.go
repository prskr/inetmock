package audit

import (
	"net"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"gitlab.com/inetmock/inetmock/internal/netutils"
	"gitlab.com/inetmock/inetmock/pkg/audit/details"
	auditv1 "gitlab.com/inetmock/inetmock/pkg/audit/v1"
)

type Details interface {
	MarshalToWireFormat() (*anypb.Any, error)
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

	var detailsEntity *anypb.Any = nil
	if e.ProtocolDetails != nil {
		if any, err := e.ProtocolDetails.MarshalToWireFormat(); err == nil {
			detailsEntity = any
		}
	}

	return &auditv1.EventEntity{
		Id:              e.ID,
		Timestamp:       timestamppb.New(e.Timestamp),
		Transport:       e.Transport,
		Application:     e.Application,
		SourceIp:        e.SourceIP,
		DestinationIp:   e.DestinationIP,
		SourcePort:      uint32(e.SourcePort),
		DestinationPort: uint32(e.DestinationPort),
		Tls:             tlsDetails,
		ProtocolDetails: detailsEntity,
	}
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
		ProtocolDetails: guessDetailsFromApp(msg.GetProtocolDetails()),
		TLS:             NewTLSDetailsFromProto(msg.GetTls()),
	}
	return
}

func guessDetailsFromApp(any *anypb.Any) Details {
	var detailsProto proto.Message
	var err error
	if detailsProto, err = any.UnmarshalNew(); err != nil {
		return nil
	}
	switch any.TypeUrl {
	case "type.googleapis.com/inetmock.audit.v1.HTTPDetailsEntity":
		return details.NewHTTPFromWireFormat(detailsProto.(*auditv1.HTTPDetailsEntity))
	case "type.googleapis.com/inetmock.audit.v1.DNSDetailsEntity":
		return details.NewDNSFromWireFormat(detailsProto.(*auditv1.DNSDetailsEntity))
	default:
		return nil
	}
}
