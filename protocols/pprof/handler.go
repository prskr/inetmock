package pprof

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/pprof"

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

func (p *pprofHandler) Start(ctx context.Context, lifecycle endpoint.Lifecycle) error {
	pprofMux := new(http.ServeMux)
	pprofMux.HandleFunc(pprofIndexPath, pprof.Index)
	pprofMux.HandleFunc(pprofCmdLinePath, pprof.Cmdline)
	pprofMux.HandleFunc(pprofProfilePath, pprof.Profile)
	pprofMux.HandleFunc(pprofSymbolPath, pprof.Symbol)
	pprofMux.HandleFunc(pprofTracePath, pprof.Trace)

	p.server = &http.Server{
		Handler:     audit.EmittingHandler(p.emitter, auditv1.AppProtocol_APP_PROTOCOL_PPROF, pprofMux),
		ConnContext: audit.StoreConnPropertiesInContext,
	}

	p.logger = p.logger.With(
		zap.String("address", lifecycle.Uplink().Addr().String()),
	)

	go p.startServer(lifecycle.Uplink().Listener)
	go p.stopServer(ctx)

	return nil
}

func (p *pprofHandler) startServer(listener net.Listener) {
	defer func() {
		if err := listener.Close(); err != nil {
			p.logger.Warn("Failed to close listeners", zap.Error(err))
		}
	}()

	if err := p.server.Serve(listener); err != nil && errors.Is(err, http.ErrServerClosed) {
		p.logger.Error("Failed to start pprof HTTP listener", zap.Error(err))
	}
}

func (p *pprofHandler) stopServer(ctx context.Context) {
	<-ctx.Done()
	p.logger.Info("Shutting down pprof HTTP protocols")
	if err := p.server.Close(); err != nil {
		p.logger.Error("Failed to shutdown pprof HTTP protocols", zap.Error(err))
	}
}
