package doh

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/soheilhy/cmux"
	"go.uber.org/zap"

	"inetmock.icb4dc0.de/inetmock/internal/endpoint"
	"inetmock.icb4dc0.de/inetmock/multiplexing"
	"inetmock.icb4dc0.de/inetmock/pkg/audit"
	auditv1 "inetmock.icb4dc0.de/inetmock/pkg/audit/v1"
	"inetmock.icb4dc0.de/inetmock/pkg/logging"
	"inetmock.icb4dc0.de/inetmock/protocols/dns"
)

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

func (d *dohHandler) Start(_ context.Context, startupSpec *endpoint.StartupSpec) error {
	var options *dns.Options
	if opts, err := dns.OptionsFromLifecycle(startupSpec); err != nil {
		return err
	} else {
		options = opts
	}

	d.logger = d.logger.With(
		zap.String("handler_name", startupSpec.Name),
		zap.String("address", startupSpec.Addr.String()),
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

	queryHandler := DNSQueryHandler(d.logger, startupSpec.Name, handler)
	emittingHandler := audit.EmittingHandler(d.emitter, auditv1.AppProtocol_APP_PROTOCOL_DNS_OVER_HTTPS, queryHandler)
	d.server = NewServer(emittingHandler)

	go d.startServer(startupSpec.Listener)
	return nil
}

func (d *dohHandler) startServer(listener net.Listener) {
	if err := endpoint.IgnoreShutdownError(d.server.Serve(listener)); err != nil {
		d.logger.Error("Failed to start DoH server", zap.Error(err))
	}
}
