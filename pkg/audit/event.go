package audit

import (
	"net"
	"strconv"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"gitlab.com/inetmock/inetmock/pkg/audit/details"
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
	ProtocolDetails Details
	TLS             *TLSDetails
}

func (e *Event) ProtoMessage() *EventEntity {
	var tlsDetails *TLSDetailsEntity = nil
	if e.TLS != nil {
		tlsDetails = e.TLS.ProtoMessage()
	}

	var detailsEntity *anypb.Any = nil
	if e.ProtocolDetails != nil {
		if any, err := e.ProtocolDetails.MarshalToWireFormat(); err == nil {
			detailsEntity = any
		}
	}

	return &EventEntity{
		Id:              e.ID,
		Timestamp:       timestamppb.New(e.Timestamp),
		Transport:       e.Transport,
		Application:     e.Application,
		SourceIP:        e.SourceIP,
		DestinationIP:   e.DestinationIP,
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

func (e *Event) SetSourceIPFromAddr(remoteAddr net.Addr) {
	ip, port := parseIPPortFromAddr(remoteAddr)
	e.SourceIP = ip
	e.SourcePort = port
}

func (e *Event) SetDestinationIPFromAddr(localAddr net.Addr) {
	ip, port := parseIPPortFromAddr(localAddr)
	e.DestinationIP = ip
	e.DestinationPort = port
}

func NewEventFromProto(msg *EventEntity) (ev Event) {
	ev = Event{
		ID:              msg.GetId(),
		Timestamp:       msg.GetTimestamp().AsTime(),
		Transport:       msg.GetTransport(),
		Application:     msg.GetApplication(),
		SourceIP:        msg.SourceIP,
		DestinationIP:   msg.DestinationIP,
		SourcePort:      uint16(msg.GetSourcePort()),
		DestinationPort: uint16(msg.GetDestinationPort()),
		ProtocolDetails: guessDetailsFromApp(msg.GetProtocolDetails()),
		TLS:             NewTLSDetailsFromProto(msg.GetTls()),
	}
	return
}

func parseIPPortFromAddr(addr net.Addr) (ip net.IP, port uint16) {
	const expectedIPPortSplitLength = 2
	if addr == nil {
		return
	}
	switch a := addr.(type) {
	case *net.TCPAddr:
		return a.IP, uint16(a.Port)
	case *net.UDPAddr:
		return a.IP, uint16(a.Port)
	case *net.UnixAddr:
		return
	default:
		ipPortSplit := strings.Split(addr.String(), ":")
		if len(ipPortSplit) != expectedIPPortSplitLength {
			return
		}

		ip = net.ParseIP(ipPortSplit[0])
		if p, err := strconv.Atoi(ipPortSplit[1]); err == nil {
			port = uint16(p)
		}
		return
	}
}

func guessDetailsFromApp(any *anypb.Any) Details {
	var detailsProto proto.Message
	var err error
	if detailsProto, err = any.UnmarshalNew(); err != nil {
		return nil
	}
	switch any.TypeUrl {
	case "type.googleapis.com/inetmock.audit.details.HTTPDetailsEntity":
		return details.NewHTTPFromWireFormat(detailsProto.(*details.HTTPDetailsEntity))
	case "type.googleapis.com/inetmock.audit.details.DNSDetailsEntity":
		return details.NewDNSFromWireFormat(detailsProto.(*details.DNSDetailsEntity))
	default:
		return nil
	}
}
