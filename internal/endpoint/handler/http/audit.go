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

	if request.TLS != nil {
		ev.TLS = &audit.TLSDetails{
			Version:     audit.TLSVersionToEntity(request.TLS.Version).String(),
			CipherSuite: tls.CipherSuiteName(request.TLS.CipherSuite),
			ServerName:  request.TLS.ServerName,
		}
	}

	ev.SetDestinationIPFromAddr(localAddr(request.Context()))
	ev.SetSourceIPFromAddr(remoteAddr(request.Context()))

	return ev
}
