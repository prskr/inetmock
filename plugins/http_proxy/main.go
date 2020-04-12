package main

import (
	"fmt"
	"github.com/baez90/inetmock/internal/config"
	"go.uber.org/zap"
	"gopkg.in/elazarl/goproxy.v1"
	"net/http"
	"sync"
)

const (
	name = "http_proxy"
)

type httpProxy struct {
	logger *zap.Logger
	proxy  *goproxy.ProxyHttpServer
	server *http.Server
}

func (h *httpProxy) Start(config config.HandlerConfig) {
	options := loadFromConfig(config.Options())
	addr := fmt.Sprintf("%s:%d", config.ListenAddress(), config.Port())
	h.server = &http.Server{Addr: addr, Handler: h.proxy}
	h.logger = h.logger.With(
		zap.String("address", addr),
	)

	proxyHandler := &proxyHttpHandler{
		options: options,
		logger:  h.logger,
	}
	h.proxy.OnRequest().Do(proxyHandler)
	go h.startProxy()
}

func (h *httpProxy) startProxy() {
	if err := h.server.ListenAndServe(); err != nil {
		h.logger.Error(
			"failed to start proxy server",
			zap.Error(err),
		)
	}
}

func (h *httpProxy) Shutdown(wg *sync.WaitGroup) {
	defer wg.Done()
	h.logger.Info("Shutting down HTTP proxy")
	if err := h.server.Close(); err != nil {
		h.logger.Error(
			"failed to shutdown proxy endpoint",
			zap.Error(err),
		)
	}
}
