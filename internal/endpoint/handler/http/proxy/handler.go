package proxy

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/soheilhy/cmux"
	"gitlab.com/inetmock/inetmock/internal/endpoint"
	imHttp "gitlab.com/inetmock/inetmock/internal/endpoint/handler/http"
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

func (h *httpProxy) Matchers() []cmux.Matcher {
	return []cmux.Matcher{cmux.HTTP1()}
}

func (h *httpProxy) Start(lifecycle endpoint.Lifecycle) (err error) {
	var opts httpProxyOptions
	if err = lifecycle.UnmarshalOptions(&opts); err != nil {
		return
	}

	h.server = &http.Server{
		Handler:     h.proxy,
		ConnContext: imHttp.StoreConnPropertiesInContext,
	}
	h.logger = h.logger.With(
		zap.String("handler_name", lifecycle.Name()),
		zap.String("address", lifecycle.Uplink().Addr().String()),
	)

	tlsConfig := lifecycle.CertStore().TLSConfig()

	proxyHandler := &proxyHttpHandler{
		handlerName: lifecycle.Name(),
		options:     opts,
		logger:      h.logger,
		emitter:     lifecycle.Audit(),
	}

	proxyHTTPSHandler := &proxyHttpsHandler{
		options:   opts,
		tlsConfig: tlsConfig,
		emitter:   lifecycle.Audit(),
	}

	h.proxy.OnRequest().Do(proxyHandler)
	h.proxy.OnRequest().HandleConnect(proxyHTTPSHandler)
	go h.startProxy(lifecycle.Uplink().Listener)
	go h.shutdownOnContextDone(lifecycle.Context())
	return
}

func (h *httpProxy) startProxy(listener net.Listener) {
	if err := h.server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		h.logger.Error(
			"failed to start proxy server",
			zap.Error(err),
		)
	}
}

func (h *httpProxy) shutdownOnContextDone(ctx context.Context) {
	<-ctx.Done()
	var err error
	h.logger.Info("Shutting down HTTP proxy")
	if err = h.server.Close(); err != nil {
		h.logger.Error(
			"failed to shutdown proxy endpoint",
			zap.Error(err),
		)
	}
	return
}
