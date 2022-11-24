package mock

import (
	"context"

	mdns "github.com/miekg/dns"
	"go.uber.org/zap"

	"inetmock.icb4dc0.de/inetmock/internal/endpoint"
	"inetmock.icb4dc0.de/inetmock/pkg/audit"
	"inetmock.icb4dc0.de/inetmock/pkg/logging"
	"inetmock.icb4dc0.de/inetmock/protocols/dns"
)

type dnsHandler struct {
	logger    logging.Logger
	emitter   audit.Emitter
	dnsServer *mdns.Server
}

func (d *dnsHandler) Start(_ context.Context, startupSpec *endpoint.StartupSpec) error {
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

	serverHandler := &Server{
		Name: startupSpec.Name,
		Handler: &dns.CacheHandler{
			Cache:    options.Cache,
			TTL:      options.TTL,
			Fallback: dns.FallbackHandler(ruleHandler, options.Default, options.TTL),
		},
		Logger:  d.logger,
		Emitter: d.emitter,
	}

	if startupSpec.IsTCP() {
		d.dnsServer = &mdns.Server{
			Listener: startupSpec.Listener,
			Handler:  serverHandler,
		}
	} else {
		d.dnsServer = &mdns.Server{
			PacketConn: startupSpec.PacketConn,
			Handler:    serverHandler,
		}
	}

	go d.startServer()
	return nil
}

func (d *dnsHandler) startServer() {
	if err := endpoint.IgnoreShutdownError(d.dnsServer.ActivateAndServe()); err != nil {
		d.logger.Error(
			"failed to start DNS server listener",
			zap.Error(err),
		)
	}
}
