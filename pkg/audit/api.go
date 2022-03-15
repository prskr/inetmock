//go:generate mockgen -source=$GOFILE -destination=./../../internal/mock/audit/audit.mock.go -package=audit_mock

package audit

import (
	"context"
	"errors"
	"io"
	"net"
	"sync"

	"gitlab.com/inetmock/inetmock/internal/netutils"
	auditv1 "gitlab.com/inetmock/inetmock/pkg/audit/v1"
)

var (
	ErrSinkAlreadyRegistered = errors.New("sink with same name already registered")
	eventPool                = sync.Pool{
		New: func() any {
			return new(Event)
		},
	}
	_ EventBuilder = (*eventBuilder)(nil)
	_ Emitter      = (*EmitterFunc)(nil)
)

type (
	EventBuilder interface {
		WithTransport(transport auditv1.TransportProtocol) EventBuilder
		WithApplication(app auditv1.AppProtocol) EventBuilder
		WithSource(ip net.IP, port uint16) EventBuilder
		WithSourceFromAddr(addr net.Addr) (EventBuilder, error)
		WithDestination(ip net.IP, port uint16) EventBuilder
		WithDestinationFromAddr(addr net.Addr) (EventBuilder, error)
		WithProtocolDetails(details Details) EventBuilder
		WithTLSDetails(details *TLSDetails) EventBuilder
		Emit()
	}

	Emitter interface {
		Emit(ev *Event)
		Builder() EventBuilder
	}

	Sink interface {
		Name() string
		OnSubscribe(evs <-chan *Event)
	}

	EventStream interface {
		io.Closer
		Emitter
		RegisterSink(ctx context.Context, s Sink) error
		Sinks() []string
		RemoveSink(name string) (exists bool)
	}

	EmitterFunc func(ev *Event)

	eventBuilder struct {
		ev      *Event
		emitter Emitter
	}
)

func (ef EmitterFunc) Builder() EventBuilder {
	return &eventBuilder{
		ev:      eventPool.Get().(*Event),
		emitter: ef,
	}
}

func (e *eventBuilder) WithTransport(transport auditv1.TransportProtocol) EventBuilder {
	e.ev.Transport = transport
	return e
}

func (e *eventBuilder) WithApplication(app auditv1.AppProtocol) EventBuilder {
	e.ev.Application = app
	return e
}

func (e *eventBuilder) WithSource(ip net.IP, port uint16) EventBuilder {
	e.ev.SourceIP = ip
	e.ev.SourcePort = port
	return e
}

func (e *eventBuilder) WithSourceFromAddr(addr net.Addr) (EventBuilder, error) {
	if ipPort, err := netutils.IPPortFromAddress(addr); err != nil {
		return e, err
	} else {
		e.ev.SourceIP = ipPort.IP
		e.ev.SourcePort = uint16(ipPort.Port)
	}
	return e, nil
}

func (e *eventBuilder) WithDestination(ip net.IP, port uint16) EventBuilder {
	e.ev.DestinationIP = ip
	e.ev.DestinationPort = port
	return e
}

func (e *eventBuilder) WithDestinationFromAddr(addr net.Addr) (EventBuilder, error) {
	if ipPort, err := netutils.IPPortFromAddress(addr); err != nil {
		return e, err
	} else {
		e.ev.DestinationIP = ipPort.IP
		e.ev.DestinationPort = uint16(ipPort.Port)
	}
	return e, nil
}

func (e *eventBuilder) WithProtocolDetails(details Details) EventBuilder {
	e.ev.ProtocolDetails = details
	return e
}

func (e *eventBuilder) WithTLSDetails(details *TLSDetails) EventBuilder {
	e.ev.TLS = details
	return e
}

func (e *eventBuilder) Emit() {
	e.emitter.Emit(e.ev)
}

func (ef EmitterFunc) Emit(ev *Event) {
	ef(ev)
}
