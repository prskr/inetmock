package config

import (
	"github.com/baez90/inetmock/pkg/logging"
	"github.com/baez90/inetmock/pkg/path"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"strings"
)

var (
	appConfig Config
)

func CreateConfig(flags *pflag.FlagSet) {
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

	appConfig = configInstance
}

func Instance() Config {
	return appConfig
}

type Config interface {
	ReadConfig(configFilePath string) error
	ReadConfigString(config, format string) error
	Viper() *viper.Viper
	TLSConfig() CertOptions
	APIConfig() RPC
	PluginsDir() string
	EndpointConfigs() map[string]MultiHandlerConfig
}

type config struct {
	cfg              *viper.Viper
	logger           logging.Logger
	TLS              CertOptions
	PluginsDirectory string
	Endpoints        map[string]MultiHandlerConfig
	API              RPC
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

func (c *config) EndpointConfigs() map[string]MultiHandlerConfig {
	return c.Endpoints
}

func (c *config) PluginsDir() string {
	return c.PluginsDirectory
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
