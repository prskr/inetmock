package test

import (
	"context"
	"errors"
	"net"
	"sync"

	"google.golang.org/grpc"
)

type APIRegistration func(registrar grpc.ServiceRegistrar)

type GRPCServer struct {
	mutex         sync.Mutex
	serverRunning chan struct{}
	cancel        context.CancelFunc
	server        *grpc.Server
	addr          *net.TCPAddr
	listener      net.Listener
}

func NewTestGRPCServer(registrations ...APIRegistration) (srv *GRPCServer, err error) {
	srv = new(GRPCServer)
	if srv.listener, err = net.Listen("tcp", "127.0.0.1:0"); err != nil {
		return
	}

	var isTCPAddr bool
	if srv.addr, isTCPAddr = srv.listener.Addr().(*net.TCPAddr); !isTCPAddr {
		err = errors.New("expected TPC addr but wasn't")
	}

	srv.server = grpc.NewServer()

	for _, reg := range registrations {
		reg(srv.server)
	}

	return
}

func (t *GRPCServer) StartServer() (err error) {
	t.mutex.Lock()
	var ctx context.Context
	ctx, t.cancel = context.WithCancel(context.Background())
	t.mutex.Unlock()
	errs := t.startServerAsync(ctx)
	select {
	case err = <-errs:
	default:
	}

	return
}

func (t *GRPCServer) Dial(ctx context.Context, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(opts, grpc.WithInsecure())

	return grpc.DialContext(ctx, t.addr.String(), opts...)
}

func (t *GRPCServer) StopServer() {
	t.cancel()
}

func (t *GRPCServer) IsRunning() bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if t.serverRunning == nil {
		return false
	}

	select {
	case _, more := <-t.serverRunning:
		return more
	default:
		return true
	}
}

func (t *GRPCServer) startServerAsync(ctx context.Context) (errs chan error) {
	errs = make(chan error)
	t.mutex.Lock()
	t.serverRunning = make(chan struct{})
	t.mutex.Unlock()
	defer close(t.serverRunning)

	go func() {
		<-ctx.Done()
		if t.IsRunning() {
			t.server.Stop()
		}
	}()

	go func() {
		if err := t.server.Serve(t.listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			errs <- err
		}
	}()
	return
}
