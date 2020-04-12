package main

import (
	"fmt"
	"github.com/baez90/inetmock/internal/config"
	"go.uber.org/zap"
	"net/http"
	"sync"
)

const (
	name = "http_mock"
)

type httpHandler struct {
	logger *zap.Logger
	router *RegexpHandler
	server *http.Server
}

func (p *httpHandler) Start(config config.HandlerConfig) {
	options := loadFromConfig(config.Options())
	addr := fmt.Sprintf("%s:%d", config.ListenAddress(), config.Port())
	p.server = &http.Server{Addr: addr, Handler: p.router}
	p.logger = p.logger.With(
		zap.String("address", addr),
	)

	for _, rule := range options.Rules {
		p.setupRoute(rule)
	}

	go p.startServer()
}

func (p *httpHandler) Shutdown(wg *sync.WaitGroup) {
	p.logger.Info("Shutting down HTTP mock")
	if err := p.server.Close(); err != nil {
		p.logger.Error(
			"failed to shutdown HTTP server",
			zap.Error(err),
		)
	}

	wg.Done()
}

func (p *httpHandler) startServer() {
	if err := p.server.ListenAndServe(); err != nil {
		p.logger.Error(
			"failed to start http listener",
			zap.Error(err),
		)
	}
}
