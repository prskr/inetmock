package http_mock

import (
	"fmt"
	"github.com/baez90/inetmock/pkg/config"
	"github.com/baez90/inetmock/pkg/logging"
	"go.uber.org/zap"
	"net/http"
)

const (
	name = "http_mock"
)

type httpHandler struct {
	logger logging.Logger
	router *RegexpHandler
	server *http.Server
}

func (p *httpHandler) Start(config config.HandlerConfig) (err error) {
	options := loadFromConfig(config.Options)
	addr := fmt.Sprintf("%s:%d", config.ListenAddress, config.Port)
	p.server = &http.Server{Addr: addr, Handler: p.router}
	p.logger = p.logger.With(
		zap.String("address", addr),
	)

	for _, rule := range options.Rules {
		p.setupRoute(rule)
	}

	go p.startServer()
	return
}

func (p *httpHandler) Shutdown() (err error) {
	p.logger.Info("Shutting down HTTP mock")
	if err = p.server.Close(); err != nil {
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
	if err := p.server.ListenAndServe(); err != nil {
		p.logger.Error(
			"failed to start http listener",
			zap.Error(err),
		)
	}
}
