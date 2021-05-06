package endpoint_test

import (
	"os"
	"testing"

	"github.com/maxatome/go-testdeep/td"

	"gitlab.com/inetmock/inetmock/internal/endpoint"
	"gitlab.com/inetmock/inetmock/internal/endpoint/handler/http/mock"
	"gitlab.com/inetmock/inetmock/pkg/cert"
	"gitlab.com/inetmock/inetmock/pkg/logging"
)

var (
	certStoreOptions = cert.Options{
		RootCACert: cert.File{
			PublicKeyPath:  "./ca/ca.pem",
			PrivateKeyPath: "./ca/ca.key",
		},
		CertCachePath:               os.TempDir(),
		Curve:                       cert.CurveTypeP256,
		Validity:                    cert.ValidityByPurpose{},
		IncludeInsecureCipherSuites: true,
		MinTLSVersion:               cert.TLSVersionTLS10,
	}
)

func Test_orchestrator_RegisterListener(t *testing.T) {
	t.Parallel()
	type args struct {
		spec endpoint.ListenerSpec
	}
	tests := []struct {
		name                 string
		args                 args
		handlerRegistrySetup func(t *testing.T) endpoint.HandlerRegistry
		wantErr              bool
		want                 interface{}
	}{
		{
			name: "Successfully register plain HTTP listener",
			args: args{
				spec: endpoint.ListenerSpec{
					Protocol: "tcp",
					Endpoints: map[string]endpoint.Spec{
						"plainHttp": {
							HandlerRef: "http_mock",
							TLS:        false,
							Options:    map[string]interface{}{},
						},
					},
					Uplink: nil,
				},
			},
			handlerRegistrySetup: func(t *testing.T) endpoint.HandlerRegistry {
				t.Helper()
				handler := endpoint.NewHandlerRegistry()
				_ = mock.AddHTTPMock(handler, logging.CreateTestLogger(t), nil, nil)
				return handler
			},
			wantErr: false,
			want:    td.Len(1),
		},
		{
			name: "Successfully register plain HTTP and HTTPS listener",
			args: args{
				spec: endpoint.ListenerSpec{
					Protocol: "tcp",
					Endpoints: map[string]endpoint.Spec{
						"plainHttp": {
							HandlerRef: "http_mock",
							TLS:        false,
							Options:    map[string]interface{}{},
						},
						"https": {
							HandlerRef: "http_mock",
							TLS:        true,
							Options:    map[string]interface{}{},
						},
					},
					Uplink: nil,
				},
			},
			handlerRegistrySetup: func(t *testing.T) endpoint.HandlerRegistry {
				t.Helper()
				handler := endpoint.NewHandlerRegistry()
				_ = mock.AddHTTPMock(handler, logging.CreateTestLogger(t), nil, nil)
				return handler
			},
			wantErr: false,
			want:    td.Len(2),
		},
		{
			name: "Fail because no matching handler registered",
			args: args{
				spec: endpoint.ListenerSpec{
					Protocol: "tcp",
					Endpoints: map[string]endpoint.Spec{
						"plainHttp": {
							HandlerRef: "http_mock",
							TLS:        false,
							Options:    map[string]interface{}{},
						},
					},
					Uplink: nil,
				},
			},
			handlerRegistrySetup: func(t *testing.T) endpoint.HandlerRegistry {
				t.Helper()
				return endpoint.NewHandlerRegistry()
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var logger = logging.CreateTestLogger(t)
			var store cert.Store
			var err error
			if store, err = cert.NewDefaultStore(certStoreOptions, logger); err != nil {
				t.Errorf("cert.NewDefaultStore() error = %v", err)
				return
			}
			orchestrator := endpoint.NewOrchestrator(store, tt.handlerRegistrySetup(t), nil, logger)
			t.Cleanup(func() {
				if uplink := tt.args.spec.Uplink; uplink != nil {
					if err := uplink.Close(); err != nil {
						t.Errorf("uplink.Close() error = %v", err)
					}
				}
			})
			if err := orchestrator.RegisterListener(tt.args.spec); err != nil {
				if !tt.wantErr {
					t.Errorf("RegisterListener() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			td.Cmp(t, orchestrator.Endpoints(), tt.want)
		})
	}
}
