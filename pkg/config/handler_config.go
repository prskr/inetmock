package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type HandlerConfig struct {
	HandlerName   string
	Port          uint16
	ListenAddress string
	Options       *viper.Viper
}

func (h HandlerConfig) ListenAddr() string {
	return fmt.Sprintf("%s:%d", h.ListenAddress, h.Port)
}
