package main

import (
	"github.com/baez90/inetmock/pkg/logging"
	"go.uber.org/zap"
)

func init() {
	logger, _ := logging.CreateLogger()
	logger = logger.With(
		zap.String("ProtocolHandler", name),
	)
}
