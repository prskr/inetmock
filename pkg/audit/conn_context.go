package audit

import (
	"context"
	"crypto/tls"
	"net"

	"github.com/soheilhy/cmux"
)

type httpContextKey string

const (
	remoteAddrKey httpContextKey = "gitlab.com/inetmock/inetmock/internal/endpoint/protocols/http/context/remoteAddr"
	localAddrKey  httpContextKey = "gitlab.com/inetmock/inetmock/internal/endpoint/protocols/http/context/localAddr"
	tlsStateKey   httpContextKey = "gitlab.com/inetmock/inetmock/internal/endpoint/protocols/http/context/tlsState"
)

func StoreConnPropertiesInContext(ctx context.Context, c net.Conn) context.Context {
	ctx = context.WithValue(ctx, remoteAddrKey, c.RemoteAddr())
	ctx = context.WithValue(ctx, localAddrKey, c.LocalAddr())
	ctx = addTLSConnectionStateToContext(ctx, c)
	return ctx
}

func addTLSConnectionStateToContext(ctx context.Context, c net.Conn) context.Context {
	switch subConn := c.(type) {
	case *tls.Conn:
		return context.WithValue(ctx, tlsStateKey, subConn.ConnectionState())
	case *cmux.MuxConn:
		return addTLSConnectionStateToContext(ctx, subConn.Conn)
	default:
		return ctx
	}
}

func TLSConnectionState(ctx context.Context) (tls.ConnectionState, bool) {
	val := ctx.Value(tlsStateKey)
	if val == nil {
		return tls.ConnectionState{}, false
	}
	return val.(tls.ConnectionState), true
}

func LocalAddr(ctx context.Context) net.Addr {
	val := ctx.Value(localAddrKey)
	if val == nil {
		return nil
	}
	return val.(net.Addr)
}

func RemoteAddr(ctx context.Context) net.Addr {
	val := ctx.Value(remoteAddrKey)
	if val == nil {
		return nil
	}
	return val.(net.Addr)
}
