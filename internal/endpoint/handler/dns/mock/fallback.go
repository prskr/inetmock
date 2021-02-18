package mock

import (
	"encoding/binary"
	"math"
	"math/rand"
	"net"
	"unsafe"

	"github.com/mitchellh/mapstructure"
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
				latestIP: ipToInt32(startIP),
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
	return uint32ToIP(i.latestIP)
}

type randomIPFallback struct {
}

func (randomIPFallback) GetIP() net.IP {
	//nolint:gosec
	return uint32ToIP(uint32(rand.Int31()))
}

func uint32ToIP(i uint32) net.IP {
	bytes := (*[4]byte)(unsafe.Pointer(&i))[:]
	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
}

func ipToInt32(ip net.IP) uint32 {
	v4 := ip.To4()
	result := binary.BigEndian.Uint32(v4)
	return result
}
