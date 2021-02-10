package mock

import (
	"context"
	"time"

	"github.com/miekg/dns"
	"gitlab.com/inetmock/inetmock/internal/endpoint"
	"gitlab.com/inetmock/inetmock/pkg/logging"
	"go.uber.org/zap"
)

type dnsHandler struct {
	logger    logging.Logger
	dnsServer *dns.Server
}

func (d *dnsHandler) Start(lifecycle endpoint.Lifecycle) (err error) {
	var options dnsOptions
	if options, err = loadFromConfig(lifecycle); err != nil {
		return
	}

	d.logger = lifecycle.Logger().With(
		zap.String("handler_name", lifecycle.Name()),
		zap.String("address", lifecycle.Uplink().Addr().String()),
	)

	handler := &regexHandler{
		handlerName:  lifecycle.Name(),
		fallback:     options.Fallback,
		logger:       lifecycle.Logger(),
		auditEmitter: lifecycle.Audit(),
	}

	for _, rule := range options.Rules {
		d.logger.Info(
			"register DNS rule",
			zap.String("pattern", rule.pattern.String()),
			zap.String("response", rule.response.String()),
		)
		handler.AddRule(rule)
	}

	if lifecycle.Uplink().Listener != nil {
		d.dnsServer = &dns.Server{
			Listener: lifecycle.Uplink().Listener,
			Handler:  handler,
		}
	} else {
		d.dnsServer = &dns.Server{
			PacketConn: lifecycle.Uplink().PacketConn,
			Handler:    handler,
		}
	}

	go d.startServer()
	return
}

func (d *dnsHandler) startServer() {
	if err := d.dnsServer.ActivateAndServe(); err != nil {
		d.logger.Error(
			"failed to start DNS server listener",
			zap.Error(err),
		)
	}
}

func (d *dnsHandler) shutdownOnEnd(ctx context.Context) {
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	if err := d.dnsServer.ShutdownContext(shutdownCtx); err != nil {
		d.logger.Error("failed to shutdown DNS server", zap.Error(err))
	}
}
