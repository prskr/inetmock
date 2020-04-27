package config

import (
	"github.com/spf13/viper"
)

type MultiHandlerConfig struct {
	Handler       string
	Ports         []uint16
	ListenAddress string
	Options       *viper.Viper
}

func (m MultiHandlerConfig) HandlerConfigs() []HandlerConfig {
	configs := make([]HandlerConfig, 0)
	for _, port := range m.Ports {
		configs = append(configs, HandlerConfig{
			HandlerName:   m.Handler,
			Port:          port,
			ListenAddress: m.ListenAddress,
			Options:       m.Options,
		})
	}
	return configs
}
