package test

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"time"
)

func HTTPClientForListener(lis net.Listener) (*http.Client, error) {
	switch l := lis.(type) {
	case *net.TCPListener:
		dialer := new(net.Dialer)
		tlsDialer := new(tls.Dialer)
		listenerAddr := l.Addr()
		return &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return dialer.DialContext(ctx, listenerAddr.Network(), listenerAddr.String())
				},
				DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return tlsDialer.DialContext(ctx, listenerAddr.Network(), listenerAddr.String())
				},
				MaxIdleConns:          5,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
		}, nil
	default:
		return nil, errors.New("not a TCP listener")
	}
}

func HTTPClientForInMemListener(lis InMemListener) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext:           lis.DialContext,
			DialTLSContext:        lis.DialContext,
			MaxIdleConns:          5,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
}
