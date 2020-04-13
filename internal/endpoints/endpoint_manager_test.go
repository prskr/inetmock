package endpoints

import (
	"github.com/baez90/inetmock/internal/config"
	"github.com/baez90/inetmock/internal/mock"
	"github.com/baez90/inetmock/internal/plugins"
	"github.com/baez90/inetmock/pkg/logging"
	"github.com/golang/mock/gomock"
	"testing"
)

func Test_endpointManager_CreateEndpoint(t *testing.T) {
	type fields struct {
		logger                   logging.Logger
		registeredEndpoints      []Endpoint
		properlyStartedEndpoints []Endpoint
		registry                 plugins.HandlerRegistry
	}
	type args struct {
		name               string
		multiHandlerConfig config.MultiHandlerConfig
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
					return mock.NewMockLogger(gomock.NewController(t))
				}(),
				registeredEndpoints:      nil,
				properlyStartedEndpoints: nil,
				registry: func() plugins.HandlerRegistry {
					registry := mock.NewMockHandlerRegistry(gomock.NewController(t))
					registry.
						EXPECT().
						HandlerForName("sampleHandler").
						MinTimes(1).
						MaxTimes(1).
						Return(mock.NewMockProtocolHandler(gomock.NewController(t)), true)
					return registry
				}(),
			},
			args: args{
				name: "sampleEndpoint",
				multiHandlerConfig: config.NewMultiHandlerConfig(
					"sampleHandler",
					[]uint16{80},
					"0.0.0.0",
					nil,
				),
			},
		},
		{
			name:          "Test add unknown handler",
			wantErr:       true,
			wantEndpoints: 0,
			fields: fields{
				logger: func() logging.Logger {
					return mock.NewMockLogger(gomock.NewController(t))
				}(),
				registeredEndpoints:      nil,
				properlyStartedEndpoints: nil,
				registry: func() plugins.HandlerRegistry {
					registry := mock.NewMockHandlerRegistry(gomock.NewController(t))
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
				multiHandlerConfig: config.NewMultiHandlerConfig(
					"sampleHandler",
					[]uint16{80},
					"0.0.0.0",
					nil,
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &endpointManager{
				logger:                   tt.fields.logger,
				registeredEndpoints:      tt.fields.registeredEndpoints,
				properlyStartedEndpoints: tt.fields.properlyStartedEndpoints,
				registry:                 tt.fields.registry,
			}

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
