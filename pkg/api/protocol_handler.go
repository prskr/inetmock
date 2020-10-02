//go:generate mockgen -source=$GOFILE -destination=./../../internal/mock/api/protocol_handler.mock.go -package=api_mock
package api

import (
	"context"
	"github.com/baez90/inetmock/pkg/config"
	"go.uber.org/zap"
)

type PluginInstanceFactory func() ProtocolHandler

type LoggingFactory func() (*zap.Logger, error)

type ProtocolHandler interface {
	Start(config config.HandlerConfig) error
	Shutdown(ctx context.Context) error
}
