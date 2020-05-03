//go:generate mockgen -source=endpoint.go -destination=./../../internal/mock/endpoints/endpoint_mock.go -package=endpoints_mock
package endpoints

import (
	"github.com/baez90/inetmock/pkg/api"
	"github.com/baez90/inetmock/pkg/config"
	"github.com/google/uuid"
)

type Endpoint interface {
	Id() uuid.UUID
	Start() error
	Shutdown() error
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

func (e *endpoint) Start() (err error) {
	return e.handler.Start(e.config)
}

func (e *endpoint) Shutdown() (err error) {
	return e.handler.Shutdown()
}
