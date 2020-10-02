package http_mock

import (
	"context"
	"errors"
	"fmt"
	"github.com/baez90/inetmock/pkg/config"
	"github.com/baez90/inetmock/pkg/logging"
	"go.uber.org/zap"
	"net/http"
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

func (p *httpHandler) Start(config config.HandlerConfig) (err error) {
	options := loadFromConfig(config.Options)
	p.logger = p.logger.With(
		zap.String("handler_name", config.HandlerName),
		zap.String("address", config.ListenAddr()),
	)

	router := &RegexpHandler{
		logger:      p.logger,
		handlerName: config.HandlerName,
	}
	p.server = &http.Server{Addr: config.ListenAddr(), Handler: router}

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
