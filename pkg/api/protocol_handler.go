//go:generate mockgen -source=protocol_handler.go -destination=./../../internal/mock/api/protocol_handler_mock.go -package=api_mock
package api

import (
	"go.uber.org/zap"
)

type PluginInstanceFactory func() ProtocolHandler

type LoggingFactory func() (*zap.Logger, error)

type ProtocolHandler interface {
	Start(config HandlerConfig) error
	Shutdown() error
}
