package dns

import (
	"time"

	mdns "github.com/miekg/dns"

	"gitlab.com/inetmock/inetmock/internal/rules"
)

type RuleHandler struct {
	resolvers []ConditionalResolver
	TTL       time.Duration
}

func (r RuleHandler) AnswerDNSQuestion(q Question) (ResourceRecord, error) {
	for idx := range r.resolvers {
		if res := r.resolvers[idx]; res.Matches(q) {
			resolvedIP := res.Lookup(q.Name)
			switch q.Qtype {
			case mdns.TypeA:
				return &mdns.A{
					A:   resolvedIP,
					Hdr: RRHeader(r.TTL, q),
				}, nil
			case mdns.TypeAAAA:
				return &mdns.AAAA{
					AAAA: resolvedIP,
					Hdr:  RRHeader(r.TTL, q),
				}, nil
			}
		}
	}

	return nil, ErrNoAnswerForQuestion
}

func (r *RuleHandler) RegisterRule(rawRule string) error {
	rule := new(rules.SingleResponsePipeline)
	if err := rules.Parse(rawRule, rule); err != nil {
		return err
	}

	var (
		conditionalResolver ConditionalResolver
		err                 error
	)

	if conditionalResolver.Predicates, err = QuestionPredicatesForRoutingRule(rule); err != nil {
		return err
	}

	if conditionalResolver.IPResolver, err = ResolverForRule(rule); err != nil {
		return err
	}

	r.resolvers = append(r.resolvers, conditionalResolver)
	return nil
}
