package main

import (
	"fmt"
	"github.com/spf13/viper"
)

const (
	targetSchemeConfigKey    = "target.scheme"
	targetIpAddressConfigKey = "target.ipAddress"
	targetPortConfigKey      = "target.port"
)

type redirectionTarget struct {
	scheme    string
	ipAddress string
	port      uint16
}

func (rt redirectionTarget) host() string {
	return fmt.Sprintf("%s:%d", rt.ipAddress, rt.port)
}

type httpProxyOptions struct {
	redirectionTarget redirectionTarget
}

func loadFromConfig(config *viper.Viper) (options httpProxyOptions) {

	config.SetDefault(targetSchemeConfigKey, "http")
	config.SetDefault(targetIpAddressConfigKey, "127.0.0.1")
	config.SetDefault(targetPortConfigKey, "80")

	options = httpProxyOptions{
		redirectionTarget{
			scheme:    config.GetString(targetSchemeConfigKey),
			ipAddress: config.GetString(targetIpAddressConfigKey),
			port:      uint16(config.GetInt(targetPortConfigKey)),
		},
	}
	return
}
