package dns_mock

import (
	"encoding/binary"
	"github.com/spf13/viper"
	"math"
	"math/rand"
	"net"
	"unsafe"
)

const (
	randomIPStrategyName      = "random"
	incrementalIPStrategyName = "incremental"
	startIPConfigKey          = "startIP"
)

var (
	fallbackStrategies map[string]ResolverFactory
)

type ResolverFactory func(conf *viper.Viper) ResolverFallback

func init() {
	fallbackStrategies = make(map[string]ResolverFactory)
	fallbackStrategies[incrementalIPStrategyName] = func(conf *viper.Viper) ResolverFallback {
		return &incrementalIPFallback{
			latestIp: ipToInt32(net.ParseIP(conf.GetString(startIPConfigKey))),
		}
	}
	fallbackStrategies[randomIPStrategyName] = func(conf *viper.Viper) ResolverFallback {
		return &randomIPFallback{}
	}
}

func CreateResolverFallback(name string, config *viper.Viper) ResolverFallback {
	if factory, ok := fallbackStrategies[name]; ok {
		return factory(config)
	} else {
		return fallbackStrategies[randomIPStrategyName](config)
	}
}

type ResolverFallback interface {
	GetIP() net.IP
}

type incrementalIPFallback struct {
	latestIp uint32
}

func (i *incrementalIPFallback) GetIP() net.IP {
	if i.latestIp < math.MaxInt32 {
		i.latestIp += 1
	}
	return uint32ToIP(i.latestIp)
}

type randomIPFallback struct {
}

func (randomIPFallback) GetIP() net.IP {
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
