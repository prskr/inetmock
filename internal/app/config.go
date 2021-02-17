package app

import (
	"net/url"
	"strings"

	"github.com/spf13/viper"
	"gitlab.com/inetmock/inetmock/internal/endpoint"
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
	APIURL() *url.URL
	AuditDataDir() string
	PCAPDataDir() string
	ListenerSpecs() map[string]endpoint.ListenerSpec
}

type apiConfig struct {
	Listen string
}

type config struct {
	cfg       *viper.Viper
	TLS       cert.CertOptions
	Listeners map[string]endpoint.ListenerSpec
	api       apiConfig
	Data      Data
}

func (c config) AuditDataDir() string {
	return c.Data.Audit
}

func (c config) PCAPDataDir() string {
	return c.Data.PCAP
}

func (c config) APIURL() *url.URL {
	if u, err := url.Parse(c.api.Listen); err != nil {
		u, _ = url.Parse("tcp://:0")
		return u
	} else {
		return u
	}
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
