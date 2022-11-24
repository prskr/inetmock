package dhcp

import (
	"net"
	"time"

	"github.com/mitchellh/mapstructure"

	"inetmock.icb4dc0.de/inetmock/internal/endpoint"
	"inetmock.icb4dc0.de/inetmock/internal/state"
)

const (
	handlerTypeRange = "range"
)

type DefaultOptions struct {
	ServerID  net.IP
	DNS       []net.IP
	Router    net.IP
	Netmask   net.IP
	LeaseTime time.Duration
}

type ProtocolOptions struct {
	Rules    []string
	Default  DefaultOptions
	Fallback DHCPv4MessageHandler
}

func LoadFromConfig(startupSpec *endpoint.StartupSpec, stateStore state.KVStore) (opts ProtocolOptions, err error) {
	var (
		composedHook       mapstructure.DecodeHookFunc
		defaultHandlerHook = endpoint.NewOptionByTypeDecoderBuilderFor(&opts.Fallback)
	)

	defaultHandlerHook.AddMappingToMapper(handlerTypeRange, rangeHandlerMappingFunc(stateStore))

	composedHook = mapstructure.ComposeDecodeHookFunc(
		defaultHandlerHook.Build(),
		mapstructure.StringToIPHookFunc(),
		mapstructure.StringToTimeDurationHookFunc(),
	)

	if err := startupSpec.UnmarshalOptions(&opts, endpoint.WithDecodeHook(composedHook)); err != nil {
		return ProtocolOptions{}, err
	}

	return
}

func rangeHandlerMappingFunc(store state.KVStore) endpoint.Mapping {
	return endpoint.MappingFunc(func(in any) (any, error) {
		h := &RangeMessageHandler{
			Store: store,
		}

		decoderCfg := &mapstructure.DecoderConfig{
			DecodeHook: mapstructure.ComposeDecodeHookFunc(
				mapstructure.StringToTimeDurationHookFunc(),
				mapstructure.StringToIPHookFunc(),
			),
			Result: h,
		}

		if decoder, err := mapstructure.NewDecoder(decoderCfg); err != nil {
			return nil, err
		} else if err := decoder.Decode(in); err != nil {
			return nil, err
		} else {
			return h, nil
		}
	})
}
