package sink

import (
	"go.uber.org/zap"

	"inetmock.icb4dc0.de/inetmock/pkg/audit"
	"inetmock.icb4dc0.de/inetmock/pkg/logging"
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

func (l logSink) OnEvent(ev *audit.Event) {
	eventLogger := l.logger

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
		logging.IP("source_ip", ev.SourceIP),
		zap.Uint16("source_port", ev.SourcePort),
		logging.IP("destination_ip", ev.DestinationIP),
		zap.Uint16("destination_port", ev.DestinationPort),
		zap.Any("details", ev.ProtocolDetails),
	)
}
