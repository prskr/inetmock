package dns

import "context"

type MockResolver struct {
	LookupAddrDelegate func(ctx context.Context, addr string) (names []string, err error)
	LookupHostDelegate func(ctx context.Context, host string) (addrs []string, err error)
}

func (r *MockResolver) LookupHost(ctx context.Context, host string) (addrs []string, err error) {
	if r == nil || r.LookupHostDelegate == nil {
		return nil, nil
	}
	return r.LookupHostDelegate(ctx, host)
}

func (r *MockResolver) LookupAddr(ctx context.Context, addr string) (names []string, err error) {
	if r == nil || r.LookupAddrDelegate == nil {
		return nil, nil
	}
	return r.LookupAddrDelegate(ctx, addr)
}
