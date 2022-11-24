package net_test

import (
	"net"
	"testing"

	net2 "inetmock.icb4dc0.de/inetmock/internal/net"
)

func TestManagedListener(t *testing.T) {
	t.Parallel()
	addr := &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1)}
	listener, err := net2.ListenTCP(
		addr,
		net2.WithFastOpen(true),
		net2.WithReusePort(true),
	)
	if err != nil {
		t.Errorf("netutils.ListenTCP() error = %v", err)
		return
	}

	t.Cleanup(func() {
		if err := listener.Close(); err != nil {
			t.Logf("Cleanup: listener.Close() error = %v", err)
		}
	})

	conn, err := net.DialTCP("tcp", nil, listener.Addr().(*net.TCPAddr))
	if err != nil {
		t.Errorf("net.DialTCP() error = %v", err)
		return
	}

	buf := make([]byte, 1024)

	if err = listener.Close(); err != nil {
		t.Errorf("listener.Close() error = %v", err)
		return
	}

	if _, err = conn.Read(buf); err == nil {
		t.Errorf("Expected error but got none")
	} else {
		t.Logf("listener.Close() error = %v", err)
	}
}
