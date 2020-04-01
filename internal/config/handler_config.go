package config

import "github.com/spf13/viper"

const (
	pluginConfigKey        = "handler"
	listenAddressConfigKey = "listenaddress"
	portConfigKey          = "port"
)

type handlerConfig struct {
	pluginName    string
	port          uint16
	listenAddress string
	options       *viper.Viper
}

func (h handlerConfig) HandlerName() string {
	return h.pluginName
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

type HandlerConfig interface {
	HandlerName() string
	ListenAddress() string
	Port() uint16
	Options() *viper.Viper
}

func CreateHandlerConfig(configMap interface{}, subConfig *viper.Viper) HandlerConfig {
	underlyingMap := configMap.(map[string]interface{})
	return &handlerConfig{
		pluginName:    underlyingMap[pluginConfigKey].(string),
		listenAddress: underlyingMap[listenAddressConfigKey].(string),
		port:          uint16(underlyingMap[portConfigKey].(int)),
		options:       subConfig,
	}
}
