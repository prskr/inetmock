package mock

import (
	"context"
	"time"

	mdns "github.com/miekg/dns"
	"go.uber.org/zap"

	"gitlab.com/inetmock/inetmock/internal/endpoint"
	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/logging"
)

const shutdownTimeout = 100 * time.Millisecond

type dnsHandler struct {
	logger    logging.Logger
	emitter   audit.Emitter
	dnsServer *mdns.Server
}

func (d *dnsHandler) Start(ctx context.Context, lifecycle endpoint.Lifecycle) error {
	var options dnsOptions
	if opts, err := loadFromConfig(lifecycle); err != nil {
		return err
	} else {
		options = opts
	}

	d.logger = d.logger.With(
		zap.String("handler_name", lifecycle.Name()),
		zap.String("address", lifecycle.Uplink().Addr().String()),
	)

	handler := &RuleHandler{
		Cache:       options.Cache,
		TTL:         options.TTL,
		HandlerName: lifecycle.Name(),
		Logger:      d.logger,
		Emitter:     d.emitter,
		Fallback:    options.Default,
	}

	for _, rule := range options.Rules {
		d.logger.Debug(
			"Register DNS rule",
			zap.String("raw", rule),
		)
		if err := handler.RegisterRule(rule); err != nil {
			return err
		}
	}

	if lifecycle.Uplink().Listener != nil {
		d.dnsServer = &mdns.Server{
			Listener: lifecycle.Uplink().Listener,
			Handler:  handler,
		}
	} else {
		d.dnsServer = &mdns.Server{
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