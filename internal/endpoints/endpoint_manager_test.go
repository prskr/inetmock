package endpoints

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	api_mock "gitlab.com/inetmock/inetmock/internal/mock/api"
	logging_mock "gitlab.com/inetmock/inetmock/internal/mock/logging"
	plugins_mock "gitlab.com/inetmock/inetmock/internal/mock/plugins"
	"gitlab.com/inetmock/inetmock/pkg/api"
	"gitlab.com/inetmock/inetmock/pkg/config"
	"gitlab.com/inetmock/inetmock/pkg/health"
	"gitlab.com/inetmock/inetmock/pkg/logging"
)

func Test_endpointManager_CreateEndpoint(t *testing.T) {
	type fields struct {
		logger   logging.Logger
		registry api.HandlerRegistry
	}
	type args struct {
		name               string
		multiHandlerConfig config.EndpointConfig
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantErr       bool
		wantEndpoints int
	}{
		{
			name:          "Test add endpoint",
			wantErr:       false,
			wantEndpoints: 1,
			fields: fields{
				logger: func() logging.Logger {
					return logging_mock.NewMockLogger(gomock.NewController(t))
				}(),
				registry: func() api.HandlerRegistry {
					registry := plugins_mock.NewMockHandlerRegistry(gomock.NewController(t))
					registry.
						EXPECT().
						HandlerForName("sampleHandler").
						MinTimes(1).
						MaxTimes(1).
						Return(api_mock.NewMockProtocolHandler(gomock.NewController(t)), true)
					return registry
				}(),
			},
			args: args{
				name: "sampleEndpoint",
				multiHandlerConfig: config.EndpointConfig{
					Handler:       "sampleHandler",
					Ports:         []uint16{80},
					ListenAddress: "0.0.0.0",
				},
			},
		},
		{
			name:          "Test add unknown handler",
			wantErr:       true,
			wantEndpoints: 0,
			fields: fields{
				logger: func() logging.Logger {
					return logging_mock.NewMockLogger(gomock.NewController(t))
				}(),
				registry: func() api.HandlerRegistry {
					registry := plugins_mock.NewMockHandlerRegistry(gomock.NewController(t))
					registry.
						EXPECT().
						HandlerForName("sampleHandler").
						MinTimes(1).
						MaxTimes(1).
						Return(nil, false)
					return registry
				}(),
			},
			args: args{
				name: "sampleEndpoint",
				multiHandlerConfig: config.EndpointConfig{
					Handler:       "sampleHandler",
					Ports:         []uint16{80},
					ListenAddress: "0.0.0.0",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewEndpointManager(tt.fields.registry, tt.fields.logger, health.New(), nil)
			if err := e.CreateEndpoint(tt.args.name, tt.args.multiHandlerConfig); (err != nil) != tt.wantErr {
				t.Errorf("CreateEndpoint() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(e.RegisteredEndpoints()) != tt.wantEndpoints {
				t.Errorf("RegisteredEndpoints() = %d, want = 1", len(e.RegisteredEndpoints()))
				return
			}

			if len(e.RegisteredEndpoints()) > 0 && e.RegisteredEndpoints()[0].Name() != tt.args.name {
				t.Errorf("Name() = %s, want = %s", e.RegisteredEndpoints()[0].Name(), tt.args.name)
			}
		})
	}
}

func Test_endpointManager_StartedEndpoints(t *testing.T) {
	type fields struct {
		logger                   logging.Logger
		registeredEndpoints      []Endpoint
		properlyStartedEndpoints []Endpoint
		registry                 api.HandlerRegistry
	}
	tests := []struct {
		name   string
		fields fields
		want   []Endpoint
	}{
		{
			name:   "",
			fields: fields{},
			want:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := endpointManager{
				logger:                   tt.fields.logger,
				registeredEndpoints:      tt.fields.registeredEndpoints,
				properlyStartedEndpoints: tt.fields.properlyStartedEndpoints,
			}
			if got := e.StartedEndpoints(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StartedEndpoints() = %v, want %v", got, tt.want)
			}
		})
	}
}
