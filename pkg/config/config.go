package config

import (
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gitlab.com/inetmock/inetmock/pkg/logging"
	"gitlab.com/inetmock/inetmock/pkg/path"
	"go.uber.org/zap"
)

func CreateConfig(flags *pflag.FlagSet) Config {
	logger, _ := logging.CreateLogger()
	configInstance := &config{
		logger: logger.Named("Config"),
		cfg:    viper.New(),
	}

	configInstance.cfg.SetConfigName("config")
	configInstance.cfg.SetConfigType("yaml")
	configInstance.cfg.AddConfigPath("/etc/inetmock/")
	configInstance.cfg.AddConfigPath("$HOME/.inetmock")
	configInstance.cfg.AddConfigPath(".")
	configInstance.cfg.SetEnvPrefix("INetMock")
	_ = configInstance.cfg.BindPFlags(flags)
	configInstance.cfg.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	configInstance.cfg.AutomaticEnv()

	for k, v := range registeredDefaults {
		configInstance.cfg.SetDefault(k, v)
	}

	for k, v := range registeredAliases {
		configInstance.cfg.RegisterAlias(k, v)
	}

	return configInstance
}

type Config interface {
	ReadConfig(configFilePath string) error
	ReadConfigString(config, format string) error
	Viper() *viper.Viper
	TLSConfig() CertOptions
	APIConfig() RPC
	EndpointConfigs() map[string]EndpointConfig
}

type config struct {
	cfg       *viper.Viper
	logger    logging.Logger
	TLS       CertOptions
	Endpoints map[string]EndpointConfig
	API       RPC
}

func (c *config) APIConfig() RPC {
	return c.API
}

func (c *config) ReadConfigString(config, format string) (err error) {
	c.cfg.SetConfigType(format)
	if err = c.cfg.ReadConfig(strings.NewReader(config)); err != nil {
		return
	}

	err = c.cfg.Unmarshal(c)
	return
}

func (c *config) EndpointConfigs() map[string]EndpointConfig {
	return c.Endpoints
}

func (c *config) TLSConfig() CertOptions {
	return c.TLS
}

func (c *config) Viper() *viper.Viper {
	return c.cfg
}

func (c *config) ReadConfig(configFilePath string) (err error) {
	if configFilePath != "" && path.FileExists(configFilePath) {
		c.logger.Info(
			"loading config from passed config file path",
			zap.String("configFilePath", configFilePath),
		)
		c.cfg.SetConfigFile(configFilePath)
	}
	if err = c.cfg.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			err = nil
			c.logger.Warn("failed to load config")
		}
	}

	err = c.cfg.Unmarshal(c)

	return
}
