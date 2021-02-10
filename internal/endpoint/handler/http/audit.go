package http

import (
	"crypto/tls"
	"net/http"

	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/audit/details"
)

func EventFromRequest(request *http.Request, app audit.AppProtocol) audit.Event {
	httpDetails := details.HTTP{
		Method:  request.Method,
		Host:    request.Host,
		URI:     request.RequestURI,
		Proto:   request.Proto,
		Headers: request.Header,
	}

	ev := audit.Event{
		Transport:       audit.TransportProtocol_TCP,
		Application:     app,
		ProtocolDetails: httpDetails,
	}

	if state, ok := tlsConnectionState(request.Context()); ok {
		ev.TLS = &audit.TLSDetails{
			Version:     audit.TLSVersionToEntity(state.Version).String(),
			CipherSuite: tls.CipherSuiteName(state.CipherSuite),
			ServerName:  state.ServerName,
		}
	}

	ev.SetDestinationIPFromAddr(localAddr(request.Context()))
	ev.SetSourceIPFromAddr(remoteAddr(request.Context()))

	return ev
}
