package pprof

import (
	"context"
	"net"
	"net/http"
	"net/http/pprof"
	"time"

	"go.uber.org/zap"

	"inetmock.icb4dc0.de/inetmock/internal/endpoint"
	"inetmock.icb4dc0.de/inetmock/pkg/audit"
	auditv1 "inetmock.icb4dc0.de/inetmock/pkg/audit/v1"
	"inetmock.icb4dc0.de/inetmock/pkg/logging"
)

const (
	name             = "go_pprof"
	pprofIndexPath   = "/debug/pprof/"
	pprofCmdLinePath = "/debug/pprof/cmdline"
	pprofProfilePath = "/debug/pprof/profile"
	pprofSymbolPath  = "/debug/pprof/symbol"
	pprofTracePath   = "/debug/pprof/trace"

	defaultReadHeaderTimeout = 100 * time.Millisecond
)

type pprofHandler struct {
	logger  logging.Logger
	emitter audit.Emitter
	server  *http.Server
}

func (p *pprofHandler) Start(_ context.Context, startupSpec *endpoint.StartupSpec) error {
	pprofMux := new(http.ServeMux)
	pprofMux.HandleFunc(pprofIndexPath, pprof.Index)
	pprofMux.HandleFunc(pprofCmdLinePath, pprof.Cmdline)
	pprofMux.HandleFunc(pprofProfilePath, pprof.Profile)
	pprofMux.HandleFunc(pprofSymbolPath, pprof.Symbol)
	pprofMux.HandleFunc(pprofTracePath, pprof.Trace)

	p.server = &http.Server{
		Handler:           audit.EmittingHandler(p.emitter, auditv1.AppProtocol_APP_PROTOCOL_PPROF, pprofMux),
		ConnContext:       audit.StoreConnPropertiesInContext,
		ReadHeaderTimeout: defaultReadHeaderTimeout,
	}

	p.logger = p.logger.With(
		zap.String("address", startupSpec.Addr.String()),
	)

	go p.startServer(startupSpec.Listener)

	return nil
}

func (p *pprofHandler) startServer(listener net.Listener) {
	if err := endpoint.IgnoreShutdownError(p.server.Serve(listener)); err != nil {
		p.logger.Error("Failed to start pprof HTTP listener", zap.Error(err))
	}
}
