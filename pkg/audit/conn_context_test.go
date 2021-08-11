package audit_test

import (
	"context"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/soheilhy/cmux"

	audit_mock "gitlab.com/inetmock/inetmock/internal/mock/audit"
	"gitlab.com/inetmock/inetmock/pkg/audit"
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
					var ctrl = gomock.NewController(tb)
					var connMock = audit_mock.NewMockConn(ctrl)

					connMock.
						EXPECT().
						LocalAddr().
						Return(new(net.TCPAddr))

					connMock.
						EXPECT().
						RemoteAddr().
						Return(new(net.TCPAddr))

					return connMock
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
					var ctrl = gomock.NewController(tb)
					var connMock = audit_mock.NewMockConn(ctrl)

					connMock.
						EXPECT().
						RemoteAddr().
						Return(nil)

					connMock.
						EXPECT().
						LocalAddr().
						Return(new(net.TCPAddr))

					return connMock
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
					var ctrl = gomock.NewController(tb)
					var connMock = audit_mock.NewMockConn(ctrl)

					connMock.
						EXPECT().
						RemoteAddr().
						Return(new(net.TCPAddr))

					connMock.
						EXPECT().
						LocalAddr().
						Return(nil)

					return connMock
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
					var ctrl = gomock.NewController(tb)
					var connMock = audit_mock.NewMockConn(ctrl)

					connMock.
						EXPECT().
						RemoteAddr().
						Return(new(net.TCPAddr))

					connMock.
						EXPECT().
						LocalAddr().
						Return(new(net.TCPAddr))
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
			var ctx = context.Background()
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
