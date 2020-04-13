package config

import "github.com/spf13/viper"

const (
	pluginConfigKey        = "handler"
	listenAddressConfigKey = "listenAddress"
	portConfigKey          = "port"
	portsConfigKey         = "ports"
)

type HandlerConfig interface {
	HandlerName() string
	ListenAddress() string
	Port() uint16
	Options() *viper.Viper
}

func NewHandlerConfig(handlerName string, port uint16, listenAddress string, options *viper.Viper) HandlerConfig {
	return &handlerConfig{
		handlerName:   handlerName,
		port:          port,
		listenAddress: listenAddress,
		options:       options,
	}
}

type handlerConfig struct {
	handlerName   string
	port          uint16
	listenAddress string
	options       *viper.Viper
}

func (h handlerConfig) HandlerName() string {
	return h.handlerName
}

func (h handlerConfig) ListenAddress() string {
	return h.listenAddress
}

func (h handlerConfig) Port() uint16 {
	return h.port
}

func (h handlerConfig) Options() *viper.Viper {
	return h.options
}
