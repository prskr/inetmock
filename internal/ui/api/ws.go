//go:build !js && !wasm
// +build !js,!wasm

package api

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/url"

	"github.com/tarndt/wasmws"
	"nhooyr.io/websocket"

	"gitlab.com/inetmock/inetmock/internal/ui/proxy"
	"gitlab.com/inetmock/inetmock/pkg/logging"
)

func RegisterGRPCWebSocketProxy(
	ctx context.Context,
	mux *http.ServeMux,
	upstream *url.URL,
	tlsCfg *tls.Config,
	logger logging.Logger,
) {
	listener := wasmws.NewWebSocketListener(ctx, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
		CompressionMode:    websocket.CompressionDisabled,
	})
	mux.Handle("/grpc-proxy", listener)
	go proxy.ListenAndProxy(ctx, listener, upstream, tlsCfg, logger)
}
