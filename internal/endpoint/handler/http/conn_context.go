package http

import (
	"context"
	"net"
)

type httpContextKey string

const (
	remoteAddrKey httpContextKey = "gitlab.com/inetmock/inetmock/internal/endpoint/handler/http/context/remoteAddr"
	localAddrKey  httpContextKey = "gitlab.com/inetmock/inetmock/internal/endpoint/handler/http/context/localAddr"
)

func StoreConnPropertiesInContext(ctx context.Context, c net.Conn) context.Context {
	ctx = context.WithValue(ctx, remoteAddrKey, c.RemoteAddr())
	ctx = context.WithValue(ctx, localAddrKey, c.LocalAddr())
	return ctx
}

func localAddr(ctx context.Context) net.Addr {
	val := ctx.Value(localAddrKey)
	if val == nil {
		return nil
	}
	return val.(net.Addr)
}

func remoteAddr(ctx context.Context) net.Addr {
	val := ctx.Value(remoteAddrKey)
	if val == nil {
		return nil
	}
	return val.(net.Addr)
}
