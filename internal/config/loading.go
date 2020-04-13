package config

import (
	"github.com/baez90/inetmock/pkg/logging"
	"github.com/baez90/inetmock/pkg/path"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"strings"
)

func CreateConfig() Config {
	logger, _ := logging.CreateLogger()
	return &config{
		logger: logger.Named("Config"),
	}
}

type Config interface {
	InitConfig(flags *pflag.FlagSet)
	ReadConfig(configFilePath string) error
}

type config struct {
	logger logging.Logger
}

func (c config) InitConfig(flags *pflag.FlagSet) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/inetmock/")
	viper.AddConfigPath("$HOME/.inetmock")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("INetMock")
	_ = viper.BindPFlags(flags)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
}

func (c *config) ReadConfig(configFilePath string) (err error) {
	if configFilePath != "" && path.FileExists(configFilePath) {
		c.logger.Info(
			"loading config from passed config file path",
			zap.String("configFilePath", configFilePath),
		)
		viper.SetConfigFile(configFilePath)
	}
	if err = viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			err = nil
			c.logger.Warn("failed to load config")
		}
	}
	return
}
