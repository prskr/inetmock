package health

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
)

func HTTPClient(cfg Config, tlsConfig *tls.Config) *http.Client {
	var roundTripper = &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			httpEndpoint := cfg.Client.HTTP
			return net.Dial("tcp", fmt.Sprintf("%s:%d", httpEndpoint.IP, httpEndpoint.Port))
		},
		DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			httpsEndpoint := cfg.Client.HTTPS
			return tls.Dial("tcp", fmt.Sprintf("%s:%d", httpsEndpoint.IP, httpsEndpoint.Port), tlsConfig)
		},
	}
	return &http.Client{
		Transport: roundTripper,
	}
}
