package mock

import (
	"errors"
	"net"
	"sync"

	mdns "github.com/miekg/dns"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"

	"gitlab.com/inetmock/inetmock/pkg/audit"
	auditv1 "gitlab.com/inetmock/inetmock/pkg/audit/v1"
	"gitlab.com/inetmock/inetmock/pkg/logging"
	"gitlab.com/inetmock/inetmock/pkg/metrics"
	"gitlab.com/inetmock/inetmock/protocols/dns"
)

const name = "dns_mock"

var (
	handlerNameLblName          = "handler_name"
	totalHandledRequestsCounter *prometheus.CounterVec
	unhandledRequestsCounter    *prometheus.CounterVec
	requestDurationHistogram    *prometheus.HistogramVec
	initLock                    sync.Locker = new(sync.Mutex)
)

func init() {
	initLock.Lock()
	defer initLock.Unlock()

	var err error
	if totalHandledRequestsCounter == nil {
		if totalHandledRequestsCounter, err = metrics.Counter(
			name,
			"handled_requests_total",
			"",
			handlerNameLblName,
		); err != nil {
			panic(err)
		}
	}

	if unhandledRequestsCounter == nil {
		if unhandledRequestsCounter, err = metrics.Counter(
			name,
			"unhandled_requests_total",
			"",
			handlerNameLblName,
		); err != nil {
			panic(err)
		}
	}

	if requestDurationHistogram == nil {
		if requestDurationHistogram, err = metrics.Histogram(
			name,
			"request_duration",
			"",
			nil,
			handlerNameLblName,
		); err != nil {
			panic(err)
		}
	}
}

type Server struct {
	Name    string
	Handler dns.Handler
	Logger  logging.Logger
	Emitter audit.Emitter
}

func (s *Server) ServeDNS(w mdns.ResponseWriter, req *mdns.Msg) {
	if requestDurationHistogram != nil {
		timer := prometheus.NewTimer(requestDurationHistogram.WithLabelValues(s.Name))
		defer timer.ObserveDuration()
	}

	s.recordRequest(req, w.LocalAddr(), w.RemoteAddr())

	resp := new(mdns.Msg)
	resp = resp.SetReply(req)

	for qIdx := range req.Question {
		question := req.Question[qIdx]
		if rr, err := s.Handler.AnswerDNSQuestion(dns.Question(question)); !errors.Is(err, nil) {
			s.Logger.Error("Error occurred while answering DNS question", zap.Error(err))
		} else {
			resp.Answer = append(resp.Answer, rr)
		}
	}

	if err := w.WriteMsg(resp); err != nil {
		s.Logger.Error("Failed to write response", zap.Error(err))
	}
}

func (s *Server) recordRequest(m *mdns.Msg, localAddr, remoteAddr net.Addr) {
	dnsDetails := &audit.DNS{
		OPCode: auditv1.DNSOpCode(m.Opcode),
	}

	for _, q := range m.Question {
		dnsDetails.Questions = append(dnsDetails.Questions, audit.DNSQuestion{
			RRType: auditv1.ResourceRecordType(q.Qtype),
			Name:   q.Name,
		})
	}

	builder := s.Emitter.Builder().
		WithTransport(guessTransportFromAddr(localAddr)).
		WithApplication(auditv1.AppProtocol_APP_PROTOCOL_DNS).
		WithProtocolDetails(dnsDetails)

	// it's considered to be okay if these details are missing
	builder, _ = builder.WithSourceFromAddr(remoteAddr)
	builder, _ = builder.WithDestinationFromAddr(localAddr)

	builder.Emit()
}

func guessTransportFromAddr(addr net.Addr) auditv1.TransportProtocol {
	switch addr.(type) {
	case *net.TCPAddr:
		return auditv1.TransportProtocol_TRANSPORT_PROTOCOL_TCP
	case *net.UDPAddr:
		return auditv1.TransportProtocol_TRANSPORT_PROTOCOL_UDP
	default:
		return auditv1.TransportProtocol_TRANSPORT_PROTOCOL_UNSPECIFIED
	}
}
