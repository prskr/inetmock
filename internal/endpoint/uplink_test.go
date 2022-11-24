package endpoint_test

import (
	"errors"
	"net"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/td"

	"inetmock.icb4dc0.de/inetmock/internal/endpoint"
)

func TestNewUplink(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		conn any
		want any
	}{
		{
			name: "nil value want empty struct",
			conn: nil,
			want: endpoint.Uplink{},
		},
		{
			name: "TCP listener",
			conn: fakeListener{FakeAddr: &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 1234}},
			want: td.Struct(endpoint.Uplink{}, td.StructFields{
				"Addr": td.Struct(new(net.TCPAddr), td.StructFields{}),
			}),
		},
		{
			name: "UDP PacketConn",
			conn: fakePacketConn{FakeLocalAddr: &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 1234}},
			want: td.Struct(endpoint.Uplink{}, td.StructFields{
				"Addr": td.Struct(new(net.UDPAddr), td.StructFields{}),
			}),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotU := endpoint.NewUplink(tt.conn)
			td.Cmp(t, gotU, tt.want)
		})
	}
}

func TestUplink_IsUDP(t *testing.T) {
	t.Parallel()
	type fields struct {
		Addr net.Addr
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "nil address",
			want: false,
		},
		{
			name: "UDP address",
			fields: fields{
				Addr: new(net.UDPAddr),
			},
			want: true,
		},
		{
			name: "TCP address",
			fields: fields{
				Addr: new(net.TCPAddr),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			u := endpoint.Uplink{Addr: tt.fields.Addr}
			if got := u.IsUDP(); got != tt.want {
				t.Errorf("IsUDP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUplink_IsTCP(t *testing.T) {
	t.Parallel()
	type fields struct {
		Addr net.Addr
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "nil address",
			want: false,
		},
		{
			name: "UDP address",
			fields: fields{
				Addr: new(net.UDPAddr),
			},
			want: false,
		},
		{
			name: "TCP address",
			fields: fields{
				Addr: new(net.TCPAddr),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			u := endpoint.Uplink{Addr: tt.fields.Addr}
			if got := u.IsTCP(); got != tt.want {
				t.Errorf("IsTCP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUplink_Close(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		uplinkSetup func(tb testing.TB) endpoint.Uplink
		wantErr     bool
	}{
		{
			name: "Nothing to close",
			uplinkSetup: func(tb testing.TB) endpoint.Uplink {
				tb.Helper()
				return endpoint.Uplink{}
			},
			wantErr: false,
		},
		{
			name: "Listener to close",
			uplinkSetup: func(tb testing.TB) endpoint.Uplink {
				tb.Helper()
				var gotClosed bool
				ul := endpoint.Uplink{
					Listener: fakeListener{
						OnClose: func() error {
							gotClosed = true
							return nil
						},
					},
				}

				tb.Cleanup(func() {
					if !gotClosed {
						tb.Error("Listener did not get closed")
					}
				})

				return ul
			},
			wantErr: false,
		},
		{
			name: "PacketConn to close",
			uplinkSetup: func(tb testing.TB) endpoint.Uplink {
				tb.Helper()
				var gotClosed [2]bool
				ul := endpoint.Uplink{
					Listener: fakeListener{
						OnClose: func() error {
							gotClosed[0] = true
							return nil
						},
					},
					PacketConn: fakePacketConn{
						OnClose: func() error {
							gotClosed[1] = true
							return nil
						},
					},
				}

				tb.Cleanup(func() {
					if !(gotClosed[0] && gotClosed[1]) {
						tb.Error("Didn't close everything")
					}
				})

				return ul
			},
			wantErr: false,
		},
		{
			name: "Close both",
			uplinkSetup: func(tb testing.TB) endpoint.Uplink {
				tb.Helper()
				var gotClosed bool
				ul := endpoint.Uplink{
					PacketConn: fakePacketConn{
						OnClose: func() error {
							gotClosed = true
							return nil
						},
					},
				}

				tb.Cleanup(func() {
					if !gotClosed {
						tb.Error("PacketConn did not get closed")
					}
				})

				return ul
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			u := tt.uplinkSetup(t)
			if err := u.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

var _ net.Listener = (*fakeListener)(nil)

type fakeListener struct {
	FakeAddr net.Addr
	OnClose  func() error
}

func (fl fakeListener) Accept() (net.Conn, error) {
	return nil, nil
}

func (fl fakeListener) Close() error {
	return fl.OnClose()
}

func (fl fakeListener) Addr() net.Addr {
	return fl.FakeAddr
}

var (
	_       net.PacketConn = (*fakePacketConn)(nil)
	errMock                = errors.New("this ain't a real PacketConn")
)

type fakePacketConn struct {
	FakeLocalAddr net.Addr
	OnClose       func() error
}

func (f fakePacketConn) ReadFrom([]byte) (n int, addr net.Addr, err error) {
	return 0, nil, errMock
}

func (f fakePacketConn) WriteTo([]byte, net.Addr) (n int, err error) {
	return 0, errMock
}

func (f fakePacketConn) Close() error {
	return f.OnClose()
}

func (f fakePacketConn) LocalAddr() net.Addr {
	return f.FakeLocalAddr
}

func (f fakePacketConn) SetDeadline(time.Time) error {
	return nil
}

func (f fakePacketConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (f fakePacketConn) SetWriteDeadline(t time.Time) error {
	return nil
}
