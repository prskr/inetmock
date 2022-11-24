package mock

import (
	"errors"
	"net"
	"sync"

	mdns "github.com/miekg/dns"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"

	"inetmock.icb4dc0.de/inetmock/pkg/audit"
	auditv1 "inetmock.icb4dc0.de/inetmock/pkg/audit/v1"
	"inetmock.icb4dc0.de/inetmock/pkg/logging"
	"inetmock.icb4dc0.de/inetmock/pkg/metrics"
	"inetmock.icb4dc0.de/inetmock/protocols"
	"inetmock.icb4dc0.de/inetmock/protocols/dns"
)

const name = "dns_mock"

var (
	handlerNameLblName             = "handler_name"
	totalProcessedQuestionsCounter *prometheus.CounterVec
	initLock                       sync.Locker = new(sync.Mutex)
)

func init() {
	initLock.Lock()
	defer initLock.Unlock()

	var err error
	if totalProcessedQuestionsCounter == nil {
		if totalProcessedQuestionsCounter, err = metrics.Counter(
			name,
			"handled_questions_total",
			"",
			handlerNameLblName,
			"answered",
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
	defer prometheus.NewTimer(protocols.RequestDurationHistogram.WithLabelValues("dns", s.Name)).ObserveDuration()

	s.recordRequest(req, w.LocalAddr(), w.RemoteAddr())

	resp := new(mdns.Msg)
	resp = resp.SetReply(req)

	for qIdx := range req.Question {
		question := req.Question[qIdx]
		if rr, err := s.Handler.AnswerDNSQuestion(dns.Question(question)); !errors.Is(err, nil) {
			if errors.Is(err, dns.ErrNoAnswerForQuestion) {
				totalProcessedQuestionsCounter.WithLabelValues(s.Name, "false")
			}
			s.Logger.Error("Error occurred while answering DNS question", zap.String("question", question.Name), zap.Error(err))
		} else {
			totalProcessedQuestionsCounter.WithLabelValues(s.Name, "true")
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
