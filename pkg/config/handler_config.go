package config

import "github.com/spf13/viper"

type HandlerConfig struct {
	HandlerName   string
	Port          uint16
	ListenAddress string
	Options       *viper.Viper
}
