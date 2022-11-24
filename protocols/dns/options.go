package dns

import (
	"errors"
	"net"
	"reflect"
	"time"

	"github.com/mitchellh/mapstructure"

	"inetmock.icb4dc0.de/inetmock/internal/endpoint"
	dnsmock "inetmock.icb4dc0.de/inetmock/internal/mock/dns"
)

const (
	inMemCacheType          = "inMemory"
	noneCacheType           = "none"
	incrementalResolverType = "incremental"
	randomResolverType      = "random"
	cidrKey                 = "cidr"
)

var (
	incrementalIPMapping endpoint.Mapping = endpoint.MappingFunc(func(in any) (any, error) {
		if m, ok := in.(map[string]any); ok {
			if cidr, ok := m[cidrKey].(string); ok {
				_, n, err := net.ParseCIDR(cidr)
				if err != nil {
					return nil, err
				}
				return NewIncrementalIPResolver(n), nil
			}
		}
		return nil, errors.New("couldn't convert to map structure")
	})
	randomIPMapping endpoint.Mapping = endpoint.MappingFunc(func(in any) (any, error) {
		if m, ok := in.(map[string]any); ok {
			if cidr, ok := m[cidrKey].(string); ok {
				_, n, err := net.ParseCIDR(cidr)
				if err != nil {
					return nil, err
				}
				return NewRandomIPResolver(n), nil
			}
		}
		return nil, errors.New("couldn't convert to map structure")
	})
	ttlCacheMapping endpoint.Mapping = endpoint.MappingFunc(func(in interface{}) (interface{}, error) {
		return GlobalCache(), nil
	})
)

type Options struct {
	Rules   []string
	Cache   ResourceRecordCache
	Default IPResolver
	TTL     time.Duration
}

func OptionsFromLifecycle(startupSpec *endpoint.StartupSpec) (*Options, error) {
	var (
		composedHook    mapstructure.DecodeHookFunc
		opts            = new(Options)
		cacheDecodeHook = endpoint.NewOptionByTypeDecoderBuilderFor(&opts.Cache)
		ipResolverHook  = endpoint.NewOptionByTypeDecoderBuilderFor(&opts.Default)
	)

	cacheDecodeHook.AddMappingToMapper(inMemCacheType, ttlCacheMapping)
	cacheDecodeHook.AddMappingToType(noneCacheType, reflect.TypeOf(dnsmock.CacheMock{}))

	ipResolverHook.AddMappingToMapper(incrementalResolverType, incrementalIPMapping)
	ipResolverHook.AddMappingToMapper(randomResolverType, randomIPMapping)

	composedHook = mapstructure.ComposeDecodeHookFunc(
		cacheDecodeHook.Build(),
		ipResolverHook.Build(),
		mapstructure.StringToTimeDurationHookFunc(),
	)

	if err := startupSpec.UnmarshalOptions(&opts, endpoint.WithDecodeHook(composedHook)); err != nil {
		return nil, err
	}

	return opts, nil
}
