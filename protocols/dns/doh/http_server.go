package doh

import (
	"context"
	"encoding/base64"
	"errors"
	"io"
	"net"
	"net/http"

	mdns "github.com/miekg/dns"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/logging"
	"gitlab.com/inetmock/inetmock/protocols/dns"
)

type Server struct {
	Handler dns.Handler
	server  *http.Server
}

func (s Server) Serve(listener net.Listener) error {
	return s.server.Serve(listener)
}

func (s Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s Server) Close() error {
	return s.server.Close()
}

func NewServer(handler http.Handler) *Server {
	mux := http.NewServeMux()
	mux.Handle("/dns-query", handler)
	return &Server{
		server: &http.Server{
			Handler:     h2c.NewHandler(mux, new(http2.Server)),
			ConnContext: audit.StoreConnPropertiesInContext,
		},
	}
}

func DNSQueryHandler(logger logging.Logger, handler dns.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		msg, err := getMsgFromRequest(request)
		if err != nil {
			logger.Error("Failed to get request from request", zap.Error(err))
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		resp := new(mdns.Msg)
		resp = resp.SetReply(msg)

		for idx := range msg.Question {
			question := msg.Question[idx]
			var rr dns.ResourceRecord
			if rr, err = handler.AnswerDNSQuestion(dns.Question(question)); !errors.Is(err, nil) {
				logger.Error("Error occurred while answering DNS question", zap.Error(err))
			} else {
				resp.Answer = append(resp.Answer, rr)
			}
		}

		var respData []byte
		if respData, err = resp.Pack(); err != nil {
			logger.Error("Failed to pack response message", zap.Error(err))
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		if _, err = writer.Write(respData); err != nil {
			logger.Error("Failed to write response", zap.Error(err))
		}
	})
}

func getMsgFromRequest(request *http.Request) (*mdns.Msg, error) {
	msg := new(mdns.Msg)
	switch request.Method {
	case http.MethodGet:
		if payload, err := base64.URLEncoding.DecodeString(request.URL.Query().Get("dns")); err != nil {
			return nil, err
		} else if err := msg.Unpack(payload); err != nil {
			return nil, err
		}
		return msg, nil
	case http.MethodPost:
		if payload, err := io.ReadAll(request.Body); err != nil {
			return nil, err
		} else if err := msg.Unpack(payload); err != nil {
			return nil, err
		}
		return msg, nil
	default:
		return nil, errors.New("unsupported HTTP method")
	}
}
