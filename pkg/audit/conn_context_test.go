package audit_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/soheilhy/cmux"

	"inetmock.icb4dc0.de/inetmock/pkg/audit"
)

func TestStoreConnPropertiesInContext(t *testing.T) {
	t.Parallel()
	type args struct {
		connSetup func(tb testing.TB) net.Conn
	}
	tests := []struct {
		name           string
		args           args
		wantLocalAddr  bool
		wantRemoteAddr bool
		wantTLSState   bool
	}{
		{
			name: "Mocked net.Conn with local and remote address",
			args: args{
				connSetup: func(tb testing.TB) net.Conn {
					tb.Helper()
					return stubConn(new(net.TCPAddr), new(net.TCPAddr))
				},
			},
			wantLocalAddr:  true,
			wantRemoteAddr: true,
			wantTLSState:   false,
		},
		{
			name: "Mocked net.Conn with local address only",
			args: args{
				connSetup: func(tb testing.TB) net.Conn {
					tb.Helper()

					return connStub{FakeLocalAddr: new(net.TCPAddr)}
				},
			},
			wantLocalAddr:  true,
			wantRemoteAddr: false,
			wantTLSState:   false,
		},
		{
			name: "Mocked net.Conn with remote address only",
			args: args{
				connSetup: func(tb testing.TB) net.Conn {
					tb.Helper()
					return stubConn(nil, new(net.TCPAddr))
				},
			},
			wantLocalAddr:  false,
			wantRemoteAddr: true,
			wantTLSState:   false,
		},
		{
			name: "Unwrap cmux connection",
			args: args{
				connSetup: func(tb testing.TB) net.Conn {
					tb.Helper()
					connMock := stubConn(new(net.TCPAddr), new(net.TCPAddr))
					return &cmux.MuxConn{
						Conn: connMock,
					}
				},
			},
			wantLocalAddr:  true,
			wantRemoteAddr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			got := audit.StoreConnPropertiesInContext(ctx, tt.args.connSetup(t))

			if present := (audit.LocalAddr(got) != nil); present != tt.wantLocalAddr {
				t.Errorf("Expected LocalAddr = %t but was %t", tt.wantLocalAddr, present)
			}

			if present := (audit.RemoteAddr(got) != nil); present != tt.wantRemoteAddr {
				t.Errorf("Expected RemoteAddr = %t but was %t", tt.wantRemoteAddr, present)
			}
			if _, present := audit.TLSConnectionState(got); present != tt.wantTLSState {
				t.Errorf("Expected TLSState = %t but was %t", tt.wantTLSState, present)
			}
		})
	}
}

var _ net.Conn = (*connStub)(nil)

func stubConn(localAddr, remoteAddr net.Addr) net.Conn {
	return connStub{
		FakeLocalAddr:  localAddr,
		FakeRemoteAddr: remoteAddr,
	}
}

type connStub struct {
	FakeLocalAddr  net.Addr
	FakeRemoteAddr net.Addr
}

func (c connStub) Read(b []byte) (n int, err error) {
	panic("implement me")
}

func (c connStub) Write(b []byte) (n int, err error) {
	panic("implement me")
}

func (c connStub) Close() error {
	panic("implement me")
}

func (c connStub) LocalAddr() net.Addr {
	if c.FakeLocalAddr != nil {
		return c.FakeLocalAddr
	}
	return nil
}

func (c connStub) RemoteAddr() net.Addr {
	if c.FakeRemoteAddr != nil {
		return c.FakeRemoteAddr
	}
	return nil
}

func (c connStub) SetDeadline(t time.Time) error {
	panic("implement me")
}

func (c connStub) SetReadDeadline(t time.Time) error {
	panic("implement me")
}

func (c connStub) SetWriteDeadline(t time.Time) error {
	panic("implement me")
}
