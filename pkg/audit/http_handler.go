package audit

import (
	"crypto/tls"
	"net/http"

	"gitlab.com/inetmock/inetmock/pkg/audit/details"
	auditv1 "gitlab.com/inetmock/inetmock/pkg/audit/v1"
)

func EmittingHandler(emitter Emitter, app auditv1.AppProtocol, delegate http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		httpDetails := details.HTTP{
			Method:  req.Method,
			Host:    req.Host,
			URI:     req.RequestURI,
			Proto:   req.Proto,
			Headers: req.Header,
		}

		ev := Event{
			Transport:       auditv1.TransportProtocol_TRANSPORT_PROTOCOL_TCP,
			Application:     app,
			ProtocolDetails: httpDetails,
		}

		if state, ok := TLSConnectionState(req.Context()); ok {
			ev.TLS = &TLSDetails{
				Version:     TLSVersionToEntity(state.Version).String(),
				CipherSuite: tls.CipherSuiteName(state.CipherSuite),
				ServerName:  state.ServerName,
			}
		}

		// it's considered to be okay if these details are missing
		_ = ev.SetDestinationIPFromAddr(LocalAddr(req.Context()))
		_ = ev.SetSourceIPFromAddr(RemoteAddr(req.Context()))

		emitter.Emit(ev)

		delegate.ServeHTTP(writer, req)
	})
}
