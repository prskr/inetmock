package mock

import (
	"inetmock.icb4dc0.de/inetmock/internal/endpoint"
	"inetmock.icb4dc0.de/inetmock/pkg/audit"
	"inetmock.icb4dc0.de/inetmock/pkg/logging"
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
