package main

import (
	"net"
	"net/http"
	"testing"

	"inetmock.icb4dc0.de/inetmock/internal/test"
	"inetmock.icb4dc0.de/inetmock/internal/test/integration"
	"inetmock.icb4dc0.de/inetmock/pkg/health"
	"inetmock.icb4dc0.de/inetmock/pkg/logging"
)

func Test_HTTP_runCheck(t *testing.T) {
	t.Parallel()
	type args struct {
		script string
	}
	tests := []struct {
		name     string
		behavior []string
		args     args
		wantErr  bool
	}{
		{
			name: "Check for 200 HTTP OK",
			behavior: []string{
				`=> Status(200)`,
			},
			args: args{
				script: `http.GET("http://stackoverflow.com/index.html") => Status(200)`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var (
				listener net.Listener
				client   *http.Client
				err      error
			)

			if listener, err = net.Listen("tcp", "127.0.0.1:0"); err != nil {
				t.Fatalf("net.Listen() err = %v", err)
			}

			testServer := integration.NewTestHTTPServer(t, tt.behavior, nil)
			go testServer.Listen(t, listener)

			client = test.HTTPClientForListener(t, listener)

			clients := health.ClientsForModuleMap{
				"http":  client,
				"http2": client,
			}

			if err := runCheck(test.Context(t), logging.CreateTestLogger(t), tt.args.script, clients, nil); (err != nil) != tt.wantErr {
				t.Errorf("runCheck() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
