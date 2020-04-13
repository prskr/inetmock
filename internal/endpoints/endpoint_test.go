package endpoints

import (
	"fmt"
	"github.com/baez90/inetmock/internal/config"
	"github.com/baez90/inetmock/internal/mock"
	"github.com/baez90/inetmock/pkg/api"
	"github.com/golang/mock/gomock"
	"testing"
)

func Test_endpoint_Name(t *testing.T) {
	type fields struct {
		name    string
		handler api.ProtocolHandler
		config  config.HandlerConfig
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "Empty Name if struct is uninitialized",
			fields: fields{},
			want:   "",
		},
		{
			name: "Expected Name if struct is initialized",
			fields: fields{
				name: "sampleHandler",
			},
			want: "sampleHandler",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := endpoint{
				name:    tt.fields.name,
				handler: tt.fields.handler,
				config:  tt.fields.config,
			}
			if got := e.Name(); got != tt.want {
				t.Errorf("Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_endpoint_Shutdown(t *testing.T) {
	type fields struct {
		name    string
		handler api.ProtocolHandler
		config  config.HandlerConfig
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Expect no error if mocked handler does not return one",
			fields: fields{
				handler: func() api.ProtocolHandler {
					handler := mock.NewMockProtocolHandler(gomock.NewController(t))
					handler.EXPECT().
						Shutdown().
						MaxTimes(1).
						Return(nil)
					return handler
				}(),
			},
			wantErr: false,
		},
		{
			name: "Expect error if mocked handler returns one",
			fields: fields{
				handler: func() api.ProtocolHandler {
					handler := mock.NewMockProtocolHandler(gomock.NewController(t))
					handler.EXPECT().
						Shutdown().
						MaxTimes(1).
						Return(fmt.Errorf(""))
					return handler
				}(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &endpoint{
				name:    tt.fields.name,
				handler: tt.fields.handler,
				config:  tt.fields.config,
			}
			if err := e.Shutdown(); (err != nil) != tt.wantErr {
				t.Errorf("Shutdown() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_endpoint_Start(t *testing.T) {

	demoHandlerConfig := config.NewHandlerConfig(
		"sampleHandler",
		80,
		"0.0.0.0",
		nil,
	)

	type fields struct {
		name    string
		handler api.ProtocolHandler
		config  config.HandlerConfig
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Expect no error if mocked handler does not return one",
			fields: fields{
				handler: func() api.ProtocolHandler {
					handler := mock.NewMockProtocolHandler(gomock.NewController(t))
					handler.EXPECT().
						Start(nil).
						MaxTimes(1).
						Return(nil)
					return handler
				}(),
			},
			wantErr: false,
		},
		{
			name: "Expect error if mocked handler returns one",
			fields: fields{
				handler: func() api.ProtocolHandler {
					handler := mock.NewMockProtocolHandler(gomock.NewController(t))
					handler.EXPECT().
						Start(nil).
						MaxTimes(1).
						Return(fmt.Errorf(""))
					return handler
				}(),
			},
			wantErr: true,
		},
		{
			name: "Expect config to be passed to Start call",
			fields: fields{
				config: demoHandlerConfig,
				handler: func() api.ProtocolHandler {
					handler := mock.NewMockProtocolHandler(gomock.NewController(t))
					handler.EXPECT().
						Start(demoHandlerConfig).
						MaxTimes(1).
						Return(fmt.Errorf(""))
					return handler
				}(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &endpoint{
				name:    tt.fields.name,
				handler: tt.fields.handler,
				config:  tt.fields.config,
			}
			if err := e.Start(); (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
