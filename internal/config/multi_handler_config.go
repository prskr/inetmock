package config

import (
	"github.com/baez90/inetmock/pkg/api"
	"github.com/spf13/viper"
)

type MultiHandlerConfig interface {
	HandlerName() string
	ListenAddress() string
	Ports() []uint16
	Options() *viper.Viper
	HandlerConfigs() []api.HandlerConfig
}

type multiHandlerConfig struct {
	handlerName   string
	ports         []uint16
	listenAddress string
	options       *viper.Viper
}

func NewMultiHandlerConfig(handlerName string, ports []uint16, listenAddress string, options *viper.Viper) MultiHandlerConfig {
	return &multiHandlerConfig{handlerName: handlerName, ports: ports, listenAddress: listenAddress, options: options}
}

func (m multiHandlerConfig) HandlerName() string {
	return m.handlerName
}

func (m multiHandlerConfig) ListenAddress() string {
	return m.listenAddress
}

func (m multiHandlerConfig) Ports() []uint16 {
	return m.ports
}

func (m multiHandlerConfig) Options() *viper.Viper {
	return m.options
}

func (m multiHandlerConfig) HandlerConfigs() []api.HandlerConfig {
	configs := make([]api.HandlerConfig, 0)
	for _, port := range m.ports {
		configs = append(configs, api.NewHandlerConfig(m.handlerName, port, m.listenAddress, m.options))
	}
	return configs
}
