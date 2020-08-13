package dns_mock

import (
	"context"
	"github.com/baez90/inetmock/pkg/config"
	"github.com/baez90/inetmock/pkg/logging"
	"github.com/miekg/dns"
	"go.uber.org/zap"
)

type dnsHandler struct {
	logger    logging.Logger
	dnsServer []*dns.Server
}

func (d *dnsHandler) Start(config config.HandlerConfig) (err error) {
	options := loadFromConfig(config.Options)

	d.logger = d.logger.With(
		zap.String("handler_name", config.HandlerName),
		zap.String("address", config.ListenAddr()),
	)

	handler := &regexHandler{
		handlerName: config.HandlerName,
		fallback:    options.Fallback,
		logger:      d.logger,
	}

	for _, rule := range options.Rules {
		d.logger.Info(
			"register DNS rule",
			zap.String("pattern", rule.pattern.String()),
			zap.String("response", rule.response.String()),
		)
		handler.AddRule(rule)
	}

	d.logger = d.logger.With(
		zap.String("address", config.ListenAddr()),
	)

	d.dnsServer = []*dns.Server{
		{
			Addr:    config.ListenAddr(),
			Net:     "udp",
			Handler: handler,
		},
		{
			Addr:    config.ListenAddr(),
			Net:     "tcp",
			Handler: handler,
		},
	}

	for _, dnsServer := range d.dnsServer {
		go d.startServer(dnsServer)
	}
	return
}

func (d *dnsHandler) startServer(dnsServer *dns.Server) {
	if err := dnsServer.ListenAndServe(); err != nil {
		d.logger.Error(
			"failed to start DNS server listener",
			zap.Error(err),
		)
	}
}

func (d *dnsHandler) Shutdown(ctx context.Context) error {
	d.logger.Info("shutting down DNS mock")
	for _, dnsServer := range d.dnsServer {
		if err := dnsServer.ShutdownContext(ctx); err != nil {
			d.logger.Error(
				"failed to shutdown server",
				zap.Error(err),
			)
		}
	}
	return nil
}
