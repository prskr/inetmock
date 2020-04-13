//go:generate mockgen -source=protocol_handler.go -destination=./../../mock/protocol_handler_mock.go -package=mock
package api

import (
	"github.com/baez90/inetmock/internal/config"
	"go.uber.org/zap"
)

type PluginInstanceFactory func() ProtocolHandler

type LoggingFactory func() (*zap.Logger, error)

type ProtocolHandler interface {
	Start(config config.HandlerConfig) error
	Shutdown() error
}
