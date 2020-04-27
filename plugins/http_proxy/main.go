package main

import (
	"fmt"
	"github.com/baez90/inetmock/pkg/api"
	"github.com/baez90/inetmock/pkg/config"
	"github.com/baez90/inetmock/pkg/logging"
	"go.uber.org/zap"
	"gopkg.in/elazarl/goproxy.v1"
	"net/http"
)

const (
	name = "http_proxy"
)

type httpProxy struct {
	logger logging.Logger
	proxy  *goproxy.ProxyHttpServer
	server *http.Server
}

func (h *httpProxy) Start(config config.HandlerConfig) (err error) {
	options := loadFromConfig(config.Options)
	addr := fmt.Sprintf("%s:%d", config.ListenAddress, config.Port)
	h.server = &http.Server{Addr: addr, Handler: h.proxy}
	h.logger = h.logger.With(
		zap.String("address", addr),
	)

	tlsConfig := api.ServicesInstance().CertStore().TLSConfig()

	proxyHandler := &proxyHttpHandler{
		options: options,
		logger:  h.logger,
	}

	proxyHttpsHandler := &proxyHttpsHandler{
		tlsConfig: tlsConfig,
		logger:    h.logger,
	}

	h.proxy.OnRequest().Do(proxyHandler)
	h.proxy.OnRequest().HandleConnect(proxyHttpsHandler)
	go h.startProxy()
	return
}

func (h *httpProxy) startProxy() {
	if err := h.server.ListenAndServe(); err != nil {
		h.logger.Error(
			"failed to start proxy server",
			zap.Error(err),
		)
	}
}

func (h *httpProxy) Shutdown() (err error) {
	h.logger.Info("Shutting down HTTP proxy")
	if err = h.server.Close(); err != nil {
		h.logger.Error(
			"failed to shutdown proxy endpoint",
			zap.Error(err),
		)

		err = fmt.Errorf(
			"failed to shutdown proxy endpoint: %w",
			err,
		)
	}
	return
}
