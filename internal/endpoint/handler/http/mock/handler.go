package mock

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"gitlab.com/inetmock/inetmock/pkg/api"
	"gitlab.com/inetmock/inetmock/pkg/config"
	"gitlab.com/inetmock/inetmock/pkg/logging"
	"go.uber.org/zap"
)

const (
	name               = "http_mock"
	handlerNameLblName = "handler_name"
	ruleMatchedLblName = "rule_matched"
)

type httpHandler struct {
	logger logging.Logger
	server *http.Server
}

func (p *httpHandler) Start(ctx api.PluginContext, config config.HandlerConfig) (err error) {
	p.logger = ctx.Logger().With(
		zap.String("protocol_handler", name),
	)

	var options httpOptions
	if options, err = loadFromConfig(config.Options); err != nil {
		return
	}

	p.logger = p.logger.With(
		zap.String("handler_name", config.HandlerName),
		zap.String("address", config.ListenAddr()),
	)

	router := &RegexpHandler{
		logger:      p.logger,
		emitter:     ctx.Audit(),
		handlerName: config.HandlerName,
	}
	p.server = &http.Server{
		Addr:        config.ListenAddr(),
		Handler:     router,
		ConnContext: StoreConnPropertiesInContext,
	}

	for _, rule := range options.Rules {
		router.setupRoute(rule)
	}

	go p.startServer()
	return
}

func (p *httpHandler) Shutdown(ctx context.Context) (err error) {
	p.logger.Info("Shutting down HTTP mock")
	if err = p.server.Shutdown(ctx); err != nil {
		p.logger.Error(
			"failed to shutdown HTTP server",
			zap.Error(err),
		)
		err = fmt.Errorf(
			"failed to shutdown HTTP server: %w",
			err,
		)
	}
	return
}

func (p *httpHandler) startServer() {
	if err := p.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		p.logger.Error(
			"failed to start http listener",
			zap.Error(err),
		)
	}
}
