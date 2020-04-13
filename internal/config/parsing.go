package config

import "github.com/spf13/viper"

func CreateMultiHandlerConfig(handlerConfig *viper.Viper) MultiHandlerConfig {
	return NewMultiHandlerConfig(
		handlerConfig.GetString(pluginConfigKey),
		portsFromConfig(handlerConfig),
		handlerConfig.GetString(listenAddressConfigKey),
		handlerConfig.Sub(OptionsKey),
	)
}

func portsFromConfig(handlerConfig *viper.Viper) (ports []uint16) {
	if portsInt := handlerConfig.GetIntSlice(portsConfigKey); len(portsInt) > 0 {
		for _, port := range portsInt {
			ports = append(ports, uint16(port))
		}
		return
	}

	if portInt := handlerConfig.GetInt(portConfigKey); portInt > 0 {
		ports = append(ports, uint16(portInt))
	}
	return
}
