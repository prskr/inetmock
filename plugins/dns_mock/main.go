package main

import (
	"fmt"
	"github.com/baez90/inetmock/internal/config"
	"github.com/miekg/dns"
	"go.uber.org/zap"
	"sync"
)

type dnsHandler struct {
	logger    *zap.Logger
	dnsServer []*dns.Server
}

func (d *dnsHandler) Run(config config.HandlerConfig) {
	options := loadFromConfig(config.Options())
	addr := fmt.Sprintf("%s:%d", config.ListenAddress(), config.Port())

	handler := &regexHandler{
		fallback: options.Fallback,
		logger:   d.logger,
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
		zap.String("address", addr),
	)

	d.dnsServer = []*dns.Server{
		{
			Addr:    addr,
			Net:     "udp",
			Handler: handler,
		},
		{
			Addr:    addr,
			Net:     "tcp",
			Handler: handler,
		},
	}

	for _, dnsServer := range d.dnsServer {
		go d.startServer(dnsServer)
	}
}

func (d *dnsHandler) startServer(dnsServer *dns.Server) {
	if err := dnsServer.ListenAndServe(); err != nil {
		d.logger.Error(
			"failed to start DNS server listener",
			zap.Error(err),
		)
	}
}

func (d *dnsHandler) Shutdown(wg *sync.WaitGroup) {
	d.logger.Info("shutting down DNS mock")
	for _, dnsServer := range d.dnsServer {
		if err := dnsServer.Shutdown(); err != nil {
			d.logger.Error(
				"failed to shutdown server",
				zap.Error(err),
			)
		}
	}
	wg.Done()
}
