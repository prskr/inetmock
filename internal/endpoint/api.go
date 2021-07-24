//go:generate mockgen -source=$GOFILE -destination=./../../internal/mock/endpoint/protocol_handler.mock.go -package=endpoint_mock

package endpoint

import (
	"context"

	"github.com/mitchellh/mapstructure"
	"github.com/soheilhy/cmux"
)

var (
	WithDecodeHook = func(decodeHook mapstructure.DecodeHookFunc) UnmarshalOption {
		return func(cfg *mapstructure.DecoderConfig) {
			cfg.DecodeHook = decodeHook
		}
	}
	WithErrorUnused = UnmarshalOption(func(cfg *mapstructure.DecoderConfig) {
		cfg.ErrorUnused = true
	})
	WithZeroFields = UnmarshalOption(func(cfg *mapstructure.DecoderConfig) {
		cfg.ZeroFields = true
	})
	WithWeaklyTypedInput = UnmarshalOption(func(cfg *mapstructure.DecoderConfig) {
		cfg.WeaklyTypedInput = true
	})
	WithSquash = UnmarshalOption(func(cfg *mapstructure.DecoderConfig) {
		cfg.Squash = true
	})
)

type UnmarshalOption func(cfg *mapstructure.DecoderConfig)

type Lifecycle interface {
	Name() string
	Uplink() Uplink
	UnmarshalOptions(cfg interface{}, opts ...UnmarshalOption) error
}

type ProtocolHandler interface {
	Start(ctx context.Context, lifecycle Lifecycle) error
}

type MultiplexHandler interface {
	ProtocolHandler
	Matchers() []cmux.Matcher
}
