package main

import (
	"bytes"
	"github.com/spf13/viper"
	"reflect"
	"regexp"
	"testing"
)

func Test_loadFromConfig(t *testing.T) {
	type args struct {
		config string
	}
	tests := []struct {
		name        string
		args        args
		wantOptions httpProxyOptions
	}{
		{
			name: "Parse proper configuration with notfound strategy",
			args: args{
				config: `
fallback: notfound,
rules:
  - pattern: ".*"
    response: ./assets/fakeFiles/default.html
`,
			},
			wantOptions: httpProxyOptions{
				FallbackStrategy: StrategyForName(notFoundStrategyName),
				Rules: []targetRule{
					{
						response: "./assets/fakeFiles/default.html",
						pattern:  regexp.MustCompile(".*"),
					},
				},
			},
		},
		{
			name: "Parse proper configuration with pass through strategy",
			args: args{
				config: `
fallback: passthrough
rules:
  - pattern: ".*"
    response: ./assets/fakeFiles/default.html
`,
			},
			wantOptions: httpProxyOptions{
				FallbackStrategy: StrategyForName(passthroughStrategyName),
				Rules: []targetRule{
					{
						response: "./assets/fakeFiles/default.html",
						pattern:  regexp.MustCompile(".*"),
					},
				},
			},
		},
		{
			name: "Parse proper configuration and preserve order of rules",
			args: args{
				config: `
fallback: notfound
rules:
  - pattern: ".*\\.(?i)txt"
    response: ./assets/fakeFiles/default.txt
  - pattern: ".*"
    response: ./assets/fakeFiles/default.html
`,
			},
			wantOptions: httpProxyOptions{
				FallbackStrategy: StrategyForName(notFoundStrategyName),
				Rules: []targetRule{
					{
						response: "./assets/fakeFiles/default.txt",
						pattern:  regexp.MustCompile(".*\\.(?i)txt"),
					},
					{
						response: "./assets/fakeFiles/default.html",
						pattern:  regexp.MustCompile(".*"),
					},
				},
			},
		},
		{
			name: "Parse configuration with non existing fallback strategy key - falling back to 'notfound'",
			args: args{
				config: `
fallback: doesNotExist
rules: []
`,
			},
			wantOptions: httpProxyOptions{
				FallbackStrategy: StrategyForName(notFoundStrategyName),
				Rules:            nil,
			},
		},
		{
			name: "Parse configuration without any fallback key",
			args: args{
				config: `
f4llb4ck: doesNotExist
rules: []
`,
			},
			wantOptions: httpProxyOptions{
				FallbackStrategy: StrategyForName(notFoundStrategyName),
				Rules:            nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := viper.New()
			config.SetConfigType("yaml")
			if err := config.ReadConfig(bytes.NewBufferString(tt.args.config)); err != nil {
				t.Errorf("failed to read config %v", err)
				return
			}
			if gotOptions := loadFromConfig(config); !reflect.DeepEqual(gotOptions, tt.wantOptions) {
				t.Errorf("loadFromConfig() = %v, want %v", gotOptions, tt.wantOptions)
			}
		})
	}
}
