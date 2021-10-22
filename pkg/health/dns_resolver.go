package health

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"

	"gitlab.com/inetmock/inetmock/pkg/health/dns"
	"gitlab.com/inetmock/inetmock/protocols/dns/client"
)

type (
	ResolverForModule interface {
		ResolverForModule(module string) (dns.Resolver, error)
	}
	ResolversForModuleMap map[string]dns.Resolver
)

func (d ResolversForModuleMap) ResolverForModule(module string) (dns.Resolver, error) {
	if resolver, ok := d[module]; !ok {
		return nil, fmt.Errorf("%w: %s", ErrNoClientForModule, module)
	} else {
		return resolver, nil
	}
}

func Resolvers(cfg Config, tlsConfig *tls.Config) ResolverForModule {
	dialer := new(net.Dialer)
	tlsDialer := tls.Dialer{
		Config: tlsConfig,
	}
	return ResolversForModuleMap{
		"dns": &client.Resolver{
			Transport: &client.TraditionalTransport{
				Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
					return dialer.DialContext(ctx, cfg.Client.DNS.Proto, fmt.Sprintf("%s:%d", cfg.Client.DNS.IP, cfg.Client.DNS.Port))
				},
			},
		},
		"dot": &client.Resolver{
			Transport: &client.TraditionalTransport{
				Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
					return tlsDialer.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", cfg.Client.DoT.IP, cfg.Client.DoT.Port))
				},
			},
		},
		"doh": &client.Resolver{
			Transport: &client.HTTPTransport{
				Packer: client.RequestPackerPOST,
				Client: HTTPClient(cfg, tlsConfig),
				Scheme: "https",
				Server: cfg.Client.DNS.IP,
			},
		},
		"doh2": &client.Resolver{
			Transport: &client.HTTPTransport{
				Packer: client.RequestPackerPOST,
				Client: HTTP2Client(cfg, tlsConfig),
				Scheme: "https",
				Server: cfg.Client.DNS.IP,
			},
		},
	}
}
