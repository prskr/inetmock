package mock

import (
	"context"
	"io/fs"
	"net"
	"net/http"
	"time"

	"github.com/soheilhy/cmux"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"inetmock.icb4dc0.de/inetmock/internal/endpoint"
	"inetmock.icb4dc0.de/inetmock/multiplexing"
	"inetmock.icb4dc0.de/inetmock/pkg/audit"
	auditv1 "inetmock.icb4dc0.de/inetmock/pkg/audit/v1"
	"inetmock.icb4dc0.de/inetmock/pkg/logging"
)

const (
	name                     = "http_mock"
	defaultReadHeaderTimeout = 100 * time.Millisecond
)

type httpHandler struct {
	logger     logging.Logger
	fakeFileFS fs.FS
	server     *http.Server
	emitter    audit.Emitter
}

func (p *httpHandler) Matchers() []cmux.Matcher {
	return []cmux.Matcher{multiplexing.HTTP()}
}

func (p *httpHandler) Start(_ context.Context, startupSpec *endpoint.StartupSpec) error {
	p.logger = p.logger.With(
		zap.String("protocol_handler", name),
	)

	var (
		options httpOptions
		err     error
	)

	if options, err = loadFromConfig(startupSpec); err != nil {
		return err
	}

	p.logger = p.logger.With(
		zap.String("address", startupSpec.Addr.String()),
	)

	router := &Router{
		HandlerName: startupSpec.Name,
		Logger:      p.logger,
		FakeFileFS:  p.fakeFileFS,
	}

	p.server = &http.Server{
		Handler:           h2c.NewHandler(audit.EmittingHandler(p.emitter, auditv1.AppProtocol_APP_PROTOCOL_HTTP, router), new(http2.Server)),
		ConnContext:       audit.StoreConnPropertiesInContext,
		ReadHeaderTimeout: defaultReadHeaderTimeout,
	}

	for idx := range options.Rules {
		rule := options.Rules[idx]
		if err = router.RegisterRule(rule); err != nil {
			p.logger.Error("failed to setup rule", zap.String("raw_rule", rule), zap.Error(err))
			return err
		}
	}

	go p.startServer(startupSpec.Listener)
	return nil
}

func (p *httpHandler) startServer(listener net.Listener) {
	if err := endpoint.IgnoreShutdownError(p.server.Serve(listener)); err != nil {
		p.logger.Error("Failed to start HTTP listener", zap.Error(err))
	}
}
