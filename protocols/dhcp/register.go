package dhcp

import (
	"inetmock.icb4dc0.de/inetmock/internal/endpoint"
	"inetmock.icb4dc0.de/inetmock/internal/state"
	"inetmock.icb4dc0.de/inetmock/pkg/audit"
	"inetmock.icb4dc0.de/inetmock/pkg/logging"
)

func New(logger logging.Logger, emitter audit.Emitter, stateStore state.KVStore) endpoint.ProtocolHandler {
	return &dhcpHandler{
		logger:     logger,
		emitter:    emitter,
		stateStore: stateStore,
	}
}

func AddDHCPMock(registry endpoint.HandlerRegistry, logger logging.Logger, emitter audit.Emitter, stateStore state.KVStore) {
	registry.RegisterHandler(name, func() endpoint.ProtocolHandler {
		return New(logger, emitter, stateStore)
	})
}
