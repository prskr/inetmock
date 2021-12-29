package proxy

import (
	"context"
	"crypto/tls"
	"io"
	"net"
	"net/url"

	"go.uber.org/zap"

	"gitlab.com/inetmock/inetmock/pkg/logging"
)

func ListenAndProxy(
	ctx context.Context,
	listener net.Listener,
	upstream *url.URL,
	tlsCfg *tls.Config,
	logger logging.Logger,
) {
	listener = tls.NewListener(listener, tlsCfg)
	for ctx.Err() == nil {
		clientConn, err := listener.Accept()
		if err != nil {
			return
		}

		go handleClientConn(clientConn, upstream, logger)
	}
}

func handleClientConn(clientConn net.Conn, upstream *url.URL, logger logging.Logger) {
	defer clientConn.Close()
	upstreamConn, err := dial(upstream)
	if err != nil {
		return
	}

	defer upstreamConn.Close()

	go func() {
		if _, err := io.Copy(clientConn, upstreamConn); err != nil {
			logger.Error("Error on proxying data to upstream", zap.Error(err))
		}
	}()

	if _, err := io.Copy(upstreamConn, clientConn); err != nil {
		logger.Error("Error on proxying data back to client", zap.Error(err))
	}
}

func dial(u *url.URL) (net.Conn, error) {
	switch u.Scheme {
	case "unix":
		return net.Dial(u.Scheme, u.Path)
	default:
		return net.Dial(u.Scheme, u.Host)
	}
}
