package health_test

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/td"

	"gitlab.com/inetmock/inetmock/internal/test"
	"gitlab.com/inetmock/inetmock/pkg/health"
	"gitlab.com/inetmock/inetmock/pkg/logging"
	"gitlab.com/inetmock/inetmock/protocols/http/mock"
)

func TestHttpClient(t *testing.T) {
	t.Parallel()
	type request struct {
		method string
		url    string
	}
	type args struct {
		serverRules []string
		request     request
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    interface{}
	}{
		{
			name: "Get StackOverflow",
			args: args{
				serverRules: []string{
					`=> Status(600)`,
				},
				request: request{
					method: http.MethodGet,
					url:    "http://stackoverflow.com/",
				},
			},
			wantErr: false,
			want: td.Struct(new(http.Response), td.StructFields{
				"StatusCode": 600,
			}),
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cfg := setupServer(t, tt.args.serverRules)
			httpClient := health.HTTPClient(cfg, new(tls.Config))

			var err error
			var req *http.Request
			ctx, cancel := context.WithTimeout(test.Context(t), 50*time.Millisecond)
			t.Cleanup(cancel)
			if req, err = http.NewRequestWithContext(ctx, tt.args.request.method, tt.args.request.url, nil); err != nil {
				t.Fatalf("http.NewRequest() - error = %v", err)
			}

			if resp, err := httpClient.Do(req); err != nil {
				if !tt.wantErr {
					t.Errorf("")
				}
				return
			} else {
				td.Cmp(t, resp, tt.want)
			}
		})
	}
}

func setupServer(tb testing.TB, rules []string) health.Config {
	tb.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		tb.Fatalf("net.Listen() error = %v", err)
	}

	tb.Cleanup(func() {
		if err = listener.Close(); err != nil {
			tb.Fatalf("listener.Close() error = %v", err)
		}
	})

	logger := logging.CreateTestLogger(tb)

	router := mock.Router{
		HandlerName: "Test",
		Logger:      logger,
	}

	for idx := range rules {
		if err := router.RegisterRule(rules[idx]); err != nil {
			tb.Fatalf("router.RegisterRule() - error = %v", err)
		}
	}

	go func(lis net.Listener, handler http.Handler) {
		//nolint:gosec // fine for testing purposes
		_ = http.Serve(lis, handler)
	}(listener, &router)

	var ok bool
	var addr *net.TCPAddr
	if addr, ok = listener.Addr().(*net.TCPAddr); !ok {
		tb.Fatalf("listener.Addr() not a TCP address but %v", listener.Addr())
	}

	srv := health.Server{
		IP:   addr.IP.String(),
		Port: uint16(addr.Port),
	}

	return health.Config{
		Client: health.ClientsConfig{
			HTTP:  srv,
			HTTPS: srv,
		},
	}
}
