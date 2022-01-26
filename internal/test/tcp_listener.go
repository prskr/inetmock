package test

import (
	"net"
	"testing"

	"github.com/maxatome/go-testdeep/td"

	"gitlab.com/inetmock/inetmock/internal/endpoint"
)

func NewTCPListener(tb testing.TB, rawAddr string) (listener net.Listener) {
	tb.Helper()
	var err error
	listener, err = net.Listen("tcp4", rawAddr)
	if !td.CmpNoError(tb, err) {
		return
	}
	listener = endpoint.AutoLingeringListener{Listener: listener}
	return
}
