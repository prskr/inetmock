package http

import (
	"crypto/tls"
	"net/http"

	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/audit/details"
	v1 "gitlab.com/inetmock/inetmock/pkg/audit/v1"
)

func EventFromRequest(request *http.Request, app v1.AppProtocol) audit.Event {
	httpDetails := details.HTTP{
		Method:  request.Method,
		Host:    request.Host,
		URI:     request.RequestURI,
		Proto:   request.Proto,
		Headers: request.Header,
	}

	ev := audit.Event{
		Transport:       v1.TransportProtocol_TRANSPORT_PROTOCOL_TCP,
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

	// it's considered to be okay if these details are missing
	_ = ev.SetDestinationIPFromAddr(localAddr(request.Context()))
	_ = ev.SetSourceIPFromAddr(remoteAddr(request.Context()))

	return ev
}
