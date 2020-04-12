package api

import (
	"github.com/baez90/inetmock/internal/config"
	"go.uber.org/zap"
	"sync"
)

type PluginInstanceFactory func() ProtocolHandler

type LoggingFactory func() (*zap.Logger, error)

type ProtocolHandler interface {
	Start(config config.HandlerConfig)
	Shutdown(wg *sync.WaitGroup)
}
