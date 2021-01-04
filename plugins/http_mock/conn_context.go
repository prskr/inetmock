package http_mock

import (
	"context"
	"net"
)

type httpContextKey string

const (
	remoteAddrKey httpContextKey = "RemoteAddr"
	localAddrKey  httpContextKey = "LocalAddr"
)

func StoreConnPropertiesInContext(ctx context.Context, c net.Conn) context.Context {
	ctx = context.WithValue(ctx, remoteAddrKey, c.RemoteAddr())
	ctx = context.WithValue(ctx, localAddrKey, c.LocalAddr())
	return ctx
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
