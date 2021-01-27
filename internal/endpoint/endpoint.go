//go:generate mockgen -source=$GOFILE -destination=./../../internal/mock/endpoints/endpoint.mock.go -package=endpoints_mock
package endpoint

import (
	"context"

	"github.com/google/uuid"
	"gitlab.com/inetmock/inetmock/pkg/api"
	"gitlab.com/inetmock/inetmock/pkg/config"
)

type Endpoint interface {
	Id() uuid.UUID
	Start(ctx api.PluginContext) error
	Shutdown(ctx context.Context) error
	Name() string
	Handler() string
	Listen() string
	Port() uint16
}

type endpoint struct {
	id      uuid.UUID
	name    string
	handler api.ProtocolHandler
	config  config.HandlerConfig
}

func (e endpoint) Id() uuid.UUID {
	return e.id
}

func (e endpoint) Name() string {
	return e.name
}

func (e endpoint) Handler() string {
	return e.config.HandlerName
}

func (e endpoint) Listen() string {
	return e.config.ListenAddress
}

func (e endpoint) Port() uint16 {
	return e.config.Port
}

func (e *endpoint) Start(ctx api.PluginContext) (err error) {
	return e.handler.Start(ctx, e.config)
}

func (e *endpoint) Shutdown(ctx context.Context) (err error) {
	return e.handler.Shutdown(ctx)
}
