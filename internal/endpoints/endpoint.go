//go:generate mockgen -source=endpoint.go -destination=./../../internal/mock/endpoints/endpoint_mock.go -package=endpoints_mock
package endpoints

import (
	"github.com/baez90/inetmock/pkg/api"
)

type Endpoint interface {
	Start() error
	Shutdown() error
	Name() string
}

type endpoint struct {
	name    string
	handler api.ProtocolHandler
	config  api.HandlerConfig
}

func (e endpoint) Name() string {
	return e.name
}

func (e *endpoint) Start() (err error) {
	return e.handler.Start(e.config)
}

func (e *endpoint) Shutdown() (err error) {
	return e.handler.Shutdown()
}
