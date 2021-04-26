//go:generate mockgen -source=$GOFILE -destination=./../../internal/mock/endpoint/protocol_handler.mock.go -package=endpoint_mock

package endpoint

import (
	"context"

	"github.com/soheilhy/cmux"
)

type Lifecycle interface {
	Name() string
	Uplink() Uplink
	UnmarshalOptions(cfg interface{}) error
}

type ProtocolHandler interface {
	Start(ctx context.Context, lifecycle Lifecycle) error
}

type MultiplexHandler interface {
	ProtocolHandler
	Matchers() []cmux.Matcher
}
