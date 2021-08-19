package mock

import (
	"math"
	"net"
	"time"

	mdns "github.com/miekg/dns"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"

	"gitlab.com/inetmock/inetmock/internal/endpoint/handler/dns"
	"gitlab.com/inetmock/inetmock/internal/rules"
	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/audit/details"
	v1 "gitlab.com/inetmock/inetmock/pkg/audit/v1"
	"gitlab.com/inetmock/inetmock/pkg/logging"
)

type RuleHandler struct {
	handlers    []ConditionHandler
	Cache       Cache
	TTL         time.Duration
	HandlerName string
	Logger      logging.Logger
	Emitter     audit.Emitter
	Fallback    dns.IPResolver
}

func (r *RuleHandler) RegisterRule(rawRule string) error {
	r.Logger.Debug("Adding routing rule", zap.String("rawRule", rawRule))
	rule := new(rules.Routing)
	if err := rules.Parse(rawRule, rule); err != nil {
		return err
	}

	var (
		conditionalHandler ConditionHandler
		err                error
	)

	if conditionalHandler.Filters, err = RequestFiltersForRoutingRule(rule); err != nil {
		return err
	}

	if conditionalHandler.IPResolver, err = HandlerForRoutingRule(rule); err != nil {
		return err
	}

	r.handlers = append(r.handlers, conditionalHandler)
	return nil
}

func (r *RuleHandler) ServeDNS(w mdns.ResponseWriter, req *mdns.Msg) {
	if requestDurationHistogram != nil {
		timer := prometheus.NewTimer(requestDurationHistogram.WithLabelValues(r.HandlerName))
		defer timer.ObserveDuration()
	}

	r.recordRequest(req, w.LocalAddr(), w.RemoteAddr())

	var (
		resp             = new(mdns.Msg)
		seconds          = r.ttlSeconds()
		matchedQuestions = make(map[int]bool)
	)
	resp = resp.SetReply(req)

questionLoop:
	for qIdx := range req.Question {
		var (
			question = req.Question[qIdx]
			ip       net.IP
		)
		switch question.Qtype {
		case mdns.TypeA, mdns.TypeAAAA:
			if ip = r.Cache.ForwardLookup(question.Name); ip != nil {
				ip = formatIP(ip, question.Qtype)
				matchedQuestions[qIdx] = true
				addARecordAnswer(resp, question, ip, question.Name, seconds)
				continue
			}
		case mdns.TypePTR:
			ip = dns.ParseInAddrArpa(question.Name)
			if host, miss := r.Cache.ReverseLookup(ip); !miss {
				matchedQuestions[qIdx] = true
				addPTRRecordAnswer(resp, question, host, seconds)
				continue
			}
		}

		for idx := range r.handlers {
			handler := r.handlers[idx]
			if handler.Matches(&question) {
				matchedQuestions[qIdx] = true
				ip = formatIP(handler.Lookup(question.Name), question.Qtype)
				r.Cache.PutRecord(question.Name, ip)
				addARecordAnswer(resp, question, ip, question.Name, seconds)
				continue questionLoop
			}
		}
	}

	for qIdx := range req.Question {
		q := req.Question[qIdx]

		if matchedQuestions[qIdx] {
			continue
		}

		ip := formatIP(r.Fallback.Lookup(q.Name), q.Qtype)
		r.Cache.PutRecord(q.Name, ip)
		addARecordAnswer(resp, q, ip, q.Name, seconds)
	}

	if err := w.WriteMsg(resp); err != nil {
		r.Logger.Error("Failed to write response", zap.Error(err))
	}
}

func (r *RuleHandler) recordRequest(m *mdns.Msg, localAddr, remoteAddr net.Addr) {
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

	r.Emitter.Emit(ev)
}

func (r RuleHandler) ttlSeconds() uint32 {
	const minTTLSeconds = 5.0
	return uint32(math.Max(minTTLSeconds, r.TTL.Seconds()))
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

func addPTRRecordAnswer(msg *mdns.Msg, q mdns.Question, host string, ttl uint32) {
	msg.Answer = append(msg.Answer, &mdns.PTR{
		Ptr: host,
		Hdr: mdns.RR_Header{
			Name:   q.Name,
			Class:  mdns.ClassINET,
			Rrtype: q.Qtype,
			Ttl:    ttl,
		},
	})
}

func addARecordAnswer(msg *mdns.Msg, q mdns.Question, ip net.IP, host string, ttl uint32) {
	msg.Answer = append(msg.Answer, &mdns.A{
		A: ip,
		Hdr: mdns.RR_Header{
			Name:   host,
			Class:  mdns.ClassINET,
			Rrtype: q.Qtype,
			Ttl:    ttl,
		},
	})
}

func formatIP(ip net.IP, qType uint16) net.IP {
	switch qType {
	case mdns.TypeA:
		return ip.To4()
	case mdns.TypeAAAA:
		return ip.To16()
	default:
		return ip
	}
}
