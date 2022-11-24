package audit

import (
	"crypto/tls"
	"net/http"

	auditv1 "inetmock.icb4dc0.de/inetmock/pkg/audit/v1"
)

func EmittingHandler(emitter Emitter, app auditv1.AppProtocol, delegate http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		httpDetails := &HTTP{
			Method:  req.Method,
			Host:    req.Host,
			URI:     req.RequestURI,
			Proto:   req.Proto,
			Headers: req.Header,
		}

		builder := emitter.Builder().
			WithTransport(auditv1.TransportProtocol_TRANSPORT_PROTOCOL_TCP).
			WithApplication(app).
			WithProtocolDetails(httpDetails)

		if state, ok := TLSConnectionState(req.Context()); ok {
			builder = builder.WithTLSDetails(&TLSDetails{
				Version:     TLSVersionToEntity(state.Version).String(),
				CipherSuite: tls.CipherSuiteName(state.CipherSuite),
				ServerName:  state.ServerName,
			})
		}

		// it's considered to be okay if these details are missing
		builder, _ = builder.WithSourceFromAddr(RemoteAddr(req.Context()))
		builder, _ = builder.WithDestinationFromAddr(LocalAddr(req.Context()))

		builder.Emit()

		delegate.ServeHTTP(writer, req)
	})
}
