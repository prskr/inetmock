package doh

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/soheilhy/cmux"
	"go.uber.org/zap"

	"gitlab.com/inetmock/inetmock/internal/endpoint"
	"gitlab.com/inetmock/inetmock/multiplexing"
	"gitlab.com/inetmock/inetmock/pkg/audit"
	v1 "gitlab.com/inetmock/inetmock/pkg/audit/v1"
	"gitlab.com/inetmock/inetmock/pkg/logging"
	"gitlab.com/inetmock/inetmock/protocols/dns"
)

const shutdownTimeout = 100 * time.Millisecond

type dohHandler struct {
	logger  logging.Logger
	emitter audit.Emitter
	server  *Server
}

func (d dohHandler) Matchers() []cmux.Matcher {
	return []cmux.Matcher{
		multiplexing.HTTPMatchAnd(func(req *multiplexing.RequestPreface) bool {
			return (req.Method == http.MethodGet || req.Method == http.MethodPost) && strings.HasPrefix(req.Path, "/dns-query")
		}),
	}
}

func (d *dohHandler) Start(ctx context.Context, lifecycle endpoint.Lifecycle) error {
	var options *dns.Options
	if opts, err := dns.OptionsFromLifecycle(lifecycle); err != nil {
		return err
	} else {
		options = opts
	}

	d.logger = d.logger.With(
		zap.String("handler_name", lifecycle.Name()),
		zap.String("address", lifecycle.Uplink().Addr().String()),
	)

	ruleHandler := &dns.RuleHandler{
		TTL: options.TTL,
	}

	for _, rule := range options.Rules {
		d.logger.Debug(
			"Register DNS rule",
			zap.String("raw", rule),
		)
		if err := ruleHandler.RegisterRule(rule); err != nil {
			return err
		}
	}

	handler := &dns.CacheHandler{
		Cache:    options.Cache,
		TTL:      options.TTL,
		Fallback: dns.FallbackHandler(ruleHandler, options.Default, options.TTL),
	}

	queryHandler := DNSQueryHandler(d.logger, handler)
	emittingHandler := audit.EmittingHandler(d.emitter, v1.AppProtocol_APP_PROTOCOL_DNS_OVER_HTTPS, queryHandler)
	d.server = NewServer(emittingHandler)

	go d.startServer(lifecycle.Uplink().Listener)
	go d.shutdownOnEnd(ctx)
	return nil
}

func (d *dohHandler) startServer(listener net.Listener) {
	if err := d.server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		d.logger.Error("Failed to start DoH server", zap.Error(err))
	}
}

func (d *dohHandler) shutdownOnEnd(ctx context.Context) {
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := d.server.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		d.logger.Error("Failed to close server", zap.Error(err))
	}
}
