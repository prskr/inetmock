package mock

import (
	"io/fs"

	"inetmock.icb4dc0.de/inetmock/internal/endpoint"
	"inetmock.icb4dc0.de/inetmock/pkg/audit"
	"inetmock.icb4dc0.de/inetmock/pkg/logging"
)

func New(logger logging.Logger, emitter audit.Emitter, fakeFileFS fs.FS) endpoint.ProtocolHandler {
	return &httpHandler{
		logger:     logger,
		fakeFileFS: fakeFileFS,
		emitter:    emitter,
	}
}

func AddHTTPMock(registry endpoint.HandlerRegistry, logger logging.Logger, emitter audit.Emitter, fakeFileFS fs.FS) {
	registry.RegisterHandler(name, func() endpoint.ProtocolHandler {
		return New(logger, emitter, fakeFileFS)
	})
}
