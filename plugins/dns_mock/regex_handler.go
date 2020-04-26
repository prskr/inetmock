package main

import (
	"github.com/baez90/inetmock/pkg/logging"
	"github.com/miekg/dns"
	"go.uber.org/zap"
)

type regexHandler struct {
	routes   []resolverRule
	fallback ResolverFallback
	logger   logging.Logger
}

func (r2 *regexHandler) AddRule(rule resolverRule) {
	r2.routes = append(r2.routes, rule)
}

func (r2 regexHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.Compress = false
	m.SetReply(r)

	switch r.Opcode {
	case dns.OpcodeQuery:
		r2.handleQuery(m)
	}
	if err := w.WriteMsg(m); err != nil {
		r2.logger.Error(
			"Failed to write DNS response message",
			zap.Error(err),
		)
	}
}

func (r2 regexHandler) handleQuery(m *dns.Msg) {
	for _, q := range m.Question {
		r2.logger.Info(
			"handling question",
			zap.String("question", q.Name),
		)
		switch q.Qtype {
		case dns.TypeA:
			for _, rule := range r2.routes {
				if rule.pattern.MatchString(q.Name) {
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
					r2.logger.Info(
						"matched DNS rule",
						zap.String("pattern", rule.pattern.String()),
						zap.String("response", rule.response.String()),
					)
					return
				}
			}
			r2.handleFallbackForMessage(m, q)
		}
	}
}

func (r2 regexHandler) handleFallbackForMessage(m *dns.Msg, q dns.Question) {
	fallbackIP := r2.fallback.GetIP()
	answer := &dns.A{
		Hdr: dns.RR_Header{
			Name:   q.Name,
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
			Ttl:    60,
		},
		A: fallbackIP,
	}
	r2.logger.Info(
		"Falling back to generated IP",
		zap.String("response", fallbackIP.String()),
	)
	m.Authoritative = true
	m.Answer = append(m.Answer, answer)
}
