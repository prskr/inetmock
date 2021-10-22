package health

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	gohttp "net/http"

	"golang.org/x/net/http2"
)

type (
	HTTPClientForModule interface {
		ClientForModule(module string) (*gohttp.Client, error)
	}
	ClientsForModuleMap map[string]*gohttp.Client
)

func (c ClientsForModuleMap) ClientForModule(module string) (*gohttp.Client, error) {
	if client, ok := c[module]; !ok {
		return nil, fmt.Errorf("%w: %s", ErrNoClientForModule, module)
	} else {
		return client, nil
	}
}

func HTTPClients(cfg Config, tlsConfig *tls.Config) HTTPClientForModule {
	return ClientsForModuleMap{
		"http":  HTTPClient(cfg, tlsConfig),
		"http2": HTTP2Client(cfg, tlsConfig),
	}
}

func HTTPClient(cfg Config, tlsConfig *tls.Config) *gohttp.Client {
	var (
		netDialer = new(net.Dialer)
		tlsDialer = &tls.Dialer{
			NetDialer: netDialer,
			Config:    tlsConfig,
		}
		httpEndpoint  = cfg.Client.HTTP
		httpsEndpoint = cfg.Client.HTTPS
	)

	roundTripper := &gohttp.Transport{
		DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
			return netDialer.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", httpEndpoint.IP, httpEndpoint.Port))
		},
		DialTLSContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
			return tlsDialer.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", httpsEndpoint.IP, httpsEndpoint.Port))
		},
	}
	return &gohttp.Client{
		Transport: roundTripper,
	}
}

func HTTP2Client(cfg Config, tlsConfig *tls.Config) *gohttp.Client {
	var (
		netDialer = new(net.Dialer)
		tlsDialer = &tls.Dialer{
			NetDialer: netDialer,
			Config:    tlsConfig,
		}
		httpsEndpoint = cfg.Client.HTTPS
	)

	http2RoundTripper := &http2.Transport{
		TLSClientConfig: tlsConfig,
		AllowHTTP:       true,
		DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
			return tlsDialer.Dial("tcp", fmt.Sprintf("%s:%d", httpsEndpoint.IP, httpsEndpoint.Port))
		},
	}
	return &gohttp.Client{
		Transport: http2RoundTripper,
	}
}
