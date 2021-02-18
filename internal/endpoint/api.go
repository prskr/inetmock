//go:generate mockgen -source=$GOFILE -destination=./../../internal/mock/endpoint/protocol_handler.mock.go -package=endpoint_mock

package endpoint

import (
	"context"

	"github.com/soheilhy/cmux"

	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/cert"
	"gitlab.com/inetmock/inetmock/pkg/logging"
)

type Lifecycle interface {
	Name() string
	Logger() logging.Logger
	CertStore() cert.Store
	Audit() audit.Emitter
	Context() context.Context
	Uplink() Uplink
	UnmarshalOptions(cfg interface{}) error
}

type ProtocolHandler interface {
	Start(ctx Lifecycle) error
}

type MultiplexHandler interface {
	ProtocolHandler
	Matchers() []cmux.Matcher
}
