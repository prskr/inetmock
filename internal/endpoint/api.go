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

func NewStartupSpec(name string, uplink Uplink, opts map[string]interface{}) *StartupSpec {
	return &StartupSpec{
		Name:    name,
		Uplink:  uplink,
		Options: opts,
	}
}

type StartupSpec struct {
	Uplink
	Name    string
	Options map[string]interface{}
}

func (s StartupSpec) UnmarshalOptions(cfg interface{}, opts ...UnmarshalOption) error {
	var (
		decoderConfig = new(mapstructure.DecoderConfig)
		decoder       *mapstructure.Decoder
	)
	for idx := range opts {
		opts[idx](decoderConfig)
	}

	decoderConfig.Result = cfg

	if d, err := mapstructure.NewDecoder(decoderConfig); err != nil {
		return err
	} else {
		decoder = d
	}

	return decoder.Decode(s.Options)
}

type (
	ProtocolHandler interface {
		Start(ctx context.Context, ss *StartupSpec) error
	}

	MultiplexHandler interface {
		ProtocolHandler
		Matchers() []cmux.Matcher
	}

	StoppableHandler interface {
		ProtocolHandler
		Stop(ctx context.Context) error
	}
)

type (
	Host interface {
		ConfiguredGroups() []GroupInfo
		ServeGroup(ctx context.Context, groupName string) error
		ServeGroups(ctx context.Context) error
		Shutdown(ctx context.Context) error
		ShutdownGroup(ctx context.Context, groupName string) error
	}

	HostBuilder interface {
		ConfigureGroup(spec ListenerSpec) (err error)
		ConfiguredGroups() []GroupInfo
	}

	GroupInfo struct {
		Name      string
		Endpoints []string
		Serving   bool
	}
)
