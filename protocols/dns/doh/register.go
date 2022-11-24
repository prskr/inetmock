package doh

import (
	"inetmock.icb4dc0.de/inetmock/internal/endpoint"
	"inetmock.icb4dc0.de/inetmock/pkg/audit"
	"inetmock.icb4dc0.de/inetmock/pkg/logging"
)

const name = "doh_mock"

func New(logger logging.Logger, emitter audit.Emitter) *dohHandler {
	return &dohHandler{
		logger:  logger,
		emitter: emitter,
	}
}

func AddDoH(registry endpoint.HandlerRegistry, logger logging.Logger, emitter audit.Emitter) {
	registry.RegisterHandler(name, func() endpoint.ProtocolHandler {
		return New(logger, emitter)
	})
}
