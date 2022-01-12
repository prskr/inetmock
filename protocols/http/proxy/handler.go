package proxy

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/elazarl/goproxy"
	"github.com/soheilhy/cmux"
	"go.uber.org/zap"

	"gitlab.com/inetmock/inetmock/internal/endpoint"
	"gitlab.com/inetmock/inetmock/multiplexing"
	"gitlab.com/inetmock/inetmock/pkg/audit"
	auditv1 "gitlab.com/inetmock/inetmock/pkg/audit/v1"
	"gitlab.com/inetmock/inetmock/pkg/cert"
	"gitlab.com/inetmock/inetmock/pkg/logging"
)

const (
	name = "http_proxy"
)

type httpProxy struct {
	logger    logging.Logger
	proxy     *goproxy.ProxyHttpServer
	certStore cert.Store
	emitter   audit.Emitter
	server    *http.Server
}

func (h *httpProxy) Matchers() []cmux.Matcher {
	return []cmux.Matcher{multiplexing.HTTP()}
}

func (h *httpProxy) Start(ctx context.Context, lifecycle endpoint.Lifecycle) error {
	var opts httpProxyOptions
	if err := lifecycle.UnmarshalOptions(&opts); err != nil {
		return err
	}

	h.server = &http.Server{
		Handler:     audit.EmittingHandler(h.emitter, auditv1.AppProtocol_APP_PROTOCOL_HTTP_PROXY, h.proxy),
		ConnContext: audit.StoreConnPropertiesInContext,
	}
	h.logger = h.logger.With(
		zap.String("handler_name", lifecycle.Name()),
		zap.String("address", lifecycle.Uplink().Addr.String()),
	)

	tlsConfig := h.certStore.TLSConfig()

	proxyHandler := &proxyHTTPHandler{
		handlerName: lifecycle.Name(),
		options:     opts,
		logger:      h.logger,
	}

	proxyHTTPSHandler := &proxyHTTPSHandler{
		options:   opts,
		tlsConfig: tlsConfig,
	}

	h.proxy.OnRequest().Do(proxyHandler)
	h.proxy.OnRequest().HandleConnect(proxyHTTPSHandler)
	go h.startProxy(lifecycle.Uplink().Listener)
	go h.shutdownOnContextDone(ctx)
	return nil
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
}
