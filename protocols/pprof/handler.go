package pprof

import (
	"context"
	"net"
	"net/http"
	"net/http/pprof"
	"time"

	"go.uber.org/zap"

	"gitlab.com/inetmock/inetmock/internal/endpoint"
	"gitlab.com/inetmock/inetmock/pkg/audit"
	auditv1 "gitlab.com/inetmock/inetmock/pkg/audit/v1"
	"gitlab.com/inetmock/inetmock/pkg/logging"
)

const (
	name             = "go_pprof"
	pprofIndexPath   = "/debug/pprof/"
	pprofCmdLinePath = "/debug/pprof/cmdline"
	pprofProfilePath = "/debug/pprof/profile"
	pprofSymbolPath  = "/debug/pprof/symbol"
	pprofTracePath   = "/debug/pprof/trace"
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
		ReadHeaderTimeout: 50 * time.Millisecond,
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
