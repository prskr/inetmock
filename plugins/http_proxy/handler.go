package http_proxy

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"gitlab.com/inetmock/inetmock/pkg/api"
	"gitlab.com/inetmock/inetmock/pkg/config"
	"gitlab.com/inetmock/inetmock/pkg/logging"
	"go.uber.org/zap"
	"gopkg.in/elazarl/goproxy.v1"
)

const (
	name = "http_proxy"
)

type httpProxy struct {
	logger logging.Logger
	proxy  *goproxy.ProxyHttpServer
	server *http.Server
}

func (h *httpProxy) Start(ctx api.PluginContext, cfg config.HandlerConfig) (err error) {
	var opts httpProxyOptions
	if err = cfg.Options.Unmarshal(&opts); err != nil {
		return
	}
	listenAddr := cfg.ListenAddr()
	h.server = &http.Server{Addr: listenAddr, Handler: h.proxy}
	h.logger = h.logger.With(
		zap.String("handler_name", cfg.HandlerName),
		zap.String("address", listenAddr),
	)

	tlsConfig := ctx.CertStore().TLSConfig()

	proxyHandler := &proxyHttpHandler{
		handlerName: cfg.HandlerName,
		options:     opts,
		logger:      h.logger,
	}

	proxyHttpsHandler := &proxyHttpsHandler{
		handlerName: cfg.HandlerName,
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
