package endpoint_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"gitlab.com/inetmock/inetmock/internal/endpoint"
	"gitlab.com/inetmock/inetmock/pkg/logging"
)

func TestEndpoint_Start(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		handlerSetup func(tb testing.TB, ctrl *gomock.Controller) endpoint.ProtocolHandler
		contextSetup func(tb testing.TB) context.Context
		wantErr      bool
	}{
		{
			name: "Successfully start endpoint",
			handlerSetup: func(tb testing.TB, ctrl *gomock.Controller) endpoint.ProtocolHandler {
				tb.Helper()
				return VerifiedProtocolHandler(t, func(ctx context.Context, lifecycle endpoint.Lifecycle) error {
					return nil
				})
			},
			contextSetup: func(tb testing.TB) context.Context {
				tb.Helper()
				ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
				tb.Cleanup(cancel)
				return ctx
			},
			wantErr: false,
		},
		{
			name: "Start fails in protocols",
			handlerSetup: func(tb testing.TB, ctrl *gomock.Controller) endpoint.ProtocolHandler {
				tb.Helper()
				return VerifiedProtocolHandler(t, func(ctx context.Context, lifecycle endpoint.Lifecycle) error {
					return errors.New("something wrong")
				})
			},
			contextSetup: func(tb testing.TB) context.Context {
				tb.Helper()
				ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
				tb.Cleanup(cancel)
				return ctx
			},
			wantErr: true,
		},
		{
			name: "Start panic recovered",
			handlerSetup: func(tb testing.TB, ctrl *gomock.Controller) endpoint.ProtocolHandler {
				tb.Helper()
				return VerifiedProtocolHandler(t, func(ctx context.Context, lifecycle endpoint.Lifecycle) error {
					panic("c'est la vie")
				})
			},
			contextSetup: func(tb testing.TB) context.Context {
				tb.Helper()
				ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
				tb.Cleanup(cancel)
				return ctx
			},
			wantErr: true,
		},
		{
			name: "Start fails due to timeout",
			handlerSetup: func(tb testing.TB, ctrl *gomock.Controller) endpoint.ProtocolHandler {
				tb.Helper()
				return VerifiedProtocolHandler(t, func(ctx context.Context, lifecycle endpoint.Lifecycle) error {
					time.Sleep(50 * time.Millisecond)
					return nil
				})
			},
			contextSetup: func(tb testing.TB) context.Context {
				tb.Helper()
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
				tb.Cleanup(cancel)
				return ctx
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			e := &endpoint.Endpoint{
				Spec: endpoint.Spec{
					Handler: tt.handlerSetup(t, ctrl),
				},
			}
			logger := logging.CreateTestLogger(t)
			lifecycle := endpoint.NewEndpointLifecycle(tt.name, endpoint.Uplink{}, nil)

			if err := e.Start(tt.contextSetup(t), logger, lifecycle); (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
