package mock

import (
	"context"
	"time"

	"github.com/miekg/dns"
	"go.uber.org/zap"

	"gitlab.com/inetmock/inetmock/internal/endpoint"
	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/logging"
)

const shutdownTimeout = 100 * time.Millisecond

type dnsHandler struct {
	logger    logging.Logger
	emitter   audit.Emitter
	dnsServer *dns.Server
}

func (d *dnsHandler) Start(ctx context.Context, lifecycle endpoint.Lifecycle) error {
	var err error
	var options dnsOptions
	if options, err = loadFromConfig(lifecycle); err != nil {
		return err
	}

	d.logger = d.logger.With(
		zap.String("handler_name", lifecycle.Name()),
		zap.String("address", lifecycle.Uplink().Addr().String()),
	)

	handler := &regexHandler{
		handlerName:  lifecycle.Name(),
		fallback:     options.Fallback,
		logger:       d.logger,
		auditEmitter: d.emitter,
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
	go d.shutdownOnEnd(ctx)
	return nil
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
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := d.dnsServer.ShutdownContext(shutdownCtx); err != nil {
		d.logger.Error("failed to shutdown DNS server", zap.Error(err))
	}
}
