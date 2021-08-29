package health

import (
	"context"
	"fmt"
	"net"
)

func DNSResolver(cfg Config) *net.Resolver {
	dialer := new(net.Dialer)
	dnsEndpoint := cfg.Client.DNS
	return &net.Resolver{
		PreferGo:     true,
		StrictErrors: true,
		Dial: func(ctx context.Context, _, _ string) (net.Conn, error) {
			return dialer.DialContext(ctx, dnsEndpoint.Proto, fmt.Sprintf("%s:%d", dnsEndpoint.IP, dnsEndpoint.Port))
		},
	}
}