package mock

import (
	"context"
	"errors"
	"io/fs"
	"net"
	"net/http"

	"github.com/soheilhy/cmux"
	"go.uber.org/zap"

	"gitlab.com/inetmock/inetmock/internal/endpoint"
	imHttp "gitlab.com/inetmock/inetmock/internal/endpoint/handler/http"
	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/logging"
)

const (
	name               = "http_mock"
	handlerNameLblName = "handler_name"
	ruleMatchedLblName = "rule_matched"
)

type httpHandler struct {
	logger     logging.Logger
	fakeFileFS fs.FS
	server     *http.Server
	emitter    audit.Emitter
}

func (p *httpHandler) Matchers() []cmux.Matcher {
	return []cmux.Matcher{cmux.HTTP1()}
}

func (p *httpHandler) Start(ctx context.Context, lifecycle endpoint.Lifecycle) error {
	p.logger = p.logger.With(
		zap.String("protocol_handler", name),
	)

	var err error
	var options httpOptions
	if options, err = loadFromConfig(lifecycle); err != nil {
		return err
	}

	p.logger = p.logger.With(
		zap.String("address", lifecycle.Uplink().Addr().String()),
	)

	router := &RegexHandler{
		handlerName: lifecycle.Name(),
		logger:      p.logger,
		emitter:     p.emitter,
		fakeFileFS:  p.fakeFileFS,
	}
	p.server = &http.Server{
		Handler:     router,
		ConnContext: imHttp.StoreConnPropertiesInContext,
	}

	for _, rule := range options.Rules {
		router.AddRouteRule(rule)
	}

	go p.startServer(lifecycle.Uplink().Listener)
	go p.shutdownOnCancel(ctx)
	return nil
}

func (p *httpHandler) shutdownOnCancel(ctx context.Context) {
	<-ctx.Done()
	p.logger.Info("Shutting down HTTP mock")
	if err := p.server.Close(); err != nil {
		p.logger.Error(
			"failed to shutdown HTTP server",
			zap.Error(err),
		)
	}
}

func (p *httpHandler) startServer(listener net.Listener) {
	defer func() {
		if err := listener.Close(); err != nil {
			p.logger.Warn("failed to close listener", zap.Error(err))
		}
	}()
	if err := p.server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		p.logger.Error(
			"failed to start http listener",
			zap.Error(err),
		)
	}
}
