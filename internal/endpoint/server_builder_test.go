package endpoint_test

import (
	"context"
	"testing"

	"github.com/maxatome/go-testdeep/td"

	"gitlab.com/inetmock/inetmock/internal/endpoint"
	"gitlab.com/inetmock/inetmock/pkg/logging"
)

func TestServerBuilder_ConfigureGroup(t *testing.T) {
	t.Parallel()
	type args struct {
		spec endpoint.ListenerSpec
	}
	tests := []struct {
		name          string
		registrySetup func(tb testing.TB) endpoint.HandlerRegistry
		args          args
		wantErr       bool
		wantEndpoints interface{}
	}{
		{
			name: "Empty listener - nothing registered",
			registrySetup: func(tb testing.TB) endpoint.HandlerRegistry {
				tb.Helper()
				return endpoint.NewHandlerRegistry()
			},
			args: args{
				spec: defaultListenerSpec,
			},
			wantErr: true,
		},
		{
			name: "Empty registry - not matching handler registered",
			registrySetup: func(tb testing.TB) endpoint.HandlerRegistry {
				tb.Helper()
				return endpoint.NewHandlerRegistry()
			},
			args: args{
				spec: endpoint.ListenerSpec{
					Protocol: "tcp",
					Port:     1234,
					Endpoints: map[string]endpoint.Spec{
						"http": {
							HandlerRef: "http_mock",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "All good",
			registrySetup: func(tb testing.TB) endpoint.HandlerRegistry {
				tb.Helper()
				registry := endpoint.NewHandlerRegistry()
				registry.RegisterHandler("http_mock", func() endpoint.ProtocolHandler {
					return ProtocolHandlerFunc(func(context.Context, *endpoint.StartupSpec) error {
						tb.Error("should not start at all")
						return nil
					})
				})

				return registry
			},
			args: args{
				spec: endpoint.ListenerSpec{
					Protocol: "tcp",
					Port:     1234,
					Endpoints: map[string]endpoint.Spec{
						"http": {
							HandlerRef: "http_mock",
						},
					},
				},
			},
			wantErr:       false,
			wantEndpoints: td.Len(1),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e := endpoint.NewServerBuilder(nil, tt.registrySetup(t), logging.CreateTestLogger(t))
			if err := e.ConfigureGroup(tt.args.spec); err != nil {
				if !tt.wantErr {
					t.Errorf("ConfigureGroup() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			td.Cmp(t, e.ConfiguredGroups(), tt.wantEndpoints)
		})
	}
}
