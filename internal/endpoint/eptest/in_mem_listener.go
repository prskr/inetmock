package eptest

import (
	"context"
	"errors"
	"net"
	"testing"
)

var (
	ErrListenerClosed = errors.New("listener closed")
)

type InMemListener interface {
	net.Listener
	Dial(network, addr string) (net.Conn, error)
	DialContext(ctx context.Context, network, addr string) (net.Conn, error)
}

func NewInMemoryListener(tb testing.TB) InMemListener {
	tb.Helper()
	listener := &inMemListener{
		state:       make(chan struct{}),
		connections: make(chan net.Conn),
	}

	tb.Cleanup(func() {
		_ = listener.Close()
	})

	return listener
}

type inMemListener struct {
	state       chan struct{}
	connections chan net.Conn
}

func (i inMemListener) Accept() (net.Conn, error) {
	select {
	case newConnection := <-i.connections:
		return newConnection, nil
	case <-i.state:
		return nil, ErrListenerClosed
	}
}

func (i *inMemListener) Close() error {
	select {
	case _, stillOpen := <-i.state:
		if stillOpen {
			close(i.state)
		}
	default:
	}

	return nil
}

func (i inMemListener) Addr() net.Addr {
	return new(inMemAddr)
}

func (i inMemListener) DialContext(_ context.Context, network, addr string) (net.Conn, error) {
	return i.Dial(network, addr)
}

func (i inMemListener) Dial(_, _ string) (net.Conn, error) {
	select {
	case _, more := <-i.state:
		if !more {
			return nil, ErrListenerClosed
		}
	default:
	}

	serverSide, clientSide := net.Pipe()
	i.connections <- serverSide
	return clientSide, nil
}

type inMemAddr struct {
}

func (i inMemAddr) Network() string {
	return "pipe"
}

func (i inMemAddr) String() string {
	return ""
}
