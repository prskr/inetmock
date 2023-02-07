package client

import (
	"context"
	"errors"
	"net"
	"sync"

	mdns "github.com/miekg/dns"
)

type TraditionalTransport struct {
	lock             sync.Mutex
	Network, Address string
	Dial             func(ctx context.Context, network, address string) (net.Conn, error)
}

func (t *TraditionalTransport) RoundTrip(ctx context.Context, question *mdns.Msg) (resp *mdns.Msg, err error) {
	var conn mdns.Conn
	if conn.Conn, err = t.dial(ctx, t.Network, t.Address); err != nil {
		return nil, err
	}

	defer func() {
		err = errors.Join(err, conn.Close())
	}()

	if err := conn.WriteMsg(question); err != nil {
		return nil, err
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return conn.ReadMsg()
}

func (t *TraditionalTransport) dial(ctx context.Context, network, address string) (net.Conn, error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if t.Dial == nil {
		dialer := new(net.Dialer)
		t.Dial = dialer.DialContext
	}

	return t.Dial(ctx, network, address)
}
