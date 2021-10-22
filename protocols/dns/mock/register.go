package mock

import (
	"gitlab.com/inetmock/inetmock/internal/endpoint"
	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/logging"
)

func New(logger logging.Logger, emitter audit.Emitter) endpoint.ProtocolHandler {
	return &dnsHandler{
		logger:  logger,
		emitter: emitter,
	}
}

func AddDNSMock(registry endpoint.HandlerRegistry, logger logging.Logger, emitter audit.Emitter) {
	registry.RegisterHandler(name, func() endpoint.ProtocolHandler {
		return New(logger, emitter)
	})
}
