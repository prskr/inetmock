//go:generate mockgen -source=protocol_handler.go -destination=./../../internal/mock/api/protocol_handler_mock.go -package=api_mock
package api

import (
	"github.com/baez90/inetmock/pkg/config"
	"go.uber.org/zap"
)

type PluginInstanceFactory func() ProtocolHandler

type LoggingFactory func() (*zap.Logger, error)

type ProtocolHandler interface {
	Start(config config.HandlerConfig) error
	Shutdown() error
}
