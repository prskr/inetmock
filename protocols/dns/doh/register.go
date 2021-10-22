package doh

import (
	"gitlab.com/inetmock/inetmock/internal/endpoint"
	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/logging"
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
