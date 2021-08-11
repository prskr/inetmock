package mock

import (
	"errors"
	"net"
	"reflect"
	"time"

	"github.com/mitchellh/mapstructure"

	"gitlab.com/inetmock/inetmock/internal/endpoint"
	"gitlab.com/inetmock/inetmock/internal/endpoint/handler/dns"
)

const (
	ttlCacheType            = "ttl"
	noneCacheType           = "none"
	incrementalResolverType = "incremental"
	randomResolverType      = "random"
	cidrKey                 = "cidr"
)

var (
	incrementalIPMapping endpoint.Mapping = endpoint.MappingFunc(func(in interface{}) (interface{}, error) {
		if m, ok := in.(map[string]interface{}); ok {
			if cidr, ok := m["cidr"].(string); ok {
				_, n, err := net.ParseCIDR(cidr)
				if err != nil {
					return nil, err
				}
				return NewIncrementalIPResolver(n), nil
			}
		}
		return nil, errors.New("couldn't convert to map structure")
	})
	randomIPMapping endpoint.Mapping = endpoint.MappingFunc(func(in interface{}) (interface{}, error) {
		if m, ok := in.(map[string]interface{}); ok {
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
		var (
			cacheOpts = &struct {
				TTL             time.Duration
				InitialCapacity int
			}{}
			decoderCfg = &mapstructure.DecoderConfig{
				DecodeHook: mapstructure.StringToTimeDurationHookFunc(),
				Result:     cacheOpts,
			}
		)

		if decoder, err := mapstructure.NewDecoder(decoderCfg); err != nil {
			return nil, err
		} else if err = decoder.Decode(in); err != nil {
			return nil, err
		} else {
			return dns.NewCache(dns.WithInitialSize(cacheOpts.InitialCapacity), dns.WithTTL(cacheOpts.TTL)), nil
		}
	})
)

type dnsOptions struct {
	Rules   []string
	Cache   Cache
	Default dns.IPResolver
	TTL     time.Duration
}

func loadFromConfig(lifecycle endpoint.Lifecycle) (dnsOptions, error) {
	var (
		opts            dnsOptions
		composedHook    mapstructure.DecodeHookFunc
		cacheDecodeHook = endpoint.NewOptionByTypeDecoderBuilderFor(&opts.Cache)
		ipResolverHook  = endpoint.NewOptionByTypeDecoderBuilderFor(&opts.Default)
	)

	cacheDecodeHook.AddMappingToMapper(ttlCacheType, ttlCacheMapping)
	cacheDecodeHook.AddMappingToType(noneCacheType, reflect.TypeOf(DelegateCache{}))

	ipResolverHook.AddMappingToMapper(incrementalResolverType, incrementalIPMapping)
	ipResolverHook.AddMappingToMapper(randomResolverType, randomIPMapping)

	composedHook = mapstructure.ComposeDecodeHookFunc(
		cacheDecodeHook.Build(),
		ipResolverHook.Build(),
		mapstructure.StringToTimeDurationHookFunc(),
	)

	if err := lifecycle.UnmarshalOptions(&opts, endpoint.WithDecodeHook(composedHook)); err != nil {
		return dnsOptions{}, err
	}

	return opts, nil
}
