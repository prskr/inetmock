package netutils_test

import (
	"net"
	"testing"

	"gitlab.com/inetmock/inetmock/internal/netutils"
)

func TestManagedListener(t *testing.T) {
	t.Parallel()
	listener, err := netutils.ListenTCP(&net.TCPAddr{IP: net.IPv4(127, 0, 0, 1)})
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

	if err = listener.Close(); err != nil {
		t.Errorf("listener.Close() error = %v", err)
		return
	}

	buf := make([]byte, 1024)
	if _, err = conn.Read(buf); err == nil {
		t.Errorf("Expected error but got none")
	} else {
		t.Logf("listener.Close() error = %v", err)
	}
}
