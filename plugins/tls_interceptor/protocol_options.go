package tls_interceptor

import (
	"fmt"
	"github.com/spf13/viper"
)

const (
	targetIpAddressConfigKey = "target.ipAddress"
	targetPortConfigKey      = "target.port"
)

type redirectionTarget struct {
	ipAddress string
	port      uint16
}

func (rt redirectionTarget) address() string {
	return fmt.Sprintf("%s:%d", rt.ipAddress, rt.port)
}

type tlsOptions struct {
	redirectionTarget redirectionTarget
}

func loadFromConfig(config *viper.Viper) tlsOptions {
	return tlsOptions{
		redirectionTarget: redirectionTarget{
			ipAddress: config.GetString(targetIpAddressConfigKey),
			port:      uint16(config.GetInt(targetPortConfigKey)),
		},
	}
}
