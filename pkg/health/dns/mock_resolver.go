package dns

import (
	"context"
	"net"
)

type MockResolver struct {
	LookupAddrDelegate func(ctx context.Context, addr string) (names []string, err error)
	LookupHostDelegate func(ctx context.Context, host string) (addrs []net.IP, err error)
}

func (r *MockResolver) ResolverForModule(string) (Resolver, error) {
	return r, nil
}

func (r *MockResolver) LookupA(ctx context.Context, host string) (addrs []net.IP, err error) {
	if r == nil || r.LookupHostDelegate == nil {
		return nil, nil
	}
	return r.LookupHostDelegate(ctx, host)
}

func (r *MockResolver) LookupPTR(ctx context.Context, addr string) (names []string, err error) {
	if r == nil || r.LookupAddrDelegate == nil {
		return nil, nil
	}
	return r.LookupAddrDelegate(ctx, addr)
}
