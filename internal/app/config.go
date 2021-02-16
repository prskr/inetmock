package app

import (
	"strings"

	"github.com/spf13/viper"
	"gitlab.com/inetmock/inetmock/internal/endpoint"
	"gitlab.com/inetmock/inetmock/internal/rpc"
	"gitlab.com/inetmock/inetmock/pkg/cert"
	"gitlab.com/inetmock/inetmock/pkg/path"
)

func CreateConfig() Config {
	configInstance := &config{
		cfg: viper.New(),
	}

	configInstance.cfg.SetConfigName("config")
	configInstance.cfg.SetConfigType("yaml")
	configInstance.cfg.AddConfigPath("/etc/inetmock/")
	configInstance.cfg.AddConfigPath("$HOME/.inetmock")
	configInstance.cfg.AddConfigPath(".")
	configInstance.cfg.SetEnvPrefix("INetMock")
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
	TLSConfig() cert.CertOptions
	APIConfig() rpc.Config
	AuditDataDir() string
	PCAPDataDir() string
	ListenerSpecs() map[string]endpoint.ListenerSpec
}

type config struct {
	cfg       *viper.Viper
	TLS       cert.CertOptions
	Listeners map[string]endpoint.ListenerSpec
	API       rpc.Config
	Data      Data
}

func (c config) AuditDataDir() string {
	return c.Data.Audit
}

func (c config) PCAPDataDir() string {
	return c.Data.PCAP
}

func (c config) APIConfig() rpc.Config {
	return c.API
}

func (c config) ListenerSpecs() map[string]endpoint.ListenerSpec {
	return c.Listeners
}

func (c config) TLSConfig() cert.CertOptions {
	return c.TLS
}

func (c *config) ReadConfigString(config, format string) (err error) {
	c.cfg.SetConfigType(format)
	if err = c.cfg.ReadConfig(strings.NewReader(config)); err != nil {
		return
	}

	err = c.cfg.Unmarshal(c)
	return
}

func (c *config) ReadConfig(configFilePath string) (err error) {
	if configFilePath != "" && path.FileExists(configFilePath) {
		c.cfg.SetConfigFile(configFilePath)
	}
	if err = c.cfg.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			err = nil
		} else {
			return
		}
	}

	if err = c.cfg.Unmarshal(c); err != nil {
		return
	}

	if err = c.Data.setup(); err != nil {
		return
	}

	return
}
