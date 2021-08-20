package integration

import (
	"errors"
	"io/fs"
	"net"
	"net/http"
	"testing"

	"gitlab.com/inetmock/inetmock/pkg/logging"
	"gitlab.com/inetmock/inetmock/protocols/http/mock"
)

func NewTestHTTPServer(tb testing.TB, rawBehavior []string, fakeFileFS fs.FS) *HTTPServer {
	tb.Helper()
	router := mock.Router{
		HandlerName: tb.Name(),
		Logger:      logging.CreateTestLogger(tb),
		FakeFileFS:  fakeFileFS,
	}

	for idx := range rawBehavior {
		if err := router.RegisterRule(rawBehavior[idx]); err != nil {
			tb.Fatalf("Failed to parse behavior: %v", err)
		}
	}

	server := &HTTPServer{
		server: &http.Server{
			Handler: &router,
		},
	}

	tb.Cleanup(func() {
		if err := server.Close(); err != nil {
			tb.Errorf("server.Close() err = %v", err)
		}
	})

	return server
}

type HTTPServer struct {
	server *http.Server
}

func (s *HTTPServer) Close() error {
	return s.server.Close()
}

func (s *HTTPServer) Listen(tb testing.TB, listener net.Listener) {
	tb.Helper()
	if err := s.server.Serve(listener); !errors.Is(err, http.ErrServerClosed) {
		tb.Errorf("server.Serve() err = %v", err)
	}
}
