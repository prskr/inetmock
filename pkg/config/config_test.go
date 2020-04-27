package config

import (
	"github.com/spf13/pflag"
	"os"
	"testing"
)

func Test_config_ReadConfig(t *testing.T) {
	type args struct {
		flags  *pflag.FlagSet
		config string
	}
	tests := []struct {
		name    string
		args    args
		matcher func(Config) bool
		wantErr bool
	}{
		{
			name: "Test endpoints config",
			args: args{
				flags: pflag.NewFlagSet("", pflag.ContinueOnError),
				config: `
endpoints:
  plainHttp:
    handler: http_mock
    listenAddress: 0.0.0.0
    ports:
    - 80
    - 8080
    options: {}
  proxy:
    handler: http_proxy
    listenAddress: 0.0.0.0
    ports:
    - 3128
    options:
      target:
        ipAddress: 127.0.0.1
        port: 80
`,
			},
			matcher: func(c Config) bool {
				if len(c.EndpointConfigs()) < 1 {
					t.Error("Expected EndpointConfigs to be set but is empty")
					return false
				}

				return true
			},
			wantErr: false,
		},
		{
			name: "Test loading of flags by adding",
			args: args{
				flags: func() *pflag.FlagSet {
					flags := pflag.NewFlagSet("", pflag.ContinueOnError)
					flags.String("plugins-directory", os.TempDir(), "")

					return flags
				}(),
				config: ``,
			},
			matcher: func(c Config) bool {
				if c.PluginsDir() != os.TempDir() {
					t.Errorf("PluginsDir() Expected = %s, Got = %s", os.TempDir(), c.PluginsDir())
				}
				return true
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CreateConfig(tt.args.flags)
			cfg := Instance()
			if err := cfg.ReadConfigString(tt.args.config, "yaml"); (err != nil) != tt.wantErr {
				t.Errorf("ReadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.matcher(cfg) {
				t.Error("matcher error")
			}
		})
	}
}
