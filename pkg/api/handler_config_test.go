package api

import (
	"github.com/spf13/viper"
	"reflect"
	"testing"
)

func Test_handlerConfig_HandlerName(t *testing.T) {
	type fields struct {
		handlerName   string
		port          uint16
		listenAddress string
		options       *viper.Viper
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "Get empty HandlerName for uninitialized struct",
			fields: fields{},
			want:   "",
		},
		{
			name: "Get expected HandlerName for initialized struct",
			fields: fields{
				handlerName: "sampleHandler",
			},
			want: "sampleHandler",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handlerConfig{
				handlerName:   tt.fields.handlerName,
				port:          tt.fields.port,
				listenAddress: tt.fields.listenAddress,
				options:       tt.fields.options,
			}
			if got := h.HandlerName(); got != tt.want {
				t.Errorf("HandlerName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handlerConfig_ListenAddress(t *testing.T) {
	type fields struct {
		handlerName   string
		port          uint16
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
			h := handlerConfig{
				handlerName:   tt.fields.handlerName,
				port:          tt.fields.port,
				listenAddress: tt.fields.listenAddress,
				options:       tt.fields.options,
			}
			if got := h.ListenAddress(); got != tt.want {
				t.Errorf("ListenAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handlerConfig_Options(t *testing.T) {
	type fields struct {
		handlerName   string
		port          uint16
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
			h := handlerConfig{
				handlerName:   tt.fields.handlerName,
				port:          tt.fields.port,
				listenAddress: tt.fields.listenAddress,
				options:       tt.fields.options,
			}
			if got := h.Options(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Options() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handlerConfig_Port(t *testing.T) {
	type fields struct {
		handlerName   string
		port          uint16
		listenAddress string
		options       *viper.Viper
	}
	tests := []struct {
		name   string
		fields fields
		want   uint16
	}{
		{
			name:   "Get empty Port for uninitialized struct",
			fields: fields{},
			want:   0,
		},
		{
			name: "Get expected Port for initialized struct",
			fields: fields{
				port: 80,
			},
			want: 80,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handlerConfig{
				handlerName:   tt.fields.handlerName,
				port:          tt.fields.port,
				listenAddress: tt.fields.listenAddress,
				options:       tt.fields.options,
			}
			if got := h.Port(); got != tt.want {
				t.Errorf("Port() = %v, want %v", got, tt.want)
			}
		})
	}
}
