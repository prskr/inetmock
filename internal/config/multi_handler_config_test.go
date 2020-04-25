package config

import (
	"github.com/baez90/inetmock/pkg/api"
	"github.com/spf13/viper"
	"reflect"
	"testing"
)

func Test_multiHandlerConfig_HandlerConfigs(t *testing.T) {
	type fields struct {
		handlerName   string
		ports         []uint16
		listenAddress string
		options       *viper.Viper
	}
	tests := []struct {
		name   string
		fields fields
		want   []api.HandlerConfig
	}{
		{
			name:   "Get empty array if no ports are set",
			fields: fields{},
			want:   make([]api.HandlerConfig, 0),
		},
		{
			name: "Get a single handler config if only one port is set",
			fields: fields{
				handlerName:   "sampleHandler",
				ports:         []uint16{80},
				listenAddress: "0.0.0.0",
				options:       nil,
			},
			want: []api.HandlerConfig{
				api.NewHandlerConfig("sampleHandler", 80, "0.0.0.0", nil),
			},
		},
		{
			name: "Get multiple handler configs if only one port is set",
			fields: fields{
				handlerName:   "sampleHandler",
				ports:         []uint16{80, 8080},
				listenAddress: "0.0.0.0",
				options:       nil,
			},
			want: []api.HandlerConfig{
				api.NewHandlerConfig("sampleHandler", 80, "0.0.0.0", nil),
				api.NewHandlerConfig("sampleHandler", 8080, "0.0.0.0", nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := multiHandlerConfig{
				handlerName:   tt.fields.handlerName,
				ports:         tt.fields.ports,
				listenAddress: tt.fields.listenAddress,
				options:       tt.fields.options,
			}
			if got := m.HandlerConfigs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandlerConfigs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_multiHandlerConfig_HandlerName(t *testing.T) {
	type fields struct {
		handlerName   string
		ports         []uint16
		listenAddress string
		options       *viper.Viper
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "Get empty handler name for uninitialized struct",
			fields: fields{},
			want:   "",
		},
		{
			name: "Get expected handler name for initialized struct",
			fields: fields{
				handlerName: "sampleHandler",
			},
			want: "sampleHandler",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := multiHandlerConfig{
				handlerName:   tt.fields.handlerName,
				ports:         tt.fields.ports,
				listenAddress: tt.fields.listenAddress,
				options:       tt.fields.options,
			}
			if got := m.HandlerName(); got != tt.want {
				t.Errorf("HandlerName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_multiHandlerConfig_ListenAddress(t *testing.T) {
	type fields struct {
		handlerName   string
		ports         []uint16
		listenAddress string
		options       *viper.Viper
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "Get empty ListenAddress for uninitialized struct",
			fields: fields{},
			want:   "",
		},
		{
			name: "Get expected ListenAddress for initialized struct",
			fields: fields{
				listenAddress: "0.0.0.0",
			},
			want: "0.0.0.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := multiHandlerConfig{
				handlerName:   tt.fields.handlerName,
				ports:         tt.fields.ports,
				listenAddress: tt.fields.listenAddress,
				options:       tt.fields.options,
			}
			if got := m.ListenAddress(); got != tt.want {
				t.Errorf("ListenAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_multiHandlerConfig_Options(t *testing.T) {
	type fields struct {
		handlerName   string
		ports         []uint16
		listenAddress string
		options       *viper.Viper
	}
	tests := []struct {
		name   string
		fields fields
		want   *viper.Viper
	}{
		{
			name:   "Get nil Options for uninitialized struct",
			fields: fields{},
			want:   nil,
		},
		{
			name: "Get expected Options for initialized struct",
			fields: fields{
				options: viper.New(),
			},
			want: viper.New(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := multiHandlerConfig{
				handlerName:   tt.fields.handlerName,
				ports:         tt.fields.ports,
				listenAddress: tt.fields.listenAddress,
				options:       tt.fields.options,
			}
			if got := m.Options(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Options() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_multiHandlerConfig_Ports(t *testing.T) {
	type fields struct {
		handlerName   string
		ports         []uint16
		listenAddress string
		options       *viper.Viper
	}
	tests := []struct {
		name   string
		fields fields
		want   []uint16
	}{
		{
			name:   "Get empty Ports for uninitialized struct",
			fields: fields{},
			want:   nil,
		},
		{
			name: "Get expected Ports for initialized struct",
			fields: fields{
				ports: []uint16{80, 8080},
			},
			want: []uint16{80, 8080},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := multiHandlerConfig{
				handlerName:   tt.fields.handlerName,
				ports:         tt.fields.ports,
				listenAddress: tt.fields.listenAddress,
				options:       tt.fields.options,
			}
			if got := m.Ports(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ports() = %v, want %v", got, tt.want)
			}
		})
	}
}
