package app

import (
	"reflect"
	"testing"

	"gitlab.com/inetmock/inetmock/internal/endpoint"
)

func Test_config_ReadConfig(t *testing.T) {
	type args struct {
		config string
	}
	type testCase struct {
		name          string
		args          args
		wantListeners map[string]endpoint.ListenerSpec
		wantErr       bool
	}
	tests := []testCase{
		{
			name: "Test endpoints config",
			args: args{
				// language=yaml
				config: `
listeners:
  tcp_80:
    name: ''
    protocol: tcp
    listenAddress: ''
    port: 80
  tcp_443:
    name: ''
    protocol: tcp
    listenAddress: ''
    port: 443
`,
			},
			wantListeners: map[string]endpoint.ListenerSpec{
				"tcp_80": {
					Name:      "",
					Protocol:  "tcp",
					Address:   "",
					Port:      80,
					Endpoints: nil,
				},
				"tcp_443": {
					Name:      "",
					Protocol:  "tcp",
					Address:   "",
					Port:      443,
					Endpoints: nil,
				},
			},
			wantErr: false,
		},
	}
	scenario := func(tt testCase) func(t *testing.T) {
		return func(t *testing.T) {
			cfg := CreateConfig()
			if err := cfg.ReadConfigString(tt.args.config, "yaml"); (err != nil) != tt.wantErr {
				t.Errorf("ReadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(tt.wantListeners, cfg.ListenerSpecs()) {
				t.Errorf("want = %v, got = %v", tt.wantListeners, cfg.ListenerSpecs())
			}
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, scenario(tt))
	}
}
