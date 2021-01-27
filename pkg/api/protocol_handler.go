//go:generate mockgen -source=$GOFILE -destination=./../../internal/mock/api/protocol_handler.mock.go -package=api_mock
package api

import (
	"context"

	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/cert"
	"gitlab.com/inetmock/inetmock/pkg/config"
	"gitlab.com/inetmock/inetmock/pkg/logging"
)

type PluginContext interface {
	Logger() logging.Logger
	CertStore() cert.Store
	Audit() audit.Emitter
}

type ProtocolHandler interface {
	Start(ctx PluginContext, config config.HandlerConfig) error
	Shutdown(ctx context.Context) error
}
