package http_proxy

import (
	"context"
	"errors"
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
	var opts httpProxyOptions
	if err = config.Options.Unmarshal(&opts); err != nil {
		return
	}
	listenAddr := config.ListenAddr()
	h.server = &http.Server{Addr: listenAddr, Handler: h.proxy}
	h.logger = h.logger.With(
		zap.String("handler_name", config.HandlerName),
		zap.String("address", listenAddr),
	)

	tlsConfig := api.ServicesInstance().CertStore().TLSConfig()

	proxyHandler := &proxyHttpHandler{
		handlerName: config.HandlerName,
		options:     opts,
		logger:      h.logger,
	}

	proxyHttpsHandler := &proxyHttpsHandler{
		handlerName: config.HandlerName,
		tlsConfig:   tlsConfig,
		logger:      h.logger,
	}

	h.proxy.OnRequest().Do(proxyHandler)
	h.proxy.OnRequest().HandleConnect(proxyHttpsHandler)
	go h.startProxy()
	return
}

func (h *httpProxy) startProxy() {
	if err := h.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		h.logger.Error(
			"failed to start proxy server",
			zap.Error(err),
		)
	}
}

func (h *httpProxy) Shutdown(ctx context.Context) (err error) {
	h.logger.Info("Shutting down HTTP proxy")
	if err = h.server.Shutdown(ctx); err != nil {
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
