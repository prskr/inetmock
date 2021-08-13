package health

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
)

func HTTPClient(cfg Config, tlsConfig *tls.Config) *http.Client {
	var (
		netDialer = new(net.Dialer)
		tlsDialer = &tls.Dialer{
			NetDialer: netDialer,
			Config:    tlsConfig,
		}
		httpEndpoint  = cfg.Client.HTTP
		httpsEndpoint = cfg.Client.HTTPS
	)

	var roundTripper = &http.Transport{
		DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
			return netDialer.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", httpEndpoint.IP, httpEndpoint.Port))
		},
		DialTLSContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
			return tlsDialer.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", httpsEndpoint.IP, httpsEndpoint.Port))
		},
	}
	return &http.Client{
		Transport: roundTripper,
	}
}
