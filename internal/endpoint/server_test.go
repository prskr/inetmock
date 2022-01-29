package endpoint_test

import (
	"context"
	"net"
	"net/http"
	"testing"
	"testing/fstest"
	"time"

	"github.com/golang/mock/gomock"
	"golang.org/x/net/context/ctxhttp"

	"gitlab.com/inetmock/inetmock/internal/endpoint"
	audit_mock "gitlab.com/inetmock/inetmock/internal/mock/audit"
	"gitlab.com/inetmock/inetmock/internal/netutils"
	"gitlab.com/inetmock/inetmock/internal/test"
	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/logging"
	"gitlab.com/inetmock/inetmock/protocols/http/mock"
)

func TestServer_ServeGroups_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockEmitter := audit_mock.NewMockEmitter(ctrl)
	mockEmitter.EXPECT().Emit(gomock.Any()).MinTimes(1)
	testLogger := logging.CreateTestLogger(t)

	listenAddr, _ := prepareServer(t, mockEmitter, testLogger)
	httpClient := test.HTTPClientForAddr(t, listenAddr)

	var resp *http.Response

	if r, err := ctxhttp.Get(context.Background(), httpClient, "http://www.stackoverflow.com/"); err != nil {
		t.Errorf("httpClient.Get() error = %v", err)
		return
	} else {
		resp = r
	}

	t.Cleanup(func() {
		if err := resp.Body.Close(); err != nil {
			t.Logf("Failed to close response body error = %v", err)
		}
	})

	if resp.StatusCode != 204 {
		t.Errorf("Status is %d expected 204", resp.StatusCode)
		return
	}
}

func prepareServer(tb testing.TB, emitter audit.Emitter, logger logging.Logger) (*net.TCPAddr, *endpoint.Server) {
	tb.Helper()
	defaultRegistry := endpoint.NewHandlerRegistry()
	mock.AddHTTPMock(defaultRegistry, logger, emitter, fstest.MapFS{})
	builder := endpoint.NewServerBuilder(nil, defaultRegistry, logger)

	var port int
	if p, err := netutils.RandomPort(); err != nil {
		tb.Fatalf("netutils.RandomPort() error = %v", err)
		return nil, nil
	} else {
		port = p
	}

	plainHTTPSpec := endpoint.ListenerSpec{
		Protocol: "tcp",
		Address:  "127.0.0.1",
		Port:     uint16(port),
		Endpoints: map[string]endpoint.Spec{
			"plain": {
				HandlerRef: "http_mock",
				Options: map[string]interface{}{
					"rules": []string{`=> Status(204)`},
				},
			},
		},
	}

	if err := builder.ConfigureGroup(plainHTTPSpec); err != nil {
		tb.Fatalf("builder.ConfigureGroup() error = %v", err)
		return nil, nil
	}

	srv := builder.Server()

	startupCtx, startupCancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	tb.Cleanup(startupCancel)
	if err := srv.ServeGroups(startupCtx); err != nil {
		tb.Fatalf("srv.ServeGroups() error = %v", err)
		return nil, nil
	}

	tb.Cleanup(func() {
		if err := srv.Shutdown(context.Background()); err != nil {
			tb.Fatalf("srv.Shutdown() error = %v", err)
		}
	})

	return &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: port}, srv
}
