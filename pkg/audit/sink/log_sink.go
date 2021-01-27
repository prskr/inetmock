package sink

import (
	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/logging"
	"go.uber.org/zap"
)

const (
	logSinkName = "logging"
)

func NewLogSink(logger logging.Logger) audit.Sink {
	return &logSink{
		logger: logger,
	}
}

type logSink struct {
	logger logging.Logger
}

func (logSink) Name() string {
	return logSinkName
}

func (l logSink) OnSubscribe(evs <-chan audit.Event) {
	go func(logger logging.Logger, evs <-chan audit.Event) {
		for ev := range evs {
			eventLogger := logger

			if ev.TLS != nil {
				eventLogger = eventLogger.With(
					zap.String("tls_server_name", ev.TLS.ServerName),
					zap.String("tls_cipher_suite", ev.TLS.CipherSuite),
					zap.String("tls_version", ev.TLS.Version),
				)
			}

			eventLogger.Info(
				"handled request",
				zap.Time("timestamp", ev.Timestamp),
				zap.String("application", ev.Application.String()),
				zap.String("transport", ev.Transport.String()),
				zap.String("source_ip", ev.SourceIP.String()),
				zap.Uint16("source_port", ev.SourcePort),
				zap.String("destination_ip", ev.DestinationIP.String()),
				zap.Uint16("destination_port", ev.DestinationPort),
				zap.Any("details", ev.ProtocolDetails),
			)
		}
	}(l.logger, evs)
}
