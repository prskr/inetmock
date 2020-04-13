package config

import (
	"bytes"
	"github.com/spf13/viper"
	"reflect"
	"testing"
)

func TestCreateMultiHandlerConfig(t *testing.T) {
	type args struct {
		handlerConfig *viper.Viper
	}
	tests := []struct {
		name string
		args args
		want MultiHandlerConfig
	}{
		{
			name: "Get simple multiHandlerConfig from config",
			args: args{
				handlerConfig: configFromString(`
handler: sampleHandler
listenAddress: 0.0.0.0
ports:
- 80
- 8080
options: {}
`),
			},
			want: &multiHandlerConfig{
				handlerName:   "sampleHandler",
				ports:         []uint16{80, 8080},
				listenAddress: "0.0.0.0",
				options:       viper.New(),
			},
		},
		{
			name: "Get more complex multiHandlerConfig from config",
			args: args{
				handlerConfig: configFromString(`
handler: sampleHandler
listenAddress: 0.0.0.0
ports:
- 80
- 8080
options:
  optionA: asdf
  optionB: as1234
`),
			},
			want: &multiHandlerConfig{
				handlerName:   "sampleHandler",
				ports:         []uint16{80, 8080},
				listenAddress: "0.0.0.0",
				options: configFromString(`
nesting:
  optionA: asdf
  optionB: as1234
`).Sub("nesting"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateMultiHandlerConfig(tt.args.handlerConfig); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateMultiHandlerConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_portsFromConfig(t *testing.T) {
	type args struct {
		handlerConfig *viper.Viper
	}
	tests := []struct {
		name      string
		args      args
		wantPorts []uint16
	}{
		{
			name: "Empty array if config value is not set",
			args: args{
				handlerConfig: viper.New(),
			},
			wantPorts: nil,
		},
		{
			name: "Array of one if `port` is set",
			args: args{
				handlerConfig: configFromString(`
port: 80
`),
			},
			wantPorts: []uint16{80},
		},
		{
			name: "Array of one if `ports` is set as array",
			args: args{
				handlerConfig: configFromString(`
ports:
- 80
`),
			},
			wantPorts: []uint16{80},
		},
		{
			name: "Array of two if `ports` is set as array",
			args: args{
				handlerConfig: configFromString(`
ports:
- 80
- 8080
`),
			},
			wantPorts: []uint16{80, 8080},
		},
		{
			name: "Array of two if `port` is set as array",
			args: args{
				handlerConfig: configFromString(`
ports:
- 80
- 8080
`),
			},
			wantPorts: []uint16{80, 8080},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotPorts := portsFromConfig(tt.args.handlerConfig); !reflect.DeepEqual(gotPorts, tt.wantPorts) {
				t.Errorf("portsFromConfig() = %v, want %v", gotPorts, tt.wantPorts)
			}
		})
	}
}

func configFromString(yaml string) (config *viper.Viper) {
	config = viper.New()
	config.SetConfigType("yaml")
	_ = config.ReadConfig(bytes.NewBufferString(yaml))
	return
}
