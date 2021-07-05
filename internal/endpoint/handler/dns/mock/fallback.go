package mock

import (
	"math"
	"math/rand"
	"net"

	"github.com/mitchellh/mapstructure"

	"gitlab.com/inetmock/inetmock/internal/endpoint/handler/dns"
)

const (
	randomIPStrategyName      = "random"
	incrementalIPStrategyName = "incremental"
)

var (
	defaultStartIPIncrementalStrategy = net.ParseIP("10.10.0.1")
	fallbackStrategies                = map[string]ResolverFactory{
		incrementalIPStrategyName: func(args map[string]interface{}) ResolverFallback {
			tmp := struct {
				StartIP string
			}{}
			var startIP net.IP
			if err := mapstructure.Decode(args, &tmp); err == nil {
				startIP = net.ParseIP(tmp.StartIP)
			}
			if startIP == nil || len(startIP) == 0 {
				startIP = defaultStartIPIncrementalStrategy
			}
			return &incrementalIPFallback{
				latestIP: dns.IPToInt32(startIP),
			}
		},
		randomIPStrategyName: func(map[string]interface{}) ResolverFallback {
			return &randomIPFallback{}
		},
	}
)

type ResolverFactory func(args map[string]interface{}) ResolverFallback

func CreateResolverFallback(name string, args map[string]interface{}) ResolverFallback {
	if factory, ok := fallbackStrategies[name]; ok {
		return factory(args)
	} else {
		return fallbackStrategies[randomIPStrategyName](args)
	}
}

type ResolverFallback interface {
	GetIP() net.IP
}

type incrementalIPFallback struct {
	latestIP uint32
}

func (i *incrementalIPFallback) GetIP() net.IP {
	if i.latestIP < math.MaxInt32 {
		i.latestIP += 1
	}
	return dns.Uint32ToIP(i.latestIP)
}

type randomIPFallback struct {
}

func (randomIPFallback) GetIP() net.IP {
	//nolint:gosec
	return dns.Uint32ToIP(uint32(rand.Int31()))
}
