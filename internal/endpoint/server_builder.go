package endpoint

import (
	"crypto/tls"
	"errors"
	"sync"

	"gitlab.com/inetmock/inetmock/pkg/logging"
)

var (
	ErrUnknownHandlerRef = errors.New("no handler for given key registered")
	ErrNoEndpoints       = errors.New("no endpoints configured in ListenerGroup")
)

func NewServerBuilder(
	tlsConfig *tls.Config,
	registry HandlerRegistry,
	logger logging.Logger,
) *ServerBuilder {
	return &ServerBuilder{
		registry: registry,
		server:   NewServer(logger, tlsConfig),
	}
}

type ServerBuilder struct {
	server   *Server
	lock     sync.Mutex
	registry HandlerRegistry
}

func (e *ServerBuilder) Server() *Server {
	return e.server
}

func (e *ServerBuilder) ConfigureGroup(spec ListenerSpec) (err error) {
	e.lock.Lock()
	defer e.lock.Unlock()

	if len(spec.Endpoints) < 1 {
		return ErrNoEndpoints
	}

	var grp *ListenerGroup
	if grp, err = NewListenerGroup(spec); err != nil {
		return err
	}

	for name, s := range spec.Endpoints {
		if handler, registered := e.registry.HandlerForName(s.HandlerRef); registered {
			grp.ConfigureEndpoint(name, NewListenerEndpoint(s, handler))
		} else {
			return ErrUnknownHandlerRef
		}
	}

	e.server.ConfigureGroup(grp)
	return nil
}

func (e *ServerBuilder) ConfiguredGroups() []GroupInfo {
	e.lock.Lock()
	defer e.lock.Unlock()

	return e.server.ConfiguredGroups()
}
