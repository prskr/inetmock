package netflow

import (
	"go.uber.org/zap"

	"inetmock.icb4dc0.de/inetmock/pkg/logging"
)

type LoggerErrorSink struct {
	logging.Logger
}

func (s LoggerErrorSink) OnError(err error) {
	s.Error("Error occurred during interface monitoring", zap.Error(err))
}
