package endpoint_test

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"testing"
	"testing/fstest"
	"time"

	"github.com/maxatome/go-testdeep/td"

	"gitlab.com/inetmock/inetmock/internal/endpoint"
	"gitlab.com/inetmock/inetmock/internal/endpoint/handler/http/mock"
	"gitlab.com/inetmock/inetmock/internal/test"
	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/cert"
	"gitlab.com/inetmock/inetmock/pkg/logging"
)

var (
	certStoreOptions = cert.Options{
		RootCACert: cert.File{
			PublicKeyPath:  "./ca/ca.pem",
			PrivateKeyPath: "./ca/ca.key",
		},
		CertCachePath: os.TempDir(),
		Curve:         cert.CurveTypeP256,
		Validity: cert.ValidityByPurpose{
			Server: cert.ValidityDuration{
				NotBeforeRelative: 1 * time.Hour,
				NotAfterRelative:  1 * time.Hour,
			},
		},
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
				registry := endpoint.NewHandlerRegistry()
				mock.AddHTTPMock(registry, logging.CreateTestLogger(t), nil, nil)
				return registry
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
				registry := endpoint.NewHandlerRegistry()
				mock.AddHTTPMock(registry, logging.CreateTestLogger(t), nil, nil)
				return registry
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
			orchestrator := endpoint.NewOrchestrator(store, tt.handlerRegistrySetup(t), logger)
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

func Test_orchestrator_StartEndpoints(t *testing.T) {
	t.Parallel()
	type request struct {
		method string
		url    string
	}
	tests := []struct {
		name                 string
		handlerRegistrySetup func(t *testing.T, emitter audit.Emitter) endpoint.HandlerRegistry
		orchestratorSetup    func(t *testing.T, orchestrator endpoint.Orchestrator, uplink *endpoint.Uplink)
		request              request
		wantErr              bool
		want                 interface{}
	}{
		{
			name: "Start single plain_http handler",
			request: request{
				method: http.MethodGet,
				url:    "http://www.inetmock.org/idx.html",
			},
			handlerRegistrySetup: func(t *testing.T, emitter audit.Emitter) endpoint.HandlerRegistry {
				t.Helper()
				registry := endpoint.NewHandlerRegistry()
				mock.AddHTTPMock(registry, logging.CreateTestLogger(t), emitter, fstest.MapFS{})
				return registry
			},
			orchestratorSetup: func(t *testing.T, orchestrator endpoint.Orchestrator, uplink *endpoint.Uplink) {
				t.Helper()
				err := orchestrator.RegisterListener(endpoint.ListenerSpec{
					Protocol: endpoint.NetProtoTCP.String(),
					Uplink:   uplink,
					Endpoints: map[string]endpoint.Spec{
						"plainHttp": {
							HandlerRef: "http_mock",
							TLS:        false,
							Options: map[string]interface{}{
								"rules": []string{
									`=> Status(204)`,
								},
							},
						},
					},
				})
				if err != nil {
					t.Fatalf("orchestrator.RegisterListener() error = %v", err)
				}
			},
			want: td.Struct(new(http.Response), td.StructFields{
				"StatusCode": 204,
			}),
			wantErr: false,
		},
		{
			name: "Start multiplexed http and https handlers on same listener",
			request: request{
				method: http.MethodGet,
				url:    "https://www.inetmock.org/idx.html",
			},
			handlerRegistrySetup: func(t *testing.T, emitter audit.Emitter) endpoint.HandlerRegistry {
				t.Helper()
				registry := endpoint.NewHandlerRegistry()
				mock.AddHTTPMock(registry, logging.CreateTestLogger(t), emitter, fstest.MapFS{})
				return registry
			},
			orchestratorSetup: func(t *testing.T, orchestrator endpoint.Orchestrator, uplink *endpoint.Uplink) {
				t.Helper()
				err := orchestrator.RegisterListener(endpoint.ListenerSpec{
					Protocol: endpoint.NetProtoTCP.String(),
					Uplink:   uplink,
					Endpoints: map[string]endpoint.Spec{
						"plainHttp": {
							HandlerRef: "http_mock",
							TLS:        false,
							Options: map[string]interface{}{
								"rules": []string{
									`=> Status(204)`,
								},
							},
						},
						"https": {
							HandlerRef: "http_mock",
							TLS:        true,
							Options: map[string]interface{}{
								"rules": []string{
									`=> Status(204)`,
								},
							},
						},
					},
				})
				if err != nil {
					t.Fatalf("orchestrator.RegisterListener() error = %v", err)
				}
			},
			want: td.Struct(new(http.Response), td.StructFields{
				"StatusCode": 204,
			}),
			wantErr: false,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var logger = logging.CreateTestLogger(t)
			var store = cert.MustDefaultStore(certStoreOptions, logger)
			var stream = audit.MustNewEventStream(logger)
			uplink, client := setupTestListener(t, store.TLSConfig())
			orchestrator := endpoint.NewOrchestrator(store, tt.handlerRegistrySetup(t, stream), logger)
			tt.orchestratorSetup(t, orchestrator, uplink)
			ctx, cancel := context.WithCancel(test.Context(t))
			t.Cleanup(cancel)
			handleStartupErrors(t, orchestrator.StartEndpoints(ctx))

			time.Sleep(500 * time.Millisecond)

			var err error
			var req *http.Request
			if req, err = http.NewRequestWithContext(ctx, tt.request.method, tt.request.url, nil); err != nil {
				t.Errorf("http.NewRequest() error = %v", err)
				return
			}

			var resp *http.Response
			if resp, err = client.Do(req); err != nil {
				t.Errorf("client.Do() error = %v", err)
				return
			}
			td.Cmp(t, resp, tt.want)
		})
	}
}

func setupTestListener(tb testing.TB, tlsConfig *tls.Config) (uplink *endpoint.Uplink, client *http.Client) {
	tb.Helper()
	var err error
	uplink = new(endpoint.Uplink)
	if uplink.Listener, err = net.Listen("tcp", "127.0.0.1:0"); err != nil {
		tb.Fatalf("net.Listen() error = %v", err)
	}

	tb.Cleanup(func() {
		if err := uplink.Listener.Close(); err != nil {
			tb.Errorf("listener.Close() error = %v", err)
		}
	})

	var addr *net.TCPAddr
	var isTCPAddr bool
	if addr, isTCPAddr = uplink.Listener.Addr().(*net.TCPAddr); !isTCPAddr {
		tb.Fatalf("address %s is not a TCP address", uplink.Listener.Addr().String())
	}

	uplink.Proto = endpoint.NetProtoTCP

	var dialer net.Dialer

	listenerAddr := fmt.Sprintf("127.0.0.1:%d", addr.Port)

	var transport = &http.Transport{
		TLSClientConfig: tlsConfig,
		DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
			return dialer.DialContext(ctx, "tcp", listenerAddr)
		},
		DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return tls.Dial("tcp", listenerAddr, tlsConfig)
		},
	}

	client = &http.Client{
		Transport: transport,
	}

	return uplink, client
}

func handleStartupErrors(tb testing.TB, errChan chan error) {
	tb.Helper()
	tb.Cleanup(func() {
		for {
			select {
			case err, more := <-errChan:
				if err != nil {
					tb.Errorf("orchestrator.StartEndpoints() error = %v", err)
				}
				if !more {
					return
				}
			default:
				return
			}
		}
	})
}
