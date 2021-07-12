package mock

import (
	"net"

	"github.com/miekg/dns"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"

	"gitlab.com/inetmock/inetmock/internal/rules"
	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/audit/details"
	v1 "gitlab.com/inetmock/inetmock/pkg/audit/v1"
	"gitlab.com/inetmock/inetmock/pkg/logging"
)

type regexHandler struct {
	handlerName  string
	auditEmitter audit.Emitter
	logger       logging.Logger
}

func (rh *regexHandler) AddRule(rawRule string) error {
	var (
		parsedRule = new(rules.Routing)
	)

	if err := rules.Parse(rawRule, parsedRule); err != nil {
		return err
	}

	return nil
}

func (rh *regexHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	timer := prometheus.NewTimer(requestDurationHistogram.WithLabelValues(rh.handlerName))
	defer func() {
		timer.ObserveDuration()
	}()

	rh.recordRequest(r, w.LocalAddr(), w.RemoteAddr())

	m := new(dns.Msg)
	m.Compress = false
	m.SetReply(r)

	if r.Opcode == dns.OpcodeQuery {
		rh.handleQuery(m)
	}
	if err := w.WriteMsg(m); err != nil {
		rh.logger.Error(
			"Failed to write DNS response message",
			zap.Error(err),
		)
	}
}

func (rh *regexHandler) handleQuery(m *dns.Msg) {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA:
			totalHandledRequestsCounter.WithLabelValues(rh.handlerName).Inc()
			/*for _, rule := range rh.routes {
				if !rule.pattern.MatchString(q.Name) {
					continue
				}
				m.Authoritative = true
				answer := &dns.A{
					Hdr: dns.RR_Header{
						Name:   q.Name,
						Rrtype: dns.TypeA,
						Class:  dns.ClassINET,
						Ttl:    60,
					},
					A: rule.response,
				}
				m.Answer = append(m.Answer, answer)
				rh.logger.Info(
					"matched DNS rule",
					zap.String("pattern", rule.pattern.String()),
					zap.String("response", rule.response.String()),
				)
				return
			}*/
			rh.handleFallbackForMessage(m, q)
		default:
			unhandledRequestsCounter.WithLabelValues(rh.handlerName).Inc()
			rh.logger.Warn(
				"Unhandled DNS question type - no response will be sent",
				zap.Uint16("question_type", q.Qtype),
			)
		}
	}
}

func (rh *regexHandler) handleFallbackForMessage(m *dns.Msg, q dns.Question) {
	// nolint:gomnd
	fallbackIP := net.IPv4(127, 0, 0, 1)
	answer := &dns.A{
		Hdr: dns.RR_Header{
			Name:   q.Name,
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
			Ttl:    60,
		},
		A: fallbackIP,
	}
	rh.logger.Info(
		"Falling back to generated IP",
		zap.String("response", fallbackIP.String()),
	)
	m.Authoritative = true
	m.Answer = append(m.Answer, answer)
}

func (rh *regexHandler) recordRequest(m *dns.Msg, localAddr, remoteAddr net.Addr) {
	dnsDetails := &details.DNS{
		OPCode: v1.DNSOpCode(m.Opcode),
	}

	for _, q := range m.Question {
		dnsDetails.Questions = append(dnsDetails.Questions, details.DNSQuestion{
			RRType: v1.ResourceRecordType(q.Qtype),
			Name:   q.Name,
		})
	}

	ev := audit.Event{
		Transport:       guessTransportFromAddr(localAddr),
		Application:     v1.AppProtocol_APP_PROTOCOL_DNS,
		ProtocolDetails: dnsDetails,
	}

	// it's considered to be okay if these details are missing
	_ = ev.SetSourceIPFromAddr(remoteAddr)
	_ = ev.SetDestinationIPFromAddr(localAddr)

	rh.auditEmitter.Emit(ev)
}

func guessTransportFromAddr(addr net.Addr) v1.TransportProtocol {
	switch addr.(type) {
	case *net.TCPAddr:
		return v1.TransportProtocol_TRANSPORT_PROTOCOL_TCP
	case *net.UDPAddr:
		return v1.TransportProtocol_TRANSPORT_PROTOCOL_UDP
	default:
		return v1.TransportProtocol_TRANSPORT_PROTOCOL_UNSPECIFIED
	}
}
