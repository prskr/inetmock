package test

import (
	"context"
	"errors"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const (
	bufconnBufferSize = 1024 * 1024
)

type APIRegistration func(registrar grpc.ServiceRegistrar)

type GRPCServer struct {
	server   *grpc.Server
	listener *bufconn.Listener
}

func NewTestGRPCServer(tb testing.TB, registrations ...APIRegistration) (srv *GRPCServer) {
	tb.Helper()
	srv = &GRPCServer{
		listener: bufconn.Listen(bufconnBufferSize),
		server:   grpc.NewServer(),
	}

	for _, reg := range registrations {
		reg(srv.server)
	}

	ctx, cancel := context.WithCancel(context.Background())
	tb.Cleanup(cancel)
	go srv.startServerLifecycle(ctx, tb)

	return
}

func (t *GRPCServer) Dial(ctx context.Context, tb testing.TB, opts ...grpc.DialOption) *grpc.ClientConn {
	tb.Helper()
	opts = append(opts,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return t.listener.Dial()
		}),
	)
	conn, err := grpc.DialContext(ctx, "", opts...)
	if err != nil {
		tb.Fatalf("failed to connect to gRPC test server - error = %v", err)
	}
	tb.Cleanup(func() {
		if err := conn.Close(); err != nil {
			tb.Errorf("Failed to close bufconn connection error = %v", err)
		}
	})
	return conn
}

func (t *GRPCServer) startServerLifecycle(ctx context.Context, tb testing.TB) {
	tb.Helper()
	go func() {
		<-ctx.Done()
		t.server.Stop()
	}()

	if err := t.server.Serve(t.listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		tb.Errorf("error occurred during running the gRPC test server - error - %v", err)
	}
}
