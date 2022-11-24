package proxy

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/elazarl/goproxy"
	"github.com/soheilhy/cmux"
	"go.uber.org/zap"

	"inetmock.icb4dc0.de/inetmock/internal/endpoint"
	"inetmock.icb4dc0.de/inetmock/multiplexing"
	"inetmock.icb4dc0.de/inetmock/pkg/audit"
	auditv1 "inetmock.icb4dc0.de/inetmock/pkg/audit/v1"
	"inetmock.icb4dc0.de/inetmock/pkg/cert"
	"inetmock.icb4dc0.de/inetmock/pkg/logging"
)

const (
	name                     = "http_proxy"
	defaultReadHeaderTimeout = 100 * time.Millisecond
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

func (h *httpProxy) Start(_ context.Context, startupSpec *endpoint.StartupSpec) error {
	var opts httpProxyOptions
	if err := startupSpec.UnmarshalOptions(&opts); err != nil {
		return err
	}

	h.server = &http.Server{
		Handler:           audit.EmittingHandler(h.emitter, auditv1.AppProtocol_APP_PROTOCOL_HTTP_PROXY, h.proxy),
		ConnContext:       audit.StoreConnPropertiesInContext,
		ReadHeaderTimeout: defaultReadHeaderTimeout,
	}
	h.logger = h.logger.With(
		zap.String("handler_name", startupSpec.Name),
		zap.String("address", startupSpec.Addr.String()),
	)

	tlsConfig := h.certStore.TLSConfig()

	proxyHandler := &proxyHTTPHandler{
		handlerName: startupSpec.Name,
		options:     opts,
		logger:      h.logger,
	}

	proxyHTTPSHandler := &proxyHTTPSHandler{
		options:   opts,
		tlsConfig: tlsConfig,
	}

	h.proxy.OnRequest().Do(proxyHandler)
	h.proxy.OnRequest().HandleConnect(proxyHTTPSHandler)
	go h.startProxy(startupSpec.Listener)
	return nil
}

func (h *httpProxy) startProxy(listener net.Listener) {
	if err := endpoint.IgnoreShutdownError(h.server.Serve(listener)); err != nil {
		h.logger.Error(
			"failed to start proxy server",
			zap.Error(err),
		)
	}
}
